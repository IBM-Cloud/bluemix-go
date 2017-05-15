package bluemix

import (
	"net/http"
	"time"

	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/endpoints"
	"github.com/IBM-Bluemix/bluemix-go/helpers"
)

//Config ...
type Config struct {
	IBMID         string
	IBMIDPassword string

	BluemixAPIKey string
	Region        string

	IAMAccessToken  *string
	IAMRefreshToken *string
	UAAAccessToken  *string
	UAARefreshToken *string

	EndpointLocator endpoints.EndpointLocator
	MaxRetries      *int
	RetryDelay      *time.Duration

	HTTPTimeout time.Duration
	Debug       bool

	HTTPClient *http.Client
	SSLDisable bool
}

//Regions ...
var Regions = map[string]struct{}{
	"us-south": struct{}{},
	"au-syd":   struct{}{},
	"eu-gb":    struct{}{},
}

//Validate  ...
func (c *Config) Validate() error {
	if (c.IBMID == "" || c.IBMIDPassword == "") && c.BluemixAPIKey == "" {
		return bmxerror.New(ErrCodeMissingCredentials, "Either IBMID/IBMIDPassword or  BluemixAPIKey must be provided")
	}

	if c.Region != "" {
		if _, ok := Regions[c.Region]; !ok {
			return bmxerror.New(ErrCodeInvalidRegion, fmt.Sprintf("Region is optional. Given region %q is invalid. Supported values are %q", c.Region, helpers.MapToKeys(Regions)))
		}
	}
	return nil
}
