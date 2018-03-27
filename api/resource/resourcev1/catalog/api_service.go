package catalog

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/authentication"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/http"
	"github.com/IBM-Bluemix/bluemix-go/rest"
	"github.com/IBM-Bluemix/bluemix-go/session"
)

//ResourceCatalogAPI is the resource client ...
type ResourceCatalogAPI interface {
	ResourceCatalog() ResourceCatalogRepository
}

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//resourceControllerService holds the client
type resourceControllerService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (ResourceCatalogAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.IAMPAPService)
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
		ep, err := config.EndpointLocator.ResourceCatalogEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}
	return &resourceControllerService{
		Client: client.New(config, bluemix.ResourceControllerService, tokenRefreher),
	}, nil
}

//ResourceCatalog API
func (a *resourceControllerService) ResourceCatalog() ResourceCatalogRepository {
	return newResourceCatalogAPI(a.Client)
}
