package accountv2

import (
	gohttp "net/http"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/authentication"
	"github.com/Mavrickk3/bluemix-go/client"
	"github.com/Mavrickk3/bluemix-go/http"
	"github.com/Mavrickk3/bluemix-go/rest"
	"github.com/Mavrickk3/bluemix-go/session"
)

// AccountServiceAPI is the accountv2 client ...
type AccountServiceAPI interface {
	Accounts() Accounts
}

// ErrCodeNoAccountExists ...
const ErrCodeNoAccountExists = "NoAccountExists"

// MccpService holds the client
type accountService struct {
	*client.Client
}

// New ...
func New(sess *session.Session) (AccountServiceAPI, error) {
	config := sess.Config.Copy()
	err := config.ValidateConfigForService(bluemix.AccountService)
	if err != nil {
		return nil, err
	}
	if config.HTTPClient == nil {
		config.HTTPClient = http.NewHTTPClient(config)
	}
	tokenRefreher, err := authentication.NewUAARepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"X-Original-User-Agent": []string{config.UserAgent},
			"User-Agent":            []string{http.UserAgent()},
		},
		HTTPClient: config.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	if config.UAAAccessToken == "" {
		err := authentication.PopulateTokens(tokenRefreher, config)
		if err != nil {
			return nil, err
		}
	}
	if config.Endpoint == nil {
		ep, err := config.EndpointLocator.AccountManagementEndpoint()
		if err != nil {
			return nil, err
		}
		config.Endpoint = &ep
	}
	return &accountService{
		Client: client.New(config, bluemix.AccountService, tokenRefreher),
	}, nil
}

// Accounts API
func (a *accountService) Accounts() Accounts {
	return newAccountAPI(a.Client)
}
