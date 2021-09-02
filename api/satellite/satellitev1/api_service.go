package satellitev1

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

//SatelliteServiceAPI is the Aramda K8s client ...
type SatelliteServiceAPI interface {
	Endpoint() Endpoint

	//TODO Add other services
}

type satService struct {
	*client.Client
}

func New(sess *session.Session) (SatelliteServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.VpcContainerService)
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
		ep, err := config.EndpointLocator.SatelliteEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &satService{
		Client: client.New(config, bluemix.VpcContainerService, tokenRefreher),
	}, nil
}

//Clusters implements Clusters API
func (c *satService) Endpoint() Endpoint {
	return newEndpointAPI(c.Client)
}
