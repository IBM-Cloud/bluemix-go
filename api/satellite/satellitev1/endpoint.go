package satellitev1

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/client"
)

type SatelliteLocationInfo struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Region            string      `json:"region"`
	ResourceGroup     string      `json:"resourceGroup"`
	ResourceGroupName string      `json:"resourceGroupName"`
	PodSubnet         string      `json:"podSubnet"`
	ServiceSubnet     string      `json:"serviceSubnet"`
	CreatedDate       string      `json:"createdDate"`
	MasterKubeVersion string      `json:"masterKubeVersion"`
	TargetVersion     string      `json:"targetVersion"`
	WorkerCount       int         `json:"workerCount"`
	Location          string      `json:"location"`
	Datacenter        string      `json:"datacenter"`
	MultiAzCapable    bool        `json:"multiAzCapable"`
	Provider          string      `json:"provider"`
	State             string      `json:"state"`
	Status            string      `json:"status"`
	VersionEOS        string      `json:"versionEOS"`
	IsPaid            bool        `json:"isPaid"`
	Entitlement       string      `json:"entitlement"`
	Type              string      `json:"type"`
	Addons            interface{} `json:"addons"`
	EtcdPort          string      `json:"etcdPort"`
	MasterURL         string      `json:"masterURL"`
	Ingress           struct {
		Hostname   string `json:"hostname"`
		SecretName string `json:"secretName"`
		Status     string `json:"status"`
		Message    string `json:"message"`
	} `json:"ingress"`
	CaCertRotationStatus struct {
		Status              string `json:"status"`
		ActionTriggerDate   string `json:"actionTriggerDate"`
		ActionCompletedDate string `json:"actionCompletedDate"`
	} `json:"caCertRotationStatus"`
	ImageSecurityEnabled bool     `json:"imageSecurityEnabled"`
	DisableAutoUpdate    bool     `json:"disableAutoUpdate"`
	Crn                  string   `json:"crn"`
	WorkerZones          []string `json:"workerZones"`
	Lifecycle            struct {
		MasterStatus             string `json:"masterStatus"`
		MasterStatusModifiedDate string `json:"masterStatusModifiedDate"`
		MasterHealth             string `json:"masterHealth"`
		MasterState              string `json:"masterState"`
		ModifiedDate             string `json:"modifiedDate"`
	} `json:"lifecycle"`
	ServiceEndpoints struct {
		PrivateServiceEndpointEnabled bool   `json:"privateServiceEndpointEnabled"`
		PrivateServiceEndpointURL     string `json:"privateServiceEndpointURL"`
		PublicServiceEndpointEnabled  bool   `json:"publicServiceEndpointEnabled"`
		PublicServiceEndpointURL      string `json:"publicServiceEndpointURL"`
	} `json:"serviceEndpoints"`
	Features struct {
		KeyProtectEnabled bool `json:"keyProtectEnabled"`
		PullSecretApplied bool `json:"pullSecretApplied"`
	} `json:"features"`
	Vpcs      interface{} `json:"vpcs"`
	CosConfig struct {
		Region          string `json:"region"`
		Bucket          string `json:"bucket"`
		Endpoint        string `json:"endpoint"`
		ServiceInstance struct {
			Crn string `json:"crn"`
		} `json:"serviceInstance"`
	} `json:"cos_config"`
	Description string `json:"description"`
	Deployments struct {
		Enabled bool   `json:"enabled"`
		Message string `json:"message"`
	} `json:"deployments"`
	Hosts struct {
		Total     int `json:"total"`
		Available int `json:"available"`
	} `json:"hosts"`
	Iaas struct {
		Provider string `json:"provider"`
		Region   string `json:"region"`
	} `json:"iaas"`
	OpenVpnServerPort int `json:"open_vpn_server_port"`
}

type SatelliteEndpointInfo struct {
	Cert             *EndpointCerts `json:"certs,omitempty"`
	ClientHost       string         `json:"client_host"`
	ClientMutualAuth bool           `json:"client_mutual_auth"`
	ClientPort       int            `json:"client_port"`
	ClientProtocol   string         `json:"client_protocol"`
	ConnType         string         `json:"conn_type"`
	ConnectorPort    int            `json:"connector_port"`
	CreatedAt        time.Time      `json:"created_at"`
	CreatedBy        string         `json:"created_by"`
	Crn              string         `json:"crn"`
	DisplayName      string         `json:"display_name"`
	HTTPTunnelOnTCP  interface{}    `json:"http_tunnel_on_tcp"`
	LastChange       time.Time      `json:"last_change"`
	LocationID       string         `json:"location_id"`
	Performance      struct {
		Connection        int           `json:"connection"`
		RxBandwidth       int           `json:"rx_bandwidth"`
		TxBandwidth       int           `json:"tx_bandwidth"`
		Bandwidth         int           `json:"bandwidth"`
		ToCloudDataRate   int           `json:"to_cloud_data_rate"`
		FromCloudDataRate int           `json:"from_cloud_data_rate"`
		TotalDataRate     int           `json:"total_data_rate"`
		Connectors        []interface{} `json:"connectors"`
	} `json:"performance"`
	Region           string        `json:"region"`
	RejectUnauth     bool          `json:"reject_unauth"`
	ServerHost       string        `json:"server_host"`
	ServerMutualAuth bool          `json:"server_mutual_auth"`
	ServerPort       int           `json:"server_port"`
	ServerProtocol   string        `json:"server_protocol"`
	ServiceName      string        `json:"service_name"`
	Sni              interface{}   `json:"sni"`
	Sources          []interface{} `json:"sources"`
	Status           string        `json:"status"`
	Timeout          int           `json:"timeout"`
	CertificateInfo  string        `json:"certificate_info"`
	EndpointID       string        `json:"endpoint_id"`
}

type CreateEndpointResponse struct {
	LocationID       string         `json:"location_id"`
	Crn              string         `json:"crn"`
	ConnType         string         `json:"conn_type"`
	DisplayName      string         `json:"display_name"`
	ServiceName      string         `json:"service_name"`
	ClientHost       interface{}    `json:"client_host"`
	ClientPort       interface{}    `json:"client_port"`
	ServerHost       string         `json:"server_host"`
	ServerPort       int            `json:"server_port"`
	ConnectorPort    int            `json:"connector_port"`
	ClientProtocol   string         `json:"client_protocol"`
	ClientMutualAuth bool           `json:"client_mutual_auth"`
	ServerProtocol   string         `json:"server_protocol"`
	ServerMutualAuth bool           `json:"server_mutual_auth"`
	RejectUnauth     bool           `json:"reject_unauth"`
	Sources          []interface{}  `json:"sources"`
	Timeout          int            `json:"timeout"`
	HTTPTunnelOnTCP  interface{}    `json:"http_tunnel_on_tcp"`
	Cert             *EndpointCerts `json:"certs,omitempty"`
	Status           string         `json:"status"`
	CreatedBy        string         `json:"created_by"`
	CreatedAt        time.Time      `json:"created_at"`
	LastChange       time.Time      `json:"last_change"`
	Performance      struct {
		Connection  int           `json:"connection"`
		RxBandwidth int           `json:"rx_bandwidth"`
		TxBandwidth int           `json:"tx_bandwidth"`
		Bandwidth   int           `json:"bandwidth"`
		Connectors  []interface{} `json:"connectors"`
	} `json:"performance"`
	Region     string `json:"region"`
	EndpointID string `json:"endpoint_id"`
}

type CreateEndpointRequest struct {
	EndpointID       string         `json:"-"`
	ConnType         string         `json:"conn_type,omitempty"`
	DisplayName      string         `json:"display_name,omitempty"`
	ServerHost       string         `json:"server_host,omitempty"`
	ServerPort       int            `json:"server_port,omitempty"`
	Sni              string         `json:"sni,omitempty"`
	ClientProtocol   string         `json:"client_protocol,omitempty"`
	ClientMutualAuth bool           `json:"client_mutual_auth,omitempty"`
	ServerProtocol   string         `json:"server_protocol,omitempty"`
	ServerMutualAuth bool           `json:"server_mutual_auth,omitempty"`
	RejectUnauth     bool           `json:"reject_unauth,omitempty"`
	Timeout          int            `json:"timeout,omitempty"`
	CreatedBy        string         `json:"created_by,omitempty"`
	Cert             *EndpointCerts `json:"certs,omitempty"`
}

type EndpointCerts struct {
	Client struct {
		Cert struct {
			Filename     string `json:"filename,omitempty"`
			FileContents string `json:"file_contents,omitempty"`
		} `json:"cert,omitempty"`
	} `json:"client,omitempty"`
	Server struct {
		Cert struct {
			Filename     string `json:"filename,omitempty"`
			FileContents string `json:"file_contents,omitempty"`
		} `json:"cert,omitempty"`
	} `json:"server,omitempty"`
	Connector struct {
		Cert struct {
			Filename     string `json:"filename,omitempty"`
			FileContents string `json:"file_contents,omitempty"`
		} `json:"cert,omitempty"`
		Key struct {
			Filename     string `json:"filename,omitempty"`
			FileContents string `json:"file_contents,omitempty"`
		} `json:"key,omitempty"`
	} `json:"connector,omitempty"`
}

type Endpoint interface {
	GetEndpointInfo(locationID, endpointID string, target containerv2.ClusterTargetHeader) (SatelliteEndpointInfo, error)
	CreateSatelliteEndpoint(params CreateEndpointRequest, target containerv2.ClusterTargetHeader) (CreateEndpointResponse, error)
	GetEndpoints(locationID string, target containerv2.ClusterTargetHeader) ([]SatelliteEndpointInfo, error)
	DeleteEndpoint(locationID, endpointID string, target containerv2.ClusterTargetHeader) error
}

type endpoint struct {
	client     *client.Client
	pathPrefix string
}

func newEndpointAPI(c *client.Client) Endpoint {
	return &endpoint{
		client: c,
	}
}

func (s *endpoint) GetEndpointInfo(locationID, endpointID string, target containerv2.ClusterTargetHeader) (SatelliteEndpointInfo, error) {
	SatEndpointInfo := SatelliteEndpointInfo{}
	rawURL := fmt.Sprintf("v1/locations/%s/endpoints/%s", locationID, endpointID)
	_, err := s.client.Get(rawURL, &SatEndpointInfo, target.ToMap())
	return SatEndpointInfo, err
}

func (s *endpoint) GetEndpoints(locationID string, target containerv2.ClusterTargetHeader) ([]SatelliteEndpointInfo, error) {
	SatEndpointInfo := []SatelliteEndpointInfo{}
	rawURL := fmt.Sprintf("v1/locations/%s/endpoints", locationID)
	_, err := s.client.Get(rawURL, &SatEndpointInfo, target.ToMap())
	if err != nil {
		return nil, err
	}
	return SatEndpointInfo, err
}

func (s *endpoint) CreateSatelliteEndpoint(params CreateEndpointRequest, target containerv2.ClusterTargetHeader) (CreateEndpointResponse, error) {
	var endpoint CreateEndpointResponse
	rawURL := fmt.Sprintf("/v1/locations/%s/endpoints", params.EndpointID)
	_, err := s.client.Post(rawURL, params, &endpoint, target.ToMap())
	return endpoint, err
}

func (s *endpoint) DeleteEndpoint(locationID, endpointID string, target containerv2.ClusterTargetHeader) error {
	rawURL := fmt.Sprintf("v1/locations/%s/endpoints/%s", locationID, endpointID)
	_, err := s.client.Delete(rawURL, target.ToMap())
	return err
}
