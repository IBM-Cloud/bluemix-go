package graphql

import (
	"log"
	"net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/client"
	bluemixHttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/session"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("VNIs", func() {
	var server *ghttp.Server

	Describe("AttachToBareMetalNode", func() {
		Context("When attaching VNI to bare metal node is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/graphql"),
						ghttp.RespondWith(http.StatusOK, `{
							"data": {
								"addVirtualNetworkInterfaceToBareMetalNode": {
									"networkAttachment": {
										"attachedTo": { "id": "worker-123" },
										"virtualNetworkInterface": {
											"externalID": "r006-vni-123",
											"name": "test-vni",
											"primaryIPAddress": "10.240.0.5",
											"macAddress": "02:00:00:00:00:01"
										},
										"vlanID": 100
									}
								}
							}
						}`),
					),
				)
			})

			It("should attach VNI successfully", func() {
				target := containerv1.ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				workerID := "worker-123"
				input := AddVNIToBareMetalNodeInput{
					VirtualNetworkInterfaceID: "r006-vni-123",
					VlanID:                    100,
					Node:                      &workerID,
				}
				resp, err := newVNI(server.URL()).AttachToBareMetalNode(input, target)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).ShouldNot(BeNil())
				Expect(resp.NetworkAttachment.AttachedTo.ID).Should(Equal("worker-123"))
				Expect(resp.NetworkAttachment.VirtualNetworkInterface.ExternalID).Should(Equal("r006-vni-123"))
				Expect(*resp.NetworkAttachment.VlanID).Should(Equal(100))
			})
		})

		Context("When attaching VNI fails", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/graphql"),
						ghttp.RespondWith(http.StatusOK, `{
							"errors": [{
								"message": "VNI not found",
								"extensions": {"code": "E0404"}
							}],
							"data": {"addVirtualNetworkInterfaceToBareMetalNode": null}
						}`),
					),
				)
			})

			It("should return error", func() {
				target := containerv1.ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				workerID := "worker-123"
				input := AddVNIToBareMetalNodeInput{
					VirtualNetworkInterfaceID: "r006-vni-123",
					VlanID:                    100,
					Node:                      &workerID,
				}
				_, err := newVNI(server.URL()).AttachToBareMetalNode(input, target)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("VNI not found"))
			})
		})
	})

	Describe("DetachFromNode", func() {
		Context("When detaching VNI is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/graphql"),
						ghttp.RespondWith(http.StatusOK, `{
							"data": {
								"removeVirtualNetworkInterfaceFromNode": {
									"cluster": { "id": "cluster-123" },
									"node": { "id": "worker-123" },
									"virtualNetworkInterface": {
										"externalID": "r006-vni-123",
										"name": "test-vni"
									}
								}
							}
						}`),
					),
				)
			})

			It("should detach VNI successfully", func() {
				target := containerv1.ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				workerID := "worker-123"
				input := RemoveVNIFromNodeInput{
					VirtualNetworkInterfaceID: "r006-vni-123",
					Node:                      &workerID,
				}
				resp, err := newVNI(server.URL()).DetachFromNode(input, target)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).ShouldNot(BeNil())
				Expect(resp.Node.ID).Should(Equal("worker-123"))
				Expect(resp.VirtualNetworkInterface.ExternalID).Should(Equal("r006-vni-123"))
			})
		})
	})

	Describe("ListAttachments", func() {
		Context("When listing VNI attachments is successful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/graphql"),
						ghttp.RespondWith(http.StatusOK, `{
							"data": {
								"node": {
									"__typename": "KubernetesCluster",
									"networkAttachments": {
										"edges": [{
											"node": {
												"attachedTo": { "id": "worker-123" },
												"virtualNetworkInterface": {
													"externalID": "r006-vni-123",
													"name": "test-vni",
													"primaryIPAddress": "10.240.0.5",
													"macAddress": "02:00:00:00:00:01"
												},
												"vlanID": 100
											}
										}],
										"pageInfo": {
											"hasNextPage": false,
											"endCursor": null
										}
									}
								}
							}
						}`),
					),
				)
			})

			It("should list VNI attachments successfully", func() {
				target := containerv1.ClusterTargetHeader{
					OrgID:     "abc",
					SpaceID:   "def",
					AccountID: "ghi",
				}
				input := ListVNIAttachmentsInput{
					NodeID: "cluster-123",
				}
				resp, err := newVNI(server.URL()).ListAttachments(input, target)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).ShouldNot(BeNil())
				Expect(resp.AttachableType).Should(Equal("KubernetesCluster"))
				Expect(len(resp.Connection.Edges)).Should(Equal(1))
				Expect(resp.Connection.Edges[0].Node.VirtualNetworkInterface.ExternalID).Should(Equal("r006-vni-123"))
			})
		})
	})
})

func newVNI(url string) VNIs {
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.MccpService,
	}
	return NewVNIAPI(&client)
}
