package cfv2

import (
	gohttp "net/http"

	"github.com/IBM-Bluemix/bluemix-go/authentication"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/http"
	"github.com/IBM-Bluemix/bluemix-go/rest"
	"github.com/IBM-Bluemix/bluemix-go/session"
)

//ErrCodeAPICreation ...
const ErrCodeAPICreation = "APICreationError"

//CfServiceAPI is the mccpv2 client ...
type CfServiceAPI interface {
	Organizations() Organizations
	Spaces() Spaces
	ServiceInstances() ServiceInstances
	ServiceKeys() ServiceKeys
	ServicePlans() ServicePlans
	ServiceOfferings() ServiceOfferings
	SpaceQuotas() SpaceQuotas
}

//CfService holds the client
type cfService struct {
	*client.Client
}

//New ...
func New(sess *session.Session, endpoints ...client.EndpointOptions) (CfServiceAPI, error) {
	config := sess.Config

	tokenRefreher, err := authentication.NewUAARepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		return nil, err
	}
	if config.UAAAccessToken == nil || config.UAARefreshToken == nil {
		err := authentication.PopulateTokens(tokenRefreher, config)
		if err != nil {
			return nil, err
		}
	}
	defaultHeader := func() gohttp.Header {
		h := gohttp.Header{}
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, *config.UAAAccessToken)
		return h
	}

	if config.HTTPClient == nil {
		config.HTTPClient = http.NewHTTPClient(config)
	}
	if config.Endpoint == nil {
		ep, err := config.EndpointLocator.CFAPIEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}

	return &cfService{
		Client: client.New(config, tokenRefreher, Paginate, defaultHeader),
	}, nil
}

//Organizations implements Organizations APIs
func (c *cfService) Organizations() Organizations {
	return newOrganizationAPI(c.Client)
}

//Spaces implements Spaces APIs
func (c *cfService) Spaces() Spaces {
	return newSpacesAPI(c.Client)
}

//ServicePlans implements ServicePlans APIs
func (c *cfService) ServicePlans() ServicePlans {
	return newServicePlanAPI(c.Client)
}

//ServiceOfferings implements ServiceOfferings APIs
func (c *cfService) ServiceOfferings() ServiceOfferings {
	return newServiceOfferingAPI(c.Client)
}

//ServiceInstances implements ServiceInstances APIs
func (c *cfService) ServiceInstances() ServiceInstances {
	return newServiceInstanceAPI(c.Client)
}

//ServiceKeys implements ServiceKey APIs
func (c *cfService) ServiceKeys() ServiceKeys {
	return newServiceKeyAPI(c.Client)
}

//SpaceQuotas implements SpaceQuota APIs
func (c *cfService) SpaceQuotas() SpaceQuotas {
	return newSpaceQuotasAPI(c.Client)
}
