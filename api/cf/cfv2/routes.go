package cfv2

import (
	"fmt"
	//"strconv"

	//"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//ErrCodeRouteDoesnotExist ...
var ErrCodeRouteDoesnotExist = "RouteDoesnotExist"

type RouteCreateRequest struct {
	Host       string `json:"host"`
	SpaceGUID  string `json:"space_guid"`
	DomainGUID string `json:"domain_guid"`
}

//Metadata ...
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

//Apps ...
type Routes interface {
	GetSharedDomains(domain string) (*Route, error)
	Find(hostname, domainGUID string) ([]Route, error)
	Create(host, domainGUID, spaceGUID string) (*RouteFields, error)
}

type route struct {
	client *client.Client
}

func newRouteAPI(c *client.Client) Routes {
	return &route{
		client: c,
	}
}

func (r *route) GetSharedDomains(domainName string) (*Route, error) {
	rawURL := "/v2/shared_domains"
	req := rest.GetRequest(rawURL).Query("q", "name:"+domainName)
	fmt.Println(req, "\n\n")
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	domain, err := r.listRouteWithPath(path)
	if err != nil {
		return nil, err
	}
	return &domain[0], nil
}

func (r *route) Find(hostname, domainGUID string) ([]Route, error) {
	rawURL := "/v2/routes?inline-relations-depth=1"
	req := rest.GetRequest(rawURL).Query("q", "host:"+hostname+";domain_guid:"+domainGUID)
	fmt.Println(req, "\n\n")
	httpReq, err := req.Build()
	fmt.Println(httpReq, "\n\n", "\n\n")
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	route, err := r.listRouteWithPath(path)
	if err != nil {
		return nil, err
	}
	return route, nil
}

func (r *route) Create(host, domainGUID, spaceGUID string) (*RouteFields, error) {
	payload := RouteCreateRequest{
		Host:       host,
		DomainGUID: domainGUID,
		SpaceGUID:  spaceGUID,
	}
	rawURL := "/v2/routes?async=true&inline-relations-depth=1"
	routeFields := RouteFields{}
	_, err := r.client.Post(rawURL, payload, &routeFields)
	if err != nil {
		return nil, err
	}
	return &routeFields, nil
}

func (r *route) listRouteWithPath(path string) ([]Route, error) {
	var route []Route
	_, err := r.client.GetPaginated(path, RouteResource{}, func(resource interface{}) bool {
		if routeResource, ok := resource.(RouteResource); ok {
			route = append(route, routeResource.ToFields())
			return true
		}
		return false
	})
	return route, err
}
