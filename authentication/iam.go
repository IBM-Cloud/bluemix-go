package authentication

import (
	"encoding/base64"
	"fmt"
	"net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM/go-sdk-core/v5/core"
)

// IAMError ...
type IAMError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	ErrorDetails string `json:"errorDetails"`
}

// Description ...
func (e IAMError) Description() string {
	if e.ErrorDetails != "" {
		return e.ErrorDetails
	}
	return e.ErrorMessage
}

// IAMTokenResponse ...
type IAMTokenResponse struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	UAAAccessToken        string `json:"uaa_token"`
	UAARefreshToken       string `json:"uaa_refresh_token"`
	DelegatedRefreshToken string `json:"delegated_refresh_token"`
	TokenType             string `json:"token_type"`
}

// IAMAuthRepository ...
type IAMAuthRepository struct {
	config   *bluemix.Config
	client   *rest.Client
	endpoint string
}

// NewIAMAuthRepository ...
func NewIAMAuthRepository(config *bluemix.Config, client *rest.Client) (*IAMAuthRepository, error) {
	var endpoint string

	if config.TokenProviderEndpoint != nil {
		endpoint = *config.TokenProviderEndpoint
	} else {
		var err error
		endpoint, err = config.EndpointLocator.IAMEndpoint()
		if err != nil {
			return nil, err
		}
	}

	return &IAMAuthRepository{
		config:   config,
		client:   client,
		endpoint: endpoint,
	}, nil
}

// AuthenticatePassword ...
func (auth *IAMAuthRepository) AuthenticatePassword(username string, password string) error {
	return auth.setTokens("bx", "bx", map[string]string{
		"response_type": "cloud_iam",
		"grant_type":    "password",
		"username":      username,
		"password":      password,
	})
}

// AuthenticateAPIKey ...
func (auth *IAMAuthRepository) AuthenticateAPIKey(apiKey string) error {
	return auth.setTokens("bx", "bx", map[string]string{
		"response_type": "cloud_iam",
		"grant_type":    "urn:ibm:params:oauth:grant-type:apikey",
		"apikey":        apiKey,
	})
}

// AuthenticateSSO ...
func (auth *IAMAuthRepository) AuthenticateSSO(passcode string) error {
	return auth.setTokens("bx", "bx", map[string]string{
		"response_type": "cloud_iam",
		"grant_type":    "urn:ibm:params:oauth:grant-type:passcode",
		"passcode":      passcode,
	})
}

// IAMAssumeAuthenticator ...
func (auth *IAMAuthRepository) AuthenticateAssume(apiKey string, trustedProfileId string) error {
	return auth.setAssumeTokens(apiKey, trustedProfileId)
}

// IAMAssumeAuthenticator ...
func (auth *IAMAuthRepository) FetchAuthorizationData(authenticator core.Authenticator) error {
	req := &http.Request{
		Header: make(http.Header),
	}
	if err := authenticator.Authenticate(req); err != nil {
		return err
	}

	auth.config.IAMAccessToken = req.Header.Get("Authorization")
	return nil
}

// GetKubeTokens fetches the kube:kube access and refresh tokens.
func (auth *IAMAuthRepository) GetKubeTokens() (string, string, error) {
	delegatedTokens, err := auth.getTokens("bx", "bx", map[string]string{
		"grant_type":                     "refresh_token",
		"refresh_token":                  auth.config.IAMRefreshToken,
		"response_type":                  "delegated_refresh_token",
		"delegated_refresh_token_expiry": "600",
		"receiver_client_ids":            "kube",
	})
	if err != nil {
		return "", "", err
	}

	kubeTokens, err := auth.getTokens("kube", "kube", map[string]string{
		"grant_type":    "urn:ibm:params:oauth:grant-type:delegated-refresh-token",
		"refresh_token": delegatedTokens.DelegatedRefreshToken,
	})
	if err != nil {
		return "", "", err
	}

	return kubeTokens.AccessToken, kubeTokens.RefreshToken, nil
}

// RefreshToken ...
func (auth *IAMAuthRepository) RefreshToken() (string, error) {
	data := map[string]string{
		"response_type": "cloud_iam",
		"grant_type":    "refresh_token",
		"refresh_token": auth.config.IAMRefreshToken,
	}

	err := auth.setTokens("bx", "bx", data)
	if err != nil {
		return "", err
	}

	return auth.config.IAMAccessToken, nil
}

// GetPasscode ...
func (auth *IAMAuthRepository) GetPasscode() (string, error) {
	request := rest.PostRequest(auth.endpoint+"/identity/passcode").
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("bx:bx"))).
		Field("grant_type", "refresh_token").
		Field("refresh_token", auth.config.IAMRefreshToken).
		Field("response_type", "cloud_iam")

	res := make(map[string]string, 0)
	var apiErr IAMError

	resp, err := auth.client.Do(request, &res, &apiErr)
	if err != nil {
		return "", err
	}

	if apiErr.ErrorCode != "" {
		if apiErr.ErrorCode == "BXNIM0407E" {
			return "", bmxerror.New(ErrCodeInvalidToken, apiErr.Description())
		}
		return "", bmxerror.NewRequestFailure(apiErr.ErrorCode, apiErr.Description(), resp.StatusCode)
	}

	return res["passcode"], nil
}

func (auth *IAMAuthRepository) getTokens(clientId, clientSecret string, data map[string]string) (IAMTokenResponse, error) {
	request := rest.PostRequest(auth.endpoint+"/identity/token").
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", clientId, clientSecret)))

	for k, v := range data {
		request.Field(k, v)
	}

	var tokens IAMTokenResponse
	var apiErr IAMError

	resp, err := auth.client.Do(request, &tokens, &apiErr)
	if err != nil {
		return tokens, err
	}

	if apiErr.ErrorCode == "" {
		return tokens, nil
	}

	if apiErr.ErrorCode == "BXNIM0407E" {
		if resp != nil && resp.Header != nil {
			return tokens, bmxerror.New(ErrCodeInvalidToken, fmt.Sprintf("Transaction-Id:%s %s", resp.Header["Transaction-Id"], apiErr.Description()))
		}
		return tokens, bmxerror.New(ErrCodeInvalidToken, apiErr.Description())
	}

	if resp != nil && resp.Header != nil {
		return tokens, bmxerror.NewRequestFailure(apiErr.ErrorCode, fmt.Sprintf("Transaction-Id:%s %s", resp.Header["Transaction-Id"], apiErr.Description()), resp.StatusCode)
	}

	return tokens, bmxerror.NewRequestFailure(apiErr.ErrorCode, apiErr.Description(), resp.StatusCode)
}

func (auth *IAMAuthRepository) setTokens(clientId, clientSecret string, data map[string]string) error {
	tokens, err := auth.getTokens(clientId, clientSecret, data)
	if err != nil {
		return err
	}

	auth.config.IAMAccessToken = fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken)
	auth.config.IAMRefreshToken = tokens.RefreshToken

	return nil
}

func (auth *IAMAuthRepository) setAssumeTokens(apiKey, trustedProfileId string) error {
	tokens, err := auth.getAssumeToken(apiKey, trustedProfileId)
	if err != nil {
		return err
	}

	auth.config.IAMAccessToken = fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken)
	return nil
}

// getAssumeTokens
func (auth *IAMAuthRepository) getAssumeToken(apiKey, trustedProfileId string) (IAMTokenResponse, error) {
	delegatedTokens, err := auth.getTokens("bx", "bx", map[string]string{
		"grant_type":    "urn:ibm:params:oauth:grant-type:apikey",
		"apikey":        apiKey,
		"response_type": "cloud_iam",
	})
	if err != nil {
		return delegatedTokens, err
	}
	assumeTokens, err := auth.getTokens("bx", "bx", map[string]string{
		"grant_type":   "urn:ibm:params:oauth:grant-type:assume",
		"access_token": delegatedTokens.AccessToken,
		"profile_id":   trustedProfileId,
	})
	if err != nil {
		return assumeTokens, err
	}
	return assumeTokens, nil
}
