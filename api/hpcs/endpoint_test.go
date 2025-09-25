package hpcs

import (
	"log"
	"net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/session"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("HpcsRepository", func() {
	var (
		server *ghttp.Server
	)

	Describe("Get()", func() {
		Context("When generating endpoint", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/instances/abc"),
						ghttp.RespondWith(http.StatusCreated, `{}`),
					),
				)
			})

			It("should return success", func() {
				_, err := newTestHpcsRepo(server.URL()).GetAPIEndpoint("abc")

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

})

func newTestHpcsRepo(url string) EndpointRepository {
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.Endpoint = &url
	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.HPCService,
	}
	return NewHpcsEndpointRepository(&client)
}
