package cisv1

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

var _ = Describe("Zones", func() {
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
                        ghttp.VerifyRequest(http.MethodPost, "/v1/crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a/zones"),
                        ghttp.RespondWith(http.StatusCreated, `
                           {
                              "result": {
                                "id": "3fefc35e7decadb111dcf85d723a4f20",
                                "name": "example.com",
                                "status": "pending",
                                "paused": false,
                                "name_servers": [
                                  "ns002.name.cloud.ibm.com",
                                  "ns007.name.cloud.ibm.com"
                                ],
                                "original_name_servers": [
                                  "ns005.name.cloud.ibm.com",
                                  "ns016.name.cloud.ibm.com"
                                ],
                                "original_registrar": null,
                                "original_dnshost": null,
                                "modified_on": "2018-05-04T14:16:28.369359Z",
                                "created_on": "2018-05-04T14:16:28.369359Z",
                                "vanity_name_servers": [],
                                "account": {
                                  "id": "796b4ef449812595ea9fe92d1e910756",
                                  "name": "b424c068-c944-4565-b0bf-b278e5ec98ed"
                                }
                              },
                              "success": true,
                              "errors": [],
                              "messages": []
                            }
                        `),
                    ),
                )
            })

            It("should return zone created", func() {
                params := v1.ZoneBody{Name: "wcpcloudnl.com"}
                myZone, err := Zones(server.URL()).CreateZone(target, params)
                Expect(err).NotTo(HaveOccurred())
                Expect(*myZone).ShouldNot(BeNil())
                Expect(*myZone.Id).Should(Equal("3fefc35e7decadb111dcf85d723a4f20"))
            })
        })
        Context("When creation is unsuccessful", func() {
            BeforeEach(func() {
                server = ghttp.NewServer()
                server.AppendHandlers(
                    ghttp.CombineHandlers(
                        ghttp.VerifyRequest(http.MethodPost, "v1/crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a/zones"),
                        ghttp.RespondWith(http.StatusInternalServerError, `Failed to create Zone`),
                    ),
                )
            })

            It("should return error during Zone creation", func() {
                params := v1.ZoneBody{Name: "wcpcloudnl.com"}
                myZone, err := newZone(server.URL()).CreateZone(target, params)
                Expect(err).To(HaveOccurred())
                Expect(*myZone).Should(BeNil())
            })
        })
    })
    //List
    Describe("List", func() {
        Context("When read of Zones is successful", func() {
            BeforeEach(func() {
                server = ghttp.NewServer()
                server.AppendHandlers(
                    ghttp.CombineHandlers(
                        ghttp.VerifyRequest(http.MethodGet, "/v1/crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a/zones"),
                        ghttp.RespondWith(http.StatusOK, `
                            {
                  "result": [
                    {
                      "id": "3fefc35e7decadb111dcf85d723a4f20",
                      "name": "example.com",
                      "status": "active",
                      "paused": false,
                      "name_servers": [
                        "ns002.name.cloud.ibm.com",
                        "ns007.name.cloud.ibm.com"
                      ],
                      "original_name_servers": [
                        "ns005.name.cloud.ibm.com",
                        "ns016.name.cloud.ibm.com"
                      ],
                      "original_registrar": null,
                      "original_dnshost": null,
                      "modified_on": "2018-10-12T06:34:35.992900Z",
                      "created_on": "2018-05-04T14:16:28.369359Z",
                      "vanity_name_servers": [],
                      "account": {
                        "id": "796b4ef449812595ea9fe92d1e910756",
                        "name": "b424c068-c944-4565-b0bf-b278e5ec98ed"
                      }
                    }
                  ],
                  "result_info": {
                    "page": 1,
                    "per_page": 20,
                    "total_pages": 1,
                    "count": 1,
                    "total_count": 1
                  },
                  "success": true,
                  "errors": [],
                  "messages": []
                }
              `),
                    ),
                )
            })

            It("should return Zone list", func() {
                target := "rn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a"
                myZones, err := Zones(server.URL()).ListZones(target)
                Expect(*myZones
                    ).ShouldNot(BeNil())
                for _, Zone := range *myZones {
                    Expect(err).NotTo(HaveOccurred())
                    Expect(Zone.Id).Should(Equal("f91adfe2-76c9-4649-939e-b01c37a3704"))
                }
            })
        })
        Context("When read of Zones is unsuccessful", func() {
            BeforeEach(func() {
                server = ghttp.NewServer()
                server.AppendHandlers(
                    ghttp.CombineHandlers(
                        ghttp.VerifyRequest(http.MethodGet, "/v1/crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a/zones"),
                        ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve Zones`),
                    ),
                )
            })

            It("should return error when Zone are retrieved", func() {
                target := "rn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a" 
                myZone, err := newZone(server.URL()).ListZones(target)
                Expect(err).To(HaveOccurred())
                Expect(*myZone).Should(BeNil())
            })
        })
    })
    //Delete
    Describe("Delete", func() {
        Context("When delete of Zone is successful", func() {
            BeforeEach(func() {
                server = ghttp.NewServer()
                server.AppendHandlers(
                    ghttp.CombineHandlers(
                        ghttp.VerifyRequest(http.MethodDelete, "/v1/crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a/zones/f91adfe2-76c9-4649-939e-b01c37a3704"),
                        ghttp.RespondWith(http.StatusOK, `{                         
                        }`),
                    ),
                )
            })

            It("should delete Zone", func() {
                target := "crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a"
                params := "f91adfe2-76c9-4649-939e-b01c37a3704"
                err := Zones(server.URL()).DeleteZone(target, params)
                Expect(err).NotTo(HaveOccurred())
            })
        })
        Context("When Zone delete has failed", func() {
            BeforeEach(func() {
                server = ghttp.NewServer()
                server.AppendHandlers(
                    ghttp.CombineHandlers(
                        ghttp.VerifyRequest(http.MethodDelete, "/v1/crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a/zones/f91adfe2-76c9-4649-939e-b01c37a3704"),
                        ghttp.RespondWith(http.StatusInternalServerError, `Failed to delete service key`),
                    ),
                )
            })

            It("should return error zone delete", func() {
                target := "crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a"
                params := "f91adfe2-76c9-4649-939e-b01c37a3704"
                err := Zones(server.URL()).DeleteZone(target, params)
                Expect(err).To(HaveOccurred())
            })
        })
    })
    //Find
    Describe("Get", func() {
        Context("When read of Zone is successful", func() {
            BeforeEach(func() {
                server = ghttp.NewServer()
                server.AppendHandlers(
                    ghttp.CombineHandlers(
                        ghttp.VerifyRequest(http.MethodGet, "/v1/crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a/zones/f91adfe2-76c9-4649-939e-b01c37a3704"),
                        ghttp.RespondWith(http.StatusOK, `
                            {
                          "result": {
                            "id": "3fefc35e7decadb111dcf85d723a4f20",
                            "name": "example.com",
                            "status": "active",
                            "paused": false,
                            "name_servers": [
                              "ns002.name.cloud.ibm.com",
                              "ns007.name.cloud.ibm.com"
                            ],
                            "original_name_servers": [
                              "ns005.name.cloud.ibm.com",
                              "ns016.name.cloud.ibm.com"
                            ],
                            "original_registrar": null,
                            "original_dnshost": null,
                            "modified_on": "2018-10-12T06:34:35.992900Z",
                            "created_on": "2018-05-04T14:16:28.369359Z",
                            "vanity_name_servers": [],
                            "account": {
                              "id": "796b4ef449812595ea9fe92d1e910756",
                              "name": "b424c068-c944-4565-b0bf-b278e5ec98ed"
                            }
                          },
                          "success": true,
                          "errors": [],
                          "messages": []
                        }
                    `),
                    ),
                )
            })

            It("should return Zone", func() {
               target := "crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a"
                params := "f91adfe2-76c9-4649-939e-b01c37a3704"
                myZone, err := Zones(server.URL()).GetZone(target, params)
                Expect(err).NotTo(HaveOccurred())
                Expect(*myZone).ShouldNot(BeNil())
                Expect(*myZone.Id).Should(Equal("f91adfe2-76c9-4649-939e-b01c37a3704"))
            })
        })
        Context("When Zone get has failed", func() {
            BeforeEach(func() {
                server = ghttp.NewServer()
                server.AppendHandlers(
                    ghttp.CombineHandlers(
                        ghttp.VerifyRequest(http.MethodGet, "/v1/crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a/zones/f91adfe2-76c9-4649-939e-b01c37a3704"),
                        ghttp.RespondWith(http.StatusInternalServerError, `Failed to retrieve Zone`),
                    ),
                )
            })

            It("should return error when Zone is retrieved", func() {
                target := "crn:v1:staging:public:iam::::apikey:ApiKey-62fefdd1-4557-4c7d-8a1c-f6da7ee2ff3a"
                params := "f91adfe2-76c9-4649-939e-b01c37a3704"
                myZone, err := Zones(server.URL()).GetZone(target, params)
                Expect(err).To(HaveOccurred())
                Expect(*myZone).Should(BeNil())
            })
        })
    })
 
})

func newServer(url string) Zones {

    sess, err := session.New()
    if err != nil {
        log.Fatal(err)
    }
    conf := sess.Config.Copy()
    conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
    conf.Endpoint = &url

    client := client.Client{
        Config:      conf,
        ServiceName: bluemix.CisService,
    }
    return newZoneAPI(&client)
}
