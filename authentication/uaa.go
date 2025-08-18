package authentication

import (
	"encoding/base64"
	"fmt"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

// UAAError ...
type UAAError struct {
	ErrorCode   string `json:"error"`
	Description string `json:"error_description"`
}

// UAATokenResponse ...
type UAATokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

// UAARepository ...
type UAARepository struct {
	config   *bluemix.Config
	client   *rest.Client
	endpoint string
}

// NewUAARepository ...
func NewUAARepository(config *bluemix.Config, client *rest.Client) (*UAARepository, error) {
	var endpoint string

	if config.TokenProviderEndpoint != nil {
		endpoint = *config.TokenProviderEndpoint
	} else {
		var err error
		endpoint, err = config.EndpointLocator.UAAEndpoint()
		if err != nil {
			return nil, err
		}
	}
	return &UAARepository{
		config:   config,
		client:   client,
		endpoint: endpoint,
	}, nil
}

// AuthenticatePassword ...
func (auth *UAARepository) AuthenticatePassword(username string, password string) error {
	return auth.setTokens("cf", "", map[string]string{
		"grant_type": "password",
		"username":   username,
		"password":   password,
	})
}

// AuthenticateSSO ...
func (auth *UAARepository) AuthenticateSSO(passcode string) error {
	return auth.setTokens("cf", "", map[string]string{
		"grant_type": "password",
		"passcode":   passcode,
	})
}

// AuthenticateAPIKey ...
func (auth *UAARepository) AuthenticateAPIKey(apiKey string) error {
	return auth.AuthenticatePassword("apikey", apiKey)
}

// GetKubeTokens fetches the kube:kube access and refresh tokens.
func (auth *UAARepository) GetKubeTokens() (string, string, error) {
	data := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": auth.config.IAMRefreshToken,
	}

	tokens, err := auth.getTokens("kube", "kube", data)
	if err != nil {
		return "", "", err
	}

	return tokens.AccessToken, tokens.RefreshToken, nil
}

// RefreshToken ...
func (auth *UAARepository) RefreshToken() (string, error) {
	err := auth.setTokens("cf", "", map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": auth.config.UAARefreshToken,
	})
	if err != nil {
		return "", err
	}

	return auth.config.UAAAccessToken, nil
}

// GetPasscode ...
func (auth *UAARepository) GetPasscode() (string, error) {
	return "", nil
}

func (auth *UAARepository) getTokens(clientId, clientSecret string, data map[string]string) (UAATokenResponse, error) {
	request := rest.PostRequest(auth.endpoint+"/oauth/token").
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", clientId, clientSecret))).
		Field("scope", "")

	for k, v := range data {
		request.Field(k, v)
	}

	var tokens UAATokenResponse
	var apiErr UAAError

	resp, err := auth.client.Do(request, &tokens, &apiErr)
	if err != nil {
		return tokens, err
	}

	if apiErr.ErrorCode == "" {
		return tokens, nil
	}

	if apiErr.ErrorCode == "invalid-token" {
		return tokens, bmxerror.NewInvalidTokenError(apiErr.Description)
	}
	return tokens, bmxerror.NewRequestFailure(apiErr.ErrorCode, apiErr.Description, resp.StatusCode)
}

func (auth *UAARepository) setTokens(clientId, clientSecret string, data map[string]string) error {
	tokens, err := auth.getTokens(clientId, clientSecret, data)
	if err != nil {
		return err
	}

	auth.config.UAAAccessToken = fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken)
	auth.config.UAARefreshToken = tokens.RefreshToken
	return nil
}
