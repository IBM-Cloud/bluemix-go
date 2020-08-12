package icdv4

import (
	"log"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/client"
	bluemixHttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/session"
)

var _ = Describe("Auto Scaling", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})
	Describe("Set", func() {
		Context("When Set Auto Scaling group is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPatch, "/v4/ibm/deployments/crn:v1:bluemix:public:databases-for-redis:us-south:a/4448261269a14562b839e0a3019ed980:58a93bd0-de14-410d-a775-8495a295e47f::/groups/member/autoscaling"),
						ghttp.RespondWith(http.StatusOK, `
                           {
							  "task": 
							  	{
									"id":"crn:v1:bluemix:public:databases-for-redis:us-south:a/4448261269a14562b839e0a3019ed980:58a93bd0-de14-410d-a775-8495a295e47f:task:e3f0d867-f6a1-4242-b753-16ba59659766",
									"description":"Updating autoscale settings.",
									"status":"running",
									"deployment_id":"crn:v1:bluemix:public:databases-for-redis:us-south:a/4448261269a14562b839e0a3019ed980:58a93bd0-de14-410d-a775-8495a295e47f::",
									"progress_percent":0,
									"created_at":"2020-08-12T11:15:45.000Z"}
								}
                            }
                        `),
					),
				)
			})
			It("should return group updated", func() {
				target1 := "crn:v1:bluemix:public:databases-for-redis:us-south:a/4448261269a14562b839e0a3019ed980:58a93bd0-de14-410d-a775-8495a295e47f::"
				target2 := "member"
				memoryReq := ASGBody{
					Scalers: ScalersBody{
						IO: &IOBody{
							Enabled:      true,
							AbovePercent: 35,
							OverPeriod:   "15m",
						},
					},
					Rate: RateBody{
						IncreasePercent:  10,
						LimitMBPerMember: 5000,
						PeriodSeconds:    900,
						Units:            "mb",
					},
				}

				groupBdy := AutoscalingGroup{
					Memory: &memoryReq,
				}
				params := AutoscalingSetGroup{
					Autoscaling: groupBdy,
				}
				myTask, err := newScalingGroup(server.URL()).SetAutoScaling(target1, target2, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(myTask).ShouldNot(BeNil())
			})
		})
		Context("When update is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPatch, "/v4/ibm/deployments/crn:v1:bluemix:public:databases-for-redis:us-south:a/4448261269a14562b839e0a3019ed980:58a93bd0-de14-410d-a775-8495a295e47f::/groups/member/autoscaling"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to update group`),
					),
				)
			})

			It("should return error during Auto Scaling group Set", func() {
				target1 := "crn:v1:bluemix:public:databases-for-redis:us-south:a/4448261269a14562b839e0a3019ed980:58a93bd0-de14-410d-a775-8495a295e47f::"
				target2 := "member"
				memoryReq := ASGBody{
					Scalers: ScalersBody{
						IO: &IOBody{
							Enabled:      true,
							AbovePercent: 35,
							OverPeriod:   "15m",
						},
					},
					Rate: RateBody{
						IncreasePercent:  10,
						LimitMBPerMember: 5000,
						PeriodSeconds:    900,
						Units:            "mb",
					},
				}
				groupBdy := AutoscalingGroup{
					Memory: &memoryReq,
					// Disk:   &diskReq,
					// CPU:    &cpuReq,
				}
				params := AutoscalingSetGroup{
					Autoscaling: groupBdy,
				}
				myTask, err := newScalingGroup(server.URL()).SetAutoScaling(target1, target2, params)
				Expect(err).To(HaveOccurred())
				Expect(myTask.Id).Should(Equal(""))
			})
		})
	})
	Describe("Get", func() {
		Context("When get is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v4/ibm/deployments/crn:v1:bluemix:public:databases-for-redis:us-south:a/4448261269a14562b839e0a3019ed980:58a93bd0-de14-410d-a775-8495a295e47f::/groups/member/autoscaling"),
						ghttp.RespondWith(http.StatusOK, `
						{
							"autoscaling":
							{
								"disk":
								{
									"scalers":
									{
										"capacity":
										{
											"enabled":false,
											"free_space_less_than_percent":10
										},
										"io_utilization":
										{
											"enabled":false,
											"over_period":"15m",
											"above_percent":90
										}
									},
									"rate":
									{
										"increase_percent":10.0,
										"period_seconds":900,
										"limit_mb_per_member":3670016,
										"units":"mb"
									}
								},
								"memory":
								{
									"scalers":
									{
										"io_utilization":
										{
											"enabled":true,
											"over_period":"15m",
											"above_percent":35
										}
									},
									"rate":
									{
										"increase_percent":10.0,
										"period_seconds":900,
										"limit_mb_per_member":5000,
										"units":"mb"
									}
								},
								"cpu":
								{
									"scalers":{},
									"rate":
									{
										"increase_percent":10.0,
										"period_seconds":900,
										"units":"count",
										"limit_count_per_member":30}
									}
								}
							}
                        `),
					),
				)
			})

			It("should return groups", func() {
				target1 := "crn:v1:bluemix:public:databases-for-redis:us-south:a/4448261269a14562b839e0a3019ed980:58a93bd0-de14-410d-a775-8495a295e47f::"
				_, err := newScalingGroup(server.URL()).GetAutoScaling(target1, "member")
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When get is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v4/ibm/deployments/crn:v1:bluemix:public:databases-for-redis:us-south:a/4448261269a14562b839e0a3019ed980:58a93bd0-de14-410d-a775-8495a295e47f::/groups/member/autoscaling"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get groups`),
					),
				)
			})

			It("should return error during AUto Scaling Group get", func() {
				target1 := "crn:v1:bluemix:public:databases-for-redis:us-south:a/4448261269a14562b839e0a3019ed980:58a93bd0-de14-410d-a775-8495a295e47f::"
				_, err := newScalingGroup(server.URL()).GetAutoScaling(target1, "member")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

func newScalingGroup(url string) AutoScaling {

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.ICDService,
	}
	return newAutoScalingAPI(&client)
}
