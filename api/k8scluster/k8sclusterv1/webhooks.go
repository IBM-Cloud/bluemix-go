package k8sclusterv1

import (
	"fmt"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
)

//WebHook is the web hook
type WebHook struct {
	Level string
	Type  string
	URL   string
}

//Webhooks interface
type Webhooks interface {
	List(clusterName string, target *ClusterTargetHeader) ([]WebHook, error)
	Add(clusterName string, params WebHook, target *ClusterTargetHeader) error
}

type webhook struct {
	client *ClusterClient
	config *bluemix.Config
}

func newWebhookAPI(c *ClusterClient) Webhooks {
	return &webhook{
		client: c,
		config: c.config,
	}
}

//WebHooks ...
func (r *webhook) List(name string, target *ClusterTargetHeader) ([]WebHook, error) {
	rawURL := fmt.Sprintf("/v1/clusters/%s/webhooks", name)
	webhooks := []WebHook{}
	_, err := r.client.get(rawURL, &webhooks, target)
	if err != nil {
		return nil, err
	}

	return webhooks, err
}

//AddWebHook ...
func (r *webhook) Add(name string, params WebHook, target *ClusterTargetHeader) error {
	rawURL := fmt.Sprintf("/v1/clusters/%s/webhooks", name)
	_, err := r.client.post(rawURL, params, nil, target)
	return err
}
