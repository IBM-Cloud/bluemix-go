package session

import (
	"net/http"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/endpoints"
	"github.com/IBM-Bluemix/bluemix-go/trace"
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

	if c.Region != "" && c.EndpointLocator == nil {
		c.EndpointLocator = endpoints.NewEndpointLocator(c.Region)
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Timeout: c.HTTPTimeout,
		}
	}

	if c.Debug {
		trace.Logger = trace.NewLogger("true")
	}

	err := c.Validate()
	return sess, err
}
