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

var _ = Describe("VPCs", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//ListVPCs
	Describe("List", func() {
		Context("When List VPCs is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getVPCs"),
						ghttp.RespondWith(http.StatusCreated, `[
							{
							  "availableIPv4AddressCount": 0,
							  "id": "string",
							  "ipv4CIDRBlock": "string",
							  "name": "string",
							  "publicGatewayID": "string",
							  "publicGatewayName": "string",
							  "vpcID": "string",
							  "vpcName": "string",
							  "zone": "string"
							}
						  ]`),
					),
				)
			})

			It("should list VPCs in a cluster", func() {
				target := ClusterTargetHeader{}

				_, err := newVPCs(server.URL()).ListVPCs(target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When list VPCs is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getVPCs"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to list VPCs`),
					),
				)
			})

			It("should return error during get VPCs", func() {
				target := ClusterTargetHeader{}
				_, err := newVPCs(server.URL()).ListVPCs(target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})

func newVPCs(url string) VPCs {
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
	return newVPCsAPI(&client)
}
