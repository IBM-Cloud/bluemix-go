package containerv2

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

var _ = Describe("Albs", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//Create
	Describe("Create", func() {
		Context("When creating alb is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/vpc/createAlb"),
						ghttp.VerifyJSON(`{"cluster":"345","type":"public","enableByDefault":true,"zone": "us-south-1", "ingressImage": "1.5.1_5367_iks"}`),
						ghttp.RespondWith(http.StatusCreated, `{"alb":"1234", "cluster":"clusterID"}`),
					),
				)
			})

			It("should create Alb in a cluster", func() {
				target := ClusterTargetHeader{}
				params := AlbCreateReq{
					Cluster: "345", Type: "public", EnableByDefault: true, ZoneAlb: "us-south-1", IngressImage: "1.5.1_5367_iks",
				}
				AlbResp, err := newAlbs(server.URL()).CreateAlb(params, target)
				Expect(AlbResp.Alb).To(Equal("1234"))
				Expect(AlbResp.Cluster).To(Equal("clusterID"))
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When creating is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/vpc/createAlb"),
						ghttp.VerifyJSON(`{"cluster":"345","type":"public","enableByDefault":true,"zone": "us-south-1","ingressImage": "1.5.1_5367_iks"}
`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to create alb`),
					),
				)
			})

			It("should return error during creating alb", func() {
				params := AlbCreateReq{
					Cluster: "345", Type: "public", EnableByDefault: true, ZoneAlb: "us-south-1", IngressImage: "1.5.1_5367_iks",
				}
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).CreateAlb(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//Enable
	Describe("Enable", func() {
		Context("When Enabling alb is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/vpc/enableAlb"),
						ghttp.VerifyJSON(`{"albBuild":"579","albID":"private-crbm64u3ed02o93vv36hb0-alb1","albType":"private","authBuild":"341","cluster":"345","createdDate":"","disableDeployment":true,"loadBalancerHostname":"","name":"","numOfInstances":"","resize":true,"state":"disabled","status":"","enable":true,"zone": "us-south-1"}`),
						ghttp.RespondWith(http.StatusCreated, `{}`),
					),
				)
			})

			It("should enable Alb in a cluster", func() {
				target := ClusterTargetHeader{}
				params := AlbConfig{
					AlbBuild: "579", AlbID: "private-crbm64u3ed02o93vv36hb0-alb1", AuthBuild: "341", Cluster: "345", CreatedDate: "", DisableDeployment: true, LoadBalancerHostname: "", AlbType: "private", Name: "", NumOfInstances: "", Resize: true, State: "disabled", Status: "", Enable: true, ZoneAlb: "us-south-1",
				}
				err := newAlbs(server.URL()).EnableAlb(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When enabling is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/vpc/enableAlb"),
						ghttp.VerifyJSON(`{"albBuild":"579","albID":"private-crbm64u3ed02o93vv36hb0-alb1","albType":"private","authBuild":"341","cluster":"345","createdDate":"","disableDeployment":true,"loadBalancerHostname":"","name":"","numOfInstances":"","resize":true,"state":"disabled","status":"","enable":true,"zone": "us-south-1"}
`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to enable alb`),
					),
				)
			})

			It("should return error during enabling alb", func() {
				params := AlbConfig{
					AlbBuild: "579", AlbID: "private-crbm64u3ed02o93vv36hb0-alb1", AlbType: "private", AuthBuild: "341", CreatedDate: "", DisableDeployment: true, Enable: true, LoadBalancerHostname: "", Name: "", NumOfInstances: "", Resize: true, State: "disabled", Status: "", Cluster: "345", ZoneAlb: "us-south-1",
				}
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).EnableAlb(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//Disable
	Describe("Disable", func() {
		Context("When Disabling alb is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/vpc/disableAlb"),
						ghttp.VerifyJSON(`{"albBuild":"579","albID":"private-crbm64u3ed02o93vv36hb0-alb1","albType":"private","authBuild":"341","cluster":"345","createdDate":"","disableDeployment":true,"loadBalancerHostname":"","name":"","numOfInstances":"","resize":true,"state":"disabled","status":"","enable":true,"zone": "us-south-1"}`),
						ghttp.RespondWith(http.StatusCreated, `{}`),
					),
				)
			})

			It("should disable Alb in a cluster", func() {
				target := ClusterTargetHeader{}
				params := AlbConfig{
					AlbBuild: "579", AlbID: "private-crbm64u3ed02o93vv36hb0-alb1", AlbType: "private", AuthBuild: "341", CreatedDate: "", DisableDeployment: true, Enable: true, LoadBalancerHostname: "", Name: "", NumOfInstances: "", Resize: true, State: "disabled", Status: "", Cluster: "345", ZoneAlb: "us-south-1",
				}
				err := newAlbs(server.URL()).DisableAlb(params, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When disabling is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/vpc/disableAlb"),
						ghttp.VerifyJSON(`{"albBuild":"579","albID":"private-crbm64u3ed02o93vv36hb0-alb1","albType":"private","authBuild":"341","cluster":"345","createdDate":"","disableDeployment":true,"loadBalancerHostname":"","name":"","numOfInstances":"","resize":true,"state":"disabled","status":"","enable":true,"zone": "us-south-1"}
`),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to disable alb`),
					),
				)
			})

			It("should return error during disabling alb", func() {
				params := AlbConfig{
					AlbBuild: "579", AlbID: "private-crbm64u3ed02o93vv36hb0-alb1", AlbType: "private", AuthBuild: "341", CreatedDate: "", DisableDeployment: true, Enable: true, LoadBalancerHostname: "", Name: "", NumOfInstances: "", Resize: true, State: "disabled", Status: "", Cluster: "345", ZoneAlb: "us-south-1",
				}
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).DisableAlb(params, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//GetAlbs
	Describe("Get", func() {
		Context("When Get Alb is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/getAlb"),
						ghttp.RespondWith(http.StatusCreated, `{"albBuild": "string","albID": "string","albType": "string","authBuild": "string","cluster": "string","createdDate": "string","disableDeployment": true,"enable": true,"loadBalancerHostname": "string","name": "string","numOfInstances": "string","resize": true,"state": "string","status": "string","zone": "string"}`),
					),
				)
			})

			It("should get Alb in a cluster", func() {
				target := ClusterTargetHeader{}

				_, err := newAlbs(server.URL()).GetAlb("aaa", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When get alb is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/getAlb"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get alb`),
					),
				)
			})

			It("should return error during get alb", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).GetAlb("aaa", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// UpdateAlb
	Describe("UpdateAlb", func() {
		Context("When Update Alb is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/updateAlb"),
						ghttp.RespondWith(http.StatusOK, `{"clusterID":"myCluster"}`),
					),
				)
			})

			It("should update Alb in a cluster", func() {
				target := ClusterTargetHeader{}
				updateReq := UpdateALBReq{
					ClusterID: "myCluster",
					ALBBuild:  "1.8.1_5384_iks",
					ALBList: []string{
						"public-crck9aaedd0p8vjmqa0asg-alb1",
					},
				}
				err := newAlbs(server.URL()).UpdateAlb(updateReq, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When update alb is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/updateAlb"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to update alb`),
					),
				)
			})

			It("should return error during get alb", func() {
				target := ClusterTargetHeader{}
				updateReq := UpdateALBReq{
					ClusterID: "myCluster",
					ALBBuild:  "1.8.1_5384_iks",
					ALBList: []string{
						"public-crck9aaedd0p8vjmqa0asg-alb1",
					},
				}
				err := newAlbs(server.URL()).UpdateAlb(updateReq, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// ListClusterAlbs
	Describe("ListClusterAlbs", func() {
		Context("When List Albs is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/getClusterAlbs"),
						ghttp.RespondWith(http.StatusOK, `{"alb":[{"albID":"private-crclusterid-alb1","cluster":"clusterid","name":"","albType":"private","enable":true,"state":"enabled","createdDate":"","numOfInstances":"2","resize":false,"zone":"testzone","disableDeployment":false,"albBuild":"1.8.1_5384_iks","authBuild":"","loadBalancerHostname":"test1-us-south.lb.test.appdomain.cloud","status":"healthy"},{"albID":"public-crclusterid-alb1","cluster":"clusterid","name":"","albType":"public","enable":true,"state":"enabled","createdDate":"","numOfInstances":"2","resize":false,"zone":"testzone","disableDeployment":false,"albBuild":"1.8.1_5384_iks","authBuild":"","loadBalancerHostname":"test2-us-south.lb.test.appdomain.cloud","status":"healthy"}]}`),
					),
				)
			})

			It("should list Albs in a cluster", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).ListClusterAlbs("clusterid", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When List Albs is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/getClusterAlbs"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to list albs`),
					),
				)
			})

			It("should return error during get alb", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).ListClusterAlbs("clusterid", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// ListAlbImages
	Describe("ListAlbImages", func() {
		Context("When List Alb images is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/getAlbImages"),
						ghttp.RespondWith(http.StatusOK, `{"defaultK8sVersion":"1.8.1_5384_iks","supportedK8sVersions":["1.5.1_5407_iks","1.6.4_5406_iks","1.8.1_5384_iks"]}`),
					),
				)
			})

			It("should list Alb images", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).ListAlbImages(target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When List Alb images is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/getAlbImages"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to list alb images`),
					),
				)
			})

			It("should return error during get alb images", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).ListAlbImages(target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// GetIngressStatus
	Describe("GetIngressStatus", func() {
		Context("When GetIngressStatus is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/getStatus"),
						ghttp.RespondWith(http.StatusOK, `{"cluster":"clusterid","enabled":true,"status":"healthy","nonTranslatedStatus":"HEALTHY","message":"All Ingress components are healthy.","StatusList":[{"component":"ingress-controller-configmap","status":"healthy","type":"cluster"},{"component":"alb-healthcheck-ingress","status":"healthy","type":"cluster"},{"component":"public-crclusterid-alb1","status":"healthy","type":"alb"},{"component":"ingressHostname.us-south.containers.appdomain.cloud","status":"healthy","type":"subdomain"},{"component":"default/ingressHostname","status":"healthy","type":"secret"},{"component":"ibm-cert-store/ingressHostname","status":"healthy","type":"secret"},{"component":"kube-system/ingressHostname","status":"healthy","type":"secret"}],"generalComponentStatus":[{"component":"ingress-controller-configmap","status":["healthy"]},{"component":"alb-healthcheck-ingress","status":["healthy"]}],"albStatus":[{"component":"public-crclusterid-alb1","status":["healthy"]}],"subdomainStatus":[{"component":"ingressHostname.us-south.containers.appdomain.cloud","status":["healthy"]}],"secretStatus":[{"component":"default/ingressHostname","status":["healthy"]},{"component":"ibm-cert-store/ingressHostname","status":["healthy"]},{"component":"kube-system/ingressHostname","status":["healthy"]}],"ignoredErrors":[]}`),
					),
				)
			})

			It("should get ingress status in a cluster", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).GetIngressStatus("clusterid", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When GetIngressStatus is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/getStatus"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get ingress status`),
					),
				)
			})

			It("should return error during get ingress status", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).GetIngressStatus("clusterid", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// SetAlbClusterHealthCheckConfig
	Describe("SetAlbClusterHealthCheckConfig", func() {
		clusterHcConfig := ALBClusterHealthCheckConfig{
			Cluster: "clusterid",
			Enable:  false,
		}

		Context("When SetAlbClusterHealthCheckConfig is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/setIngressClusterHealthcheck"),
						ghttp.RespondWith(http.StatusOK, `{"cluster":"clusterid"}`),
					),
				)
			})

			It("should set ingress in cluster health checker config", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).SetAlbClusterHealthCheckConfig(clusterHcConfig, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When SetAlbClusterHealthCheckConfig is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/setIngressClusterHealthcheck"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to set ingress in cluster health checkher config`),
					),
				)
			})

			It("should return error during set ingress in cluster health checkher config", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).SetAlbClusterHealthCheckConfig(clusterHcConfig, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// GetIngressStatus
	Describe("GetIngressStatus", func() {
		Context("When GetIngressStatus is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/getIngressClusterHealthcheck"),
						ghttp.RespondWith(http.StatusOK, `{ "cluster": "ckdo7li2075cidkk3pu0", "enable": true }`),
					),
				)
			})

			It("should get ingress in cluster health checker config", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).GetAlbClusterHealthCheckConfig("clusterid", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When GetIngressStatus is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/getIngressClusterHealthcheck"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get ingress in cluster health checker`),
					),
				)
			})

			It("should return error during get ingress in cluster health checker", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).GetAlbClusterHealthCheckConfig("clusterid", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// GetIgnoredIngressStatusErrors
	Describe("GetIgnoredIngressStatusErrors", func() {
		Context("When GetIgnoredIngressStatusErrors is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/listIgnoredIngressStatusErrors"),
						ghttp.RespondWith(http.StatusOK, `{ "cluster": "mycluster", "ignoredErrors": ["ERRAHCF"] }`),
					),
				)
			})

			It("should get ignored error codes", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).GetIgnoredIngressStatusErrors("mycluster", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When GetIgnoredIngressStatusErrors is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/alb/listIgnoredIngressStatusErrors"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to list ignored error codes`),
					),
				)
			})

			It("should return error during get ignored error codes", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).GetIgnoredIngressStatusErrors("mycluster", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// AddIgnoredIngressStatusErrors
	Describe("AddIgnoredIngressStatusErrors", func() {
		ignoredErrors := IgnoredIngressStatusErrors{
			Cluster: "mycluster",
			IgnoredErrors: []string{
				"ERRAHCF",
			},
		}
		Context("When AddIgnoredIngressStatusErrors is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/addIgnoredIngressStatusErrors"),
						ghttp.RespondWith(http.StatusOK, `{ "cluster": "mycluster" }`),
					),
				)
			})

			It("should set ingnored error codes", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).AddIgnoredIngressStatusErrors(ignoredErrors, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When AddIgnoredIngressStatusErrors is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/addIgnoredIngressStatusErrors"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to set ignored error codes`),
					),
				)
			})

			It("should return error during set ignored error codes", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).AddIgnoredIngressStatusErrors(ignoredErrors, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// RemoveIgnoredIngressStatusErrors
	Describe("RemoveIgnoredIngressStatusErrors", func() {
		ignoredErrors := IgnoredIngressStatusErrors{
			Cluster: "mycluster",
			IgnoredErrors: []string{
				"ERRAHCF",
			},
		}
		Context("When RemoveIgnoredIngressStatusErrors is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v2/alb/removeIgnoredIngressStatusErrors"),
						ghttp.RespondWith(http.StatusOK, `{ "cluster": "mycluster" }`),
					),
				)
			})

			It("should remove ingnored error codes", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).RemoveIgnoredIngressStatusErrors(ignoredErrors, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When RemoveIgnoredIngressStatusErrors is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/v2/alb/removeIgnoredIngressStatusErrors"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to remove ignored error codes`),
					),
				)
			})

			It("should return error during remove ignored error codes", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).RemoveIgnoredIngressStatusErrors(ignoredErrors, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// SetIngressStatusState
	Describe("SetIngressStatusState", func() {
		ingressStatus := IngressStatusState{
			Cluster: "mycluster",
			Enable:  false,
		}
		Context("When SetIngressStatusState is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/setIngressStatusState"),
						ghttp.RespondWith(http.StatusOK, `{ "cluster": "mycluster" }`),
					),
				)
			})

			It("should set ingress status", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).SetIngressStatusState(ingressStatus, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When SetIngressStatusState is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/v2/alb/setIngressStatusState"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to set ingress status`),
					),
				)
			})

			It("should return error during set ingress status", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).SetIngressStatusState(ingressStatus, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// GetIngressLoadBalancerConfig
	Describe("GetIngressLoadBalancerConfig", func() {
		Context("When GetIngressLoadBalancerConfig is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/ingress/v2/load-balancer/configuration"),
						ghttp.RespondWith(http.StatusOK, `{ "cluster": "clusterid", "proxyProtocol": { "enable": true, "headerTimeout": 0 }, "type": "public" }`),
					),
				)
			})

			It("should get lb config", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).GetIngressLoadBalancerConfig("clusterid", "public", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When GetIngressLoadBalancerConfig is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/ingress/v2/load-balancer/configuration"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get lb config`),
					),
				)
			})

			It("should return error during get lb config", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).GetIngressLoadBalancerConfig("clusterid", "public", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// UpdateIngressLoadBalancerConfig
	Describe("UpdateIngressLoadBalancerConfig", func() {
		lbConf := ALBLBConfig{
			Cluster: "clusterid",
			Type:    "public",
			ProxyProtocol: &ALBLBProxyProtocolConfig{
				Enable: true,
			},
		}
		Context("When UpdateIngressLoadBalancerConfig is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPatch, "/ingress/v2/load-balancer/configuration"),
						ghttp.RespondWith(http.StatusOK, `{ "cluster": "clusterid", "proxyProtocol": { "enable": true, "headerTimeout": 0 }, "type": "public" }`),
					),
				)
			})

			It("should update lb config", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).UpdateIngressLoadBalancerConfig(lbConf, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When UpdateIngressLoadBalancerConfig is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPatch, "/ingress/v2/load-balancer/configuration"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to update lb config`),
					),
				)
			})

			It("should return error during update lb config", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).UpdateIngressLoadBalancerConfig(lbConf, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// GetALBAutoscaleConfiguration
	Describe("GetALBAutoscaleConfiguration", func() {
		Context("When GetALBAutoscaleConfiguration is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/ingress/v2/clusters/clusterid/albs/public-crclusterid-alb1/autoscale"),
						ghttp.RespondWith(http.StatusOK, `{"config":{"minReplicas":2,"maxReplicas":4,"cpuAverageUtilization":600}}`),
					),
				)
			})

			It("should get autoscale config", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).GetALBAutoscaleConfiguration("clusterid", "public-crclusterid-alb1", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When GetALBAutoscaleConfiguration is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/ingress/v2/clusters/clusterid/albs/public-crclusterid-alb1/autoscale"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get autoscale config`),
					),
				)
			})

			It("should return error during get autoscale config", func() {
				target := ClusterTargetHeader{}
				_, err := newAlbs(server.URL()).GetALBAutoscaleConfiguration("clusterid", "public-crclusterid-alb1", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// SetALBAutoscaleConfiguration
	Describe("SetALBAutoscaleConfiguration", func() {
		autoscaleConf := AutoscaleDetails{
			Config: &AutoscaleConfig{
				MinReplicas:           2,
				MaxReplicas:           2,
				CPUAverageUtilization: 600,
			},
		}
		Context("When SetALBAutoscaleConfiguration is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, "/ingress/v2/clusters/clusterid/albs/public-crclusterid-alb1/autoscale"),
						ghttp.RespondWith(http.StatusOK, ``),
					),
				)
			})

			It("should set autoscale config", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).SetALBAutoscaleConfiguration("clusterid", "public-crclusterid-alb1", autoscaleConf, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When SetALBAutoscaleConfiguration is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, "/ingress/v2/clusters/clusterid/albs/public-crclusterid-alb1/autoscale"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to set autoscale config`),
					),
				)
			})

			It("should return error during set autoscale config", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).SetALBAutoscaleConfiguration("clusterid", "public-crclusterid-alb1", autoscaleConf, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// RemoveALBAutoscaleConfiguration
	Describe("RemoveALBAutoscaleConfiguration", func() {
		Context("When RemoveALBAutoscaleConfiguration is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/ingress/v2/clusters/clusterid/albs/public-crclusterid-alb1/autoscale"),
						ghttp.RespondWith(http.StatusNoContent, `{"config":{"minReplicas":2,"maxReplicas":4,"cpuAverageUtilization":600}}`),
					),
				)
			})

			It("should delete autoscale config", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).RemoveALBAutoscaleConfiguration("clusterid", "public-crclusterid-alb1", target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When RemoveALBAutoscaleConfiguration is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/ingress/v2/clusters/clusterid/albs/public-crclusterid-alb1/autoscale"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to delete autoscale config`),
					),
				)
			})

			It("should return error during delete autoscale config", func() {
				target := ClusterTargetHeader{}
				err := newAlbs(server.URL()).RemoveALBAutoscaleConfiguration("clusterid", "public-crclusterid-alb1", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})

func newAlbs(url string) Alb {

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
	return newAlbAPI(&client)
}
