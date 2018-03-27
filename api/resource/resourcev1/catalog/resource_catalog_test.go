package catalog

import (
	"log"
	"net/http"

	"github.com/IBM-Bluemix/bluemix-go"

	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/models"
	"github.com/IBM-Bluemix/bluemix-go/session"
	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResourceCatalogRepository", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	Describe("FindByName", func() {
		Context("When there is no resource", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v1/", "q=test"),
						ghttp.RespondWith(http.StatusOK, `{"resources":[]}`),
					),
				)
			})
			It("should return not found error", func() {
				repo := newTestResourceCatalogRepo(server.URL())
				resources, err := repo.FindByName("test", false)

				Expect(err).Should(HaveOccurred())
				Expect(resources).Should(BeEmpty())
			})
		})

		Context("When there is one resource", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v1/", "q=global-geo-test-service&include=*"),
						ghttp.RespondWith(http.StatusOK, `
							{"offset":0,"limit":50,"count":1,"resource_count":1,"first":"https://resource-catalog.stage1.ng.bluemix.net/api/v1?include=%2A&languages=en-US%2Cen%3Bq%3D0.8%2Czh-CN%3Bq%3D0.6%2Czh%3Bq%3D0.4&q=global-geo-test-service","resources":[{"children_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/b275d042-8376-11e7-bb31-be2e44b06b34/%2A","created":"2017-08-17T20:16:54.137Z","catalog_crn":"crn:v1:staging:public:resource-catalog::a/ac5bc3a3a4f1dd065f5c92b444be3d01::service:b275d042-8376-11e7-bb31-be2e44b06b34","geo_tags":["eu-gb","us-south"],"id":"b275d042-8376-11e7-bb31-be2e44b06b34","images":{"feature_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant64.png","image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant50.png","medium_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant32.png","small_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant24.png"},"kind":"service","metadata":{"callbacks":{"broker_url":"https://dev-resource-catalog.stage1.ng.bluemix.net/api/v1/b275d042-8376-11e7-bb31-be2e44b06b34?include=%2A&languages=en-US%2Cen%3Bq%3D0.8"},"service":{"active":true,"bindable":true,"cf_plan_guid":{"stage1.ng.bluemix.net":"1bcf9cc3-3d26-481f-a72c-f9ad40495b31"},"iam_compatible":true,"plan_updateable":true,"rc_compatible":true,"rc_provisionable":true,"service_check_enabled":true,"state":"","test_check_interval":0,"unique_api_key":false},"ui":{"allow_internal_users":false,"urls":{}}},"name":"global-geo-test-service","overview_ui":{"en":{"description":"global-geo-test-service.","display_name":"global-geo-test-service","long_description":"global-geo-test-service"}},"provider":{"email":"sales@cloudant.com","name":"IBM"},"tags":["ibm_created"],"updated":"2017-09-08T23:26:28.981Z","url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/b275d042-8376-11e7-bb31-be2e44b06b34?include=%2A&languages=en-US%2Cen%3Bq%3D0.8%2Czh-CN%3Bq%3D0.6%2Czh%3Bq%3D0.4","visibility":{"exclude":{},"include":{"accounts":{"79801189f2ab72d269b633cdea173b20":"global","ac5bc3a3a4f1dd065f5c92b444be3d01":"global","bd087e8d604cef993f957c859e37f283":"global"}},"owner":"a/ac5bc3a3a4f1dd065f5c92b444be3d01","restrict":false}}]}`),
					),
				)
			})
			It("should return one resourxce", func() {
				repo := newTestResourceCatalogRepo(server.URL())
				resources, err := repo.FindByName("global-geo-test-service", true)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resources).Should(Equal([]models.Service{{
					ID:       "b275d042-8376-11e7-bb31-be2e44b06b34",
					CRN:      "crn:v1:staging:public:resource-catalog::a/ac5bc3a3a4f1dd065f5c92b444be3d01::service:b275d042-8376-11e7-bb31-be2e44b06b34",
					Name:     "global-geo-test-service",
					Kind:     "service",
					URL:      "https://resource-catalog.stage1.ng.bluemix.net/api/v1/b275d042-8376-11e7-bb31-be2e44b06b34?include=%2A&languages=en-US%2Cen%3Bq%3D0.8%2Czh-CN%3Bq%3D0.6%2Czh%3Bq%3D0.4",
					Metadata: []byte(`{"callbacks":{"broker_url":"https://dev-resource-catalog.stage1.ng.bluemix.net/api/v1/b275d042-8376-11e7-bb31-be2e44b06b34?include=%2A&languages=en-US%2Cen%3Bq%3D0.8"},"service":{"active":true,"bindable":true,"cf_plan_guid":{"stage1.ng.bluemix.net":"1bcf9cc3-3d26-481f-a72c-f9ad40495b31"},"iam_compatible":true,"plan_updateable":true,"rc_compatible":true,"rc_provisionable":true,"service_check_enabled":true,"state":"","test_check_interval":0,"unique_api_key":false},"ui":{"allow_internal_users":false,"urls":{}}}`),
				}}))
			})
		})

		Context("When there is multiple resources", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v1/", "q=Automation%20test&include=*"),
						ghttp.RespondWith(http.StatusOK, `
							{"offset":0,"limit":50,"count":5,"resource_count":5,"first":"https://resource-catalog.stage1.ng.bluemix.net/api/v1?include=%2A&languages=en-US%2Cen%3Bq%3D0.8%2Czh-CN%3Bq%3D0.6%2Czh%3Bq%3D0.4&q=Automation+test","resources":[{"children_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/edca1a82-7326-11e7-8cf7-a6006ad3dba0/%2A","created":"2017-08-08T21:24:31.731Z","catalog_crn":"crn:v1:staging:public:resource-catalog::a/ac5bc3a3a4f1dd065f5c92b444be3d01::service:edca1a82-7326-11e7-8cf7-a6006ad3dba0","geo_tags":["eu-gb","us-south"],"id":"edca1a82-7326-11e7-8cf7-a6006ad3dba0","images":{"feature_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant64.png","image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant50.png","medium_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant32.png","small_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant24.png"},"kind":"service","metadata":{"callbacks":{"broker_url":"https://dev-resource-catalog.stage1.ng.bluemix.net/api/v1/edca1a82-7326-11e7-8cf7-a6006ad3dba0?include=%2A&languages=en-US%2Cen%3Bq%3D0.8"},"service":{"active":false,"bindable":false,"cf_plan_guid":{"stage1.eu-gb.bluemix.net":"9bf099a2-7b22-4b32-a165-42e99d0286e8","stage1.ng.bluemix.net":"b1a78e15-b43f-40c5-821a-787dc9523409"},"iam_compatible":true,"plan_updateable":false,"rc_compatible":true,"rc_provisionable":true,"service_check_enabled":false,"state":"","test_check_interval":0,"unique_api_key":false},"ui":{"allow_internal_users":false,"urls":{}}},"name":"Async CF Automation Test","overview_ui":{"en":{"description":"Async CF  broker.","display_name":"Async CF  broker","long_description":"Async CF  broker"}},"pricing_tags":["free","paygo"],"provider":{"email":"sales@cloudant.com","name":"IBM"},"tags":["ibm_created"],"updated":"2017-08-09T19:22:11.866Z","url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/edca1a82-7326-11e7-8cf7-a6006ad3dba0?include=%2A&languages=en-US%2Cen%3Bq%3D0.8%2Czh-CN%3Bq%3D0.6%2Czh%3Bq%3D0.4","visibility":{"exclude":{},"include":{"accounts":{"79801189f2ab72d269b633cdea173b20":"global","ac5bc3a3a4f1dd065f5c92b444be3d01":"global","bd087e8d604cef993f957c859e37f283":"global"}},"owner":"a/ac5bc3a3a4f1dd065f5c92b444be3d01","restrict":false}},{"children_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/297ffc2e-732a-11e7-8cf7-a6006ad3dba0/%2A","created":"2017-08-08T20:45:02.148Z","catalog_crn":"crn:v1:staging:public:resource-catalog::a/3d0b3905d569f2c9c231b2a700c22ef4::service:297ffc2e-732a-11e7-8cf7-a6006ad3dba0","geo_tags":["eu-gb","us-south"],"id":"297ffc2e-732a-11e7-8cf7-a6006ad3dba0","images":{"feature_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant64.png","image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant50.png","medium_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant32.png","small_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant24.png"},"kind":"service","metadata":{"callbacks":{"broker_url":"https://dev-resource-catalog.stage1.ng.bluemix.net/api/v1/297ffc2e-732a-11e7-8cf7-a6006ad3dba0?include=%2A&languages=en-US%2Cen%3Bq%3D0.8"},"service":{"active":false,"bindable":false,"cf_plan_guid":{"stage1.eu-gb.bluemix.net":"9bf099a2-7b22-4b32-a165-42e99d0286e8","stage1.ng.bluemix.net":"b1a78e15-b43f-40c5-821a-787dc9523409"},"iam_compatible":true,"plan_updateable":false,"rc_compatible":true,"rc_provisionable":true,"service_check_enabled":false,"state":"","test_check_interval":0,"unique_api_key":false},"ui":{"allow_internal_users":false,"urls":{}}},"name":"Async Rc Automation test","overview_ui":{"en":{"description":"Async RC  broker.","display_name":"Async RC  broker","long_description":"Async RC  broker"}},"pricing_tags":["free","paygo"],"provider":{"email":"sales@cloudant.com","name":"IBM"},"tags":["ibm_created"],"updated":"2017-09-05T01:02:35.715Z","url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/297ffc2e-732a-11e7-8cf7-a6006ad3dba0?include=%2A&languages=en-US%2Cen%3Bq%3D0.8%2Czh-CN%3Bq%3D0.6%2Czh%3Bq%3D0.4","visibility":{"exclude":{},"include":{"accounts":{"3d0b3905d569f2c9c231b2a700c22ef4":"global","79801189f2ab72d269b633cdea173b20":"global","bd087e8d604cef993f957c859e37f283":"global"}},"owner":"a/3d0b3905d569f2c9c231b2a700c22ef4","restrict":false}},{"children_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/a802af12-7094-11e7-8cf7-a6006ad3dba0/%2A","created":"2017-08-08T20:03:51.77Z","catalog_crn":"crn:v1:staging:public:resource-catalog::a/ac5bc3a3a4f1dd065f5c92b444be3d01::service:a802af12-7094-11e7-8cf7-a6006ad3dba0","geo_tags":["eu-gb","us-south"],"id":"a802af12-7094-11e7-8cf7-a6006ad3dba0","images":{"feature_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant64.png","image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant50.png","medium_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant32.png","small_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant24.png"},"kind":"service","metadata":{"callbacks":{"broker_url":"https://dev-resource-catalog.stage1.ng.bluemix.net/api/v1/a802af12-7094-11e7-8cf7-a6006ad3dba0?include=%2A&languages=en-US%2Cen%3Bq%3D0.8"},"service":{"active":false,"bindable":false,"cf_plan_guid":{"stage1.eu-gb.bluemix.net":"9bf099a2-7b22-4b32-a165-42e99d0286e8","stage1.ng.bluemix.net":"b1a78e15-b43f-40c5-821a-787dc9523409"},"iam_compatible":true,"plan_updateable":false,"rc_compatible":true,"rc_provisionable":true,"service_check_enabled":false,"state":"","test_check_interval":0,"unique_api_key":false},"ui":{"allow_internal_users":false,"urls":{}}},"name":"Automation test","overview_ui":{"en":{"description":"rc compatible broker.","display_name":"rc compatible broker","long_description":"rc compatible broker"}},"pricing_tags":["free","paygo"],"provider":{"email":"sales@cloudant.com","name":"IBM"},"tags":["ibm_created"],"updated":"2017-09-13T00:54:15.522Z","url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/a802af12-7094-11e7-8cf7-a6006ad3dba0?include=%2A&languages=en-US%2Cen%3Bq%3D0.8%2Czh-CN%3Bq%3D0.6%2Czh%3Bq%3D0.4","visibility":{"exclude":{},"include":{"accounts":{"79801189f2ab72d269b633cdea173b20":"global","ac5bc3a3a4f1dd065f5c92b444be3d01":"global","bd087e8d604cef993f957c859e37f283":"global"}},"owner":"a/ac5bc3a3a4f1dd065f5c92b444be3d01","restrict":false}},{"children_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/bc4b35b4-70cd-11e7-8cf7-a6006ad3dba0/%2A","created":"2017-08-08T19:51:22.336Z","catalog_crn":"crn:v1:staging:public:resource-catalog::a/ac5bc3a3a4f1dd065f5c92b444be3d01::service:bc4b35b4-70cd-11e7-8cf7-a6006ad3dba0","geo_tags":["eu-gb","us-south"],"id":"bc4b35b4-70cd-11e7-8cf7-a6006ad3dba0","images":{"feature_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant64.png","image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant50.png","medium_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant32.png","small_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant24.png"},"kind":"service","metadata":{"callbacks":{"broker_url":"https://dev-resource-catalog.stage1.ng.bluemix.net/api/v1/bc4b35b4-70cd-11e7-8cf7-a6006ad3dba0?include=%2A&languages=en-US%2Cen%3Bq%3D0.8"},"service":{"active":false,"bindable":false,"cf_plan_guid":{"stage1.eu-gb.bluemix.net":"9bf099a2-7b22-4b32-a165-42e99d0286e8","stage1.ng.bluemix.net":"b1a78e15-b43f-40c5-821a-787dc9523409"},"iam_compatible":false,"plan_updateable":false,"rc_compatible":false,"rc_provisionable":true,"service_check_enabled":false,"state":"","test_check_interval":0,"unique_api_key":false},"ui":{"allow_internal_users":false,"urls":{}}},"name":"CF Automation Test","overview_ui":{"en":{"description":"cf compatible broker.","display_name":"cf compatible broker","long_description":"cf compatible broker"}},"pricing_tags":["free","paygo"],"provider":{"email":"sales@cloudant.com","name":"IBM"},"tags":["ibm_created"],"updated":"2017-08-09T19:07:27.852Z","url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/bc4b35b4-70cd-11e7-8cf7-a6006ad3dba0?include=%2A&languages=en-US%2Cen%3Bq%3D0.8%2Czh-CN%3Bq%3D0.6%2Czh%3Bq%3D0.4","visibility":{"exclude":{},"include":{"accounts":{"79801189f2ab72d269b633cdea173b20":"global","ac5bc3a3a4f1dd065f5c92b444be3d01":"global","bd087e8d604cef993f957c859e37f283":"global"}},"owner":"a/ac5bc3a3a4f1dd065f5c92b444be3d01","restrict":false}},{"children_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/d2619a98-86c9-11e7-bb31-be2e44b06b34/%2A","created":"2017-09-08T22:01:03.548Z","catalog_crn":"crn:v1:staging:public:resource-catalog::a/ac5bc3a3a4f1dd065f5c92b444be3d01::service:d2619a98-86c9-11e7-bb31-be2e44b06b34","geo_tags":["eu-gb","us-south"],"id":"d2619a98-86c9-11e7-bb31-be2e44b06b34","images":{"feature_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant64.png","image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant50.png","medium_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant32.png","small_image":"https://cloudantbroker.stage1.ng.bluemix.net/cloudant24.png"},"kind":"service","metadata":{"callbacks":{"broker_url":"https://dev-resource-catalog.stage1.ng.bluemix.net/api/v1/d2619a98-86c9-11e7-bb31-be2e44b06b34?include=%2A&languages=en-US%2Cen%3Bq%3D0.8"},"service":{"active":false,"bindable":false,"iam_compatible":true,"plan_updateable":false,"rc_compatible":true,"rc_provisionable":true,"service_check_enabled":false,"state":"","test_check_interval":0,"unique_api_key":false},"ui":{"allow_internal_users":false,"urls":{}}},"name":"RcAutomation-Test-Validation","overview_ui":{"en":{"description":"rc compatible broker.","display_name":"rc compatible broker","long_description":"rc compatible broker"}},"pricing_tags":["free"],"provider":{"email":"sales@cloudant.com","name":"IBM"},"tags":["ibm_created"],"updated":"2017-09-08T22:02:55.265Z","url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/d2619a98-86c9-11e7-bb31-be2e44b06b34?include=%2A&languages=en-US%2Cen%3Bq%3D0.8%2Czh-CN%3Bq%3D0.6%2Czh%3Bq%3D0.4","visibility":{"exclude":{},"include":{"accounts":{"79801189f2ab72d269b633cdea173b20":"global","ac5bc3a3a4f1dd065f5c92b444be3d01":"global","bd087e8d604cef993f957c859e37f283":"global"}},"owner":"a/ac5bc3a3a4f1dd065f5c92b444be3d01","restrict":false}}]}
							`),
					),
				)
			})
			It("should return not found error", func() {
				repo := newTestResourceCatalogRepo(server.URL())
				resources, err := repo.FindByName("Automation test", true)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resources).Should(Equal([]models.Service{{
					ID:       "a802af12-7094-11e7-8cf7-a6006ad3dba0",
					CRN:      "crn:v1:staging:public:resource-catalog::a/ac5bc3a3a4f1dd065f5c92b444be3d01::service:a802af12-7094-11e7-8cf7-a6006ad3dba0",
					Name:     "Automation test",
					Kind:     "service",
					URL:      "https://resource-catalog.stage1.ng.bluemix.net/api/v1/a802af12-7094-11e7-8cf7-a6006ad3dba0?include=%2A&languages=en-US%2Cen%3Bq%3D0.8%2Czh-CN%3Bq%3D0.6%2Czh%3Bq%3D0.4",
					Metadata: []byte(`{"callbacks":{"broker_url":"https://dev-resource-catalog.stage1.ng.bluemix.net/api/v1/a802af12-7094-11e7-8cf7-a6006ad3dba0?include=%2A&languages=en-US%2Cen%3Bq%3D0.8"},"service":{"active":false,"bindable":false,"cf_plan_guid":{"stage1.eu-gb.bluemix.net":"9bf099a2-7b22-4b32-a165-42e99d0286e8","stage1.ng.bluemix.net":"b1a78e15-b43f-40c5-821a-787dc9523409"},"iam_compatible":true,"plan_updateable":false,"rc_compatible":true,"rc_provisionable":true,"service_check_enabled":false,"state":"","test_check_interval":0,"unique_api_key":false},"ui":{"allow_internal_users":false,"urls":{}}}`),
				}}))
			})
		})

		Context("When there is error", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v1/", "q=test"),
						ghttp.RespondWith(http.StatusUnauthorized, `Token expired`),
					),
				)
			})
			It("should return not found error", func() {
				repo := newTestResourceCatalogRepo(server.URL())
				resources, err := repo.FindByName("test", false)

				Expect(err).Should(HaveOccurred())
				Expect(resources).Should(BeEmpty())
			})
		})

		Context("When there are two resources with same name", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v1/", "q=cloud-object-storage&include=*"),
						ghttp.RespondWith(http.StatusOK, `
							{"offset":0,"limit":50,"count":1,"resource_count":1,"first":"https://resource-catalog.stage1.ng.bluemix.net/api/v1?include=%2A&languages=en_US%2Cen&q=cloud-object-storage","resources":[{"children":[{"children_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/dff97f5c-bc5e-4455-b470-411c3edbe49c/%2A","created":"2017-09-21T18:25:05.176Z","catalog_crn":"crn:v1:staging:public:resource-catalog::a/9e16d1fed8aa7e1bd73e7a9d23434a5a::iaas:dff97f5c-bc5e-4455-b470-411c3edbe49c","geo_tags":["us-south-aus01","us-south-dal10","us-south-dal12","us-south-dal13"],"id":"dff97f5c-bc5e-4455-b470-411c3edbe49c","images":{"image":"https://sl-catalogapi-opsconsoledev.stage1.ng.bluemix.net/cache/99d-1554309718/assets/icons/object-storage.png"},"kind":"iaas","metadata":{"callbacks":{"broker_proxy_url":"https://console.ng.bluemix.net","broker_url":"https://console.ng.bluemix.net"},"other":{"dashboard":{"serviceCustomNavigationItems":[{"i18n":{"en":{"label":"Buckets and objects"}},"id":"manage","subItems":[{"i18n":{"en":{"label":"Bucket overview"}},"id":"bucket_overview","url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Bucket permissions"}},"id":"bucket_permissions","url":"http://console.stage1.bluemix.net/objectstorage/"}],"url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Endpoint"}},"id":"endpoint","url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Usage details"}},"id":"usage","url":"http://console.stage1.bluemix.net/objectstorage/"}],"serviceNavigationOrder":["gettingStarted","manage","endpoint","credentials","connectedObjects","usage","plan"]}},"pricing":{"costs":[{"currencies":[{"amount":{"USD":"0.00"},"country":"USA"}],"part_number":"","tier_model":"","tier_quantity":"1.000","unit":"HOURS_PER_MONTH","unit_id":"HOURS_PER_MONTH","unit_quantity":""}],"ibm_pricing":false},"service":{"active":false,"async_provisioning_supported":false,"async_unprovisioning_supported":false,"bindable":false,"cf_service_name":"cos-temp-name","custom_create_page_hybrid_enabled":false,"extension":null,"iam_compatible":true,"parameters":null,"plan_updateable":true,"rc_compatible":true,"rc_provisionable":false,"service_check_enabled":false,"service_key_supported":true,"state":"","test_check_interval":0,"type":"","unique_api_key":false,"user_defined_service":null},"sla":{"dr":{"dr":false},"tenancy":"single_tenant","terms":"http://term.condition.url"},"ui":{"side_by_side_index":1,"strings":{"en":{"bullets":[{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/api.svg","title":"S3 API (limited tooling)"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/regions.svg","title":"Regional and Cross Regional Resiliency"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/plans.svg","title":"Lite & Premium Plans"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/data.svg","title":"Flexible data classes cost optimized for your workload needs"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/security.svg","title":"Built in Security for data at rest and in transit"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/binding.svg","title":"Bluemix Application Binding"}]}},"urls":{"create_url":"https://console.stage1.bluemix.net/catalog/infrastructure/cloud-object-storage","doc_url":"https://www.stage1.ng.bluemix.net/docs/services/cloud-object-storage/getting-started.html","instructions_url":"https://www.stage1.ng.bluemix.net/docs/services/cloud-object-storage/getting-started.html","terms_url":"https://www.ibm.com/software/sla/sladb.nsf/sla/bm-7230-01"}}},"name":"cloud-object-storage","overview_ui":{"en":{"description":"Our most recommended service constantly adding robust features","display_name":"Cloud Object Storage","long_description":"Access your unstructured data from anywhere in the world with industry standard RESTful APIs, or via a self-service web-based portal. Enterprise availability and security lets your applications scale with confidence. You can deploy stand-alone services, or seamlessly integrate with other IBM Cloud and Bluemix services, including analytics, computing, and cognitive applications."}},"parent_id":"object-storage-group","parent_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/object-storage-group","pricing_tags":["paid_only","paygo"],"provider":{"name":"IBM"},"tags":["g2","ibm_created","lux","storage"],"updated":"2017-09-25T22:47:37.079Z","url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/dff97f5c-bc5e-4455-b470-411c3edbe49c?include=%2A&languages=en_US%2Cen","visibility":{"exclude":{},"include":{"accounts":{"11bf145efcc3081d1706ffa5734f0672":"global","249d926cb9bb46119e3bda8ed7f7f618":"global","265d9d22597d4ee589138929093f1246":"global","2a024d99093d4c3fc409e25d553ecb21":"global","2fdf0c082d324f4684b532e0c92fffd8":"global","3d9e32cc8e2c8c33026a8717256616e6":"global","49abfb9d0abddaa5f15d6e1f37e1ac58":"global","50860b903d7d78346258b14aa2667a23":"global","999c0f832ad0cd5b7aad7a26fcb92680":"global","9e16d1fed8aa7e1bd73e7a9d23434a5a":"global","a00c731cff012576d0612c0c811772e4":"global","a252ff3406dd29d7f77d5cf5807990d7":"global","a6032862a2dd327dedae0a2159f818cd":"global","b09edf5642ebfad587c594f4d4a354b0":"global","c32adbc8c683c3f39b3e643ae19993f5":"global","ccc74502c2d2c438763eef92e0291d24":"global","d4a70ce3b153a2e5a0e1bbf5569867ba":"global","d7026b374f39f570d20984c19c6ecf63":"global","d939d87aaebc9ee4339e36e4ccc67a0a":"global","df395fb648a606d256db06262000901c":"global","f3d1407c99e9fd93246225c90067903f":"global","fab6b5b09c41c4076a01d016f0e85b47":"global"}},"owner":"a/9e16d1fed8aa7e1bd73e7a9d23434a5a","restrict":false}}],"children_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/object-storage-group/%2A","created":"2017-09-21T18:21:33.96Z","catalog_crn":"crn:v1:staging:public:resource-catalog::a/2a024d99093d4c3fc409e25d553ecb21::iaas:object-storage-group","geo_tags":["us-south-aus01","us-south-dal10","us-south-dal12","us-south-dal13"],"group":true,"id":"object-storage-group","images":{"image":"https://sl-catalogapi-opsconsoledev.stage1.ng.bluemix.net/cache/99d-1554309718/assets/icons/object-storage.png"},"kind":"iaas","metadata":{"pricing":{"costs":[{"currencies":null,"part_number":"","tier_model":"","tier_quantity":"","unit":"","unit_id":"","unit_quantity":""}],"ibm_pricing":false},"ui":{"urls":{"create_url":"https://tbd.sidebyside.comparision"}}},"name":"object-storage-group","overview_ui":{"en":{"description":"Provides flexible, cost-effective, and scalable cloud storage for unstructured data.","display_name":"Cloud Object Storage","long_description":"Access your unstructured data from anywhere in the world with industry standard RESTful APIs, or via a self-service web-based portal. Enterprise availability and security lets your applications scale with confidence. You can deploy stand-alone services, or seamlessly integrate with other IBM Cloud and Bluemix services, including analytics, computing and cognitive applications."}},"pricing_tags":["paygo"],"provider":{"name":"IBM"},"tags":["group","ibm_created","storage"],"updated":"2017-09-21T18:23:31.427Z","url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/object-storage-group?include=%2A&languages=en_US%2Cen","visibility":{"exclude":{},"include":{"accounts":{"2a024d99093d4c3fc409e25d553ecb21":"global","a00c731cff012576d0612c0c811772e4":"global","d939d87aaebc9ee4339e36e4ccc67a0a":"global"}},"owner":"a/2a024d99093d4c3fc409e25d553ecb21","restrict":false}},{"children":[{"children_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/dff97f5c-bc5e-4455-b470-411c3edbe49c/%2A","created":"2017-09-21T18:25:05.176Z","catalog_crn":"crn:v1:staging:public:resource-catalog::a/9e16d1fed8aa7e1bd73e7a9d23434a5a::iaas:dff97f5c-bc5e-4455-b470-411c3edbe49c","geo_tags":["us-south-aus01","us-south-dal10","us-south-dal12","us-south-dal13"],"id":"dff97f5c-bc5e-4455-b470-411c3edbe49d","images":{"image":"https://sl-catalogapi-opsconsoledev.stage1.ng.bluemix.net/cache/99d-1554309718/assets/icons/object-storage.png"},"kind":"iaas","metadata":{"callbacks":{"broker_proxy_url":"https://console.ng.bluemix.net","broker_url":"https://console.ng.bluemix.net"},"other":{"dashboard":{"serviceCustomNavigationItems":[{"i18n":{"en":{"label":"Buckets and objects"}},"id":"manage","subItems":[{"i18n":{"en":{"label":"Bucket overview"}},"id":"bucket_overview","url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Bucket permissions"}},"id":"bucket_permissions","url":"http://console.stage1.bluemix.net/objectstorage/"}],"url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Endpoint"}},"id":"endpoint","url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Usage details"}},"id":"usage","url":"http://console.stage1.bluemix.net/objectstorage/"}],"serviceNavigationOrder":["gettingStarted","manage","endpoint","credentials","connectedObjects","usage","plan"]}},"pricing":{"costs":[{"currencies":[{"amount":{"USD":"0.00"},"country":"USA"}],"part_number":"","tier_model":"","tier_quantity":"1.000","unit":"HOURS_PER_MONTH","unit_id":"HOURS_PER_MONTH","unit_quantity":""}],"ibm_pricing":false},"service":{"active":false,"async_provisioning_supported":false,"async_unprovisioning_supported":false,"bindable":false,"cf_service_name":"cos-temp-name","custom_create_page_hybrid_enabled":false,"extension":null,"iam_compatible":true,"parameters":null,"plan_updateable":true,"rc_compatible":true,"rc_provisionable":false,"service_check_enabled":false,"service_key_supported":true,"state":"","test_check_interval":0,"type":"","unique_api_key":false,"user_defined_service":null},"sla":{"dr":{"dr":false},"tenancy":"single_tenant","terms":"http://term.condition.url"},"ui":{"side_by_side_index":1,"strings":{"en":{"bullets":[{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/api.svg","title":"S3 API (limited tooling)"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/regions.svg","title":"Regional and Cross Regional Resiliency"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/plans.svg","title":"Lite & Premium Plans"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/data.svg","title":"Flexible data classes cost optimized for your workload needs"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/security.svg","title":"Built in Security for data at rest and in transit"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/binding.svg","title":"Bluemix Application Binding"}]}},"urls":{"create_url":"https://console.stage1.bluemix.net/catalog/infrastructure/cloud-object-storage","doc_url":"https://www.stage1.ng.bluemix.net/docs/services/cloud-object-storage/getting-started.html","instructions_url":"https://www.stage1.ng.bluemix.net/docs/services/cloud-object-storage/getting-started.html","terms_url":"https://www.ibm.com/software/sla/sladb.nsf/sla/bm-7230-01"}}},"name":"cloud-object-storage","overview_ui":{"en":{"description":"Our most recommended service constantly adding robust features","display_name":"Cloud Object Storage","long_description":"Access your unstructured data from anywhere in the world with industry standard RESTful APIs, or via a self-service web-based portal. Enterprise availability and security lets your applications scale with confidence. You can deploy stand-alone services, or seamlessly integrate with other IBM Cloud and Bluemix services, including analytics, computing, and cognitive applications."}},"parent_id":"object-storage-group","parent_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/object-storage-group","pricing_tags":["paid_only","paygo"],"provider":{"name":"IBM"},"tags":["g2","ibm_created","lux","storage"],"updated":"2017-09-25T22:47:37.079Z","url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/dff97f5c-bc5e-4455-b470-411c3edbe49c?include=%2A&languages=en_US%2Cen","visibility":{"exclude":{},"include":{"accounts":{"11bf145efcc3081d1706ffa5734f0672":"global","249d926cb9bb46119e3bda8ed7f7f618":"global","265d9d22597d4ee589138929093f1246":"global","2a024d99093d4c3fc409e25d553ecb21":"global","2fdf0c082d324f4684b532e0c92fffd8":"global","3d9e32cc8e2c8c33026a8717256616e6":"global","49abfb9d0abddaa5f15d6e1f37e1ac58":"global","50860b903d7d78346258b14aa2667a23":"global","999c0f832ad0cd5b7aad7a26fcb92680":"global","9e16d1fed8aa7e1bd73e7a9d23434a5a":"global","a00c731cff012576d0612c0c811772e4":"global","a252ff3406dd29d7f77d5cf5807990d7":"global","a6032862a2dd327dedae0a2159f818cd":"global","b09edf5642ebfad587c594f4d4a354b0":"global","c32adbc8c683c3f39b3e643ae19993f5":"global","ccc74502c2d2c438763eef92e0291d24":"global","d4a70ce3b153a2e5a0e1bbf5569867ba":"global","d7026b374f39f570d20984c19c6ecf63":"global","d939d87aaebc9ee4339e36e4ccc67a0a":"global","df395fb648a606d256db06262000901c":"global","f3d1407c99e9fd93246225c90067903f":"global","fab6b5b09c41c4076a01d016f0e85b47":"global"}},"owner":"a/9e16d1fed8aa7e1bd73e7a9d23434a5a","restrict":false}}],"children_url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/object-storage-group/%2A","created":"2017-09-21T18:21:33.96Z","catalog_crn":"crn:v1:staging:public:resource-catalog::a/2a024d99093d4c3fc409e25d553ecb21::iaas:object-storage-group","geo_tags":["us-south-aus01","us-south-dal10","us-south-dal12","us-south-dal13"],"group":true,"id":"object-storage-group","images":{"image":"https://sl-catalogapi-opsconsoledev.stage1.ng.bluemix.net/cache/99d-1554309718/assets/icons/object-storage.png"},"kind":"iaas","metadata":{"pricing":{"costs":[{"currencies":null,"part_number":"","tier_model":"","tier_quantity":"","unit":"","unit_id":"","unit_quantity":""}],"ibm_pricing":false},"ui":{"urls":{"create_url":"https://tbd.sidebyside.comparision"}}},"name":"object-storage-group2","overview_ui":{"en":{"description":"Provides flexible, cost-effective, and scalable cloud storage for unstructured data.","display_name":"Cloud Object Storage","long_description":"Access your unstructured data from anywhere in the world with industry standard RESTful APIs, or via a self-service web-based portal. Enterprise availability and security lets your applications scale with confidence. You can deploy stand-alone services, or seamlessly integrate with other IBM Cloud and Bluemix services, including analytics, computing and cognitive applications."}},"pricing_tags":["paygo"],"provider":{"name":"IBM"},"tags":["group","ibm_created","storage"],"updated":"2017-09-21T18:23:31.427Z","url":"https://resource-catalog.stage1.ng.bluemix.net/api/v1/object-storage-group?include=%2A&languages=en_US%2Cen","visibility":{"exclude":{},"include":{"accounts":{"2a024d99093d4c3fc409e25d553ecb21":"global","a00c731cff012576d0612c0c811772e4":"global","d939d87aaebc9ee4339e36e4ccc67a0a":"global"}},"owner":"a/2a024d99093d4c3fc409e25d553ecb21","restrict":false}}]}`),
					),
				)
			})
			It("should return two resources", func() {
				repo := newTestResourceCatalogRepo(server.URL())
				resources, err := repo.FindByName("cloud-object-storage", true)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resources).Should(Equal([]models.Service{{
					ID:       "dff97f5c-bc5e-4455-b470-411c3edbe49c",
					CRN:      "crn:v1:staging:public:resource-catalog::a/9e16d1fed8aa7e1bd73e7a9d23434a5a::iaas:dff97f5c-bc5e-4455-b470-411c3edbe49c",
					Name:     "cloud-object-storage",
					Kind:     "iaas",
					URL:      "https://resource-catalog.stage1.ng.bluemix.net/api/v1/dff97f5c-bc5e-4455-b470-411c3edbe49c?include=%2A&languages=en_US%2Cen",
					Metadata: []byte(`{"callbacks":{"broker_proxy_url":"https://console.ng.bluemix.net","broker_url":"https://console.ng.bluemix.net"},"other":{"dashboard":{"serviceCustomNavigationItems":[{"i18n":{"en":{"label":"Buckets and objects"}},"id":"manage","subItems":[{"i18n":{"en":{"label":"Bucket overview"}},"id":"bucket_overview","url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Bucket permissions"}},"id":"bucket_permissions","url":"http://console.stage1.bluemix.net/objectstorage/"}],"url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Endpoint"}},"id":"endpoint","url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Usage details"}},"id":"usage","url":"http://console.stage1.bluemix.net/objectstorage/"}],"serviceNavigationOrder":["gettingStarted","manage","endpoint","credentials","connectedObjects","usage","plan"]}},"pricing":{"costs":[{"currencies":[{"amount":{"USD":"0.00"},"country":"USA"}],"part_number":"","tier_model":"","tier_quantity":"1.000","unit":"HOURS_PER_MONTH","unit_id":"HOURS_PER_MONTH","unit_quantity":""}],"ibm_pricing":false},"service":{"active":false,"async_provisioning_supported":false,"async_unprovisioning_supported":false,"bindable":false,"cf_service_name":"cos-temp-name","custom_create_page_hybrid_enabled":false,"extension":null,"iam_compatible":true,"parameters":null,"plan_updateable":true,"rc_compatible":true,"rc_provisionable":false,"service_check_enabled":false,"service_key_supported":true,"state":"","test_check_interval":0,"type":"","unique_api_key":false,"user_defined_service":null},"sla":{"dr":{"dr":false},"tenancy":"single_tenant","terms":"http://term.condition.url"},"ui":{"side_by_side_index":1,"strings":{"en":{"bullets":[{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/api.svg","title":"S3 API (limited tooling)"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/regions.svg","title":"Regional and Cross Regional Resiliency"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/plans.svg","title":"Lite & Premium Plans"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/data.svg","title":"Flexible data classes cost optimized for your workload needs"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/security.svg","title":"Built in Security for data at rest and in transit"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/binding.svg","title":"Bluemix Application Binding"}]}},"urls":{"create_url":"https://console.stage1.bluemix.net/catalog/infrastructure/cloud-object-storage","doc_url":"https://www.stage1.ng.bluemix.net/docs/services/cloud-object-storage/getting-started.html","instructions_url":"https://www.stage1.ng.bluemix.net/docs/services/cloud-object-storage/getting-started.html","terms_url":"https://www.ibm.com/software/sla/sladb.nsf/sla/bm-7230-01"}}}`),
				},
					{
						ID:       "dff97f5c-bc5e-4455-b470-411c3edbe49d",
						CRN:      "crn:v1:staging:public:resource-catalog::a/9e16d1fed8aa7e1bd73e7a9d23434a5a::iaas:dff97f5c-bc5e-4455-b470-411c3edbe49c",
						Name:     "cloud-object-storage",
						Kind:     "iaas",
						URL:      "https://resource-catalog.stage1.ng.bluemix.net/api/v1/dff97f5c-bc5e-4455-b470-411c3edbe49c?include=%2A&languages=en_US%2Cen",
						Metadata: []byte(`{"callbacks":{"broker_proxy_url":"https://console.ng.bluemix.net","broker_url":"https://console.ng.bluemix.net"},"other":{"dashboard":{"serviceCustomNavigationItems":[{"i18n":{"en":{"label":"Buckets and objects"}},"id":"manage","subItems":[{"i18n":{"en":{"label":"Bucket overview"}},"id":"bucket_overview","url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Bucket permissions"}},"id":"bucket_permissions","url":"http://console.stage1.bluemix.net/objectstorage/"}],"url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Endpoint"}},"id":"endpoint","url":"http://console.stage1.bluemix.net/objectstorage/"},{"i18n":{"en":{"label":"Usage details"}},"id":"usage","url":"http://console.stage1.bluemix.net/objectstorage/"}],"serviceNavigationOrder":["gettingStarted","manage","endpoint","credentials","connectedObjects","usage","plan"]}},"pricing":{"costs":[{"currencies":[{"amount":{"USD":"0.00"},"country":"USA"}],"part_number":"","tier_model":"","tier_quantity":"1.000","unit":"HOURS_PER_MONTH","unit_id":"HOURS_PER_MONTH","unit_quantity":""}],"ibm_pricing":false},"service":{"active":false,"async_provisioning_supported":false,"async_unprovisioning_supported":false,"bindable":false,"cf_service_name":"cos-temp-name","custom_create_page_hybrid_enabled":false,"extension":null,"iam_compatible":true,"parameters":null,"plan_updateable":true,"rc_compatible":true,"rc_provisionable":false,"service_check_enabled":false,"service_key_supported":true,"state":"","test_check_interval":0,"type":"","unique_api_key":false,"user_defined_service":null},"sla":{"dr":{"dr":false},"tenancy":"single_tenant","terms":"http://term.condition.url"},"ui":{"side_by_side_index":1,"strings":{"en":{"bullets":[{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/api.svg","title":"S3 API (limited tooling)"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/regions.svg","title":"Regional and Cross Regional Resiliency"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/plans.svg","title":"Lite & Premium Plans"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/data.svg","title":"Flexible data classes cost optimized for your workload needs"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/security.svg","title":"Built in Security for data at rest and in transit"},{"icon":"https://resource-catalog.stage1.ng.bluemix.net/static/images/g2/binding.svg","title":"Bluemix Application Binding"}]}},"urls":{"create_url":"https://console.stage1.bluemix.net/catalog/infrastructure/cloud-object-storage","doc_url":"https://www.stage1.ng.bluemix.net/docs/services/cloud-object-storage/getting-started.html","instructions_url":"https://www.stage1.ng.bluemix.net/docs/services/cloud-object-storage/getting-started.html","terms_url":"https://www.ibm.com/software/sla/sladb.nsf/sla/bm-7230-01"}}}`),
					}}))
			})
		})

	})
})

func newTestResourceCatalogRepo(url string) ResourceCatalogRepository {
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.ResourceCatalogrService,
	}

	return newResourceCatalogAPI(&client)
}
