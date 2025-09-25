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

var _ = Describe("dedicatedhostflavor", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	Describe("List", func() {
		Context("When list dedicatedhostflavors is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/getDedicatedHostFlavors"),
						ghttp.RespondWith(http.StatusCreated, `[
							{
								"deprecated":false,
								"flavorClass":"flavorclass1",
								"id":"flavorid1",
								"instanceStorage":[
									{
										"count": 0,
										"size":	0
									}
								],
								"maxMemory": 52,
								"maxVCPUs":	12,
								"region": "region1",
								"zone": "zone1"
							  }
						  ]`),
					),
				)
			})

			It("should list dedicatedhostflavors in a zone", func() {
				target := ClusterTargetHeader{}

				ldhf, err := newDedicatedHostFlavor(server.URL()).ListDedicatedHostFlavors("zone1", target)
				Expect(err).NotTo(HaveOccurred())
				expectedDedicatedHostFlavors := GetDedicatedHostFlavors{
					GetDedicatedHostFlavor{
						ID:          "flavorid1",
						FlavorClass: "flavorclass1",
						Region:      "region1",
						Zone:        "zone1",
						Deprecated:  false,
						MaxVCPUs:    12,
						MaxMemory:   52,
						InstanceStorage: []InstanceStorage{
							InstanceStorage{
								Count: 0,
								Size:  0,
							},
						},
					},
				}
				Expect(ldhf).To(BeEquivalentTo(expectedDedicatedHostFlavors))
			})
		})
		Context("When list dedicatedhostflavors is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/getDedicatedHostFlavors"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to list dedicatedhostflavors`),
					),
				)
			})

			It("should return error during get dedicatedhosts", func() {
				target := ClusterTargetHeader{}
				_, err := newDedicatedHostFlavor(server.URL()).ListDedicatedHostFlavors("zone1", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})

func newDedicatedHostFlavor(url string) DedicatedHostFlavor {

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
	return newDedicatedHostFlavorAPI(&client)
}
