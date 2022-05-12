package containerv2

import (
	"log"
	"net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/client"
	bluemixHttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/session"

	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clusters", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//List
	Describe("List", func() {
		Context("When read of clusters is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, ContainSubstring("/v2/vpc/getClusters")),
						ghttp.RespondWith(http.StatusOK, `[{
              "CreatedDate": "",
			  "DataCenter": "dal10",
			  "Entitlement": "",
              "ID": "f91adfe2-76c9-4649-939e-b01c37a3704",
              "IngressHostname": "",
              "IngressSecretName": "",
              "Location": "",
              "MasterKubeVersion": "1.8.1",
              "Prefix": "worker",
              "ModifiedDate": "",
              "Name": "test",
              "Region": "abc",
              "ServerURL": "",
              "State": "normal",
              "IsPaid": false,
              "IsTrusted": true,
              "WorkerCount": 1
              }]`),
					),

					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, ContainSubstring("/v2/satellite/getClusters")),
						ghttp.RespondWith(http.StatusOK, `[{
	              "CreatedDate": "",
				  "DataCenter": "dal10",
				  "Entitlement": "",
	              "ID": "d91adfe2-76c9-4649-939e-b01c37a3704",
	              "IngressHostname": "",
	              "IngressSecretName": "",
	              "Location": "",
	              "MasterKubeVersion": "1.8.1",
	              "Prefix": "worker",
	              "ModifiedDate": "",
	              "Name": "test",
	              "Region": "abc",
	              "ServerURL": "",
	              "State": "normal",
	              "IsPaid": false,
	              "IsTrusted": true,
	              "WorkerCount": 1
	              }]`),
					),
				)
			})

			It("should return cluster list", func() {
				target := ClusterTargetHeader{}
				myCluster, err := newCluster(server.URL()).List(target)
				Expect(myCluster).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
				Expect(len(myCluster)).Should(Equal(2))
				Expect(myCluster[0].ID).Should(Equal("f91adfe2-76c9-4649-939e-b01c37a3704"))
				Expect(myCluster[0].WorkerCount).Should(Equal(1))
				Expect(myCluster[0].MasterKubeVersion).Should(Equal("1.8.1"))
				Expect(myCluster[1].ID).Should(Equal("d91adfe2-76c9-4649-939e-b01c37a3704"))
				Expect(myCluster[1].WorkerCount).Should(Equal(1))
				Expect(myCluster[1].MasterKubeVersion).Should(Equal("1.8.1"))
			})
		})

		Context("When provider is satellite", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, ContainSubstring("/v2/satellite/getClusters")),
						ghttp.RespondWith(http.StatusOK, `[{
							"CreatedDate": "",
				"DataCenter": "dal10",
				"Entitlement": "",
							"ID": "d91adfe2-76c9-4649-939e-b01c37a3704",
							"IngressHostname": "",
							"IngressSecretName": "",
							"Location": "",
							"MasterKubeVersion": "1.8.1",
							"Prefix": "worker",
							"ModifiedDate": "",
							"Name": "test-satellite",
							"Region": "abc",
							"ServerURL": "",
							"State": "normal",
							"IsPaid": false,
							"IsTrusted": true,
							"WorkerCount": 1
							}]`),
					),
				)
			})

			It("should return only satellite cluster list", func() {
				target := ClusterTargetHeader{}
				target.Provider = "satellite"
				myCluster, err := newCluster(server.URL()).List(target)
				Expect(myCluster).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
				Expect(len(myCluster)).Should(Equal(1))
				Expect(myCluster[0].ID).Should(Equal("d91adfe2-76c9-4649-939e-b01c37a3704"))
				Expect(myCluster[0].Name).Should(Equal("test-satellite"))
				Expect(myCluster[0].WorkerCount).Should(Equal(1))
				Expect(myCluster[0].MasterKubeVersion).Should(Equal("1.8.1"))
			})
		})

		Context("When provider is vpc-classic", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getClusters", "provider=vpc-classic"),
						ghttp.RespondWith(http.StatusOK, `[{
									"CreatedDate": "",
						"DataCenter": "dal10",
						"Entitlement": "",
									"ID": "c91adfe2-76c9-4649-939e-b01c37a3704",
									"IngressHostname": "",
									"IngressSecretName": "",
									"Location": "",
									"MasterKubeVersion": "1.8.1",
									"Prefix": "worker",
									"ModifiedDate": "",
									"Name": "test-vpc-classic",
									"Region": "abc",
									"ServerURL": "",
									"State": "normal",
									"IsPaid": false,
									"IsTrusted": true,
									"WorkerCount": 1
									}]`),
					),
				)
			})

			It("should return only vpc-classic cluster list", func() {
				target := ClusterTargetHeader{}
				target.Provider = "vpc-classic"
				myCluster, err := newCluster(server.URL()).List(target)
				Expect(myCluster).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
				Expect(len(myCluster)).Should(Equal(1))
				Expect(myCluster[0].ID).Should(Equal("c91adfe2-76c9-4649-939e-b01c37a3704"))
				Expect(myCluster[0].Name).Should(Equal("test-vpc-classic"))
				Expect(myCluster[0].WorkerCount).Should(Equal(1))
				Expect(myCluster[0].MasterKubeVersion).Should(Equal("1.8.1"))
			})
		})
		Context("When read of clusters is unsuccessful", func() {
			BeforeEach(func() {

				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, ContainSubstring("/v2/")),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve clusters`),
					),
				)
			})

			It("should return error when cluster are retrieved", func() {
				target := ClusterTargetHeader{}
				myCluster, err := newCluster(server.URL()).List(target)
				Expect(err).To(HaveOccurred())
				Expect(myCluster).Should(BeNil())
			})
		})
	})

	//Create
	Describe("Create", func() {
		Context("When creation is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/vpc/createCluster"),
						ghttp.VerifyJSON(`{"disablePublicServiceEndpoint": false, "defaultWorkerPoolEntitlement": "", "kubeVersion": "", "podSubnet": "podnet", "provider": "abc", "serviceSubnet": "svcnet", "name": "abcd", "cosInstanceCRN": "", "workerPool": {"flavor": "", "hostPoolID": "hostpoolid", "name": "", "vpcID": "", "workerCount": 0, "zones": null, "entitlement": ""}}`),
						ghttp.RespondWith(http.StatusCreated, `{
							 "clusterID": "f91adfe2-76c9-4649-939e-b01c37a3704c"
						}`),
					),
				)
			})

			It("should return cluster created", func() {
				WPools := WorkerPoolConfig{
					Flavor: "", WorkerCount: 0, VpcID: "", Name: "",
					HostPoolID: "hostpoolid",
				}
				params := ClusterCreateRequest{
					DisablePublicServiceEndpoint: false, KubeVersion: "", PodSubnet: "podnet", Provider: "abc", ServiceSubnet: "svcnet", Name: "abcd", WorkerPools: WPools, CosInstanceCRN: "",
				}
				target := ClusterTargetHeader{}
				myCluster, err := newCluster(server.URL()).Create(params, target)
				Expect(err).NotTo(HaveOccurred())
				Expect(myCluster).ShouldNot(BeNil())
				Expect(myCluster.ID).Should(Equal("f91adfe2-76c9-4649-939e-b01c37a3704c"))
			})
		})
		Context("When creation is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/vpc/createCluster"),
						ghttp.VerifyJSON(`{"disablePublicServiceEndpoint": false, "defaultWorkerPoolEntitlement": "", "kubeVersion": "", "podSubnet": "podnet", "provider": "abc", "serviceSubnet": "svcnet", "name": "abcd", "cosInstanceCRN": "", "workerPool": {"flavor": "", "name": "", "vpcID": "", "workerCount": 0, "zones": null, "entitlement": ""}}`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to create cluster`),
					),
				)
			})
			It("should return error during cluster creation", func() {
				WPools := WorkerPoolConfig{
					Flavor: "", WorkerCount: 0, VpcID: "", Name: "", Entitlement: "",
				}
				params := ClusterCreateRequest{
					DisablePublicServiceEndpoint: false, KubeVersion: "", PodSubnet: "podnet", Provider: "abc", ServiceSubnet: "svcnet", Name: "abcd", WorkerPools: WPools, DefaultWorkerPoolEntitlement: "", CosInstanceCRN: "",
				}
				target := ClusterTargetHeader{}
				myCluster, err := newCluster(server.URL()).Create(params, target)
				Expect(err).To(HaveOccurred())
				Expect(myCluster.ID).Should(Equal(""))
			})
		})
		Context("When creating with kms enabled", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/vpc/createCluster"),
						ghttp.VerifyJSON(`{"disablePublicServiceEndpoint": false, "defaultWorkerPoolEntitlement": "", "kubeVersion": "", "podSubnet": "podnet", "provider": "abc", "serviceSubnet": "svcnet", "name": "abcd", "cosInstanceCRN": "", "workerPool": {"flavor": "", "name": "", "vpcID": "", "workerCount": 0, "zones": null, "entitlement": "", "workerVolumeEncryption": {"kmsInstanceID": "kmsid", "workerVolumeCRKID": "rootkeyid"}}}`),
						ghttp.RespondWith(http.StatusCreated, `{
							 "clusterID": "f91adfe2-76c9-4649-939e-b01c37a3704c"
						}`),
					),
				)
			})

			It("should return cluster created", func() {
				WVE := WorkerVolumeEncryption{KmsInstanceID: "kmsid", WorkerVolumeCRKID: "rootkeyid"}
				WPools := WorkerPoolConfig{
					Flavor: "", WorkerCount: 0, VpcID: "", Name: "", WorkerVolumeEncryption: &WVE,
				}
				params := ClusterCreateRequest{
					DisablePublicServiceEndpoint: false, KubeVersion: "", PodSubnet: "podnet", Provider: "abc", ServiceSubnet: "svcnet", Name: "abcd", WorkerPools: WPools, CosInstanceCRN: "",
				}
				target := ClusterTargetHeader{}
				myCluster, err := newCluster(server.URL()).Create(params, target)
				Expect(err).NotTo(HaveOccurred())
				Expect(myCluster).ShouldNot(BeNil())
				Expect(myCluster.ID).Should(Equal("f91adfe2-76c9-4649-939e-b01c37a3704c"))
			})
		})
	})
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
		ServiceName: bluemix.VpcContainerService,
	}
	return newClusterAPI(&client)
}
