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

var _ = Describe("Kms", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//Enable
	Describe("Enable", func() {
		Context("When Enabling kms is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/enableKMS"),
						ghttp.VerifyJSON(`{"cluster":"bs65tjud0j4njc8pu30g","instance_id":"12043812-757f-4e1e-8436-6af3245e6a69","crk_id":"0792853c-b9f9-4b35-9d9e-ffceab51d3c1","private_endpoint":false}`),
						ghttp.RespondWith(http.StatusCreated, `{}`),
					),
				)
			})

			It("should enable Kms in a cluster", func() {
				target := ClusterHeader{}
				params := KmsEnableReq{
					Cluster: "bs65tjud0j4njc8pu30g", Kms: "12043812-757f-4e1e-8436-6af3245e6a69", Crk: "0792853c-b9f9-4b35-9d9e-ffceab51d3c1", PrivateEndpoint: false,
				}
				err := newKms(server.URL()).EnableKms(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When enabling is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/enableKMS"),
						ghttp.VerifyJSON(`{"cluster":"bs65tjud0j4njc8pu30g","instance_id":"12043812-757f-4e1e-8436-6af3245e6a69","crk_id":"0792853c-b9f9-4b35-9d9e-ffceab51d3c1","private_endpoint":false}`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to enable kms`),
					),
				)
			})

			It("should return error during enabling kms", func() {
				params := KmsEnableReq{
					Cluster: "bs65tjud0j4njc8pu30g", Kms: "12043812-757f-4e1e-8436-6af3245e6a69", Crk: "0792853c-b9f9-4b35-9d9e-ffceab51d3c1", PrivateEndpoint: false,
				}
				target := ClusterHeader{}
				err := newKms(server.URL()).EnableKms(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

func newKms(url string) Kms {

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
	return newKmsAPI(&client)
}
