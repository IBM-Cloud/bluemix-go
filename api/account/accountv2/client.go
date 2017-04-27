package accountv2

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net"
	gohttp "net/http"
	"net/url"
	"path"
	"reflect"
	"strings"

	"golang.org/x/net/websocket"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/authentication"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/http"
	"github.com/IBM-Bluemix/bluemix-go/rest"
	"github.com/IBM-Bluemix/bluemix-go/session"
)

//AuthorizationHeader ...
const AuthorizationHeader = "Authorization"

//Client is the accountv2 client ...
type Client interface {
	Accounts() Accounts
}

//ErrCodeNoAccountExists ...
var ErrCodeNoAccountExists = "NoAccountExists"

//PaginatedResources ...
type PaginatedResources struct {
	NextURL        string          `json:"next_url"`
	ResourcesBytes json.RawMessage `json:"resources"`
	resourceType   reflect.Type
}

//NewPaginatedResources ...
func NewPaginatedResources(resource interface{}) PaginatedResources {
	return PaginatedResources{
		resourceType: reflect.TypeOf(resource),
	}
}

//Resources ...
func (pr PaginatedResources) Resources() ([]interface{}, error) {
	slicePtr := reflect.New(reflect.SliceOf(pr.resourceType))
	err := json.Unmarshal([]byte(pr.ResourcesBytes), slicePtr.Interface())
	slice := reflect.Indirect(slicePtr)

	contents := make([]interface{}, 0, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		contents = append(contents, slice.Index(i).Interface())
	}
	return contents, err
}

//URLGetter ...
type URLGetter func() string

//ErrHandler ...
type ErrHandler func(statusCode int, rawResponse []byte) error

//BeforeHandler ...
type BeforeHandler func(*rest.Request) error

//TokenRefresher ...
type TokenRefresher interface {
	RefreshToken() (string, error)
}

//AccountClient ...
type AccountClient struct {
	UAATokenRefresher TokenRefresher
	BaseURL           URLGetter
	OnError           ErrHandler
	Before            BeforeHandler

	config     *bluemix.Config
	HTTPClient *gohttp.Client
}

//NewClient ...
func NewClient(s *session.Session) (*AccountClient, error) {
	config := s.Config.Copy()

	_, err := config.EndpointLocator.AccountManagementEndpoint()
	if err != nil {
		return nil, err
	}
	baseURL := func() string {
		ep, _ := config.EndpointLocator.AccountManagementEndpoint()
		return ep
	}

	httpClient := http.NewHTTPClient(config)

	tokenRefreher, err := authentication.NewUAARepository(config, &rest.Client{
		HTTPClient: httpClient,
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})

	if err != nil {
		return nil, err
	}
	client := &AccountClient{
		BaseURL:           baseURL,
		UAATokenRefresher: tokenRefreher,
		config:            config,
		HTTPClient:        httpClient,
	}
	return client, nil
}

//Accounts API
func (c *AccountClient) Accounts() Accounts {
	return newAccountAPI(c)
}

//SendRequest ...
func (c *AccountClient) sendRequest(r *rest.Request, respV interface{}) (*gohttp.Response, error) {
	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = gohttp.DefaultClient
	}

	restClient := &rest.Client{
		DefaultHeader: http.DefaultHeader(c.config),
		HTTPClient:    httpClient,
	}

	if c.Before != nil {
		err := c.Before(r)
		if err != nil {
			return new(gohttp.Response), err
		}
	}
	//TODO
	resp, err := restClient.Do(r, respV, nil)

	// The response returned by go HTTP client.Do() could be nil if request timeout.
	// For convenience, we ensure that response returned by this method is always not nil.
	if resp == nil {
		return new(gohttp.Response), err
	}

	if err != nil {
		err = WrapNetworkErrors(resp.Request.URL.Host, err)
	}

	// if token is invalid, refresh and try again
	if resp.StatusCode == 401 && c.UAATokenRefresher != nil {
		newToken, uaaErr := c.UAATokenRefresher.RefreshToken()
		switch uaaErr.(type) {
		case nil:
			restClient.DefaultHeader.Set(AuthorizationHeader, newToken)
			resp, err = restClient.Do(r, respV, nil)
		case *bmxerror.InvalidTokenError:
			return resp, bmxerror.NewRequestFailure("InvalidToken", fmt.Sprintf("%v", uaaErr), 401)
		default:
			return resp, fmt.Errorf("Authentication failed, Unable to refresh auth token: %v. Try again later", uaaErr)
		}
	}

	if errResp, ok := err.(bmxerror.RequestFailure); ok && c.OnError != nil {
		err = c.OnError(errResp.StatusCode(), []byte(errResp.Description()))
	}

	return resp, err
}

//WrapNetworkErrors ...
func WrapNetworkErrors(host string, err error) error {
	var innerErr error
	switch typedErr := err.(type) {
	case *url.Error:
		innerErr = typedErr.Err
	case *websocket.DialError:
		innerErr = typedErr.Err
	}

	if innerErr != nil {
		switch typedInnerErr := innerErr.(type) {
		case x509.UnknownAuthorityError:
			return bmxerror.NewInvalidSSLCert(host, "unknown authority")
		case x509.HostnameError:
			return bmxerror.NewInvalidSSLCert(host, "not valid for the requested host")
		case x509.CertificateInvalidError:
			return bmxerror.NewInvalidSSLCert(host, "")
		case *net.OpError:
			if typedInnerErr.Op == "dial" {
				return fmt.Errorf("%s\n%s", err.Error(), "TIP: If you are behind a firewall and require an HTTP proxy, verify the https_proxy environment variable is correctly set. Else, check your network connection.")
			}
		}
	}

	return err
}

//Get ...
func (c *AccountClient) get(path string, respV interface{}) (*gohttp.Response, error) {
	return c.sendRequest(rest.GetRequest(c.url(path)), respV)
}

//Put ...
func (c *AccountClient) put(path string, data interface{}, respV interface{}) (*gohttp.Response, error) {
	return c.sendRequest(rest.PutRequest(c.url(path)).Body(data), respV)
}

//Patch ...
func (c *AccountClient) patch(path string, data interface{}, respV interface{}) (*gohttp.Response, error) {
	return c.sendRequest(rest.PatchRequest(c.url(path)).Body(data), respV)
}

//Post ...
func (c *AccountClient) post(path string, data interface{}, respV interface{}) (*gohttp.Response, error) {
	return c.sendRequest(rest.PostRequest(c.url(path)).Body(data), respV)
}

//Delete ...
func (c *AccountClient) delete(path string) (*gohttp.Response, error) {
	return c.sendRequest(rest.DeleteRequest(c.url(path)), nil)
}

//GetPaginated ...
func (c *AccountClient) getPaginated(path string, resource interface{}, cb func(interface{}) bool) (resp *gohttp.Response, err error) {
	for path != "" {
		paginatedResources := NewPaginatedResources(resource)

		resp, err = c.get(path, &paginatedResources)
		if err != nil {
			return
		}

		var resources []interface{}
		resources, err = paginatedResources.Resources()
		if err != nil {
			err = fmt.Errorf("%s: Error parsing JSON", err.Error())
			return
		}

		for _, resource := range resources {
			if !cb(resource) {
				return
			}
		}

		path = paginatedResources.NextURL
	}
	return
}

func (c *AccountClient) url(path string) string {
	if c.BaseURL == nil {
		return path
	}

	return c.BaseURL() + cleanPath(path)
}

//SetHTTPClient ...
func (c *AccountClient) SetHTTPClient(httpClient *gohttp.Client) {
	c.HTTPClient = httpClient
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
