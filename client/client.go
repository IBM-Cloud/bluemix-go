//Package client provides a generic client to be used by all services
package client

import (
	"fmt"
	"log"
	"path"
	"strings"

	gohttp "net/http"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/http"

	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//TokenProvider ...
type TokenProvider interface {
	RefreshToken() (string, error)
	AuthenticatePassword(string, string) error
	AuthenticateAPIKey(string) error
}

//HandlePagination ...
type HandlePagination func(c *Client, path string, resource interface{}, cb func(interface{}) bool) (resp *gohttp.Response, err error)

//Client is the base client for all service api client
type Client struct {
	Config           *bluemix.Config
	DefaultHeader    gohttp.Header
	ServiceName      bluemix.ServiceName
	TokenRefresher   TokenProvider
	HandlePagination HandlePagination
}

//Config stores any generic service client configurations
type Config struct {
	Config   *bluemix.Config
	Endpoint string
}

//New ...
func New(c *bluemix.Config, serviceName bluemix.ServiceName, refresher TokenProvider, pagination HandlePagination) *Client {
	config := c.Copy()
	return &Client{
		Config:           config,
		ServiceName:      serviceName,
		TokenRefresher:   refresher,
		HandlePagination: pagination,
		DefaultHeader:    getDefaultAuthHeaders(serviceName, c),
	}
}

//SendRequest ...
func (c *Client) SendRequest(r *rest.Request, respV interface{}) (*gohttp.Response, error) {
	httpClient := c.Config.HTTPClient
	if httpClient == nil {
		httpClient = gohttp.DefaultClient
	}

	restClient := &rest.Client{
		DefaultHeader: c.DefaultHeader,
		HTTPClient:    httpClient,
	}

	resp, err := restClient.Do(r, respV, nil)

	// The response returned by go HTTP client.Do() could be nil if request timeout.
	// For convenience, we ensure that response returned by this method is always not nil.
	if resp == nil {
		return new(gohttp.Response), err
	}

	if err != nil {
		err = bmxerror.WrapNetworkErrors(resp.Request.URL.Host, err)
	}

	// if token is invalid, refresh and try again
	if resp.StatusCode == 401 && c.TokenRefresher != nil {
		_, err := c.TokenRefresher.RefreshToken()
		switch err.(type) {
		case nil:
			restClient.DefaultHeader = getDefaultAuthHeaders(c.ServiceName, c.Config)
			resp, err = restClient.Do(r, respV, nil)
		case *bmxerror.InvalidTokenError:
			return resp, bmxerror.NewRequestFailure("InvalidToken", fmt.Sprintf("%v", err), 401)
		default:
			return resp, fmt.Errorf("Authentication failed, Unable to refresh auth token: %v. Try again later", err)
		}
	}

	return resp, err
}

//Get ...
func (c *Client) Get(path string, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.GetRequest(c.URL(path))
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.SendRequest(r, respV)
}

//Put ...
func (c *Client) Put(path string, data interface{}, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PutRequest(c.URL(path)).Body(data)
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.SendRequest(r, respV)
}

//Patch ...
func (c *Client) Patch(path string, data interface{}, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PatchRequest(c.URL(path)).Body(data)
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.SendRequest(r, respV)
}

//Post ...
func (c *Client) Post(path string, data interface{}, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PostRequest(c.URL(path)).Body(data)
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.SendRequest(r, respV)
}

//Delete ...
func (c *Client) Delete(path string, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.DeleteRequest(c.URL(path))
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.SendRequest(r, nil)
}

func addToRequestHeader(h interface{}, r *rest.Request) {
	switch v := h.(type) {
	case map[string]string:
		for key, value := range v {
			r.Set(key, value)
		}
	}
}

//GetPaginated ...
func (c *Client) GetPaginated(path string, resource interface{}, cb func(interface{}) bool) (resp *gohttp.Response, err error) {
	return c.HandlePagination(c, path, resource, cb)
}

//URL ...
func (c *Client) URL(path string) string {
	return *c.Config.Endpoint + cleanPath(path)
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return path.Clean(p)
}

const (
	userAgentHeader      = "User-Agent"
	authorizationHeader  = "Authorization"
	uaaAccessTokenHeader = "X-Auth-Uaa-Token"

	iamRefreshTokenHeader = "X-Auth-Refresh-Token"
)

func getDefaultAuthHeaders(serviceName bluemix.ServiceName, c *bluemix.Config) gohttp.Header {
	h := gohttp.Header{}
	switch serviceName {
	case bluemix.CfService, bluemix.AccountService:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, c.UAAAccessToken)

	case bluemix.ClusterService:
		h.Set(userAgentHeader, http.UserAgent())
		h.Set(authorizationHeader, c.IAMAccessToken)
		h.Set(iamRefreshTokenHeader, c.IAMRefreshToken)
		h.Set(uaaAccessTokenHeader, c.UAAAccessToken)
	default:
		log.Println("Unknown service")
	}
	return h
}
