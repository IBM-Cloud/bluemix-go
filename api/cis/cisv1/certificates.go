package cisv1

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type Certificate struct {
	Id                string   `json:"id"`
	Hostnames         []string `json:"hostnames"`
	RequestType       string   `json:"request_type"`
	RequestedValidity int      `json:"requested_validity"`
	ExpiresOn         string   `json:"expires_on"`
	Certificate       string   `json:"certificate"`
	PrivateKey        string   `json:"private_key"`
	Csr               string   `json:"csr"`
}

type CertificateResults struct {
	CertificateList []Certificate `json:"result"`
	ResultsInfo     ResultsCount  `json:"result_info"`
	Success         bool          `json:"success"`
	Errors          []Error       `json:"errors"`
}

type CertificateResult struct {
	Certificate Certificate `json:"result"`
	Success     bool        `json:"success"`
	Errors      []Error     `json:"errors"`
	Messages    []string    `json:"messages"`
}

type CertificateBody struct {
	Hostnames         []string `json:"hostnames"`
	RequestType       string   `json:"request_type"`
	RequestedValidity int      `json:"requested_validity"`
}

type CertificateDelete struct {
	Result struct {
		CertificateId string
	} `json:"result"`
	Success  bool     `json:"success"`
	Errors   []Error  `json:"errors"`
	Messages []string `json:"messages"`
}

type Certificates interface {
	ListCertificates(cisId string, zoneId string) ([]Certificate, error)
	GetCertificate(cisId string, zoneId string, certId string) (*Certificate, error)
	CreateCertificate(cisId string, zoneId string, certBody CertificateBody) (*Certificate, error)
	DeleteCertificate(cisId string, zoneId string, certId string) error
	UpdateCertificate(cisId string, zoneId string, certId string, certBody CertificateBody) (*Certificate, error)
}

type certificates struct {
	client *client.Client
}

func newCertificatesApi(c *client.Client) Certificates {
	return &certificates{
		client: c,
	}
}

func (r *certificates) ListCertificates(cisId string, zoneId string) ([]Certificate, error) {
	certResults := CertificateResults{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/origin_certificates", cisId, zoneId)
	_, err := r.client.Get(rawURL, &certResults)
	if err != nil {
		return nil, err
	}
	return certResults.CertificateList, err
}

func (r *certificates) GetCertificate(cisId string, zoneId string, certId string) (*Certificate, error) {
	certResult := CertificateResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/origin_certificates/%s", cisId, zoneId, certId)
	_, err := r.client.Get(rawURL, &certResult)
	if err != nil {
		return nil, err
	}
	return &certResult.Certificate, nil
}

func (r *certificates) DeleteCertificate(cisId string, zoneId string, certId string) error {
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/origin_certificates/%s", cisId, zoneId, certId)
	_, err := r.client.Delete(rawURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *certificates) CreateCertificate(cisId string, zoneId string, certBody CertificateBody) (*Certificate, error) {
	certResult := CertificateResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/origin_certificates", cisId, zoneId)
	_, err := r.client.Post(rawURL, &certBody, &certResult)
	if err != nil {
		return nil, err
	}
	return &certResult.Certificate, nil
}

func (r *certificates) UpdateCertificate(cisId string, zoneId string, certId string, certBody CertificateBody) (*Certificate, error) {
	certResult := CertificateResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/origin_certificates/%s", cisId, zoneId, certId)
	_, err := r.client.Put(rawURL, &certBody, &certResult)
	if err != nil {
		return nil, err
	}
	return &certResult.Certificate, nil
}
