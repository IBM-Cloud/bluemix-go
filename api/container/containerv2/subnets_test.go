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

var _ = Describe("Subnets", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//ListSubnets
	Describe("List", func() {
		Context("When List subnets is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getSubnets"),
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

			It("should list subnets in a cluster", func() {
				target := ClusterTargetHeader{}

				_, err := newSubnets(server.URL()).ListSubnets("aaa", "bbb", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When list subnets is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getSubnets"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to list subnets`),
					),
				)
			})

			It("should return error during get subnets", func() {
				target := ClusterTargetHeader{}
				_, err := newSubnets(server.URL()).ListSubnets("aaa", "bbb", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})

func newSubnets(url string) Subnets {
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
	return newSubnetsAPI(&client)
}
