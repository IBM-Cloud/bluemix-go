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

var _ = Describe("Logging", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//ListLoggingInstances
	Describe("ListLoggingInstances", func() {
		Context("When read of logging instances is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, ContainSubstring("/v2/observe/logging/getConfigs")),
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

			It("should return loggin isnatnces list", func() {
				target := LoggingTargetHeader{}
				myLogging, err := newLogging(server.URL()).ListLoggingInstances("clusterName", target)
				Expect(myLogging).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
				Expect(len(myLogging)).Should(Equal(1))
				Expect(myLogging[0].InstanceID).Should(Equal("ec4f0886-edc4-409e-8720-574035538f91"))
				Expect(myLogging[0].CRN).Should(Equal("crn:v1:bluemix:public:sysdig-monitor:us-south:a/fcdb764102154c7ea8e1b79d3a64afe0:ec4f0886-edc4-409e-8720-574035538f91::"))
				Expect(myLogging[0].PrivateEndpoint).Should(Equal(true))
			})
		})

		Context("When read of logging instances is unsuccessful", func() {
			BeforeEach(func() {

				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, ContainSubstring("/v2/")),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve logging istances`),
					),
				)
			})

			It("should return error when ogging istances are retrieved", func() {
				target := LoggingTargetHeader{}
				myLogging, err := newLogging(server.URL()).ListLoggingInstances("DragonBoat-cluster", target)
				Expect(err).To(HaveOccurred())
				Expect(myLogging).Should(BeNil())
			})
		})
	})

	//Create
	Describe("CreateLoggingConfig", func() {
		Context("When creation is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/observe/logging/createConfig"),
						ghttp.VerifyJSON(`{"cluster": "DragonBoat-cluster", "instance": "ec4f0886-edc4-409e-8720-574035538f91"}`),
						ghttp.RespondWith(http.StatusCreated, `{
							 "instanceId": "ec4f0886-edc4-409e-8720-574035538f91"
						}`),
					),
				)
			})

			It("should return Logging instances created", func() {
				params := LoggingCreateRequest{
					Cluster: "DragonBoat-cluster", LoggingInstance: "ec4f0886-edc4-409e-8720-574035538f91",
				}
				target := LoggingTargetHeader{}
				myLogging, err := newLogging(server.URL()).CreateLoggingConfig(params, target)
				Expect(err).NotTo(HaveOccurred())
				Expect(myLogging).ShouldNot(BeNil())
				Expect(myLogging.InstanceID).Should(Equal("ec4f0886-edc4-409e-8720-574035538f91"))
			})
		})
		Context("When creation is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/observe/logging/createConfig"),
						ghttp.VerifyJSON(`{"cluster": "DragonBoat-cluster", "instance": "ec4f0886-edc4-409e-8720-574035538f91"}`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to create sysdig logging instance`),
					),
				)
			})
			It("should return error during logging instance creation", func() {

				params := LoggingCreateRequest{
					Cluster: "DragonBoat-cluster", LoggingInstance: "ec4f0886-edc4-409e-8720-574035538f91",
				}
				target := LoggingTargetHeader{}
				_, err := newLogging(server.URL()).CreateLoggingConfig(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

func newLogging(url string) Logging {

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
	return newLoggingAPI(&client)
}
