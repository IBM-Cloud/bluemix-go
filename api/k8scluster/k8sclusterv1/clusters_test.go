package k8sclusterv1

import (
	"log"
	"net/http"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/client"
	bluemixHttp "github.com/IBM-Bluemix/bluemix-go/http"
	"github.com/IBM-Bluemix/bluemix-go/session"

	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clusters", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})
	Describe("Create", func() {
		Context("When creation is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v1/clusters"),
						ghttp.VerifyJSON(`{"Billing":"","Datacenter":"dal10","Isolation":"","MachineType":"free","Name":"testservice","PrivateVlan":"vlan","PublicVlan":"vlan","WorkerNum":0,"NoSubnet":false}
`),
						ghttp.RespondWith(http.StatusCreated, `{							 	
							 "id": "f91adfe2-76c9-4649-939e-b01c37a3704c"
						}`),
					),
				)
			})

			It("should return cluster created", func() {
				params := &ClusterCreateRequest{
					Name: "testservice", Datacenter: "dal10", MachineType: "free", PublicVlan: "vlan", PrivateVlan: "vlan",
				}
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				myCluster, err := newCluster(server.URL()).Create(params, target)
				Expect(err).NotTo(HaveOccurred())
				Expect(myCluster).ShouldNot(BeNil())
				Expect(myCluster.ID).Should(Equal("f91adfe2-76c9-4649-939e-b01c37a3704c"))
			})
		})
		Context("When creation is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v1/clusters"),
						ghttp.VerifyJSON(`{"Billing":"","Datacenter":"dal10","Isolation":"","MachineType":"free","Name":"testservice","PrivateVlan":"vlan","PublicVlan":"vlan","WorkerNum":0,"NoSubnet":false}
`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to create cluster`),
					),
				)
			})

			It("should return error during cluster creation", func() {
				params := &ClusterCreateRequest{
					Name: "testservice", Datacenter: "dal10", MachineType: "free", PublicVlan: "vlan", PrivateVlan: "vlan",
				}
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				myCluster, err := newCluster(server.URL()).Create(params, target)
				Expect(err).To(HaveOccurred())
				Expect(myCluster.ID).Should(Equal(""))
			})
		})
	})
	//List
	Describe("List", func() {
		Context("When read of clusters is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/clusters"),
						ghttp.RespondWith(http.StatusOK, `[{							 	
							 "GUID": "f91adfe2-76c9-4649-939e-b01c37a3704c",
              "CreatedDate": "",
              "DataCenter": "dal10",
              "ID": "f91adfe2-76c9-4649-939e-b01c37a3704",
              "IngressHostname": "",
              "IngressSecretName": "",
              "Location": "",
              "MasterKubeVersion": "",
              "ModifiedDate": "",
              "Name": "test",
              "Region": "abc",
              "ServerURL": "",
              "State": "normal",
              "IsPaid": false,
              "WorkerCount": 1
              }]`),
					),
				)
			})

			It("should return cluster list", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				myCluster, err := newCluster(server.URL()).List(target)
				Expect(myCluster).ShouldNot(BeNil())
				for _, cluster := range myCluster {
					Expect(err).NotTo(HaveOccurred())
					Expect(cluster.ID).Should(Equal("f91adfe2-76c9-4649-939e-b01c37a3704"))
					Expect(cluster.WorkerCount).Should(Equal(1))
				}
			})
		})
		Context("When read of clusters is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/clusters"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve clusters`),
					),
				)
			})

			It("should return error when cluster are retrieved", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				myCluster, err := newCluster(server.URL()).List(target)
				Expect(err).To(HaveOccurred())
				Expect(myCluster).Should(BeNil())
			})
		})
	})
	//Delete
	Describe("Delete", func() {
		Context("When delete of cluster is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v1/clusters/test"),
						ghttp.RespondWith(http.StatusOK, `{							
						}`),
					),
				)
			})

			It("should delete cluster", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				err := newCluster(server.URL()).Delete("test", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When cluster delete is failed", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v1/clusters/test"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to delete service key`),
					),
				)
			})

			It("should return error service key delete", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				err := newCluster(server.URL()).Delete("test", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})
	//Find
	Describe("Find", func() {
		Context("When read of cluster is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/clusters/test"),
						ghttp.RespondWith(http.StatusOK, `{							 	
							 "GUID": "f91adfe2-76c9-4649-939e-b01c37a3704c",
              "CreatedDate": "",
              "DataCenter": "dal10",
              "ID": "f91adfe2-76c9-4649-939e-b01c37a3704",
              "IngressHostname": "",
              "IngressSecretName": "",
              "Location": "",
              "MasterKubeVersion": "",
              "ModifiedDate": "",
              "Name": "test",
              "Region": "abc",
              "ServerURL": "",
              "State": "normal",
              "IsPaid": false,
              "WorkerCount": 1}`),
					),
				)
			})

			It("should return cluster", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				myCluster, err := newCluster(server.URL()).Find("test", target)
				Expect(err).NotTo(HaveOccurred())
				Expect(myCluster).ShouldNot(BeNil())
				Expect(myCluster.GUID).Should(Equal("f91adfe2-76c9-4649-939e-b01c37a3704c"))

			})
		})
		Context("When cluster retrieve is failed", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/clusters/test"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve cluster`),
					),
				)
			})

			It("should return error when cluster is retrieved", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				myCluster, err := newCluster(server.URL()).Find("test", target)
				Expect(err).To(HaveOccurred())
				Expect(myCluster.ID).Should(Equal(""))
			})
		})
	})
	//set credentials
	Describe("set credentials", func() {
		Context("When credential set is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v1/credentials"),
						ghttp.RespondWith(http.StatusOK, `{}`),
					),
				)
			})

			It("should set credentials", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				err := newCluster(server.URL()).SetCredentials("test", "abcdef-df-fg", target)
				Expect(err).NotTo(HaveOccurred())

			})
		})
		Context("When credential set is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v1/credentials"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to set credentials`),
					),
				)
			})

			It("should throw error when setting credentials", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				err := newCluster(server.URL()).SetCredentials("test", "abcdef-df-fg", target)
				Expect(err).To(HaveOccurred())

			})
		})
	})
	//Unset credentials
	Describe("unset credentials", func() {
		Context("When unset credential is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v1/credentials"),
						ghttp.RespondWith(http.StatusOK, `{}`),
					),
				)
			})

			It("should set credentials", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				err := newCluster(server.URL()).UnsetCredentials(target)
				Expect(err).NotTo(HaveOccurred())

			})
		})
		Context("When unset credential is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v1/credentials"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to unset credentials`),
					),
				)
			})

			It("should set credentials", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				err := newCluster(server.URL()).UnsetCredentials(target)
				Expect(err).To(HaveOccurred())

			})
		})
	})
	//Bind service
	Describe("Bind service", func() {
		Context("When bind service is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v1/clusters/test/services"),
						ghttp.RespondWith(http.StatusOK, `{}`),
					),
				)
			})

			It("should bind service to a cluster", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				params := &ServiceBindRequest{
					ClusterNameOrID: "test", SpaceGUID: "ffed-ret-534-ghrk", ServiceInstanceNameOrID: "cloudantDB", NamespaceID: "default"}
				serviceResp, err := newCluster(server.URL()).BindService(params, target)
				Expect(err).NotTo(HaveOccurred())
				Expect(serviceResp).ShouldNot(BeNil())
			})
		})
		Context("When bind service is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v1/clusters/test/services"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to set credentials`),
					),
				)
			})

			It("should throw error when binding service to a cluster", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				params := &ServiceBindRequest{
					ClusterNameOrID: "test", SpaceGUID: "ffed-ret-534-ghrk", ServiceInstanceNameOrID: "cloudantDB", NamespaceID: "default"}
				serviceResp, err := newCluster(server.URL()).BindService(params, target)
				Expect(err).To(HaveOccurred())
				Expect(serviceResp.ServiceInstanceGUID).Should(Equal(""))
				Expect(serviceResp.SecretName).Should(Equal(""))
			})
		})
	})
	//Unbind service
	Describe("UnBind service", func() {
		Context("When bind service is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v1/clusters/test/services/default/cloudantDB"),
						ghttp.RespondWith(http.StatusOK, `{}`),
					),
				)
			})

			It("should bind service to a cluster", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				err := newCluster(server.URL()).UnBindService("test", "default", "cloudantDB", target)
				Expect(err).NotTo(HaveOccurred())

			})
		})
		Context("When unbind service is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v1/clusters/test/services/default/cloudantDB"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to unbind service`),
					),
				)
			})

			It("should set credentials", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				err := newCluster(server.URL()).UnBindService("test", "default", "cloudantDB", target)
				Expect(err).To(HaveOccurred())

			})
		})
	})
	//List bound services
	Describe("ListServicesBoundToCluster", func() {
		Context("When read of cluster services is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/clusters/test/services/default"),
						ghttp.RespondWith(http.StatusOK, `[{							 	
              "ServiceName": "testService",
              "ServiceID": "f91adfe2-76c9-4649-939e-b01c37a3704",
              "ServiceKeyName": "kube-testService",
              "Namespace": "default"
              }]`),
					),
				)
			})

			It("should return cluster service list", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				boundServices, err := newCluster(server.URL()).ListServicesBoundToCluster("test", "default", target)
				Expect(boundServices).ShouldNot(BeNil())
				for _, service := range boundServices {
					Expect(err).NotTo(HaveOccurred())
					Expect(service.ServiceName).Should(Equal("testService"))
					Expect(service.ServiceID).Should(Equal("f91adfe2-76c9-4649-939e-b01c37a3704"))
				}
			})
		})
		Context("When read of cluster services is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/clusters/test/services/default"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve clusters`),
					),
				)
			})

			It("should return error when cluster services are retrieved", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				service, err := newCluster(server.URL()).ListServicesBoundToCluster("test", "default", target)
				Expect(err).To(HaveOccurred())
				Expect(service).Should(BeNil())
			})
		})
	})
	//Find Cluster service
	Describe("FindServiceBoundToClusters", func() {
		Context("When read a service bound to cluster is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/clusters/test/services/default"),
						ghttp.RespondWith(http.StatusOK, `[{							 	
              "ServiceName": "testService",
              "ServiceID": "f91adfe2-76c9-4649-939e-b01c37a3704",
              "ServiceKeyName": "kube-testService",
              "Namespace": "default"
              }]`),
					),
				)
			})

			It("should return cluster service list", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				boundService, err := newCluster(server.URL()).FindServiceBoundToCluster("test", "f91adfe2-76c9-4649-939e-b01c37a3704", "default", target)
				Expect(boundService).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
				Expect(boundService.ServiceName).Should(Equal("testService"))
				Expect(boundService.ServiceID).Should(Equal("f91adfe2-76c9-4649-939e-b01c37a3704"))
			})
		})
		Context("When read of cluster services is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v1/clusters/test/services/default"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve clusters`),
					),
				)
			})

			It("should return error when cluster services are retrieved", func() {
				target := &ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				_, err := newCluster(server.URL()).FindServiceBoundToCluster("test", "f91adfe2-76c9-4649-939e-b01c37a3704", "default", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})
	//
})

func newCluster(url string) Clusters {

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.CfService,
	}
	return newClusterAPI(&client)
}
