package cfv2

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//Space ...
type Space struct {
	GUID           string
	Name           string
	OrgGUID        string
	SpaceQuotaGUID string
	AllowSSH       bool
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

//Spaces ...
type Spaces interface {
	ListSpacesInOrg(orgGUID string) ([]Space, error)
	FindByNameInOrg(orgGUID string, name string) (*Space, error)
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
