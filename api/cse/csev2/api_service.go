package csev2

import (
	gohttp "net/http"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/authentication"
	"github.com/Mavrickk3/bluemix-go/client"
	"github.com/Mavrickk3/bluemix-go/http"
	"github.com/Mavrickk3/bluemix-go/rest"
	"github.com/Mavrickk3/bluemix-go/session"
)

//const ErrCodeAPICreation = "APICreationError"

type CseServiceAPI interface {
	ServiceEndpoints() ServiceEndpoints
}

type cseService struct {
	*client.Client
}

func New(sess *session.Session) (CseServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.CseService)
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
		ep, err := config.EndpointLocator.CseEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &cseService{
		Client: client.New(config, bluemix.CseService, tokenRefreher),
	}, nil
}

func (c *cseService) ServiceEndpoints() ServiceEndpoints {
	return newServiceEndpointsAPI(c.Client)
}
