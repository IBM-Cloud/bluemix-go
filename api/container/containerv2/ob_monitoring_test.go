package containerv2

import (
	"log"
	"net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/client"
	bluemixHttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/session"

	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Monitoring", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//ListAllMonitors
	Describe("ListAllMonitors", func() {
		Context("When read of monitors is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, ContainSubstring("v2/observe/monitoring/getConfigs")),
						ghttp.RespondWith(http.StatusOK, `[{
              "AgentKey": "",
			  "AgentNamespace": "",
			  "CRN": "crn:v1:bluemix:public:sysdig-monitor:us-south:a/fcdb764102154c7ea8e1b79d3a64afe0:ec4f0886-edc4-409e-8720-574035538f91::",
              "DaemonsetName": "",
              "DiscoveredAgent": true,
              "InstanceID": "ec4f0886-edc4-409e-8720-574035538f91",
              "InstanceName": "ns",
              "Namespace": "",
              "PrivateEndpoint": true
              }]`),
					),
				)
			})

			It("should return sysdig monitor list", func() {
				target := MonitoringTargetHeader{}
				myMonitor, err := newMonitoring(server.URL()).ListAllMonitors("clusterName", target)
				Expect(myMonitor).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
				Expect(len(myMonitor)).Should(Equal(1))
				Expect(myMonitor[0].InstanceID).Should(Equal("ec4f0886-edc4-409e-8720-574035538f91"))
				Expect(myMonitor[0].CRN).Should(Equal("crn:v1:bluemix:public:sysdig-monitor:us-south:a/fcdb764102154c7ea8e1b79d3a64afe0:ec4f0886-edc4-409e-8720-574035538f91::"))
				Expect(myMonitor[0].PrivateEndpoint).Should(Equal(true))
			})
		})

		Context("When read of monitors is unsuccessful", func() {
			BeforeEach(func() {

				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, ContainSubstring("/v2/")),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve monitors`),
					),
				)
			})

			It("should return error when monitors are retrieved", func() {
				target := MonitoringTargetHeader{}
				myMonitor, err := newMonitoring(server.URL()).ListAllMonitors("DragonBoat-cluster", target)
				Expect(err).To(HaveOccurred())
				Expect(myMonitor).Should(BeNil())
			})
		})
	})

	//Create
	Describe("CreateMonitorConfig", func() {
		Context("When creation is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/observe/monitoring/createConfig"),
						ghttp.VerifyJSON(`{"cluster": "DragonBoat-cluster", "instance": "ec4f0886-edc4-409e-8720-574035538f91"}`),
						ghttp.RespondWith(http.StatusCreated, `{
							 "instanceId": "ec4f0886-edc4-409e-8720-574035538f91"
						}`),
					),
				)
			})

			It("should return monitor created", func() {
				params := MonitoringCreateRequest{
					Cluster: "DragonBoat-cluster", SysidigInstance: "ec4f0886-edc4-409e-8720-574035538f91",
				}
				target := MonitoringTargetHeader{}
				myMonitor, err := newMonitoring(server.URL()).CreateMonitoringConfig(params, target)
				Expect(err).NotTo(HaveOccurred())
				Expect(myMonitor).ShouldNot(BeNil())
				Expect(myMonitor.InstanceID).Should(Equal("ec4f0886-edc4-409e-8720-574035538f91"))
			})
		})
		Context("When creation is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/observe/monitoring/createConfig"),
						ghttp.VerifyJSON(`{"cluster": "DragonBoat-cluster",  "instance": "ec4f0886-edc4-409e-8720-574035538f91"}`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to create sysdig monitor`),
					),
				)
			})
			It("should return error during monitor creation", func() {

				params := MonitoringCreateRequest{
					Cluster: "DragonBoat-cluster", SysidigInstance: "ec4f0886-edc4-409e-8720-574035538f91",
				}
				target := MonitoringTargetHeader{}
				_, err := newMonitoring(server.URL()).CreateMonitoringConfig(params, target)
				Expect(err).To(HaveOccurred())
				//Expect(myMonitor.InstanceID).Should(Equal("ec4f0886-edc4-409e-8720-574035538f91"))
			})
		})
	})
})

func newMonitoring(url string) Monitoring {

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.VpcContainerService,
	}
	return newMonitoringAPI(&client)
}
