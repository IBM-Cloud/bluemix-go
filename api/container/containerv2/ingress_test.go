package containerv2

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

var _ = Describe("Ingress Secrets", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//Enable
	Describe("Create", func() {
		Context("When creating ingress secret successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/ingress/v2/secret/createSecret"),
						ghttp.VerifyJSON(`{"cluster":"bugi52rf0rtfgadjfso0","name":"testabc2","namespace":"","crn":"crn:v1:bluemix:public:cloudcerts:us-south:a/883079c85357a1f3f85d968780e56518:b65b5b7f-e904-4d2b-bd87-f0ccd57e76ba:certificate:333d8673f4d03c148ff81192b9edaafc","persistence":true}`),
						ghttp.RespondWith(http.StatusCreated, `{}`),
					),
				)
			})

			It("should create Ingress Secret in a cluster", func() {
				target := ClusterTargetHeader{}
				params := SecretCreateConfig{
					Cluster:     "bugi52rf0rtfgadjfso0",
					Name:        "testabc2",
					CRN:         "crn:v1:bluemix:public:cloudcerts:us-south:a/883079c85357a1f3f85d968780e56518:b65b5b7f-e904-4d2b-bd87-f0ccd57e76ba:certificate:333d8673f4d03c148ff81192b9edaafc",
					Persistence: true,
				}
				_, err := newIngresses(server.URL()).CreateIngressSecret(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When creating is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/ingress/v2/secret/createSecret"),
						ghttp.VerifyJSON(`{"cluster":"bugi52rf0rtfgadjfso0","name":"testabc2","namespace":"","crn":"crn:v1:bluemix:public:cloudcerts:us-south:a/883079c85357a1f3f85d968780e56518:b65b5b7f-e904-4d2b-bd87-f0ccd57e76ba:certificate:333d8673f4d03c148ff81192b9edaafc","persistence":true}
`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to enable ingress`),
					),
				)
			})

			It("should return error during creating ingress secret", func() {
				params := SecretCreateConfig{
					Cluster:     "bugi52rf0rtfgadjfso0",
					Name:        "testabc2",
					CRN:         "crn:v1:bluemix:public:cloudcerts:us-south:a/883079c85357a1f3f85d968780e56518:b65b5b7f-e904-4d2b-bd87-f0ccd57e76ba:certificate:333d8673f4d03c148ff81192b9edaafc",
					Persistence: true,
				}
				target := ClusterTargetHeader{}
				_, err := newIngresses(server.URL()).CreateIngressSecret(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//Disable
	Describe("Destroy", func() {
		Context("When deleting ingress secret successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/ingress/v2/secret/deleteSecret"),
						ghttp.VerifyJSON(`{"cluster":"bugi52rf0rtfgadjfso0","name":"testabc2","namespace":"default"}`),
						ghttp.RespondWith(http.StatusCreated, `{}`),
					),
				)
			})

			It("should destroy Ingress Secret in a cluster", func() {
				target := ClusterTargetHeader{}
				params := SecretDeleteConfig{
					Cluster:   "bugi52rf0rtfgadjfso0",
					Name:      "testabc2",
					Namespace: "default",
				}
				err := newIngresses(server.URL()).DeleteIngressSecret(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When Destroying is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/ingress/v2/secret/deleteSecret"),
						ghttp.VerifyJSON(`{"cluster":"bugi52rf0rtfgadjfso0","name":"testabc2","namespace":"default"}`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to disable ingress`),
					),
				)
			})

			It("should return error during destroying ingress", func() {
				params := SecretDeleteConfig{
					Cluster:   "bugi52rf0rtfgadjfso0",
					Name:      "testabc2",
					Namespace: "default",
				}
				target := ClusterTargetHeader{}
				err := newIngresses(server.URL()).DeleteIngressSecret(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//GetIngress Secrets
	Describe("Get", func() {
		Context("When Get Ingress Secret is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/ingress/v2/secret/getSecret"),
						ghttp.RespondWith(http.StatusCreated, `{"cluster":"bugi52rf0rtfgadjfso0","name":"testabc2","namespace":"default","domain":"*.mytestclustercb8f-dce1dcf4a47f9ff42332256e6c4eb998-0000.us-south.containers.appdomain.cloud","crn":"crn:v1:bluemix:public:cloudcerts:us-south:a/883079c85357a1f3f85d968780e56518:b65b5b7f-e904-4d2b-bd87-f0ccd57e76ba:certificate:333d8673f4d03c148ff81192b9edaafc","expiresOn":"2021-01-27T00:18:56+0000","status":"created","userManaged":true,"persistence":true}`),
					),
				)
			})

			It("should get Ingress Secret in a cluster", func() {
				target := ClusterTargetHeader{}

				_, err := newIngresses(server.URL()).GetIngressSecret("bugi52rf0rtfgadjfso0", "testabc2", "default", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When get ingress secret unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/ingress/v2/secret/getSecret"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get ingress`),
					),
				)
			})

			It("should return error during get ingress", func() {
				target := ClusterTargetHeader{}
				_, err := newIngresses(server.URL()).GetIngressSecret("bugi52rf0rtfgadjfso0", "testabc2", "default", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})

func newIngresses(url string) Ingress {

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
	return newIngressAPI(&client)
}
