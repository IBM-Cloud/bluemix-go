package k8sclusterv1

import (
	"log"
	"net/http"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/client"
	bluemixHttp "github.com/IBM-Bluemix/bluemix-go/http"
	"github.com/IBM-Bluemix/bluemix-go/session"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Subnets", func() {
	var server *ghttp.Server
	Describe("Add", func() {
		Context("When adding a subnet is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, "/v1/clusters/test/subnets/1109876"),
						ghttp.RespondWith(http.StatusCreated, `{}`),
					),
				)
			})

			It("should return subnet added to cluster", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				err := newSubnet(server.URL()).AddSubnet("test", "1109876", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When adding subnet is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, "/v1/clusters/test/subnets/1109876"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to add subnet to cluster`),
					),
				)
			})

			It("should return error during add subnet to cluster", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				err := newSubnet(server.URL()).AddSubnet("test", "1109876", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})
	//List
	Describe("List", func() {
		Context("When retrieving available subnets is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/subnets"),
						ghttp.RespondWith(http.StatusOK, `[{
						"ID": "535642",
						"Type": "private",     
						"VlanID": "1565297",
						"IPAddresses": ["10.98.25.2","10.98.25.3","10.98.25.4"],
						"Properties": {
						"CIDR": "26 ",
						"NetworkIdentifier":"10.130.229.64",
						"Note": "",
						"SubnetType":"additional_primary",
						"DisplayLabel":"",
						"Gateway":""
						}
						}]`),
					),
				)
			})

			It("should return available subnets ", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				subnets, err := newSubnet(server.URL()).List(target)
				Expect(err).NotTo(HaveOccurred())
				Expect(subnets).ShouldNot(BeNil())
				for _, sObj := range subnets {
					Expect(sObj).ShouldNot(BeNil())
					Expect(sObj.ID).Should(Equal("535642"))
					Expect(sObj.Type).Should(Equal("private"))
				}
			})
		})
		Context("When retrieving available subnets is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/subnets"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve subnets`),
					),
				)
			})

			It("should return error during retrieveing subnets", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				subnets, err := newSubnet(server.URL()).List(target)
				Expect(err).To(HaveOccurred())
				Expect(subnets).Should(BeNil())
			})
		})
	})

})

func newSubnet(url string) Subnets {

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.CfService,
	}
	return newSubnetAPI(&client)
}
