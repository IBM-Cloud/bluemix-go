package catalog

import (
	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
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
	err := config.ValidateConfigForService(bluemix.ResourceCatalogrService)
	if err != nil {
		return nil, err
	}
	if config.HTTPClient == nil {
		config.HTTPClient = http.NewHTTPClient(config)
	}
	tokenRefreher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"X-Original-User-Agent": []string{config.UserAgent},
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
