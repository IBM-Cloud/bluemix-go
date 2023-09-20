package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type AlbCreateReq struct {
	Cluster         string `json:"cluster"`
	EnableByDefault bool   `json:"enableByDefault"`
	Type            string `json:"type"`
	ZoneAlb         string `json:"zone"`
	IngressImage    string `json:"ingressImage"`
}

type AlbCreateResp struct {
	Cluster string `json:"cluster"`
	Alb     string `json:"alb"`
}

type ClusterALB struct {
	ID                      string      `json:"id"`
	Region                  string      `json:"region"`
	DataCenter              string      `json:"dataCenter"`
	IsPaid                  bool        `json:"isPaid"`
	PublicIngressHostname   string      `json:"publicIngressHostname"`
	PublicIngressSecretName string      `json:"publicIngressSecretName"`
	ALBs                    []AlbConfig `json:"alb"`
}
type AlbConfig struct {
	AlbBuild             string `json:"albBuild"`
	AlbID                string `json:"albID"`
	AlbType              string `json:"albType"`
	AuthBuild            string `json:"authBuild"`
	Cluster              string `json:"cluster"`
	CreatedDate          string `json:"createdDate"`
	DisableDeployment    bool   `json:"disableDeployment"`
	Enable               bool   `json:"enable"`
	LoadBalancerHostname string `json:"loadBalancerHostname"`
	Name                 string `json:"name"`
	NumOfInstances       string `json:"numOfInstances"`
	Resize               bool   `json:"resize"`
	State                string `json:"state"`
	Status               string `json:"status"`
	ZoneAlb              string `json:"zone"`
}

// UpdateALBReq is the body of the v2 Update ALB API endpoint
type UpdateALBReq struct {
	ClusterID string   `json:"cluster" description:"The ID of the cluster on which the update ALB action shall be performed"`
	ALBBuild  string   `json:"albBuild" description:"The version of the build to which the ALB should be updated"`
	ALBList   []string `json:"albList" description:"The list of ALBs that should be updated to the requested albBuild"`
}

// AlbUpdateResp
type AlbUpdateResp struct {
	ClusterID string `json:"clusterID"`
}

type AlbImageVersions struct {
	DefaultK8sVersion    string   `json:"defaultK8sVersion"`
	SupportedK8sVersions []string `json:"supportedK8sVersions"`
}

// IngressStatus struct for the top level ingress status, for a cluster
type IngressStatus struct {
	Cluster                string `json:"cluster"`
	Status                 string `json:"status"`
	NonTranslatedStatus    string `json:"nonTranslatedStatus"`
	Message                string `json:"message"`
	StatusList             []IngressComponentStatus
	GeneralComponentStatus []V2IngressComponentStatus `json:"generalComponentStatus,omitempty"`
	ALBStatus              []V2IngressComponentStatus `json:"albStatus,omitempty"`
	RouterStatus           []V2IngressComponentStatus `json:"routerStatus,omitempty"`
	SubdomainStatus        []V2IngressComponentStatus `json:"subdomainStatus,omitempty"`
	SecretStatus           []V2IngressComponentStatus `json:"secretStatus,omitempty"`
	IgnoredErrors          []string                   `json:"ignoredErrors" description:"list of error codes that the user wants to ignore"`
}

// IngressComponentStatus status of individual ingress component
type IngressComponentStatus struct {
	Component string `json:"component"`
	Status    string `json:"status"`
	Type      string `json:"type"`
}

// V2IngressComponentStatus status of individual ingress component
type V2IngressComponentStatus struct {
	Component string   `json:"component,omitempty"`
	Status    []string `json:"status,omitempty"`
}

// ALBClusterHealthCheckConfig configuration for ALB in-cluster health check
type ALBClusterHealthCheckConfig struct {
	Cluster string `json:"cluster"`
	Enable  bool   `json:"enable"`
}

// IgnoredIngressStatusErrors
type IgnoredIngressStatusErrors struct {
	Cluster       string   `json:"cluster" description:"the ID or name of the cluster"`
	IgnoredErrors []string `json:"ignoredErrors" description:"list of error codes that the user wants to ignore"`
}

// IngressStatusState
type IngressStatusState struct {
	Cluster string `json:"cluster" description:"the ID or name of the cluster"`
	Enable  bool   `json:"enable" description:"true or false to enable or disable ingress status"`
}

type alb struct {
	client *client.Client
}

// Clusters interface
type Alb interface {
	CreateAlb(albCreateReq AlbCreateReq, target ClusterTargetHeader) (AlbCreateResp, error)
	DisableAlb(disableAlbReq AlbConfig, target ClusterTargetHeader) error
	EnableAlb(enableAlbReq AlbConfig, target ClusterTargetHeader) error
	UpdateAlb(updateAlbReq UpdateALBReq, target ClusterTargetHeader) error
	GetAlb(albid string, target ClusterTargetHeader) (AlbConfig, error)
	ListClusterAlbs(clusterNameOrID string, target ClusterTargetHeader) ([]AlbConfig, error)
	ListAlbImages(target ClusterTargetHeader) (AlbImageVersions, error)
	GetIngressStatus(clusterNameOrID string, target ClusterTargetHeader) (IngressStatus, error)
	GetAlbClusterHealthCheckConfig(clusterNameOrID string, target ClusterTargetHeader) (ALBClusterHealthCheckConfig, error)
	SetAlbClusterHealthCheckConfig(albHealthCheckReq ALBClusterHealthCheckConfig, target ClusterTargetHeader) error
	GetIgnoredIngressStatusErrors(clusterNameOrID string, target ClusterTargetHeader) (IgnoredIngressStatusErrors, error)
	AddIgnoredIngressStatusErrors(ignoredErrorsReq IgnoredIngressStatusErrors, target ClusterTargetHeader) error
	RemoveIgnoredIngressStatusErrors(ignoredErrorsReq IgnoredIngressStatusErrors, target ClusterTargetHeader) error
	SetIngressStatusState(ingressStatusStateReq IngressStatusState, target ClusterTargetHeader) error
}

func newAlbAPI(c *client.Client) Alb {
	return &alb{
		client: c,
	}
}

func (r *alb) CreateAlb(albCreateReq AlbCreateReq, target ClusterTargetHeader) (AlbCreateResp, error) {
	var successV AlbCreateResp
	_, err := r.client.Post("/v2/alb/vpc/createAlb", albCreateReq, &successV, target.ToMap())
	return successV, err
}

func (r *alb) DisableAlb(disableAlbReq AlbConfig, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/vpc/disableAlb", disableAlbReq, nil, target.ToMap())
	return err
}

func (r *alb) EnableAlb(enableAlbReq AlbConfig, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/vpc/enableAlb", enableAlbReq, nil, target.ToMap())
	return err
}

func (r *alb) GetAlb(albID string, target ClusterTargetHeader) (AlbConfig, error) {
	var successV AlbConfig
	_, err := r.client.Get(fmt.Sprintf("/v2/alb/getAlb?albID=%s", albID), &successV, target.ToMap())
	return successV, err
}

func (r *alb) UpdateAlb(updateAlbReq UpdateALBReq, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/updateAlb", updateAlbReq, nil, target.ToMap())
	return err
}

// ListClusterALBs returns the list of albs available for cluster
func (r *alb) ListClusterAlbs(clusterNameOrID string, target ClusterTargetHeader) ([]AlbConfig, error) {
	var successV ClusterALB
	rawURL := fmt.Sprintf("v2/alb/getClusterAlbs?cluster=%s", clusterNameOrID)
	_, err := r.client.Get(rawURL, &successV, target.ToMap())
	return successV.ALBs, err
}

// ListAlbImages lists the default and the supported ALB image versions
func (r *alb) ListAlbImages(target ClusterTargetHeader) (AlbImageVersions, error) {
	var successV AlbImageVersions
	_, err := r.client.Get("v2/alb/getAlbImages", &successV, target.ToMap())
	return successV, err
}

func (r *alb) GetIngressStatus(clusterNameOrID string, target ClusterTargetHeader) (IngressStatus, error) {
	var successV IngressStatus
	_, err := r.client.Get(fmt.Sprintf("/v2/alb/getStatus?cluster=%s", clusterNameOrID), &successV, target.ToMap())
	return successV, err
}

// GetAlbClusterHealthCheckConfig returns the ALB in-cluster healthcheck config
func (r *alb) GetAlbClusterHealthCheckConfig(clusterNameOrID string, target ClusterTargetHeader) (ALBClusterHealthCheckConfig, error) {
	var successV ALBClusterHealthCheckConfig
	_, err := r.client.Get(fmt.Sprintf("/v2/alb/getIngressClusterHealthcheck?cluster=%s", clusterNameOrID), &successV, target.ToMap())
	return successV, err
}

// SetAlbClusterHealthCheckConfig configure the ALB in-cluster healthcheck
func (r *alb) SetAlbClusterHealthCheckConfig(albHealthCheckReq ALBClusterHealthCheckConfig, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/setIngressClusterHealthcheck", albHealthCheckReq, nil, target.ToMap())
	return err
}

// GetIgnoredIngressStatusErrors lists of error codes that the user wants to ignore
func (r *alb) GetIgnoredIngressStatusErrors(clusterNameOrID string, target ClusterTargetHeader) (IgnoredIngressStatusErrors, error) {
	var successV IgnoredIngressStatusErrors
	_, err := r.client.Get(fmt.Sprintf("/v2/alb/listIgnoredIngressStatusErrors?cluster=%s", clusterNameOrID), &successV, target.ToMap())
	return successV, err
}

// AddIgnoredIngressStatusErrors
func (r *alb) AddIgnoredIngressStatusErrors(ignoredErrorsReq IgnoredIngressStatusErrors, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/addIgnoredIngressStatusErrors", ignoredErrorsReq, nil, target.ToMap())
	return err
}

// RemoveIgnoredIngressStatusErrors
func (r *alb) RemoveIgnoredIngressStatusErrors(ignoredErrorsReq IgnoredIngressStatusErrors, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Delete("/v2/alb/removeIgnoredIngressStatusErrors", ignoredErrorsReq, nil, target.ToMap())
	return err
}

// SetIngressStatusState
func (r *alb) SetIngressStatusState(ingressStatusStateReq IngressStatusState, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := r.client.Post("/v2/alb/setIngressStatusState", ingressStatusStateReq, nil, target.ToMap())
	return err
}
