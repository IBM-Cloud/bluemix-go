package cfv2

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//ServiceBindingRequest ...
type ServiceBindingRequest struct {
	ServiceInstanceGUID string `json:"service_instance_guid"`
	AppGUID             string `json:"app_guid"`
	Parameters          string `json:"parameters,omitempty"`
}

//ServiceBindingMetadata ...
type ServiceBindingMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//ServiceBindingEntity ...
type ServiceBindingEntity struct {
	ServiceInstanceGUID string                 `json:"service_instance_guid"`
	AppGUID             string                 `json:"app_guid"`
	Credentials         map[string]interface{} `json:"credentials"`
}

//ServiceBindingResource ...
type ServiceBindingResource struct {
	Resource
	Entity ServiceBindingEntity
}

//ServiceBindingFields ...
type ServiceBindingFields struct {
	Metadata ServiceBindingMetadata
	Entity   ServiceBindingEntity
}

//ServiceBinding model
type ServiceBinding struct {
	GUID                string
	ServiceInstanceGUID string
	AppGUID             string
	Credentials         map[string]interface{}
}

//ToFields ..
func (resource ServiceBindingResource) ToFields() ServiceBinding {
	entity := resource.Entity

	return ServiceBinding{
		GUID:                resource.Metadata.GUID,
		ServiceInstanceGUID: entity.ServiceInstanceGUID,
		AppGUID:             entity.AppGUID,
		Credentials:         entity.Credentials,
	}
}

//ServiceBindings ...
type ServiceBindings interface {
	Create(req ServiceBindingRequest) (*ServiceBindingFields, error)
	Get(guid string) (*ServiceBindingFields, error)
	Delete(guid string, async bool) error
}

type serviceBinding struct {
	client *client.Client
}

func newServiceBindingAPI(c *client.Client) ServiceBindings {
	return &serviceBinding{
		client: c,
	}
}

func (r *serviceBinding) Get(sbGUID string) (*ServiceBindingFields, error) {
	rawURL := fmt.Sprintf("/v2/service_bindings/%s", sbGUID)
	sbFields := ServiceBindingFields{}
	_, err := r.client.Get(rawURL, &sbFields, nil)
	if err != nil {
		return nil, err
	}
	return &sbFields, nil
}

func (r *serviceBinding) Create(req ServiceBindingRequest) (*ServiceBindingFields, error) {
	rawURL := "/v2/service_bindings"
	sbFields := ServiceBindingFields{}
	_, err := r.client.Post(rawURL, req, &sbFields)
	if err != nil {
		return nil, err
	}
	return &sbFields, nil
}

func (r *serviceBinding) Delete(guid string, async bool) error {
	rawURL := fmt.Sprintf("/v2/service_bindings/%s", guid)
	req := rest.GetRequest(rawURL).Query("recursive", "true")
	if async {
		req.Query("async", "true")
	}
	httpReq, err := req.Build()
	if err != nil {
		return err
	}
	path := httpReq.URL.String()
	_, err = r.client.Delete(path)
	return err
}
