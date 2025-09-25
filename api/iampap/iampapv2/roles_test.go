package iampapv2

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

var _ = Describe("RoleRepository", func() {
	var (
		server *ghttp.Server
	)

	Describe("Create()", func() {
		Context("When create one role", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/roles"),
						ghttp.RespondWith(http.StatusOK, `
						{
							"id": "12345678abcd1a2ba1b21234567890ab",
							"crn": "crn:v1:bluemix:public:iam-access-management::::customRole:Example",
							"name": "abc",
							"display_name": "abc",
							"description": "test role",
							"service_name": "kms",
							"account_id": "acc",
							"actions": [
							  "kms.secrets.readmetadata"
							],
							"created_at": "2018-08-30T14:09:09.907Z",
							"created_by_id": "USER_ID",
							"last_modified_at": "2018-08-30T14:09:09.907Z",
							"last_modified_by_id": "USER_ID",
							"href": "https://iam.cloud.ibm.com/v2/roles/12345678abcd1a2ba1b21234567890ab"
						  }`),
					),
				)
			})

			It("should return success", func() {
				response, err := newTestRoleRepo(server.URL()).Create(CreateRoleRequest{
					Name:        "abc",
					ServiceName: "kms",
					AccountID:   "acc",
					Description: "test role",
					DisplayName: "abc",
					Actions:     []string{"kms.secrets.readmetadata"}})
				Expect(err).ShouldNot(HaveOccurred())

				Expect(response.Name).Should(Equal("abc"))
				Expect(response.ServiceName).Should(Equal("kms"))
			})

		})
	})

	Describe("Remove()", func() {
		Context("When role is deleted", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v2/roles/abc"),
						ghttp.RespondWith(http.StatusNoContent, ""),
					),
				)
			})

			It("should return success", func() {
				err := newTestRoleRepo(server.URL()).Delete("abc")

				Expect(err).Should(Succeed())
			})
		})

		Context("When role is not found", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v2/roles/abc"),
						ghttp.RespondWith(http.StatusNotFound, `{
							"StatusCode": 404,
							"code": "not_found",
							"message": "role abc is not found"
						}`),
					),
				)
			})

			It("should return not found error", func() {
				err := newTestRoleRepo(server.URL()).Delete("abc")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("not_found"))
			})
		})
	})

})

func newTestRoleRepo(url string) RoleRepository {
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.Endpoint = &url
	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.IAMPAPServicev2,
	}
	return NewRoleRepository(&client)
}
