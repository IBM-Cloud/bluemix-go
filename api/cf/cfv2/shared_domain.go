package cfv2

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//ErrCodeSharedDomainDoesnotExist ...
var ErrCodeSharedDomainDoesnotExist = "SharedDomainDoesnotExist"

//SharedDomainRequest ...
type SharedDomainRequest struct {
	Name            string `json:"name"`
	RouterGroupGUID string `json:"router_group_guid,omitempty"`
}

//SharedDomaineMetadata ...
type SharedDomainMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//SharedDomainEntity ...
type SharedDomainEntity struct {
	Name            string `json:"name"`
	RouterGroupGUID string `json:"router_group_guid"`
	RouterGroupType string `json:"router_group_type"`
}

//SharedDomainResource ...
type SharedDomainResource struct {
	Resource
	Entity SharedDomainEntity
}

//SharedDomainFields ...
type SharedDomainFields struct {
	Metadata SharedDomainMetadata
	Entity   SharedDomainEntity
}

//ToFields ..
func (resource SharedDomainResource) ToFields() SharedDomain {
	entity := resource.Entity

	return SharedDomain{
		GUID:            resource.Metadata.GUID,
		Name:            entity.Name,
		RouterGroupGUID: entity.RouterGroupGUID,
		RouterGroupType: entity.RouterGroupType,
	}
}

//SharedDomain model
type SharedDomain struct {
	GUID            string
	Name            string
	RouterGroupGUID string
	RouterGroupType string
}

//SharedDomains ...
type SharedDomains interface {
	FindByName(domainName string) (*SharedDomain, error)
	Create(req SharedDomainRequest) (*SharedDomainFields, error)
	Get(sharedDomainGUID string) (*SharedDomainFields, error)
	Delete(sharedDomainGUID string, async bool) error
}

type sharedDomain struct {
	client *client.Client
}

func newSharedDomainAPI(c *client.Client) SharedDomains {
	return &sharedDomain{
		client: c,
	}
}

func (d *sharedDomain) FindByName(domainName string) (*SharedDomain, error) {
	rawURL := "/v2/shared_domains"
	req := rest.GetRequest(rawURL).Query("q", "name:"+domainName)
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	domain, err := listSharedDomainWithPath(d.client, path)
	if err != nil {
		return nil, err
	}
	if len(domain) == 0 {
		return nil, bmxerror.New(ErrCodeSharedDomainDoesnotExist, fmt.Sprintf("Shared Domain: %q doesn't exist", domainName))
	}
	return &domain[0], nil
}

func listSharedDomainWithPath(c *client.Client, path string) ([]SharedDomain, error) {
	var sharedDomain []SharedDomain
	_, err := c.GetPaginated(path, SharedDomainResource{}, func(resource interface{}) bool {
		if sharedDomainResource, ok := resource.(SharedDomainResource); ok {
			sharedDomain = append(sharedDomain, sharedDomainResource.ToFields())
			return true
		}
		return false
	})
	return sharedDomain, err
}

func (d *sharedDomain) Create(req SharedDomainRequest) (*SharedDomainFields, error) {
	rawURL := "/v2/shared_domains"
	sharedDomainFields := SharedDomainFields{}
	_, err := d.client.Post(rawURL, req, &sharedDomainFields)
	if err != nil {
		return nil, err
	}
	return &sharedDomainFields, nil
}

func (d *sharedDomain) Get(sharedDomainGUID string) (*SharedDomainFields, error) {
	rawURL := fmt.Sprintf("/v2/shared_domains/%s", sharedDomainGUID)
	sharedDomainFields := SharedDomainFields{}
	_, err := d.client.Get(rawURL, &sharedDomainFields, nil)
	if err != nil {
		return nil, err
	}
	return &sharedDomainFields, nil
}

func (d *sharedDomain) Delete(sharedDomainGUID string, async bool) error {
	rawURL := fmt.Sprintf("/v2/shared_domains/%s", sharedDomainGUID)
	req := rest.GetRequest(rawURL).Query("recursive", "true")
	if async {
		req.Query("async", "true")
	}
	httpReq, err := req.Build()
	if err != nil {
		return err
	}
	path := httpReq.URL.String()
	_, err = d.client.Delete(path)
	return err
}
