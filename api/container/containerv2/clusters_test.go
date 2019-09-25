package containerv2

import (
	"log"
	"net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/client"
	bluemixHttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/session"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"

	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clusters", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//List
	Describe("List", func() {
		Context("When read of clusters is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getClusters"),
						ghttp.RespondWith(http.StatusOK, `[{
              "CreatedDate": "",
              "DataCenter": "dal10",
              "ID": "f91adfe2-76c9-4649-939e-b01c37a3704",
              "IngressHostname": "",
              "IngressSecretName": "",
              "Location": "",
              "MasterKubeVersion": "1.8.1",
              "Prefix": "worker",
              "ModifiedDate": "",
              "Name": "test",
              "Region": "abc",
              "ServerURL": "",
              "State": "normal",
              "IsPaid": false,
              "IsTrusted": true,
              "WorkerCount": 1
              }]`),
					),
				)
			})

			It("should return cluster list", func() {
				target := v1.ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				myCluster, err := newCluster(server.URL()).List(target)
				Expect(myCluster).ShouldNot(BeNil())
				for _, cluster := range myCluster {
					Expect(err).NotTo(HaveOccurred())
					Expect(cluster.ID).Should(Equal("f91adfe2-76c9-4649-939e-b01c37a3704"))
					Expect(cluster.WorkerCount).Should(Equal(1))
					Expect(cluster.MasterKubeVersion).Should(Equal("1.8.1"))
				}
			})
		})
		Context("When read of clusters is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getClusters"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve clusters`),
					),
				)
			})

			It("should return error when cluster are retrieved", func() {
				target := v1.ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				myCluster, err := newCluster(server.URL()).List(target)
				Expect(err).To(HaveOccurred())
				Expect(myCluster).Should(BeNil())
			})
		})
	})

	//
})

func newCluster(url string) Clusters {

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.MccpService,
	}
	return newClusterAPI(&client)
}
