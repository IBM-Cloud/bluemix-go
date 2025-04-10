package containerv2

import (
	"log"
	"net/http"

	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/client"
	bluemixHttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/session"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Workers", func() {
	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})

	//ListByWorkerpool
	Describe("List", func() {
		Context("When List worker is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getWorkers"),
						ghttp.RespondWith(http.StatusCreated, `[
							{
							  "flavor": "string",
							  "health": {
								"message": "string",
								"state": "string"
							  },
							  "id": "string",
							  "kubeVersion": {
								"actual": "string",
								"desired": "string",
								"eos": "string",
								"masterEOS": "string",
								"target": "string"
							  },
							  "lifecycle": {
								"actualState": "string",
								"desiredState": "string",
								"message": "string",
								"messageDate": "string",
								"messageDetails": "string",
								"messageDetailsDate": "string",
								"pendingOperation": "string",
								"reasonForDelete": "string"
							  },
							  "location": "string",
							  "networkInterfaces": [
								{
								  "cidr": "string",
								  "ipAddress": "string",
								  "primary": true,
								  "subnetID": "string"
								}
							  ],
							  "poolID": "string",
							  "poolName": "string"
							}
						  ]`),
					),
				)
			})

			It("should list workers in a cluster", func() {
				target := ClusterTargetHeader{}

				_, err := newWorker(server.URL()).ListByWorkerPool("aaa", "bbb", true, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When list worker is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getWorkers"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to list worker`),
					),
				)
			})

			It("should return error during get worker", func() {
				target := ClusterTargetHeader{}
				_, err := newWorker(server.URL()).ListByWorkerPool("aaa", "bbb", true, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// ListAllWorkers
	Describe("ListAllWorkers", func() {
		Context("When ListAllWorkers is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						// https://containers.cloud.ibm.com/global/swagger-global-api/#/v2/getWorkers
						ghttp.VerifyRequest(http.MethodGet, "/v2/getWorkers"),
						ghttp.RespondWith(http.StatusOK, `[
  {
    "dedicatedHostId": "string",
    "dedicatedHostPoolId": "string",
    "flavor": "string",
    "health": {
      "message": "string",
      "state": "string"
    },
    "id": "string",
    "kubeVersion": {
      "actual": "string",
      "desired": "string",
      "eos": "string",
      "masterEOS": "string",
      "target": "string"
    },
    "lifecycle": {
      "actualOperatingSystem": "string",
      "actualState": "string",
      "desiredOperatingSystem": "string",
      "desiredState": "string",
      "message": "string",
      "messageDate": "string",
      "messageDetails": "string",
      "messageDetailsDate": "string",
      "pendingOperation": "string",
      "reasonForDelete": "string"
    },
    "location": "string",
    "networkInformation": {
      "privateIP": "string",
      "privateVLAN": "string",
      "publicIP": "string",
      "publicVLAN": "string"
    },
    "poolID": "string",
    "poolName": "string"
  }
]`),
					),
				)
			})

			It("should list workers in a cluster", func() {
				target := ClusterTargetHeader{}

				_, err := newWorker(server.URL()).ListAllWorkers("aaa", false, target)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("When ListAllWorkers is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/getWorkers"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to list worker`),
					),
				)
			})

			It("should return error during get worker", func() {
				target := ClusterTargetHeader{}
				_, err := newWorker(server.URL()).ListAllWorkers("aaa", false, target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//Get
	Describe("Get", func() {
		Context("When Get worker is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getWorker"),
						ghttp.RespondWith(http.StatusCreated, `{
							  "dedicatedHostId": "dedicatedhostid1",
							  "dedicatedHostPoolId": "dedicatedhostpoolid1",
							  "flavor": "string",
							  "health": {
								"message": "string",
								"state": "string"
							  },
							  "id": "string",
							  "kubeVersion": {
								"actual": "string",
								"desired": "string",
								"eos": "string",
								"masterEOS": "string",
								"target": "string"
							  },
							  "lifecycle": {
								"actualState": "string",
								"desiredState": "string",
								"message": "string",
								"messageDate": "string",
								"messageDetails": "string",
								"messageDetailsDate": "string",
								"pendingOperation": "string",
								"reasonForDelete": "string"
							  },
							  "location": "string",
							  "networkInterfaces": [
								{
								  "cidr": "string",
								  "ipAddress": "string",
								  "primary": true,
								  "subnetID": "string"
								}
							  ],
							  "poolID": "string",
							  "poolName": "string"
							}`),
					),
				)
			})

			It("should get workers in a cluster", func() {
				target := ClusterTargetHeader{}

				w, err := newWorker(server.URL()).Get("test", "kube-bmrtar0d0st4h9b09vm0-myclustervp-default-0000013", target)
				Expect(err).NotTo(HaveOccurred())
				Expect(w.HostID).To(BeEquivalentTo("dedicatedhostid1"))
				Expect(w.HostPoolID).To(BeEquivalentTo("dedicatedhostpoolid1"))
			})
		})
		Context("When get worker is unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/v2/vpc/getWorker"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to get worker`),
					),
				)
			})

			It("should return error during get worker", func() {
				target := ClusterTargetHeader{}
				_, err := newWorker(server.URL()).Get("test", "kube-bmrtar0d0st4h9b09vm0-myclustervp-default-0000013", target)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})

func newWorker(url string) Workers {

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
	return newWorkerAPI(&client)
}
