package cfv2

import (
	"fmt"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
)

//ErrCodeServiceKeyDoesNotExist ...
const ErrCodeServiceKeyDoesNotExist = "erviceKeyDoesNotExist"

//ServiceKeyRequest ...
type ServiceKeyRequest struct {
	Name                string                 `json:"name"`
	ServiceInstanceGUID string                 `json:"service_instance_guid"`
	Params              map[string]interface{} `json:"parameters,omitempty"`
}

//ServiceKey  model...
type ServiceKey struct {
	GUID                string
	Name                string                 `json:"name"`
	ServiceInstanceGUID string                 `json:"service_instance_guid"`
	ServiceInstanceURL  string                 `json:"service_instance_url"`
	Credentials         map[string]interface{} `json:"credentials"`
}

//ServiceKeyFields ...
type ServiceKeyFields struct {
	Metadata ServiceInstanceMetadata
	Entity   ServiceKey
}

//ServiceKeyMetadata ...
type ServiceKeyMetadata struct {
	GUID      string `json:"guid"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

//ServiceKeyResource ...
type ServiceKeyResource struct {
	Resource
	Entity ServiceKeyEntity
}

//ServiceKeyEntity ...
type ServiceKeyEntity struct {
	Name                string                 `json:"name"`
	ServiceInstanceGUID string                 `json:"service_instance_guid"`
	ServiceInstanceURL  string                 `json:"service_instance_url"`
	Credentials         map[string]interface{} `json:"credentials"`
}

//ToModel ...
func (resource ServiceKeyResource) ToModel() ServiceKey {

	entity := resource.Entity

	return ServiceKey{
		GUID:                resource.Metadata.GUID,
		Name:                entity.Name,
		ServiceInstanceGUID: entity.ServiceInstanceGUID,
		ServiceInstanceURL:  entity.ServiceInstanceURL,
		Credentials:         entity.Credentials,
	}
}

//ServiceKeys ...
type ServiceKeys interface {
	Create(serviceInstanceGUID string, keyName string, params map[string]interface{}) (*ServiceKeyFields, error)
	FindByName(serviceInstanceGUID string, keyName string) (*ServiceKey, error)
	Get(serviceKeyGUID string) (*ServiceKeyFields, error)
	Delete(serviceKeyGUID string) error
}

type serviceKey struct {
	client *cfAPIClient
	config *bluemix.Config
}

func newServiceKeyAPI(c *cfAPIClient) ServiceKeys {
	return &serviceKey{
		client: c,
		config: c.config,
	}
}

func (r *serviceKey) Create(serviceInstanceGUID string, keyName string, params map[string]interface{}) (*ServiceKeyFields, error) {
	serviceKeyFields := ServiceKeyFields{}
	reqParam := ServiceKeyRequest{
		ServiceInstanceGUID: serviceInstanceGUID,
		Name:                keyName,
		Params:              params,
	}
	_, err := r.client.post("/v2/service_keys", reqParam, &serviceKeyFields)
	if err != nil {
		return nil, err
	}
	return &serviceKeyFields, nil
}

func (r *serviceKey) Delete(serviceKeyGUID string) error {
	rawURL := fmt.Sprintf("/v2/service_keys/%s", serviceKeyGUID)
	_, err := r.client.delete(rawURL)
	return err
}

func (r *serviceKey) Get(guid string) (*ServiceKeyFields, error) {
	rawURL := fmt.Sprintf("/v2/service_keys/%s", guid)
	serviceKeyFields := ServiceKeyFields{}
	_, err := r.client.get(rawURL, &serviceKeyFields)
	if err != nil {
		return nil, err
	}

	return &serviceKeyFields, err
}

func (r *serviceKey) FindByName(serviceInstanceGUID string, keyName string) (*ServiceKey, error) {
	rawURL := fmt.Sprintf("/v2/service_instances/%s/service_keys", serviceInstanceGUID)
	req := rest.GetRequest(rawURL)
	if keyName != "" {
		req.Query("q", "name:"+keyName)
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	serviceKeys, err := r.listServiceKeysWithPath(path)
	if err != nil {
		return nil, err
	}
	if len(serviceKeys) == 0 {
		return nil, bmxerror.New(ErrCodeServiceKeyDoesNotExist,
			fmt.Sprintf("Given service key %q doesn't exist for the given service instance  %q", keyName, serviceInstanceGUID))
	}
	return &serviceKeys[0], nil
}

func (r *serviceKey) listServiceKeysWithPath(path string) ([]ServiceKey, error) {
	var serviceKeys []ServiceKey
	_, err := r.client.getPaginated(path, ServiceKeyResource{}, func(resource interface{}) bool {
		if serviceKeyResource, ok := resource.(ServiceKeyResource); ok {
			serviceKeys = append(serviceKeys, serviceKeyResource.ToModel())
			return true
		}
		return false
	})
	return serviceKeys, err
}
