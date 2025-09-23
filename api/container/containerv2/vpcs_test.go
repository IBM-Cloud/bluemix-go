package containerv2

import (
	"log"
	"net/http"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/client"
	bluemixHttp "github.com/Mavrickk3/bluemix-go/http"
	"github.com/Mavrickk3/bluemix-go/session"

	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VPCs", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	// ListVPCs
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

	// SetOutboundTrafficProtection
	Describe("SetOutboundTrafficProtection", func() {
		Context("When SetOutboundTrafficProtection is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/network/v2/outbound-traffic-protection"),
						ghttp.VerifyJSON(`{"cluster":"testCluster","operation":"enable-outbound-protection"}`),
						ghttp.RespondWith(http.StatusOK, ""),
					),
				)
			})

			It("should return with 200 OK", func() {
				target := ClusterTargetHeader{}

				err := newVPCs(server.URL()).SetOutboundTrafficProtection("testCluster", true, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When SetOutboundTrafficProtection is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/network/v2/outbound-traffic-protection"),
						ghttp.VerifyJSON(`{"cluster":"testCluster","operation":"disable-outbound-protection"}`),
						ghttp.RespondWith(http.StatusInternalServerError, ""),
					),
				)
			})

			It("should return with 500 Internal server error", func() {
				target := ClusterTargetHeader{}
				err := newVPCs(server.URL()).SetOutboundTrafficProtection("testCluster", false, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// Enable secure by default
	Describe("Enable Secure by Default", func() {
		Context("When EnableSecureByDefault is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/network/v2/secure-by-default/enable"),
						ghttp.VerifyJSON(`{"cluster":"testCluster","disableOutboundTrafficProtection":true}`),
						ghttp.RespondWith(http.StatusOK, ""),
					),
				)
			})

			It("should return with 200 OK", func() {
				target := ClusterTargetHeader{}

				err := newVPCs(server.URL()).EnableSecureByDefault("testCluster", true, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When EnableSecureByDefault is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/network/v2/secure-by-default/enable"),
						ghttp.VerifyJSON(`{"cluster":"testCluster"}`),
						ghttp.RespondWith(http.StatusInternalServerError, ""),
					),
				)
			})

			It("should return with 500 Internal server error", func() {
				target := ClusterTargetHeader{}
				err := newVPCs(server.URL()).EnableSecureByDefault("testCluster", false, target)
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
