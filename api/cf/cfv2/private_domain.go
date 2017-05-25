package cfv2

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//ErrCodePrivateDomainDoesnotExist ...
var ErrCodePrivateDomainDoesnotExist = "PrivateDomainDoesnotExist"

//PrivateDomainRequest ...
type PrivateDomainRequest struct {
	Name    string `json:"name,omitempty"`
	OrgGUID string `json:"owning_organization_guid,omitempty"`
}

//PrivateDomaineMetadata ...
type PrivateDomainMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

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

//PrivateDomainFields ...
type PrivateDomainFields struct {
	Metadata PrivateDomainMetadata
	Entity   PrivateDomainEntity
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
	FindByNameInOrg(orgGUID, domainName string) (*PrivateDomain, error)
	FindByName(domainName string) (*PrivateDomain, error)
	Create(req PrivateDomainRequest) (*PrivateDomainFields, error)
	Get(privateDomainGUID string) (*PrivateDomainFields, error)
	Delete(privateDomainGUID string, async bool) error
}

type privateDomain struct {
	client *client.Client
}

func newPrivateDomainAPI(c *client.Client) PrivateDomains {
	return &privateDomain{
		client: c,
	}
}

func (d *privateDomain) FindByNameInOrg(orgGUID, domainName string) (*PrivateDomain, error) {
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

func (d *privateDomain) FindByName(domainName string) (*PrivateDomain, error) {
	rawURL := fmt.Sprintf("/v2/private_domains")
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

func (d *privateDomain) Create(req PrivateDomainRequest) (*PrivateDomainFields, error) {
	rawURL := "/v2/private_domains"
	privateDomainFields := PrivateDomainFields{}
	_, err := d.client.Post(rawURL, req, &privateDomainFields)
	if err != nil {
		return nil, err
	}
	return &privateDomainFields, nil
}

func (d *privateDomain) Get(privateDomainGUID string) (*PrivateDomainFields, error) {
	rawURL := fmt.Sprintf("/v2/private_domains/%s", privateDomainGUID)
	privateDomainFields := PrivateDomainFields{}
	_, err := d.client.Get(rawURL, &privateDomainFields, nil)
	if err != nil {
		return nil, err
	}
	return &privateDomainFields, nil
}

func (d *privateDomain) Delete(privateDomainGUID string, async bool) error {
	rawURL := fmt.Sprintf("/v2/private_domains/%s", privateDomainGUID)
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
