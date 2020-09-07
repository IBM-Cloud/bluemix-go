package containerv1

import (
	"log"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/client"
	bluemixHttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/session"
)

var _ = Describe("AddOns", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})
	//Configure
	Describe("Configure", func() {
		Context("When configuring addons is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPatch, "/v1/clusters/testcluster/addons"),
						ghttp.VerifyJSON(`{
							"addons": [
							  {
								"deprecated": false,
								"name": "istio",
								"vlan_spanning_required": false
							  }
							],
							"enable": false,
							"update": false
						  }`),
						ghttp.RespondWith(http.StatusCreated, `{}`),
					),
				)
			})

			It("should configure addon to a cluster", func() {
				target := ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
					Region:    "eu-de",
				}
				clustername := "testcluster"
				params := ConfigureAddOns{
					AddonsList: []AddOn{},
					Enable:     false,
					Update:     false,
				}
				var addOn = AddOn{
					Name: "istio",
				}
				params.AddonsList = append(params.AddonsList, addOn)
				_, err := newAddOns(server.URL()).ConfigureAddons(clustername, &params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When configuring addons is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPatch, "/v1/clusters/testcluster/addons"),
						ghttp.VerifyJSON(`{
							"addons": [
							  {
								"deprecated": false,
								"name": "istio",
								"vlan_spanning_required": false
							  }
							],
							"enable": false,
							"update": false
						  }`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to configure addons`),
					),
				)
			})

			It("should return error during configuring addons", func() {

				params := ConfigureAddOns{
					AddonsList: []AddOn{},
					Enable:     false,
					Update:     false,
				}
				var addOn = AddOn{
					Name: "istio",
				}
				params.AddonsList = append(params.AddonsList, addOn)
				clustername := "testcluster"
				target := ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
					Region:    "eu-de",
				}
				_, err := newAddOns(server.URL()).ConfigureAddons(clustername, &params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//GetAlb
	Describe("Get cluster addons", func() {
		Context("When read of cluster addons is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/clusters/testcluster/addons"),
						ghttp.RespondWith(http.StatusOK, `[
							{
							  "name": "istio",
							  "version": "1.6",
							  "targetVersion": "1.6",
							  "healthStatus": "Enabling",
							  "allowed_upgrade_versions": [
								"1.7"
							  ]
							}
						  ]
						  `),
					),
				)
			})

			It("should return addons", func() {
				target := ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
					Region:    "eu-de",
				}
				addons, err := newAddOns(server.URL()).GetAddons("testcluster", target)
				Expect(addons).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When read of cluster addon is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/clusters/testcluster/addons"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve addons.`),
					),
				)
			})

			It("should return error when addons are retrieved", func() {
				target := ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
					Region:    "eu-de",
				}
				_, err := newAddOns(server.URL()).GetAddons("testcluster", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})

func newAddOns(url string) AddOns {

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
	return newAddOnsAPI(&client)
}
