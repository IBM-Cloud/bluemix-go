package satellitev1

import (
	gohttp "net/http"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/authentication"
	"github.com/Mavrickk3/bluemix-go/client"
	"github.com/Mavrickk3/bluemix-go/http"
	"github.com/Mavrickk3/bluemix-go/rest"
	"github.com/Mavrickk3/bluemix-go/session"
)

// ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

// SatelliteServiceAPI is the Aramda K8s client ...
type SatelliteServiceAPI interface {
	Endpoint() Endpoint
	Source() Source

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
			"X-Original-User-Agent": []string{config.UserAgent},
			"User-Agent":            []string{http.UserAgent()},
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

// Endpoint implements Endpoint API
func (c *satService) Endpoint() Endpoint {
	return newEndpointAPI(c.Client)
}

// Source implements Source API
func (c *satService) Source() Source {
	return newSourceAPI(c.Client)
}
