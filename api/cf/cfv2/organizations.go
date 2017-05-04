package cfv2

import (
	"fmt"

	"github.ibm.com/ashishth/bluemix-go/rest"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/client"
)

//ErrCodeOrgDoesnotExist ...
var ErrCodeOrgDoesnotExist = "OrgDoesnotExist"

//Metadata ...
type Metadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//Resource ...
type Resource struct {
	Metadata Metadata
}

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
	client *client.Client
}

func newOrganizationAPI(c *client.Client) Organizations {
	return &organization{
		client: c,
	}
}

//FindByName ...
func (o *organization) FindByName(name string) (*Organization, error) {
	path, err := o.urlOfOrgWithName(name, false)
	if err != nil {
		return nil, err
	}

	var org Organization
	var found bool
	err = o.listOrgResourcesWithPath(path, func(orgResource OrgResource) bool {
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
		fmt.Sprintf("Given org %q doesn't exist", name))

}

func (o *organization) listOrgResourcesWithPath(path string, cb func(OrgResource) bool) error {
	_, err := o.client.GetPaginated(path, OrgResource{}, func(resource interface{}) bool {
		if orgResource, ok := resource.(OrgResource); ok {
			return cb(orgResource)
		}
		return false
	})
	return err
}

func (o *organization) urlOfOrgWithName(name string, inline bool) (string, error) {
	req := rest.GetRequest("/v2/organizations").Query("q", fmt.Sprintf("name:%s", name))

	if inline {
		req.Query("inline-relations-depth", "1")
	}
	return o.url(req)
}

func (o *organization) url(req *rest.Request) (string, error) {
	httpReq, err := req.Build()
	if err != nil {
		return "", err
	}
	return httpReq.URL.String(), nil
}
