package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM-Cloud/bluemix-go"
	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/api/satellite/satellitev1"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

const sourceName = "portworx-service"

func main() {

	c := new(bluemix.Config)

	var locationID string
	flag.StringVar(&locationID, "locationID", "", "Location ID")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	if locationID == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}
	target := v2.ClusterTargetHeader{}

	satelliteClient, err := satellitev1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	satSourceAPI := satelliteClient.Source()
	satEndpointAPI := satelliteClient.Endpoint()

	// 1. Create Endpoint Source.
	createSourceReq := satellitev1.CreateSatelliteEndpointSourceRequest{
		LocationID: locationID,
		Type:       "service",
		SourceName: sourceName,
		Addresses:  []string{"10.0.0.0/12"},
	}
	createSourceResp, err := satSourceAPI.CreateSatelliteEndpointSource(createSourceReq, target)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Create Satellite endpoint source response: %+v\n", createSourceResp)

	// 2. Verify Endpoint Source is created
	out, err := satSourceAPI.ListSatelliteEndpointSources(locationID, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("List of sources added to satellite location [%s]: %+v\n", locationID, out)

	// 3. Update existing source
	var sourceID string
	for _, satSource := range out.Sources {
		if strings.EqualFold(satSource.SourceName, sourceName) {
			sourceID = satSource.SourceID
		}
	}
	updateSatSourceReq := satellitev1.UpdateSatelliteEndpointSourceRequest{
		SourceName: sourceName,
		Addresses:  []string{"10.8.0.0/24"},
	}
	updatedSatSource, err := satSourceAPI.UpdateSatelliteEndpointSources(locationID, sourceID, updateSatSourceReq, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated satellite source details are: %+v\n", updatedSatSource)

	// 4. List existing endpoints
	endpoints, err := satEndpointAPI.GetEndpoints(locationID, target)

	if err != nil {
		log.Fatal(err)
	}

	var openShiftAPIEndpointID string

	for _, ep := range endpoints.Endpoint {
		if strings.Contains(ep.DisplayName, "openshift-api") {
			openShiftAPIEndpointID = ep.EndpointID
		}
	}

	// 4. Add endpoint source created in step 1 to `openshift-api-xxx-xxxx` endpoint

	endpointSources := satellitev1.EndpointSource{
		SourceID: createSourceResp.SourceID,
		Enabled:  true,
	}

	addSourcesToEndpointReq := satellitev1.AddSourcesRequest{
		Sources: []satellitev1.EndpointSource{endpointSources},
	}

	err = satEndpointAPI.AddSourcesToEndpoint(locationID, openShiftAPIEndpointID, addSourcesToEndpointReq, target)

	if err != nil {
		log.Fatal(err)
	}
}
