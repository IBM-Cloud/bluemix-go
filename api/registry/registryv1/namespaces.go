package registryv1

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type NamespaceTargetHeader struct {
	AccountID string
}

//ToMap ...
func (c NamespaceTargetHeader) ToMap() map[string]string {
	m := make(map[string]string, 1)
	m[accountIDHeader] = c.AccountID
	return m
}

//Subnets interface
type Namespaces interface {
	GetNamespaces(target NamespaceTargetHeader) ([]string, error)
	AddNamespace(target NamespaceTargetHeader, namespace string) (PutNamespaceResponse, error)
	DeleteNamespace(target NamespaceTargetHeader, namespace string) error
}

type namespaces struct {
	client *client.Client
}

func newNamespaceAPI(c *client.Client) Namespaces {
	return &namespaces{
		client: c,
	}
}

type PutNamespaceResponse struct {
	Namespace string `json:"namespace,omitempty"`
}

//Create ...
func (r *namespaces) GetNamespaces(target NamespaceTargetHeader) ([]string, error) {

	var retVal []string
	req := rest.GetRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/api/v1/namespaces"))

	for key, value := range target.ToMap() {
		req.Set(key, value)
	}

	_, err := r.client.SendRequest(req, &retVal)
	return retVal, err
}

//Add ...
func (r *namespaces) AddNamespace(target NamespaceTargetHeader, namespace string) (PutNamespaceResponse, error) {

	var retVal PutNamespaceResponse
	req := rest.PutRequest(helpers.GetFullURL(*r.client.Config.Endpoint, fmt.Sprintf("/api/v1/namespaces/%s", namespace)))

	for key, value := range target.ToMap() {
		req.Set(key, value)
	}

	_, err := r.client.SendRequest(req, &retVal)
	return retVal, err
}

//Delete...
func (r *namespaces) DeleteNamespace(target NamespaceTargetHeader, namespace string) error {

	req := rest.DeleteRequest(helpers.GetFullURL(*r.client.Config.Endpoint, fmt.Sprintf("/api/v1/namespaces/%s", namespace)))

	for key, value := range target.ToMap() {
		req.Set(key, value)
	}

	_, err := r.client.SendRequest(req, nil)
	return err
}
