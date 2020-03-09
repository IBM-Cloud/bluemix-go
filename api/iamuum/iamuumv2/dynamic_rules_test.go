package iamuumv2

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

var _ = Describe("DynamicRuleRepository", func() {
	var (
		server *ghttp.Server
	)

	Describe("List()", func() {
		Context("When API error 403 returns", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/groups/def/rules"),
						ghttp.RespondWith(http.StatusForbidden, `
						{
							"message": "The provided access token does not have the proper authority to access this operation."
						}`),
					),
				)
			})

			It("should return API 403 error", func() {
				_, err := newTestDynamicRuleRepo(server.URL()).List("def")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("Request failed with status code: 403"))
			})
		})

		Context("When other JSON error returns", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/groups/def/rules"),
						ghttp.RespondWith(http.StatusBadGateway, `{
							"message": "other json error"
						}`),
					),
				)
			})

			It("should return server error", func() {
				_, err := newTestDynamicRuleRepo(server.URL()).List("def")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("other json error"))
			})
		})

	})

	Describe("Create()", func() {
		Context("When create one rule", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/groups/def/rules"),
						ghttp.RespondWith(http.StatusOK, `
						{
            "id": "ClaimRule-039dbce6-d5ee-4ddf-80a3-70be45f85f3b",
            "name": "test rule name 3",
            "expiration": 24,
            "realm_name": "test-idp.com",
            "access_group_id": "AccessGroupId-eac7839a-2a15-4f1b-a08e-58df17354845",
            "account_id": "faf6addbf6bf476896f5e342a5bdd702",
            "conditions": [
                {
                    "claim": "blueGroups",
                    "operator": "CONTAINS",
                    "value": "\"test-bluegroup-saml\""
                }
            ],
            "created_at": "2020-03-09T08:46:46Z",
            "created_by_id": "IBMid-550003NN7D",
            "last_modified_at": "2020-03-09T08:46:46Z",
            "last_modified_by_id": "IBMid-550003NN7D"
        }`),
					),
				)
			})

			It("should return success", func() {
				response, err := newTestDynamicRuleRepo(server.URL()).Create("def", CreateRuleRequest{
					Name:       "abc",
					Expiration: 24,
					RealmName:  "test-idp.com",
					Conditions: []Condition{
						{
							Claim:    "blueGroups",
							Operator: "CONTAINS",
							Value:    "\"test-bluegroup-saml\"",
						},
					},
				})
				Expect(err).ShouldNot(HaveOccurred())

				Expect(response.Name).Should(Equal("test rule name 3"))
				Expect(response.Expiration).Should(Equal(24))
			})

		})
	})

	Describe("Remove()", func() {
		Context("When rule is deleted", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v2/groups/abc/rules/def"),
						ghttp.RespondWith(http.StatusNoContent, ""),
					),
				)
			})

			It("should return success", func() {
				err := newTestDynamicRuleRepo(server.URL()).Delete("abc", "def")

				Expect(err).Should(Succeed())
			})
		})

		Context("When group is not found", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v2/groups/abc/rules/def"),
						ghttp.RespondWith(http.StatusNotFound, `{
							"StatusCode": 404,
							"code": "not_found",
							"message": "Group abc is not found"
						}`),
					),
				)
			})

			It("should return not found error", func() {
				err := newTestDynamicRuleRepo(server.URL()).Delete("abc", "def")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("not_found"))
			})
		})
	})

})

func newTestDynamicRuleRepo(url string) DynamicRuleRepository {
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.Endpoint = &url
	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.IAMUUMServicev2,
	}
	return NewDynamicRuleRepository(&client)
}
