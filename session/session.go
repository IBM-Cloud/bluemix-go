package session

import (
	"fmt"
	"time"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/trace"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/common/rest"
	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/authentication"
	"github.com/IBM-Bluemix/bluemix-go/endpoints"
	"github.com/IBM-Bluemix/bluemix-go/helpers"
	"github.com/IBM-Bluemix/bluemix-go/http"
)

//Session ...
type Session struct {
	Config *bluemix.Config
}

//New ...
func New(configs ...*bluemix.Config) (*Session, error) {
	var c *bluemix.Config

	if len(configs) == 0 {
		c = &bluemix.Config{}
	} else {
		c = configs[0]
	}
	sess := &Session{
		Config: c,
	}

	if len(c.IBMID) == 0 {
		c.IBMID = helpers.EnvFallBack([]string{"IBMID"}, "")
	}
	if len(c.IBMIDPassword) == 0 {
		c.IBMIDPassword = helpers.EnvFallBack([]string{"IBMID_PASSWORD"}, "")
	}
	if len(c.Region) == 0 {
		c.Region = helpers.EnvFallBack([]string{"BM_REGION", "BLUEMIX_REGION"}, "us-south")
	}
	if c.MaxRetries == nil {
		c.MaxRetries = helpers.Int(3)
	}
	if c.HTTPTimeout == 0 {
		c.HTTPTimeout = 180 * time.Second
		timeout := helpers.EnvFallBack([]string{"BM_TIMEOUT", "BLUEMIX_TIMEOUT"}, "180")
		timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%ss", timeout))
		if err != nil {
			fmt.Printf("BM_TIMEOUT or BLUEMIX_TIMEOUT has invalid time format. Default timeout will be set to %q", c.HTTPTimeout)
		}
		if err == nil {
			c.HTTPTimeout = timeoutDuration
		}
	}

	if c.RetryDelay == nil {
		c.RetryDelay = helpers.Duration(30 * time.Second)
	}
	if c.EndpointLocator == nil {
		c.EndpointLocator = endpoints.NewEndpointLocator(c.Region)
	}

	if len(c.IAMAccessToken) == 0 {
		c.IAMAccessToken = helpers.EnvFallBack([]string{"IBMCLOUD_IAM_TOKEN"}, "")
	}
	if len(c.IAMRefreshToken) == 0 {
		c.IAMRefreshToken = helpers.EnvFallBack([]string{"IBMCLOUD_IAM_REFRESH_TOKEN"}, "")
	}
	if len(c.UAAAccessToken) == 0 {
		c.UAAAccessToken = helpers.EnvFallBack([]string{"IBMCLOUD_UAA_TOKEN"}, "")
	}
	if len(c.UAARefreshToken) == 0 {
		c.UAARefreshToken = helpers.EnvFallBack([]string{"IBMCLOUD_UAA_REFRESH_TOKEN"}, "")
	}

	if c.Debug {
		trace.Logger = trace.NewLogger("true")
	}

	err := c.ValidateConfig()

	if err != nil {
		return sess, err
	}

	if c.UAAAccessToken == "" || c.UAARefreshToken == "" {
		restClient := rest.NewClient()
		restClient.HTTPClient = http.NewHTTPClient(c)
		iam, err := authentication.NewIAMAuthRepository(c, restClient)
		if err != nil {
			return sess, err
		}
		err = iam.AuthenticatePassword(c.IBMID, c.IBMIDPassword)
		if err != nil {
			return sess, err
		}
	}

	return sess, nil
}
