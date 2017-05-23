package cfv2

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//ErrCodePrivateDomainDoesnotExist ...
var ErrCodePrivateDomainDoesnotExist = "PrivateDomainDoesnotExist"

//PrivateDomainEntity ...
type PrivateDomainEntity struct {
	Name                   string `json:"name"`
	OwningOrganizationGUID string `json:"owning_organization_guid"`
	OwningOrganizationURL  string `json:"owning_organization_url"`
	SharedOrganizationURL  string `json:"shared_organizations_url"`
}

//PrivateDomainResource ...
type PrivateDomainResource struct {
	Resource
	Entity PrivateDomainEntity
}

//ToFields ..
func (resource PrivateDomainResource) ToFields() PrivateDomain {
	entity := resource.Entity

	return PrivateDomain{
		GUID: resource.Metadata.GUID,
		Name: entity.Name,
		OwningOrganizationGUID: entity.OwningOrganizationGUID,
		OwningOrganizationURL:  entity.OwningOrganizationURL,
		SharedOrganizationURL:  entity.OwningOrganizationURL,
	}
}

//PrivateDomain model
type PrivateDomain struct {
	GUID                   string
	Name                   string
	OwningOrganizationGUID string
	OwningOrganizationURL  string
	SharedOrganizationURL  string
}

//PrivateDomains ...
type PrivateDomains interface {
	FindByName(orgGUID, domainName string) (*PrivateDomain, error)
}

type privateDomain struct {
	client *client.Client
}

func newPrivateDomainAPI(c *client.Client) PrivateDomains {
	return &privateDomain{
		client: c,
	}
}

func (d *privateDomain) FindByName(orgGUID, domainName string) (*PrivateDomain, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/private_domains", orgGUID)
	req := rest.GetRequest(rawURL).Query("q", "name:"+domainName)
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	domain, err := listPrivateDomainWithPath(d.client, path)
	if err != nil {
		return nil, err
	}
	if len(domain) == 0 {
		return nil, bmxerror.New(ErrCodePrivateDomainDoesnotExist, fmt.Sprintf("Private Domain: %q doesn't exist", domainName))
	}
	return &domain[0], nil
}

func listPrivateDomainWithPath(c *client.Client, path string) ([]PrivateDomain, error) {
	var privateDomain []PrivateDomain
	_, err := c.GetPaginated(path, PrivateDomainResource{}, func(resource interface{}) bool {
		if privateDomainResource, ok := resource.(PrivateDomainResource); ok {
			privateDomain = append(privateDomain, privateDomainResource.ToFields())
			return true
		}
		return false
	})
	return privateDomain, err
}
