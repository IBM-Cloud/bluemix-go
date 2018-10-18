package registryv1

import (
	gohttp "net/http"

	ibmcloud "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
)

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//RegistryServiceAPI is the IBM Cloud Registry client ...
type RegistryServiceAPI interface {
	Builds() Builds
	/*Auth() Auth

	Images() Images
	Messages() Messages
	Namespaces() Namespaces
	Plans() Plans
	Quotas() Quotas
	Tokens() Tokens*/
}

//RegistryService holds the client
type rsService struct {
	*client.Client
}

//New ...
func New(sess *session.Session) (RegistryServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(ibmcloud.ContainerRegistryService)
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
		ep, err := config.EndpointLocator.ContainerRegistryEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &rsService{
		Client: client.New(config, ibmcloud.ContainerRegistryService, tokenRefreher),
	}, nil
}



//Builds implements builds API
func (c *rsService) Builds() Builds {
	return newBuildAPI(c.Client)
}
/*
//Auth implement auth API
func (c *csService) Auth() Auth {
	return newAuthAPI(c.Client)
}

//Images implements Images API
func (c *csService) Images() Images {
	return newImageAPI(c.Client)
}

//Messages implements Messages API
func (c *csService) Messages() Messages {
	return newMessageAPI(c.Client)
}

//Namespaces implements Namespaces API
func (c *csService) Namespaces() Namespaces {
	return newNamespaceAPI(c.Client)
}

//Plans implements Plans API
func (c *csService) Plans() Plans {
	return newPlanAPI(c.Client)
}

//Quotas implements Quotas API
func (c *csService) Quotas() Quotas {
	return newQuotaAPI(c.Client)
}

//Tokens implements Tokens API
func (c *csService) Tokens() Tokens {
	return newTokenAPI(c.Client)
}
*/