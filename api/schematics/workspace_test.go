package schematics

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

var _ = Describe("workspaces", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//getworkspace
	Describe("Get", func() {
		Context("When Get workspace is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/workspaces/myworkspaceptab-bc54176d-bbcb-42"),
						ghttp.RespondWith(http.StatusCreated, `{
							"catalog_ref": {
							  "item_icon_url": "string",
							  "item_id": "string",
							  "item_name": "string",
							  "item_readme_url": "string",
							  "item_url": "string",
							  "launch_url": "string",
							  "offering_version": "string"
							},
							"created_at": "2019-11-12T17:56:31.081Z",
							"created_by": "string",
							"description": "string",
							"id": "string",
							"last_health_check_at": "2019-11-12T17:56:31.081Z",
							"location": "string",
							"name": "string",
							"resource_group": "string",
							"runtime_data": [
							  {
								"engine_cmd": "string",
								"engine_name": "string",
								"engine_version": "string",
								"id": "string",
								"log_store_url": "string",
								"output_values": {
								  "additionalProp1": "string",
								  "additionalProp2": "string",
								  "additionalProp3": "string"
								},
								"resources": [
								  [
									{}
								  ]
								],
								"state_store_url": "string"
							  }
							],
							"shared_data": {
							  "cluster_id": "string",
							  "cluster_name": "string",
							  "entitlement_keys": [
								{}
							  ],
							  "namespace": "string",
							  "region": "string",
							  "resource_group_id": "string"
							},
							"status": "string",
							"tags": [
							  "string"
							],
							"template_data": [
							  {
								"env_values": [
								  {
									"hidden": true,
									"name": "string",
									"secure": true,
									"value": "string"
								  }
								],
								"folder": "string",
								"id": "string",
								"type": "string",
								"uninstall_script_name": "string",
								"values": "string",
								"values_metadata": [
								  {}
								],
								"values_url": "string",
								"variablestore": [
								  {
									"description": "string",
									"name": "string",
									"type": "string",
									"value": "string"
								  }
								]
							  }
							],
							"template_ref": "string",
							"template_repo": {
							  "branch": "string",
							  "release": "string",
							  "repo_url": "string",
							  "url": "string"
							},
							"type": [
							  "string"
							],
							"updated_at": "2019-11-12T17:56:31.081Z",
							"updated_by": "string",
							"workspace_status": {
							  "frozen": true,
							  "frozen_at": "2019-11-12T17:56:31.081Z",
							  "frozen_by": "string",
							  "locked": true,
							  "locked_by": "string",
							  "locked_time": "2019-11-12T17:56:31.081Z"
							},
							"workspace_status_msg": {
							  "status_code": "string",
							  "status_msg": "string"
							}
						  }`),
					),
				)
			})

			It("should get Workspace", func() {

				_, err := newWorkspace(server.URL()).GetWorkspaceByID("myworkspaceptab-bc54176d-bbcb-42")
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When get workspace is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/workspaces/myworkspaceptab-bc54176d-bbcb-42"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get workerpool`),
					),
				)
			})

			It("should return error during get workspace", func() {
				_, err := newWorkspace(server.URL()).GetWorkspaceByID("myworkspaceptab-bc54176d-bbcb-42")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//getStateStore
	Describe("GetState", func() {
		Context("When Get State store is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/workspaces/myworkspaceptab-bc54176d-bbcb-42/runtime_data/29291199-ca08-46/state_store"),
						ghttp.RespondWith(http.StatusCreated, `{}`),
					),
				)
			})

			It("should get State", func() {

				_, err := newWorkspace(server.URL()).GetStateStore("myworkspaceptab-bc54176d-bbcb-42", "29291199-ca08-46")
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When get state is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/workspaces/myworkspaceptab-bc54176d-bbcb-42/runtime_data/29291199-ca08-46/state_store"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get state`),
					),
				)
			})

			It("should return error during get state", func() {
				_, err := newWorkspace(server.URL()).GetStateStore("myworkspaceptab-bc54176d-bbcb-42", "29291199-ca08-46")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//getStateStore
	Describe("GetOutput Values", func() {
		Context("When Get output is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/workspaces/myworkspaceptab-bc54176d-bbcb-42/output_values"),
						ghttp.RespondWith(http.StatusCreated, `[{
						"folder": "string",
						"id": "string",
						"output_values": [
						  {}
						],
						"type": "string"
					  }]`),
					),
				)
			})

			It("should get output values", func() {

				_, err := newWorkspace(server.URL()).GetOutputValues("myworkspaceptab-bc54176d-bbcb-42")
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When get output is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/workspaces/myworkspaceptab-bc54176d-bbcb-42/output_values"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get state`),
					),
				)
			})

			It("should return error during get state", func() {
				_, err := newWorkspace(server.URL()).GetOutputValues("myworkspaceptab-bc54176d-bbcb-42")
				Expect(err).To(HaveOccurred())
			})
		})
	})

})

func newWorkspace(url string) Workspaces {

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.SchematicsService,
	}
	return newWorkspaceAPI(&client)
}
