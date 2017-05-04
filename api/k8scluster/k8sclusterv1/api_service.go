package k8sclusterv1

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/authentication"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/http"
	"github.com/IBM-Bluemix/bluemix-go/rest"
	"github.com/IBM-Bluemix/bluemix-go/session"
)

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//ClusterServiceAPI is the Aramda K8s client ...
type ClusterServiceAPI interface {
	Clusters() Clusters
	Workers() Workers
	WebHooks() Webhooks
	Subnets() Subnets
}

//ClusterService holds the client
type csService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (ClusterServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.ClusterService)
	if err != nil {
		return nil, err
	}
	tokenRefreher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if config.IAMAccessToken == "" {
		authentication.PopulateTokens(tokenRefreher, config)
	}

	if err != nil {
		return nil, err
	}

	if config.HTTPClient == nil {
		config.HTTPClient = http.NewHTTPClient(config)
	}
	if config.Endpoint == nil {
		ep, err := config.EndpointLocator.ClusterEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &csService{
		Client: client.New(config, bluemix.ClusterService, tokenRefreher, nil),
	}, nil
}

//Clusters implements Clusters API
func (c *csService) Clusters() Clusters {
	return newClusterAPI(c.Client)
}

//Workers implements Cluster Workers API
func (c *csService) Workers() Workers {
	return newWorkerAPI(c.Client)
}

//Subnets implements Cluster Subnets API
func (c *csService) Subnets() Subnets {
	return newSubnetAPI(c.Client)
}

//Webhooks implements Cluster WebHooks API
func (c *csService) WebHooks() Webhooks {
	return newWebhookAPI(c.Client)
}
