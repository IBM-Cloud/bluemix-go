package csev2

import (
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/client"
)

type ServiceCSE struct {
	Srvid            string   `json:"srvid"`
	ServiceName      string   `json:"service"`
	CustomerName     string   `json:"customer"`
	ServiceAddresses []string `json:"serviceAddresses"`
	EstadoProto      string   `json:"estadoProto"`
	EstadoPort       int      `json:"estadoPort"`
	EstadoPath       string   `json:"estadoPath"`
	TCPPorts         []int    `json:"tcpports"`
	UDPPorts         []int    `json:"udpports"`
	TCPRange         string   `json:"tcpportrange"`
	UDPRange         string   `json:"udpportrange"`
	Region           string   `json:"region"`
	DataCenters      []string `json:"dataCenters"`
	ACL              []string `json:"acl"`
	MaxSpeed         string   `json:"maxSpeed"`
	URL              string   `json:"url"`
	Dedicated        int      `json:"dedicated"`
	MultiTenant      int      `json:"multitenant"`
}

type ServiceEndpoint struct {
	Seid          string `json:"seid"`
	StaticAddress string `json:"staticAddress"`
	Netmask       string `json:"netmask"`
	DNSStatus     string `json:"dnsStatus"`
	DataCenter    string `json:"dataCenter"`
	Status        string `json:"status"`
}

type ServiceObject struct {
	Service   ServiceCSE        `json:"service"`
	Endpoints []ServiceEndpoint `json:"endpoints"`
}

type ServiceEndpoints interface {
	GetServiceEndpoint(srvID string) (*ServiceObject, error)
	CreateServiceEndpoint(payload map[string]interface{}) (string, error)
	UpdateServiceEndpoint(srvID string, payload map[string]interface{}) error
	DeleteServiceEndpoint(srvID string) error
}

type serviceendpoints struct {
	client *client.Client
}

func newServiceEndpointsAPI(c *client.Client) ServiceEndpoints {
	return &serviceendpoints{
		client: c,
	}
}

func (r *serviceendpoints) GetServiceEndpoint(srvID string) (*ServiceObject, error) {
	srvObj := ServiceObject{}
	rawURL := fmt.Sprintf("/v2/serviceendpoint/%s", srvID)
	_, err := r.client.Get(rawURL, &srvObj, nil)
	if err != nil {
		return nil, err
	}
	return &srvObj, nil
}

func (r *serviceendpoints) DeleteServiceEndpoint(srvID string) error {
	rawURL := fmt.Sprintf("/v2/serviceendpoint/%s", srvID)
	_, err := r.client.Delete(rawURL)
	if err != nil {
		return err
	}
	return nil
}

/*
payload map includes:
  KEY                 VALUETYPE        DESCRIPTION                  MANDATORY
  service             string           service name                 yes
  customer            string           customer name                yes
  serviceAddresses    []string         service backend addresses    yes
  estadoProto         string           estado check protocol        no
                                       tcp, http or https

  estadoPort          int              estado check port            no
  estadoPath          string           estado check path            no
  tcpports            []int            tcp access ports             yes
  udpports            []int            udp access ports             no
  tcpportrange        string           tcp access portrange         no
  udpportrange        string           udp access portrange         no
  region              string           region where mb is put       yes
  dataCenters         []string         datacenter where mb is put   yes
  acl                 []string         white networks for access    no
  maxSpeed            string           1g or 20g                    yes
  dedicated           int              indicate use dedicated       no
                                       mb or not
  multitenant         int              indicate to share mb with    no
                                       other servicenedpoints or
                                       not
*/
func (r *serviceendpoints) CreateServiceEndpoint(payload map[string]interface{}) (string, error) {
	rawURL := "/v2/serviceendpoint"
	result := make(map[string]interface{})
	_, err := r.client.Post(rawURL, &payload, &result)
	if err != nil {
		return "", err
	}

	return result["serviceid"].(string), nil
}

/*
For update, only below fields are supported:
   serviceAddresses,
   dataCenters,
   tcpports,
   udpports,
   tcpportrange,
   udpportrange,
   estadoProto,
   estadoPort,
   estadoPath,
   acl
*/
func (r *serviceendpoints) UpdateServiceEndpoint(srvID string, payload map[string]interface{}) error {
	rawURL := fmt.Sprintf("/v2/serviceendpointtf/%s", srvID)
	result := make(map[string]interface{})
	_, err := r.client.Put(rawURL, &payload, &result)
	if err != nil {
		return err
	}

	return nil
}
