package cfv2

import (
	"fmt"
	"strconv"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//SpaceCreateRequest ...
type SpaceCreateRequest struct {
	Name           string `json:"name"`
	OrgGUID        string `json:"organization_guid"`
	SpaceQuotaGUID string `json:"space_quota_definition_guid,omitempty"`
}

//SpaceUpdateRequest ...
type SpaceUpdateRequest struct {
	Name *string `json:"name,omitempty"`
}

//Space ...
type Space struct {
	GUID           string
	Name           string
	OrgGUID        string
	SpaceQuotaGUID string
	AllowSSH       bool
}

//SpaceFields ...
type SpaceFields struct {
	Metadata SpaceMetadata
	Entity   SpaceEntity
}

//SpaceMetadata ...
type SpaceMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//ErrCodeSpaceDoesnotExist ...
const ErrCodeSpaceDoesnotExist = "SpaceDoesnotExist"

//SpaceResource ...
type SpaceResource struct {
	Resource
	Entity SpaceEntity
}

//SpaceEntity ...
type SpaceEntity struct {
	Name           string `json:"name"`
	OrgGUID        string `json:"organization_guid"`
	SpaceQuotaGUID string `json:"space_quota_definition_guid"`
	AllowSSH       bool   `json:"allow_ssh"`
}

//ToFields ...
func (resource *SpaceResource) ToFields() Space {
	entity := resource.Entity

	return Space{
		GUID:           resource.Metadata.GUID,
		Name:           entity.Name,
		OrgGUID:        entity.OrgGUID,
		SpaceQuotaGUID: entity.SpaceQuotaGUID,
		AllowSSH:       entity.AllowSSH,
	}
}

type RouteFilter struct {
	DomainGUID string
	Host       *string
	Path       *string
	Port       *int
}

//Spaces ...
type Spaces interface {
	ListSpacesInOrg(orgGUID string) ([]Space, error)
	FindByNameInOrg(orgGUID string, name string) (*Space, error)
	Create(req SpaceCreateRequest) (*SpaceFields, error)
	Update(spaceGUID string, req SpaceUpdateRequest) (*SpaceFields, error)
	Delete(spaceGUID string) error
	Get(spaceGUID string) (*SpaceFields, error)
	ListRoutes(spaceGUID string, req RouteFilter) ([]Route, error)
}

type spaces struct {
	client *client.Client
}

func newSpacesAPI(c *client.Client) Spaces {
	return &spaces{
		client: c,
	}
}

func (r *spaces) FindByNameInOrg(orgGUID string, name string) (*Space, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/spaces", orgGUID)
	req := rest.GetRequest(rawURL).Query("q", "name:"+name)

	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()

	spaces, err := r.listSpacesWithPath(path)

	if err != nil {
		return nil, err
	}
	if len(spaces) == 0 {
		return nil, bmxerror.New(ErrCodeSpaceDoesnotExist,
			fmt.Sprintf("Given space:  %q doesn't exist in given org: %q", name, orgGUID))

	}
	return &spaces[0], nil
}

func (r *spaces) ListSpacesInOrg(orgGUID string) ([]Space, error) {
	rawURL := fmt.Sprintf("v2/organizations/%s/spaces", orgGUID)
	req := rest.GetRequest(rawURL)

	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()

	return r.listSpacesWithPath(path)
}

func (r *spaces) listSpacesWithPath(path string) ([]Space, error) {
	var spaces []Space
	_, err := r.client.GetPaginated(path, SpaceResource{}, func(resource interface{}) bool {
		if spaceResource, ok := resource.(SpaceResource); ok {
			spaces = append(spaces, spaceResource.ToFields())
			return true
		}
		return false
	})
	return spaces, err
}
func (r *spaces) Create(req SpaceCreateRequest) (*SpaceFields, error) {
	rawURL := "/v2/spaces?accepts_incomplete=true&async=true"
	spaceFields := SpaceFields{}
	_, err := r.client.Post(rawURL, req, &spaceFields)
	if err != nil {
		return nil, err
	}
	return &spaceFields, nil
}

func (r *spaces) Get(spaceGUID string) (*SpaceFields, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s", spaceGUID)
	spaceFields := SpaceFields{}
	_, err := r.client.Get(rawURL, &spaceFields)
	if err != nil {
		return nil, err
	}

	return &spaceFields, err
}

func (r *spaces) Update(spaceGUID string, req SpaceUpdateRequest) (*SpaceFields, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s?accepts_incomplete=true&async=true", spaceGUID)
	spaceFields := SpaceFields{}
	_, err := r.client.Put(rawURL, req, &spaceFields)
	if err != nil {
		return nil, err
	}
	return &spaceFields, nil
}

func (r *spaces) Delete(spaceGUID string) error {
	rawURL := fmt.Sprintf("/v2/spaces/%s", spaceGUID)
	_, err := r.client.Delete(rawURL)
	return err
}

func (r *spaces) ListRoutes(spaceGUID string, routeFilter RouteFilter) ([]Route, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s/routes", spaceGUID)
	req := rest.GetRequest(rawURL)
	var query string
	if routeFilter.DomainGUID != "" {
		query = "domain_guid:" + routeFilter.DomainGUID + ";"
	}
	if routeFilter.Host != nil {
		query += "host:" + *routeFilter.Host + ";"
	}
	if routeFilter.Path != nil {
		query += "path:" + *routeFilter.Path + ";"
	}
	if routeFilter.Port != nil {
		query += "port:" + strconv.Itoa(*routeFilter.Port) + ";"
	}

	if len(query) > 0 {
		req.Query("q", query)
	}

	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	route, err := listRouteWithPath(r.client, path)
	if err != nil {
		return nil, err
	}
	return route, nil
}
