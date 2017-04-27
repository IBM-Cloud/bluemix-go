package bluemix

import (
	"net/http"
	"time"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/endpoints"
)

//Config ...
type Config struct {
	IBMID         string
	IBMIDPassword string

	IAMAccessToken  string
	IAMRefreshToken string
	UAAAccessToken  string
	UAARefreshToken string

	Region          string
	EndpointLocator endpoints.EndpointLocator
	MaxRetries      *int
	RetryDelay      *time.Duration

	HTTPTimeout time.Duration

	Debug bool

	HTTPClient *http.Client

	SSLDisable bool
}

//Copy performs a shallow copy of the config
func (c *Config) Copy() *Config {
	out := new(Config)
	*out = *c
	return out
}

//ValidateConfig ...
func (c *Config) ValidateConfig() error {
	if (c.IBMID == "" || c.IBMIDPassword == "") &&
		(c.UAAAccessToken == "" || c.UAARefreshToken == "") {
		return bmxerror.New(ErrInvalidConfigurationCode, "Either IBMID and IBMID_PASSWORD or IBMCLOUD_UAA_TOKEN and IBMCLOUD_UAA_REFRESH_TOKEN should be provided in the environment variable or they should be set in the Config")
	}
	if c.Region == "" {
		return bmxerror.New(ErrMissingRegionCode, "Either BM_REGION or  BLUEMIX_REGION should be provided in the environment variable or they should be set in the Config")
	}
	return nil
}
