package functions

import (
	"log"
	"net/http"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/client"
	"github.com/Mavrickk3/bluemix-go/session"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Functions", func() {
	var (
		server *ghttp.Server
	)

	Describe("Remove()", func() {
		Context("When namespace is deleted", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/api/v1/namespaces/abc"),
						ghttp.RespondWith(http.StatusNoContent, ""),
					),
				)
			})

			It("should return success", func() {
				_, err := newTestNamespace(server.URL()).DeleteNamespace("/api/v1/namespaces/abc")

				Expect(err).Should(Succeed())
			})
		})

		Context("When namespace is not found", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/api/v1/namespaces/abc"),
						ghttp.RespondWith(http.StatusNotFound, `{
							"StatusCode": 404,
							"code": "not_found",
							"message": "namespace abc is not found"
						}`),
					),
				)
			})

			It("should return not found error", func() {
				_, err := newTestNamespace(server.URL()).DeleteNamespace("abc")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("not_found"))
			})
		})
	})

})

func newTestNamespace(url string) Functions {
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.Endpoint = &url
	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.FunctionsService,
	}
	return newFunctionsAPI(&client)
}
