package satellitev1

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/client"
)

type CreateSatelliteEndpointSourceRequest struct {
	LocationID string   `json:"-"`
	Type       string   `json:"type"` // service
	SourceName string   `json:"source_name"`
	Addresses  []string `json:"addresses"`
}

type CreateSatelliteEndpointSourceResponse struct {
	LocationID string   `json:"location_id"`
	SourceName string   `json:"source_name"`
	Type       string   `json:"type"`
	Addresses  []string `json:"addresses"`
	CreatedAT  string   `json:"created_at"`
	LastChange string   `json:"last_change"`
	SourceID   string   `json:"source_id"`
}

type source struct {
	client     *client.Client
	pathPrefix string
}

type SatelliteSources struct {
	Sources []SatelliteSourceInfo `json:"sources"`
}

type SatelliteSourceInfo struct {
	Addresses  []string `json:"addresses"`
	CreatedAt  string   `json:"created_at"`
	LastChange string   `json:"last_change"`
	LocationID string   `json:"location_id"`
	SourceName string   `json:"source_name"`
	Type       string   `json:"type"`
	SourceID   string   `json:"source_id"`
}
type Source interface {
	//GetEndpointInfo(locationID, endpointID string, target containerv2.ClusterTargetHeader) (SatelliteEndpointInfo, error)
	CreateSatelliteEndpointSource(params CreateSatelliteEndpointSourceRequest,
		target containerv2.ClusterTargetHeader) (CreateSatelliteEndpointSourceResponse, error)
	ListSatelliteEndpointSources(locationID string,
		target containerv2.ClusterTargetHeader) (*SatelliteSources, error)
}

func newSourceAPI(c *client.Client) Source {
	return &source{
		client: c,
	}
}

func (s *source) CreateSatelliteEndpointSource(params CreateSatelliteEndpointSourceRequest,
	target containerv2.ClusterTargetHeader) (CreateSatelliteEndpointSourceResponse, error) {

	var source CreateSatelliteEndpointSourceResponse
	rawURL := fmt.Sprintf("/v1/locations/%s/sources", params.LocationID)
	_, err := s.client.Post(rawURL, params, &source, target.ToMap())
	return source, err
}

func (s *source) ListSatelliteEndpointSources(locationID string,
	target containerv2.ClusterTargetHeader) (*SatelliteSources, error) {

	SatSourceInfo := new(SatelliteSources)
	rawURL := fmt.Sprintf("v1/locations/%s/sources", locationID)
	_, err := s.client.Get(rawURL, &SatSourceInfo, target.ToMap())
	if err != nil {
		return nil, err
	}
	return SatSourceInfo, nil
}
