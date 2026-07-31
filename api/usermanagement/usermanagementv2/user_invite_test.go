package usermanagementv2

import (
	"log"
	"net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/client"
	bluemixHttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/session"

	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserInvite", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//List
	Describe("GetUsers", func() {
		Context("When read of users is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/accounts/1234567890abcdefgh/users"),
						ghttp.RespondWith(http.StatusOK, `{"TotalUsers": 1, "Limit": 0, "FistURL": "/v2/accounts/1234567890abcdefgh/users", "Resources":[{"ID" :"12ab34cd56", "IamID" :"xxxxxxxx"}]}`),
					),
				)
			})

			It("should return invited users list", func() {
				usersList, err := newUserInviteHandler(server.URL()).GetUsers("1234567890abcdefgh")
				Expect(usersList).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When read of users is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/accounts/1234567890abcdefgh/users"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve users info`),
					),
				)
			})

			It("should return error when users are retrieved", func() {
				usersList, err := newUserInviteHandler(server.URL()).GetUsers("1234567890abcdefgh")
				Expect(err).To(HaveOccurred())
				Expect(usersList.Resources).Should(BeNil())
			})
		})
	})

	//UserProfile
	Describe("GetUserProfiles", func() {
		Context("When read of user profile is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/accounts/1234567890abcdefgh/users/xxxxxxxx"),
						ghttp.RespondWith(http.StatusOK, `{"ID" :"12ab34cd56", "IamID" :"xxxxxxxx", "Realm": "yyyyyyy", "Firstname": "abcdef", "Lastname": "ijklmnop", "State":"Pending"}`),
					),
				)
			})

			It("should return user profile info", func() {
				usersProfile, err := newUserInviteHandler(server.URL()).GetUserProfile("1234567890abcdefgh", "xxxxxxxx")
				Expect(usersProfile).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When read of user profile is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/accounts/1234567890abcdefgh/users/xxxxxxxx"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve user profile info`),
					),
				)
			})

			It("should return error when users are retrieved", func() {
				userProfile, err := newUserInviteHandler(server.URL()).GetUserProfile("1234567890abcdefgh", "xxxxxxxx")
				Expect(err).To(HaveOccurred())
				Expect(userProfile.Phonenumber).Should(Equal(""))
			})
		})
	})

	//InviteUser
	Describe("InviteUser", func() {
		Context("When user invite is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/accounts/1234567890abcdefgh/users"),
						ghttp.VerifyJSON(`{"users": [
							{
							  "email": "test@in.ibm.com",
							  "account_role": "MEMBER"
							}
						  ]}`),
						ghttp.RespondWith(http.StatusCreated, `{ "resources": [{"id": "12000000000000000000000000000001", "email": "test@in.ibm.com", "state": "PROCESSING"}
  ]}`),
					),
				)
			})

			It("should return user is invited created", func() {
				users := make([]models.User, 0)
				users = append(users, models.User{Email: "test@in.ibm.com", AccountRole: "MEMBER"})
				userInvitePayload := models.UserInvite{Users: users}
				userInviteResp, err := newUserInviteHandler(server.URL()).InviteUsers("1234567890abcdefgh", userInvitePayload)
				Expect(err).NotTo(HaveOccurred())
				Expect(userInviteResp).ShouldNot(BeNil())
			})
		})
		Context("When user invitation is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/accounts/1234567890abcdefgh/users"),
						ghttp.VerifyJSON(`{"users": [
							{
							  "email": "test@in.ibm.com",
							  "account_role": "MEMBER"
							}
						  ]}`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to invite user`),
					),
				)
			})
			It("should return error during cluster creation", func() {
				users := make([]models.User, 0)
				users = append(users, models.User{Email: "test@in.ibm.com", AccountRole: "MEMBER"})
				userInvitePayload := models.UserInvite{Users: users}
				_, err := newUserInviteHandler(server.URL()).InviteUsers("1234567890abcdefgh", userInvitePayload)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//RemoveUsers
	Describe("RemoveUsers", func() {
		Context("When remove of user is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v2/accounts/1234567890abcdefgh/users/xxxxxxxx"),
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)
			})

			It("should return user profile info", func() {
				err := newUserInviteHandler(server.URL()).RemoveUsers("1234567890abcdefgh", "xxxxxxxx")
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When remove of user is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v2/accounts/1234567890abcdefgh/users/xxxxxxxx"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to reemove the user`),
					),
				)
			})

			It("should return error when tried remove the user", func() {
				err := newUserInviteHandler(server.URL()).RemoveUsers("1234567890abcdefgh", "xxxxxxxx")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//List
	Describe("GetUsers", func() {
		Context("When read of users is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/accounts/1234567890abcdefgh/users"),
						ghttp.RespondWith(http.StatusOK, `{"TotalUsers": 1, "Limit": 0, "FistURL": "/v2/accounts/1234567890abcdefgh/users", "Resources":[{"ID" :"12ab34cd56", "IamID" :"xxxxxxxx"}]}`),
					),
				)
			})

			It("should return invited users list", func() {
				usersList, err := newUserInviteHandler(server.URL()).GetUsers("1234567890abcdefgh")
				Expect(usersList).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When read of users is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/accounts/1234567890abcdefgh/users"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve users info`),
					),
				)
			})

			It("should return error when users are retrieved", func() {
				usersList, err := newUserInviteHandler(server.URL()).GetUsers("1234567890abcdefgh")
				Expect(err).To(HaveOccurred())
				Expect(usersList.Resources).Should(BeNil())
			})
		})
	})

})

func newUserInviteHandler(url string) Users {
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.UserManagement,
	}
	return NewUserInviteHandler(&client)
}
