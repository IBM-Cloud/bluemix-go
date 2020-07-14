package controllerv2

import (
	"log"
	"net/http"

	"github.com/IBM-Cloud/bluemix-go"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceInstances", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	Describe("ListInstances()", func() {
		Context("When there is no service instance", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/resource_instances"),
						ghttp.RespondWith(http.StatusOK, `{"resources":[]}`),
					),
				)
			})
			It("should return zero service instance", func() {
				repo := newTestServiceInstanceRepo(server.URL())
				instances, err := repo.ListInstances(ServiceInstanceQuery{
					ResourceGroupID: "resource_group_id",
				})

				Expect(err).ShouldNot(HaveOccurred())
				Expect(instances).Should(BeEmpty())
			})
		})
		Context("When there is one service instance", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/resource_instances"),
						ghttp.RespondWith(http.StatusOK, `{
																	"rows_count":1,
																	"resources":[{
																		"id":"foo",
																		"guid":"a83261db-5cd1-46ae-8cfb-8ebbcc7c0184",
																		"url":"/v1/resource_instances/a83261db-5cd1-46ae-8cfb-8ebbcc7c0184",
																		"created_at":"2017-07-31T06:19:45.16112535Z",
																		"updated_at":null,
																		"deleted_at":null,
																		"name":"test-instance",
																		"target_crn":"crn:v1:d_att288:dedicated::us-south::::d_att288-us-south",
																		"account_id":"560df2058b1e7c402303cc598b3e5540",
																		"resource_plan_id":"rc-pb-28c24fccc-4ca6-4ddd-a3be-7746cdce9912",
																		"resource_group_id":"resource_group_id",
																		"create_time":0,"crn":"",
																		"state":"inactive",
																		"type":"service_instance",
																		"resource_id":"fake-resource-id",
																		"dashboard_url":null,
																		"last_operation":null,
																		"account_url":"/v1/accounts/560df2058b1e7c402303cc598b3e5540",
																		"resource_plan_url":"/v1/catalog/regions/ibm:ys1:us-south/plans/rc-pb-28c24fccc-4ca6-4ddd-a3be-7746cdce9912",
																		"resource_bindings_url":"/v1/resource_instances/a83261db-5cd1-46ae-8cfb-8ebbcc7c0184/resource_bindings",
																		"resource_aliases_url":"/v1/resource_instances/a83261db-5cd1-46ae-8cfb-8ebbcc7c0184/resource_aliases",
																		"siblings_url":"/v1/resource_instances/a83261db-5cd1-46ae-8cfb-8ebbcc7c0184/siblings"}]}`),
					),
				)
			})
			It("should return one service instance", func() {
				repo := newTestServiceInstanceRepo(server.URL())
				instances, err := repo.ListInstances(ServiceInstanceQuery{
					ResourceGroupID: "resource_group_id",
				})

				Expect(err).ShouldNot(HaveOccurred())

				Expect(instances).Should(HaveLen(1))
				instance := instances[0]
				Expect(instance.ID).Should(Equal("foo"))
				Expect(instance.ServiceID).Should(Equal("fake-resource-id"))
				Expect(instance.Name).Should(Equal("test-instance"))
				Expect(instance.State).Should(Equal("inactive"))
				Expect(instance.Type).Should(Equal("service_instance"))
			})
		})

		Context("When there are multiple service instances", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/resource_instances"),
						ghttp.RespondWith(http.StatusOK, `{
						"rows_count":3,
						"resources":[{
							"id":"foo",
							"guid":"a83261db-5cd1-46ae-8cfb-8ebbcc7c0184",
							"url":"/v1/resource_instances/a83261db-5cd1-46ae-8cfb-8ebbcc7c0184",
							"created_at":"2017-07-31T06:19:45.16112535Z",
							"updated_at":null,
							"deleted_at":null,
							"name":"test-instance",
							"target_crn":"crn:v1:d_att288:dedicated::us-south::::d_att288-us-south",
							"account_id":"560df2058b1e7c402303cc598b3e5540",
							"resource_plan_id":"rc-pb-28c24fccc-4ca6-4ddd-a3be-7746cdce9912",
							"resource_group_id":"",
							"create_time":0,"crn":"",
							"state":"active",
							"type":"service_instance",
							"resource_id":"fake-resource-id",
							"dashboard_url":null,
							"last_operation":null,
							"account_url":"/v1/accounts/560df2058b1e7c402303cc598b3e5540",
							"resource_plan_url":"/v1/catalog/regions/ibm:ys1:us-south/plans/rc-pb-28c24fccc-4ca6-4ddd-a3be-7746cdce9912",
							"resource_bindings_url":"/v1/resource_instances/a83261db-5cd1-46ae-8cfb-8ebbcc7c0184/resource_bindings",
							"resource_aliases_url":"/v1/resource_instances/a83261db-5cd1-46ae-8cfb-8ebbcc7c0184/resource_aliases",
							"siblings_url":"/v1/resource_instances/a83261db-5cd1-46ae-8cfb-8ebbcc7c0184/siblings"
						},
						{
							"id":"foo1",
							"guid":"dea23694-1a2c-45e4-bfa8-7be5226d998c",
							"url":"/v1/resource_instances/dea23694-1a2c-45e4-bfa8-7be5226d998c",
							"created_at":"2017-07-31T06:20:14.592704474Z",
							"updated_at":null,
							"deleted_at":null,
							"name":"test-instance1",
							"target_crn":"crn:v1:d_att288:dedicated::us-south::::d_att288-us-south",
							"account_id":"560df2058b1e7c402303cc598b3e5540",
							"resource_plan_id":"rc-pb-28c24fccc-4ca6-4ddd-a3be-7746cdce9912",
							"resource_group_id":"",
							"create_time":0,
							"crn":"",
							"state":"inactive",
							"type":"service_instance",
							"resource_id":"fake-resource-id1",
							"dashboard_url":null,
							"last_operation":null,
							"account_url":"/v1/accounts/560df2058b1e7c402303cc598b3e5540",
							"resource_plan_url":"/v1/catalog/regions/ibm:ys1:us-south/plans/rc-pb-28c24fccc-4ca6-4ddd-a3be-7746cdce9912",
							"resource_bindings_url":"/v1/resource_instances/dea23694-1a2c-45e4-bfa8-7be5226d998c/resource_bindings",
							"resource_aliases_url":"/v1/resource_instances/dea23694-1a2c-45e4-bfa8-7be5226d998c/resource_aliases",
							"siblings_url":"/v1/resource_instances/dea23694-1a2c-45e4-bfa8-7be5226d998c/siblings"
						},
						{
							"id":"foo2",
							"guid":"50312f63-f43b-4a67-aa32-8626da609adb",
							"url":"/v1/resource_instances/50312f63-f43b-4a67-aa32-8626da609adb",
							"created_at":"2017-07-31T06:27:46.215093281Z",
							"updated_at":"2017-07-31T07:34:07.740506169Z",
							"deleted_at":null,
							"name":"test-instance2",
							"target_crn":"crn:v1:d_att288:dedicated::us-south::::d_att288-us-south",
							"account_id":"560df2058b1e7c402303cc598b3e5540",
							"resource_plan_id":"rc-pb-28c24fccc-4ca6-4ddd-a3be-7746cdce9912",
							"resource_group_id":"",
							"create_time":0,
							"crn":"",
							"state":"active",
							"type":"service_instance",
							"resource_id":"fake-resource-id2",
							"dashboard_url":null,
							"last_operation":null,
							"account_url":"/v1/accounts/560df2058b1e7c402303cc598b3e5540",
							"resource_plan_url":"/v1/catalog/regions/ibm:ys1:us-south/plans/rc-pb-28c24fccc-4ca6-4ddd-a3be-7746cdce9912",
							"resource_bindings_url":"/v1/resource_instances/50312f63-f43b-4a67-aa32-8626da609adb/resource_bindings",
							"resource_aliases_url":"/v1/resource_instances/50312f63-f43b-4a67-aa32-8626da609adb/resource_aliases",
							"siblings_url":"/v1/resource_instances/50312f63-f43b-4a67-aa32-8626da609adb/siblings"
						}
						]}`),
					),
				)
			})
			It("should return all of them", func() {
				repo := newTestServiceInstanceRepo(server.URL())
				instances, err := repo.ListInstances(ServiceInstanceQuery{
					ResourceGroupID: "resource_group_id",
				})

				Expect(err).ShouldNot(HaveOccurred())

				Expect(instances).Should(HaveLen(3))
				instance := instances[0]
				Expect(instance.ID).Should(Equal("foo"))
				Expect(instance.ServiceID).Should(Equal("fake-resource-id"))
				Expect(instance.Name).Should(Equal("test-instance"))
				Expect(instance.State).Should(Equal("active"))
				Expect(instance.Type).Should(Equal("service_instance"))

				instance = instances[1]
				Expect(instance.ID).Should(Equal("foo1"))
				Expect(instance.ServiceID).Should(Equal("fake-resource-id1"))
				Expect(instance.Name).Should(Equal("test-instance1"))
				Expect(instance.State).Should(Equal("inactive"))
				Expect(instance.Type).Should(Equal("service_instance"))

				instance = instances[2]
				Expect(instance.ID).Should(Equal("foo2"))
				Expect(instance.ServiceID).Should(Equal("fake-resource-id2"))
				Expect(instance.Name).Should(Equal("test-instance2"))
				Expect(instance.State).Should(Equal("active"))
				Expect(instance.Type).Should(Equal("service_instance"))
			})
		})

	})

})

func newTestServiceInstanceRepo(url string) ResourceServiceInstanceRepository {

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.ResourceManagementServicev2,
	}

	return newResourceServiceInstanceAPI(&client)
}
