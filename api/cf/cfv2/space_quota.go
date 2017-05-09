package cfv2

import (
	"encoding/json"
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//SpaceQuotaCreateRequest ...
type SpaceQuotaCreateRequest struct {
	Name                    string `json:"name"`
	OrgGUID                 string `json:"organization_guid"`
	MemoryLimitInMB         int64  `json:"memory_limit,,omitempty"`
	InstanceMemoryLimitInMB int64  `json:"instance_memory_limit,,omitempty"`
	RoutesLimit             int    `json:"total_routes,omitempty"`
	ServicesLimit           int    `json:"total_services,omitempty"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
}

//SpaceQuotaUpdateRequest ...
type SpaceQuotaUpdateRequest struct {
	Name                    string `json:"name"`
	OrgGUID                 string `json:"organization_guid"`
	MemoryLimitInMB         int64  `json:"memory_limit,,omitempty"`
	InstanceMemoryLimitInMB int64  `json:"instance_memory_limit,,omitempty"`
	RoutesLimit             int    `json:"total_routes,omitempty"`
	ServicesLimit           int    `json:"total_services,omitempty"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
}

type SpaceQuota struct {
	GUID                    string
	Name                    string
	NonBasicServicesAllowed bool
	ServicesLimit           int
	RoutesLimit             int
	MemoryLimitInMB         int64
	InstanceMemoryLimitInMB int64
	TrialDBAllowed          bool
	AppInstanceLimit        int
	PrivateDomainsLimit     int
	AppTaskLimit            int
}

//SpaceQuotaFields ...
type SpaceQuotaFields struct {
	Metadata SpaceQuotaMetadata
	Entity   SpaceQuotaEntity
}

//SpaceQuotaMetadata ...
type SpaceQuotaMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//ErrCodeSpaceQuotaDoesnotExist ...
const ErrCodeSpaceQuotaDoesnotExist = "SpaceQuotaDoesnotExist"

type SpaceQuotaResource struct {
	Resource
	Entity SpaceQuotaEntity
}

type SpaceQuotaEntity struct {
	Name                    string      `json:"name"`
	NonBasicServicesAllowed bool        `json:"non_basic_services_allowed"`
	ServicesLimit           int         `json:"total_services"`
	RoutesLimit             int         `json:"total_routes"`
	MemoryLimitInMB         int64       `json:"memory_limit"`
	InstanceMemoryLimitInMB int64       `json:"instance_memory_limit"`
	TrialDBAllowed          bool        `json:"trial_db_allowed"`
	AppInstanceLimit        json.Number `json:"app_instance_limit"`
	PrivateDomainsLimit     json.Number `json:"total_private_domains"`
	AppTaskLimit            json.Number `json:"app_task_limit"`
}

func (resource SpaceQuotaResource) ToFields() SpaceQuota {
	entity := resource.Entity

	return SpaceQuota{
		GUID: resource.Metadata.GUID,
		Name: entity.Name,
		NonBasicServicesAllowed: entity.NonBasicServicesAllowed,
		ServicesLimit:           entity.ServicesLimit,
		RoutesLimit:             entity.RoutesLimit,
		MemoryLimitInMB:         entity.MemoryLimitInMB,
		InstanceMemoryLimitInMB: entity.InstanceMemoryLimitInMB,
		TrialDBAllowed:          entity.TrialDBAllowed,
		AppInstanceLimit:        NumberToInt(entity.AppInstanceLimit, -1),
		PrivateDomainsLimit:     NumberToInt(entity.PrivateDomainsLimit, -1),
		AppTaskLimit:            NumberToInt(entity.AppTaskLimit, -1),
	}
}

//SpaceQuotas ...
type SpaceQuotas interface {
	FindByName(name, orgGUID string) (*SpaceQuota, error)
	Create(name, orgGUID string, memoryLimit, instanceMemoryLimit int64, routesLimit, servicesLimit int, serviceAllowed bool) (*SpaceQuotaFields, error)
	Update(newName, spaceQuotaGUID, orgGUID string, memoryLimit, instanceMemoryLimit int64, routesLimit, servicesLimit int, serviceAllowed bool) (*SpaceQuotaFields, error)
	Delete(spaceQuotaGUID string) error
	Get(spaceQuotaGUID string) (*SpaceQuotaFields, error)
}

type spaceQuota struct {
	client *client.Client
}

func newSpaceQuotasAPI(c *client.Client) SpaceQuotas {
	return &spaceQuota{
		client: c,
	}
}

func (r *spaceQuota) FindByName(name, orgGUID string) (*SpaceQuota, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/space_quota_definitions", orgGUID)
	req := rest.GetRequest(rawURL)

	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()

	spaceQuotas, err := r.listSpaceQuotaWithPath(path)
	if err != nil {
		return nil, err
	}

	if len(spaceQuotas) == 0 {
		return nil, bmxerror.New(ErrCodeSpaceQuotaDoesnotExist,
			fmt.Sprintf("Given space quota  %q doesn't  exist for the organization %q", name, orgGUID))
	}

	for _, q := range spaceQuotas {
		if q.Name == name {
			return &q, nil
		}

	}
	return nil, bmxerror.New(ErrCodeSpaceQuotaDoesnotExist,
		fmt.Sprintf("Given space quota  %q doesn't  exist for the organization %q", name, orgGUID))
}

func (r *spaceQuota) listSpaceQuotaWithPath(path string) ([]SpaceQuota, error) {
	var spaceQuota []SpaceQuota
	_, err := r.client.GetPaginated(path, SpaceQuotaResource{}, func(resource interface{}) bool {
		if spaceQuotaResource, ok := resource.(SpaceQuotaResource); ok {
			spaceQuota = append(spaceQuota, spaceQuotaResource.ToFields())
			return true
		}
		return false
	})
	return spaceQuota, err
}

func (r *spaceQuota) Create(name, orgGUID string, memoryLimit, instanceMemoryLimit int64, routesLimit, servicesLimit int, serviceAllowed bool) (*SpaceQuotaFields, error) {
	payload := SpaceQuotaCreateRequest{
		Name:                    name,
		OrgGUID:                 orgGUID,
		MemoryLimitInMB:         memoryLimit,
		InstanceMemoryLimitInMB: instanceMemoryLimit,
		RoutesLimit:             routesLimit,
		ServicesLimit:           servicesLimit,
		NonBasicServicesAllowed: serviceAllowed,
	}
	rawURL := "/v2/space_quota_definitions?accepts_incomplete=true&async=true"
	spaceQuotaFields := SpaceQuotaFields{}
	_, err := r.client.Post(rawURL, payload, &spaceQuotaFields)
	if err != nil {
		return nil, err
	}
	return &spaceQuotaFields, nil
}

func (r *spaceQuota) Get(spaceQuotaGUID string) (*SpaceQuotaFields, error) {
	rawURL := fmt.Sprintf("/v2/space_quota_definitions/%s", spaceQuotaGUID)
	spaceQuotaFields := SpaceQuotaFields{}
	_, err := r.client.Get(rawURL, &spaceQuotaFields)
	if err != nil {
		return nil, err
	}

	return &spaceQuotaFields, err
}

func (r *spaceQuota) Update(newName, spaceQuotaGUID, orgGUID string, memoryLimit, instanceMemoryLimit int64, routesLimit, servicesLimit int, serviceAllowed bool) (*SpaceQuotaFields, error) {
	payload := SpaceQuotaUpdateRequest{
		Name:                    newName,
		OrgGUID:                 orgGUID,
		MemoryLimitInMB:         memoryLimit,
		InstanceMemoryLimitInMB: instanceMemoryLimit,
		RoutesLimit:             routesLimit,
		ServicesLimit:           servicesLimit,
		NonBasicServicesAllowed: serviceAllowed,
	}
	rawURL := fmt.Sprintf("/v2/space_quota_definitions/%s?accepts_incomplete=true&async=true", spaceQuotaGUID)
	spaceQuotaFields := SpaceQuotaFields{}
	_, err := r.client.Put(rawURL, payload, &spaceQuotaFields)
	if err != nil {
		return nil, err
	}
	return &spaceQuotaFields, nil
}

func (r *spaceQuota) Delete(spaceQuotaGUID string) error {
	rawURL := fmt.Sprintf("/v2/space_quota_definitions/%s", spaceQuotaGUID)
	_, err := r.client.Delete(rawURL)
	return err
}

func NumberToInt(number json.Number, defaultValue int) int {
	if number != "" {
		i, err := number.Int64()
		if err == nil {
			return int(i)
		}
	}
	return defaultValue
}
