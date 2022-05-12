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

var _ = Describe("dedicatedhostpools", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	Describe("Create", func() {
		Context("When creating dedicatedhostpool is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/createDedicatedHostPool"),
						ghttp.VerifyJSON(`{
							"flavorClass": "flavorclass1",
							"metro": "metro1",
							"name": "name1"
						  }`),
						ghttp.RespondWith(http.StatusCreated, `{
							"id":"dedicatedhostpoolid1"
						}`),
					),
				)
			})

			It("should create Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}
				params := CreateDedicatedHostPoolRequest{
					FlavorClass: "flavorclass1",
					Metro:       "metro1",
					Name:        "name1",
				}
				dh, err := newDedicatedHostPool(server.URL()).CreateDedicatedHostPool(params, target)
				Expect(err).NotTo(HaveOccurred())
				Expect(dh.ID).Should(BeEquivalentTo("dedicatedhostpoolid1"))
			})
		})
		Context("When creating dedicatedhostpool is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/createDedicatedHostPool"),
						ghttp.VerifyJSON(`{
							"flavorClass": "flavorclass1",
							"metro": "metro1",
							"name": "name1"
						  }`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to create dedicatedhostpool`),
					),
				)
			})

			It("should return error during creating dedicatedhostpool", func() {
				params := CreateDedicatedHostPoolRequest{
					FlavorClass: "flavorclass1",
					Metro:       "metro1",
					Name:        "name1",
				}
				target := ClusterTargetHeader{}
				_, err := newDedicatedHostPool(server.URL()).CreateDedicatedHostPool(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Get", func() {
		Context("When Get dedicatedhostpool is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/getDedicatedHostPool"),
						ghttp.RespondWith(http.StatusCreated, `{
							"flavorClass": "flavorclass1",
							"hostCount": 3,
							"id": "dedicatedhostpool1",
							"metro": "metro1",
							"name": "name1",
							"state": "state1",
							"workerPools": [
							  {
								"clusterID": "cluster1",
								"workerPoolID": "workerpool1"
							  }
							],
							"zones": [
							  {
								"capacity": {
								  "memoryBytes": 12,
								  "vcpu": 1
								},
								"hostCount": 3,
								"zone": "zone1"
							  }
							]
						  }`),
					),
				)
			})

			It("should get Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}
				dh, err := newDedicatedHostPool(server.URL()).GetDedicatedHostPool("dedicatedhostpoolid1", target)
				Expect(err).NotTo(HaveOccurred())
				expectedDedicatedHostPool := GetDedicatedHostPoolResponse{
					FlavorClass: "flavorclass1",
					HostCount:   3,
					ID:          "dedicatedhostpool1",
					Metro:       "metro1",
					Name:        "name1",
					State:       "state1",
					WorkerPools: []DedicatedHostPoolWorkerPool{
						DedicatedHostPoolWorkerPool{
							ClusterID:    "cluster1",
							WorkerPoolID: "workerpool1",
						},
					},
					Zones: []DedicatedHostZoneResources{
						DedicatedHostZoneResources{
							Capacity: DedicatedHostResource{
								MemoryBytes: 12,
								VCPU:        1,
							},
							HostCount: 3,
							Zone:      "zone1",
						},
					},
				}
				Expect(dh).Should(BeEquivalentTo(expectedDedicatedHostPool))
			})
		})
		Context("When get dedicatedhostpool is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/getDedicatedHostPool"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get dedicatedhostpool`),
					),
				)
			})

			It("should return error during get dedicatedhostpool", func() {
				target := ClusterTargetHeader{}
				_, err := newDedicatedHostPool(server.URL()).GetDedicatedHostPool("dedicatedhostpoolid1", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("List", func() {
		Context("When list dedicatedhostpool is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/getDedicatedHostPools"),
						ghttp.RespondWith(http.StatusCreated, `[
							{
								"flavorClass": "flavorclass1",
								"hostCount": 3,
								"id": "dedicatedhostpool1",
								"metro": "metro1",
								"name": "name1",
								"state": "state1",
								"workerPools": [
								  {
									"clusterID": "cluster1",
									"workerPoolID": "workerpool1"
								  }
								],
								"zones": [
								  {
									"capacity": {
									  "memoryBytes": 12,
									  "vcpu": 1
									},
									"hostCount": 3,
									"zone": "zone1"
								  }
								]
							  }
						  ]`),
					),
				)
			})

			It("should list dedicatedhostpools in a cluster", func() {
				target := ClusterTargetHeader{}

				ldh, err := newDedicatedHostPool(server.URL()).ListDedicatedHostPools(target)
				Expect(err).NotTo(HaveOccurred())
				expectedDedicatedHostPools := []GetDedicatedHostPoolResponse{GetDedicatedHostPoolResponse{
					FlavorClass: "flavorclass1",
					HostCount:   3,
					ID:          "dedicatedhostpool1",
					Metro:       "metro1",
					Name:        "name1",
					State:       "state1",
					WorkerPools: []DedicatedHostPoolWorkerPool{
						DedicatedHostPoolWorkerPool{
							ClusterID:    "cluster1",
							WorkerPoolID: "workerpool1",
						},
					},
					Zones: []DedicatedHostZoneResources{
						DedicatedHostZoneResources{
							Capacity: DedicatedHostResource{
								MemoryBytes: 12,
								VCPU:        1,
							},
							HostCount: 3,
							Zone:      "zone1",
						},
					},
				}}
				Expect(ldh).To(BeEquivalentTo(expectedDedicatedHostPools))
			})
		})
		Context("When list dedicatedhostpool is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/getDedicatedHostPools"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to list dedicatedhostpool`),
					),
				)
			})

			It("should return error during get dedicatedhostpools", func() {
				target := ClusterTargetHeader{}
				_, err := newDedicatedHostPool(server.URL()).ListDedicatedHostPools(target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Remove", func() {
		Context("When removing dedicatedhostpool is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/removeDedicatedHostPool"),
						ghttp.VerifyJSON(`{
							"hostPool": "hostpoolid1"
						  }`),
						ghttp.RespondWith(http.StatusCreated, nil),
					),
				)
			})

			It("should remove dedicatedhostpool", func() {
				target := ClusterTargetHeader{}
				params := RemoveDedicatedHostPoolRequest{
					HostPoolID: "hostpoolid1",
				}
				err := newDedicatedHostPool(server.URL()).RemoveDedicatedHostPool(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When removing dedicatedhostpool is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/removeDedicatedHostPool"),
						ghttp.VerifyJSON(`{
							"hostPool": "hostpoolid1"
						  }`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to remove dedicatedhostpool`),
					),
				)
			})

			It("should return error during creating dedicatedhostpool", func() {
				target := ClusterTargetHeader{}
				params := RemoveDedicatedHostPoolRequest{
					HostPoolID: "hostpoolid1",
				}
				err := newDedicatedHostPool(server.URL()).RemoveDedicatedHostPool(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})

func newDedicatedHostPool(url string) DedicatedHostPool {

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
	return newDedicatedHostPoolAPI(&client)
}
