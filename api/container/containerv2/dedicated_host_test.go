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

var _ = Describe("dedicatedhosts", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	Describe("Create", func() {
		Context("When creating dedicatedhost is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/createDedicatedHost"),
						ghttp.VerifyJSON(`{
							"flavor": "flavor1",
							"hostPoolID": "hostpoolid1",
							"zone": "zone1"
						  }`),
						ghttp.RespondWith(http.StatusCreated, `{
							"id":"dedicatedhostid1"
						}`),
					),
				)
			})

			It("should create Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}
				params := CreateDedicatedHostRequest{
					Flavor:     "flavor1",
					HostPoolID: "hostpoolid1",
					Zone:       "zone1",
				}
				dh, err := newDedicatedHost(server.URL()).CreateDedicatedHost(params, target)
				Expect(err).NotTo(HaveOccurred())
				Expect(dh.ID).Should(BeEquivalentTo("dedicatedhostid1"))
			})
		})
		Context("When creating dedicatedhost is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/createDedicatedHost"),
						ghttp.VerifyJSON(`{
							"flavor": "flavor1",
							"hostPoolID": "hostpoolid1",
							"zone": "zone1"
						  }`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to create dedicatedhost`),
					),
				)
			})

			It("should return error during creating dedicatedhost", func() {
				params := CreateDedicatedHostRequest{
					Flavor:     "flavor1",
					HostPoolID: "hostpoolid1",
					Zone:       "zone1",
				}
				target := ClusterTargetHeader{}
				_, err := newDedicatedHost(server.URL()).CreateDedicatedHost(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Get", func() {
		Context("When Get dedicatedhost is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/getDedicatedHost"),
						ghttp.RespondWith(http.StatusCreated, `{
							"flavor": "flavor1",
							"id": "dedicatedhostid1",
							"lifecycle": {
							  "actualState": "actualstate1",
							  "desiredState": "desiredstate1"
							},
							"placementEnabled": true,
							"resources": {
							  "capacity": {
								"memoryBytes": 12000,
								"vcpu": 2
							  },
							  "consumed": {
								"memoryBytes": 1,
								"vcpu": 1
							  }
							},
							"workers": [
							  {
								"clusterID": "clusterid1",
								"flavor": "flavor2",
								"workerID": "workerid1",
								"workerPoolID": "workerpoolid1"
							  }
							],
							"zone": "zone1"
						  }`),
					),
				)
			})

			It("should get Workerpool in a cluster", func() {
				target := ClusterTargetHeader{}
				dh, err := newDedicatedHost(server.URL()).GetDedicatedHost("dedicatedhostid1", "dedicatedhostpoolid1", target)
				Expect(err).NotTo(HaveOccurred())
				expectedDedicatedHost := GetDedicatedHostResponse{
					Flavor: "flavor1",
					ID:     "dedicatedhostid1",
					Lifecycle: DedicatedHostLifecycle{
						ActualState:  "actualstate1",
						DesiredState: "desiredstate1",
					},
					PlacementEnabled: true,
					Resources: DedicatedHostResources{
						Capacity: DedicatedHostResource{
							MemoryBytes: 12000,
							VCPU:        2,
						},
						Consumed: DedicatedHostResource{
							MemoryBytes: 1,
							VCPU:        1,
						},
					},
					Workers: []DedicatedHostWorker{
						DedicatedHostWorker{
							ClusterID:    "clusterid1",
							Flavor:       "flavor2",
							WorkerID:     "workerid1",
							WorkerPoolID: "workerpoolid1",
						},
					},
					Zone: "zone1",
				}
				Expect(dh).Should(BeEquivalentTo(expectedDedicatedHost))
			})
		})
		Context("When get dedicatedhost is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/getDedicatedHost"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get dedicatedhost`),
					),
				)
			})

			It("should return error during get dedicatedhost", func() {
				target := ClusterTargetHeader{}
				_, err := newDedicatedHost(server.URL()).GetDedicatedHost("dedicatedhostid1", "dedicatedhostpoolid1", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("List", func() {
		Context("When list dedicatedhost is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/getDedicatedHosts"),
						ghttp.RespondWith(http.StatusCreated, `[
							{
								"flavor": "flavor1",
								"id": "dedicatedhostid1",
								"lifecycle": {
								  "actualState": "actualstate1",
								  "desiredState": "desiredstate1"
								},
								"placementEnabled": true,
								"resources": {
								  "capacity": {
									"memoryBytes": 12000,
									"vcpu": 2
								  },
								  "consumed": {
									"memoryBytes": 1,
									"vcpu": 1
								  }
								},
								"workers": [
								  {
									"clusterID": "clusterid1",
									"flavor": "flavor2",
									"workerID": "workerid1",
									"workerPoolID": "workerpoolid1"
								  }
								],
								"zone": "zone1"
							  }
						  ]`),
					),
				)
			})

			It("should list dedicatedhosts in a cluster", func() {
				target := ClusterTargetHeader{}

				ldh, err := newDedicatedHost(server.URL()).ListDedicatedHosts("dedicatedhostpoolid1", target)
				Expect(err).NotTo(HaveOccurred())
				expectedDedicatedHosts := []GetDedicatedHostResponse{GetDedicatedHostResponse{
					Flavor: "flavor1",
					ID:     "dedicatedhostid1",
					Lifecycle: DedicatedHostLifecycle{
						ActualState:  "actualstate1",
						DesiredState: "desiredstate1",
					},
					PlacementEnabled: true,
					Resources: DedicatedHostResources{
						Capacity: DedicatedHostResource{
							MemoryBytes: 12000,
							VCPU:        2,
						},
						Consumed: DedicatedHostResource{
							MemoryBytes: 1,
							VCPU:        1,
						},
					},
					Workers: []DedicatedHostWorker{
						DedicatedHostWorker{
							ClusterID:    "clusterid1",
							Flavor:       "flavor2",
							WorkerID:     "workerid1",
							WorkerPoolID: "workerpoolid1",
						},
					},
					Zone: "zone1",
				}}
				Expect(ldh).To(BeEquivalentTo(expectedDedicatedHosts))
			})
		})
		Context("When list dedicatedhost is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/getDedicatedHosts"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to list dedicatedhost`),
					),
				)
			})

			It("should return error during get dedicatedhosts", func() {
				target := ClusterTargetHeader{}
				_, err := newDedicatedHost(server.URL()).ListDedicatedHosts("dedicatedhostpoolid1", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Remove", func() {
		Context("When removing dedicatedhost is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/removeDedicatedHost"),
						ghttp.VerifyJSON(`{
							"host": "host1",
							"hostPool": "hostpoolid1"
						  }`),
						ghttp.RespondWith(http.StatusCreated, nil),
					),
				)
			})

			It("should remove dedicatedhost", func() {
				target := ClusterTargetHeader{}
				params := RemoveDedicatedHostRequest{
					HostID:     "host1",
					HostPoolID: "hostpoolid1",
				}
				err := newDedicatedHost(server.URL()).RemoveDedicatedHost(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When removing dedicatedhost is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/removeDedicatedHost"),
						ghttp.VerifyJSON(`{
							"host": "host1",
							"hostPool": "hostpoolid1"
						  }`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to remove dedicatedhost`),
					),
				)
			})

			It("should return error during creating dedicatedhost", func() {
				target := ClusterTargetHeader{}
				params := RemoveDedicatedHostRequest{
					HostID:     "host1",
					HostPoolID: "hostpoolid1",
				}
				err := newDedicatedHost(server.URL()).RemoveDedicatedHost(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Update", func() {
		Context("When enabling placement on a dedicatedhost is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/enableDedicatedHostPlacement"),
						ghttp.VerifyJSON(`{
							"hostID": "host1",
							"hostPoolID": "hostpoolid1"
						  }`),
						ghttp.RespondWith(http.StatusCreated, nil),
					),
				)
			})

			It("should enable dedicatedhost placement", func() {
				target := ClusterTargetHeader{}
				params := UpdateDedicatedHostPlacementRequest{
					HostID:     "host1",
					HostPoolID: "hostpoolid1",
				}
				err := newDedicatedHost(server.URL()).EnableDedicatedHostPlacement(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When disabling placement on a dedicatedhost is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/disableDedicatedHostPlacement"),
						ghttp.VerifyJSON(`{
							"hostID": "host1",
							"hostPoolID": "hostpoolid1"
						  }`),
						ghttp.RespondWith(http.StatusCreated, nil),
					),
				)
			})

			It("should disable dedicatedhost placement", func() {
				target := ClusterTargetHeader{}
				params := UpdateDedicatedHostPlacementRequest{
					HostID:     "host1",
					HostPoolID: "hostpoolid1",
				}
				err := newDedicatedHost(server.URL()).DisableDedicatedHostPlacement(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When enabling placement on a dedicatedhost is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/enableDedicatedHostPlacement"),
						ghttp.VerifyJSON(`{
							"hostID": "host1",
							"hostPoolID": "hostpoolid1"
						  }`),
						ghttp.RespondWith(http.StatusInternalServerError, `Bang`),
					),
				)
			})

			It("should enable dedicatedhost placement", func() {
				target := ClusterTargetHeader{}
				params := UpdateDedicatedHostPlacementRequest{
					HostID:     "host1",
					HostPoolID: "hostpoolid1",
				}
				err := newDedicatedHost(server.URL()).EnableDedicatedHostPlacement(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When disabling placement on a dedicatedhost is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/disableDedicatedHostPlacement"),
						ghttp.VerifyJSON(`{
							"hostID": "host1",
							"hostPoolID": "hostpoolid1"
						  }`),
						ghttp.RespondWith(http.StatusInternalServerError, `Bang`),
					),
				)
			})

			It("should disable dedicatedhost placement", func() {
				target := ClusterTargetHeader{}
				params := UpdateDedicatedHostPlacementRequest{
					HostID:     "host1",
					HostPoolID: "hostpoolid1",
				}
				err := newDedicatedHost(server.URL()).DisableDedicatedHostPlacement(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

func newDedicatedHost(url string) DedicatedHost {

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
	return newDedicatedHostAPI(&client)
}
