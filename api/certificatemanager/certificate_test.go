package certificatemanager

import (
	"log"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/client"
	bluemixHttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/session"
)

var _ = Describe("Certificate_Manager", func() {

	//CertID := "qwertyuiopasdfghjklmnbvcxzsdfghyhbn"
	data := models.Data{
		Content:                 "asdfghjkl",
		Privatekey:              "",
		IntermediateCertificate: "",
	}

	importdata := models.CertificateImportData{
		Name:        "Test",
		Description: "Test Certificate",
		Data:        data,
	}
	updatemetaData := models.CertificateMetadataUpdate{
		Name:        "Test",
		Description: "Test Certificate",
	}

	reimportData := models.CertificateReimportData{
		Content:                 "asdfghjk",
		Privatekey:              "qwertyuijhgfd",
		IntermediateCertificate: "sdfghjmnbvcdf",
	}
	orderdata := models.CertificateOrderData{
		Name:                   "string",
		Description:            "string",
		Domains:                []string{"asdfg"},
		DomainValidationMethod: "dns-01",
		DNSProviderInstanceCrn: "string",
		Issuer:                 "Let's Encrypt",
		Algorithm:              "sha256WithRSAEncryption",
		KeyAlgorithm:           "rsaEncryption 2048 bit",
		AutoRenewEnabled:       false,
	}

	var server *ghttp.Server
	AfterEach(func() {
		server.Close()
	})
	// import
	Describe("Import Certificate", func() {
		Context("Server returns response", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/api/v3/zxcvbnmmkjhgfdsaqwertyuiokjhgfdsa/certificates/import"),
						ghttp.RespondWith(http.StatusOK, `{
							"_id": "crn:v1:bluemix:public:cloudcerts:us-south:a/4448261269a14562b839e0a3019ed980:f352ce16-97c6-436c-a7b2-0f9bbe3fecf1:certificate:bd06f268c1b9ebf4d9aa403b556a132f",
							"name": "aaa",
							"issuer": "etcd ca",
							"domains": [
							  "ca"
							],
							"begins_on": 1510768920000,
							"expires_on": 1826128920000,
							"imported": true,
							"status": "valid",
							"algorithm": "sha256WithRSAEncryption",
							"key_algorithm": "rsaEncryption 2048 bit",
							"description": "aaaa",
							"has_previous": false
						  }`),
					),
				)
			})
			It("Should return a response", func() {
				resp, err := newCert(server.URL()).ImportCertificate("zxcvbnmmkjhgfdsaqwertyuiokjhgfdsa", importdata)
				Expect(resp).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/api/v3/zxcvbnmmkjhgfdsaqwertyuiokjhgfdsa/certificates/import"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to Import`),
					),
				)
			})
			It("should return error ", func() {
				_, err := newCert(server.URL()).ImportCertificate("zxcvbnmmkjhgfdsaqwertyuiokjhgfdsa", importdata)
				Expect(err).To(HaveOccurred())

			})

		})
	})

	// get cert data
	Describe("Get certificate data", func() {
		Context("Server returns response", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v2/certificate/qwertyuiolkjhgfdsazxcvbnm"),
						ghttp.RespondWith(http.StatusOK, `{
								"_id": "crn:v1:bluemix:public:cloudcerts:us-south:a/4448261269a14562b839e0a3019ed980:f352ce16-97c6-436c-a7b2-0f9bbe3fecf1:certificate:bd06f268c1b9ebf4d9aa403b556a132f",
								"name": "aaa",
								"issuer": "etcd ca",
								"domains": [
								  "ca"
								],
								"begins_on": 1510768920000,
								"expires_on": 1826128920000,
								"imported": true,
								"status": "valid",
								"algorithm": "sha256WithRSAEncryption",
								"key_algorithm": "rsaEncryption 2048 bit",
								"description": "aaaa",
								"has_previous": false,
								"data": {
									"content": "-----BEGIN CERTIFICATE-----\r\nMIIDsTCCApmgAwIBAgIUWzsBehxAkgLLYBUZEUpSjHkIaMowDQYJKoZIhvcNAQEL\r\nBQAwbzEMMAoGA1UEBhMDVVNBMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQH\r\nEw1TYW4gRnJhbmNpc2NvMQ0wCwYDVQQKEwRldGNkMRYwFAYDVQQLEw1ldGNkIFNl\r\nY3VyaXR5MQswCQYDVQQDEwJjYTAeFw0xNzExMTUxODAyMDBaFw0yNzExMTMxODAy\r\nMDBaMG8xDDAKBgNVBAYTA1VTQTETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UE\r\nBxMNU2FuIEZyYW5jaXNjbzENMAsGA1UEChMEZXRjZDEWMBQGA1UECxMNZXRjZCBT\r\nZWN1cml0eTELMAkGA1UEAxMCY2EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\r\nAoIBAQCxjHVNtcCSCz1w9AiN7zAql0ZsPN6MNQWJ2j3iPCvmy9oi0wqSfYXTs+xw\r\nY4Q+j0dfA54+PcyIOSBQCZBeLLIwCaXN+gLkMxYEWCCVgWYUa6UY+NzPKRCfkbwG\r\noE2Ilv3R1FWIpMqDVE2rLmTb3YxSiw460Ruv4l16kodEzfs4BRcqrEiobBwaIMLd\r\n0rDJju7Q2TcioNji+HFoXV2aLN58LDgKO9AqszXxW88IKwUspfGBcsA4Zti/OHr+\r\nW+i/VxsxnQSJiAoKYbv9SkS8fUWw2hQ9SBBCKqE3jLzI71HzKgjS5TiQVZJaD6oK\r\ncw8FjexOELZd4r1+/p+nQdKqwnb5AgMBAAGjRTBDMA4GA1UdDwEB/wQEAwIBBjAS\r\nBgNVHRMBAf8ECDAGAQH/AgECMB0GA1UdDgQWBBRLfPxmhlZix1eTdBMAzMVlAnOV\r\ngTANBgkqhkiG9w0BAQsFAAOCAQEAeT2NfOt3WsBLUVcnyGMeVRQ0gXazxJXD/Z+3\r\n2RF3KClqBLuGmPUZVl0FU841J6hLlwNjS33mye7k2OHrjJcouElbV3Olxsgh/EV0\r\nJ7b7Wf4zWYHFNZz/VxwGHunsEZ+SCXUzU8OiMrEcHkOVzhtbC2veVPJzrESqd88z\r\nm1MseGW636VIcrg4fYRS9EebRPFvlwfymMd+bqLky9KsUbjNupYd/TlhpAudrIzA\r\nwO9ZUDb/0P44iOo+xURCoodxDTM0vvfZ8eJ6VZ/17HIf/a71kvk1oMqEhf060nmF\r\nIxnbK6iUqqhV8DLE1869vpFvgbDdOxP7BeabN5FXEnZFDTLDqg==\r\n-----END CERTIFICATE-----"
								  },
								  "last_modified": 1575527437586,
								  "data_key_id": "data_key",
								  "has_previous": true
								

							  }`),
					),
				)
			})
			It("Should return a response", func() {
				resp, err := newCert(server.URL()).GetCertData("qwertyuiolkjhgfdsazxcvbnm")
				Expect(resp).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v2/certificate/qwertyuiolkjhgfdsazxcvbnm"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to Get certificate data`),
					),
				)
			})
			It("should return error ", func() {
				_, err := newCert(server.URL()).GetCertData("qwertyuiolkjhgfdsazxcvbnm")
				Expect(err).To(HaveOccurred())

			})

		})
	})

	// get meta  data
	Describe("Get certificate meta data", func() {
		Context("Server returns response", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v1/certificate/qwertyuiolkjhgfdsazxcvbnm/metadata"),
						ghttp.RespondWith(http.StatusOK, `{
							"_id": "crn:v1:bluemix:public:cloudcerts:us-south:a/4448261269a14562b839e0a3019ed980:f352ce16-97c6-436c-a7b2-0f9bbe3fecf1:certificate:bd06f268c1b9ebf4d9aa403b556a132f",
							"name": "aaa",
							"issuer": "etcd ca",
							"domains": [
							  "ca"
							],
							"begins_on": 1510768920000,
							"expires_on": 1826128920000,
							"imported": true,
							"status": "valid",
							"algorithm": "sha256WithRSAEncryption",
							"key_algorithm": "rsaEncryption 2048 bit",
							"description": "aaaa",
							"has_previous": false
						  }`),
					),
				)
			})
			It("Should return a response", func() {
				resp, err := newCert(server.URL()).GetMetaData("qwertyuiolkjhgfdsazxcvbnm")
				Expect(resp).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, "/api/v1/certificate/qwertyuiolkjhgfdsazxcvbnm/metadata"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to Get certificate meta data`),
					),
				)
			})
			It("should return error ", func() {
				_, err := newCert(server.URL()).GetMetaData("qwertyuiolkjhgfdsazxcvbnm")
				Expect(err).To(HaveOccurred())

			})

		})
	})

	// Delete Certificate
	Describe("Delete certificate", func() {
		Context("Server returns response", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/api/v2/certificate/qwertyuiolkjhgfdsazxcvbnm"),
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)
			})
			It("Should return a response", func() {
				err := newCert(server.URL()).DeleteCertificate("qwertyuiolkjhgfdsazxcvbnm")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, "/api/v2/certificate/qwertyuiolkjhgfdsazxcvbnm"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to Delete certificate data`),
					),
				)
			})
			It("should return error ", func() {
				err := newCert(server.URL()).DeleteCertificate("qwertyuiolkjhgfdsazxcvbnm")
				Expect(err).To(HaveOccurred())

			})

		})
	})

	// Update Certificate Meta Data
	Describe("Update certificate", func() {
		Context("Server returns response", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/api/v3/certificate/qwertyuiolkjhgfdsazxcvbnm"),
						ghttp.RespondWith(http.StatusOK, nil),
					),
				)
			})
			It("Should return a response", func() {
				err := newCert(server.URL()).UpdateCertificateMetaData("qwertyuiolkjhgfdsazxcvbnm", updatemetaData)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/api/v3/certificate/qwertyuiolkjhgfdsazxcvbnm"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to Update certificate metadata`),
					),
				)
			})
			It("should return error ", func() {
				err := newCert(server.URL()).UpdateCertificateMetaData("qwertyuiolkjhgfdsazxcvbnm", updatemetaData)
				Expect(err).To(HaveOccurred())

			})

		})
	})

	// Reimport Certificate
	Describe("ReImport Certificate", func() {
		Context("Server returns response", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, "/api/v1/certificate/zxcvbnmmkjhgfdsaqwertyuiokjhgfdsa"),
						ghttp.RespondWith(http.StatusOK, `{
							"_id": "string",
							"name": "string",
							"issuer": "string",
							"domains": [
							  "string"
							],
							"begins_on": 0,
							"expires_on": 0,
							"imported": true,
							"status": "active",
							"algorithm": "string",
							"key_algorithm": "string",
							"description": "string",
							"has_previous": true,
							"issuance_info": {
							  "ordered_on": 0,
							  "status": "pending",
							  "code": "string",
							  "additional_info": "string",
							  "auto": true
							},
							"order_policy": {
							  "auto_renew_enabled": false
							},
							"downloaded": true
						  }`),
					),
				)
			})
			It("Should return a response", func() {
				resp, err := newCert(server.URL()).ReimportCertificate("zxcvbnmmkjhgfdsaqwertyuiokjhgfdsa", reimportData)
				Expect(resp).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, "/api/v1/certificate/zxcvbnmmkjhgfdsaqwertyuiokjhgfdsa"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to ReImport`),
					),
				)
			})
			It("should return error ", func() {
				_, err := newCert(server.URL()).ReimportCertificate("zxcvbnmmkjhgfdsaqwertyuiokjhgfdsa", reimportData)
				Expect(err).To(HaveOccurred())

			})

		})
	})

	// Order Certificate
	Describe("Order Certificate", func() {
		Context("Server returns response", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/api/v1/asdfghjmnbvcdfgh/certificates/order"),
						ghttp.RespondWith(http.StatusOK, `{
							"_id": "string",
							"name": "string",
							"issuer": "string",
							"domains": [
							  "string"
							],
							"begins_on": 0,
							"expires_on": 0,
							"imported": true,
							"status": "expired",
							"algorithm": "string",
							"key_algorithm": "string",
							"description": "string",
							"has_previous": true,
							"issuance_info": {
							  "ordered_on": 0,
							  "status": "pending",
							  "code": "string",
							  "additional_info": "string",
							  "auto": true
							},
							"order_policy": {
							  "auto_renew_enabled": false
							}
						  }`),
					),
				)
			})
			It("Should return a response", func() {
				resp, err := newCert(server.URL()).OrderCertificate("asdfghjmnbvcdfgh", orderdata)
				Expect(resp).ShouldNot(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When unsuccessful", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.SetAllowUnhandledRequests(true)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, "/api/v1/asdfghjmnbvcdfgh/certificates/order"),
						ghttp.RespondWith(http.StatusInternalServerError, `Failed to Order`),
					),
				)
			})
			It("should return error ", func() {
				_, err := newCert(server.URL()).OrderCertificate("asdfghjmnbvcdfgh", orderdata)
				Expect(err).To(HaveOccurred())

			})

		})
	})

})

func newCert(url string) Certificate {

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}
	conf := sess.Config.Copy()
	conf.HTTPClient = bluemixHttp.NewHTTPClient(conf)
	conf.Endpoint = &url

	client := client.Client{
		Config:      conf,
		ServiceName: bluemix.CertificateManager,
	}

	return newCertificateAPI(&client)
}
