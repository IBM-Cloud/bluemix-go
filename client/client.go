//Package client provides a generic client to be used by all services
package client

import (
	"fmt"
	"path"
	"strings"

	gohttp "net/http"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"

	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//EndpointOptions ...
type EndpointOptions struct {
	APIEndpoint           string
	TokenProviderEndpoint string
}

//TokenRefresher ...
type TokenRefresher interface {
	RefreshToken() (string, error)
}

//HandlePagination ...
type HandlePagination func(c *Client, path string, resource interface{}, cb func(interface{}) bool) (resp *gohttp.Response, err error)

//DefaultHeader ...
type DefaultHeader interface {
	DefaultHeader() gohttp.Header
}

//Client is the base client for all service api client
type Client struct {
	Config           *bluemix.Config
	DefaultHeader    func() gohttp.Header
	TokenRefresher   TokenRefresher
	HandlePagination HandlePagination
	Endpoint         string
}

//New ...
func New(c *bluemix.Config, refresher TokenRefresher, pagination HandlePagination, defaultHeader func() gohttp.Header, endpoint string) *Client {
	return &Client{
		Config:           c,
		TokenRefresher:   refresher,
		HandlePagination: pagination,
		DefaultHeader:    defaultHeader,
		Endpoint:         endpoint,
	}
}

func (c *Client) sendRequest(r *rest.Request, respV interface{}) (*gohttp.Response, error) {
	httpClient := c.Config.HTTPClient
	if httpClient == nil {
		httpClient = gohttp.DefaultClient
	}

	restClient := &rest.Client{
		DefaultHeader: c.DefaultHeader(),
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
			restClient.DefaultHeader = c.DefaultHeader()
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
	r := rest.GetRequest(c.url(path))
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.sendRequest(r, respV)
}

//Put ...
func (c *Client) Put(path string, data interface{}, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PutRequest(c.url(path)).Body(data)
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.sendRequest(r, respV)
}

//Patch ...
func (c *Client) Patch(path string, data interface{}, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PatchRequest(c.url(path)).Body(data)
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.sendRequest(r, respV)
}

//Post ...
func (c *Client) Post(path string, data interface{}, respV interface{}, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.PostRequest(c.url(path)).Body(data)
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.sendRequest(r, respV)
}

//Delete ...
func (c *Client) Delete(path string, extraHeader ...interface{}) (*gohttp.Response, error) {
	r := rest.DeleteRequest(c.url(path))
	for _, t := range extraHeader {
		addToRequestHeader(t, r)
	}
	return c.sendRequest(r, nil)
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

func (c *Client) url(path string) string {
	return c.Endpoint + cleanPath(path)
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
