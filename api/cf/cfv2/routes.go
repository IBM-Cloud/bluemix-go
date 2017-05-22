package cfv2

import (

	//"strconv"

	//"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//ErrCodeRouteDoesnotExist ...
var ErrCodeRouteDoesnotExist = "RouteDoesnotExist"

//RouteRequest ...
type RouteRequest struct {
	Host       string `json:"host,omitempty"`
	SpaceGUID  string `json:"space_guid"`
	DomainGUID string `json:"domain_guid,omitempty"`
	Path       string `json:"path,omitempty"`
	Port       int    `json:"port,omitempty"`
}

//RouteMetadata ...
type RouteMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//RouteEntity ...
type RouteEntity struct {
	Name string `json:"name"`
}

//RouteResource ...
type RouteResource struct {
	Resource
	Entity RouteEntity
}

//RouteFields ...
type RouteFields struct {
	Metadata RouteMetadata
	Entity   RouteEntity
}

//ToFields ..
func (resource RouteResource) ToFields() Route {
	entity := resource.Entity

	return Route{
		GUID: resource.Metadata.GUID,
		Name: entity.Name,
	}
}

//Route model
type Route struct {
	GUID string
	Name string
}

//Routes ...
type Routes interface {
	GetSharedDomains(domain string) (*Route, error)
	Find(hostname, domainGUID string) ([]Route, error)
	Create(req RouteRequest) (*RouteFields, error)
	Get(routeGUID string) (*RouteFields, error)
	Update(routeGUID string, req RouteRequest) (*RouteFields, error)
	Delete(routeGUID string, async bool) error
}

type route struct {
	client *client.Client
}

func newRouteAPI(c *client.Client) Routes {
	return &route{
		client: c,
	}
}

func (r *route) Get(routeGUID string) (*RouteFields, error) {
	rawURL := fmt.Sprintf("/v2/routes/%s", routeGUID)
	routeFields := RouteFields{}
	_, err := r.client.Get(rawURL, &routeFields, nil)
	if err != nil {
		return nil, err
	}
	return &routeFields, nil
}

func (r *route) GetSharedDomains(domainName string) (*Route, error) {
	rawURL := "/v2/shared_domains"
	req := rest.GetRequest(rawURL).Query("q", "name:"+domainName)
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	domain, err := listRouteWithPath(r.client, path)
	if err != nil {
		return nil, err
	}
	return &domain[0], nil
}

func (r *route) Find(hostname, domainGUID string) ([]Route, error) {
	rawURL := "/v2/routes?inline-relations-depth=1"
	req := rest.GetRequest(rawURL).Query("q", "host:"+hostname+";domain_guid:"+domainGUID)
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	route, err := listRouteWithPath(r.client, path)
	if err != nil {
		return nil, err
	}
	return route, nil
}

func (r *route) Create(req RouteRequest) (*RouteFields, error) {
	rawURL := "/v2/routes?async=true&inline-relations-depth=1"
	routeFields := RouteFields{}
	_, err := r.client.Post(rawURL, req, &routeFields)
	if err != nil {
		return nil, err
	}
	return &routeFields, nil
}

func (r *route) Update(routeGUID string, req RouteRequest) (*RouteFields, error) {
	rawURL := fmt.Sprintf("/v2/routes/%s", routeGUID)
	routeFields := RouteFields{}
	_, err := r.client.Put(rawURL, req, &routeFields)
	if err != nil {
		return nil, err
	}
	return &routeFields, nil
}

func (r *route) Delete(routeGUID string, async bool) error {
	rawURL := fmt.Sprintf("/v2/route/%s", routeGUID)
	req := rest.GetRequest(rawURL).Query("recursive", "true")
	if async {
		req.Query("async", "true")
	}
	httpReq, err := req.Build()
	if err != nil {
		return err
	}
	path := httpReq.URL.String()
	_, err = r.client.Delete(path)
	return err
}

func listRouteWithPath(c *client.Client, path string) ([]Route, error) {
	var route []Route
	_, err := c.GetPaginated(path, RouteResource{}, func(resource interface{}) bool {
		if routeResource, ok := resource.(RouteResource); ok {
			route = append(route, routeResource.ToFields())
			return true
		}
		return false
	})
	return route, err
}
