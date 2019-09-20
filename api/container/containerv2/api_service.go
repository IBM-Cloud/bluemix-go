package containerv2

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//ContainerServiceAPI is the Aramda K8s client ...
type ContainerServiceAPI interface {
	Clusters() Clusters
	//TODO Add other services
}

//ContainerService holds the client
type csService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (ContainerServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.ContainerService)
	if err != nil {
		return nil, err
	}
	if config.HTTPClient == nil {
		config.HTTPClient = http.NewHTTPClient(config)
	}
	tokenRefreher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
		HTTPClient: config.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	if config.IAMAccessToken == "" {
		err := authentication.PopulateTokens(tokenRefreher, config)
		if err != nil {
			return nil, err
		}
	}
	if config.Endpoint == nil {
		ep, err := config.EndpointLocator.ContainerEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &csService{
		Client: client.New(config, bluemix.ContainerService, tokenRefreher),
	}, nil
}

//Clusters implements Clusters API
func (c *csService) Clusters() Clusters {
	return newClusterAPI(c.Client)
}
