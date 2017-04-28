package cfv2

import (
	"fmt"

	bluemix "github.com/IBM-Bluemix/bluemix-go"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
)

//ErrCodeOrgDoesnotExist ...
var ErrCodeOrgDoesnotExist = "OrgDoesnotExist"

//OrgResource ...
type OrgResource struct {
	Resource
	Entity OrgEntity
}

//OrgEntity ...
type OrgEntity struct {
	Name           string `json:"name"`
	Region         string `json:"region"`
	BillingEnabled bool   `json:"billing_enabled"`
}

//ToFields ..
func (resource OrgResource) ToFields() Organization {
	entity := resource.Entity

	return Organization{
		GUID:           resource.Metadata.GUID,
		Name:           entity.Name,
		Region:         entity.Region,
		BillingEnabled: entity.BillingEnabled,
	}
}

//Organization model
type Organization struct {
	GUID           string
	Name           string
	Region         string
	BillingEnabled bool
}

//Organizations ...
type Organizations interface {
	FindByName(orgName string) (*Organization, error)
}

type organization struct {
	client *cfAPIClient
	config *bluemix.Config
}

func newOrganizationAPI(c *cfAPIClient) Organizations {
	return &organization{
		client: c,
		config: c.config,
	}
}

//FindByName ...
func (r *organization) FindByName(name string) (*Organization, error) {
	region := r.config.Region
	path, err := r.urlOfOrgWithName(name, false)
	if err != nil {
		return nil, err
	}

	var org Organization
	var found bool
	err = r.listOrgResourcesWithPath(path, func(orgResource OrgResource) bool {
		org = orgResource.ToFields()
		found = true
		return false
	})

	if err != nil {
		return nil, err
	}

	if found {
		return &org, err
	}

	//May not be found and no error
	return nil, bmxerror.New(ErrCodeOrgDoesnotExist,
		fmt.Sprintf("Given org %q doesn't exist in the given region %q", name, region))

}

func (r *organization) listOrgResourcesWithPath(path string, cb func(OrgResource) bool) error {
	_, err := r.client.getPaginated(path, OrgResource{}, func(resource interface{}) bool {
		if orgResource, ok := resource.(OrgResource); ok {
			return cb(orgResource)
		}
		return false
	})
	return err
}

func (r *organization) urlOfOrgWithName(name string, inline bool) (string, error) {
	req := rest.GetRequest("/v2/organizations").Query("q", fmt.Sprintf("name:%s", name))

	if inline {
		req.Query("inline-relations-depth", "1")
	}
	return r.url(req)
}

func (r *organization) url(req *rest.Request) (string, error) {
	httpReq, err := req.Build()
	if err != nil {
		return "", err
	}
	return httpReq.URL.String(), nil
}
