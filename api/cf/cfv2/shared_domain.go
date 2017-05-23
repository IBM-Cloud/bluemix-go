package cfv2

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//ErrCodeSharedDomainDoesnotExist ...
var ErrCodeSharedDomainDoesnotExist = "SharedDomainDoesnotExist"

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

//SharedDomain ...
type SharedDomains interface {
	FindByName(domainName string) (*SharedDomain, error)
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
		return nil, fmt.Errorf("Shared Domain: %q doesn't exist", domainName)
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
