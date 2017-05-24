package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/session"
	"github.com/IBM-Bluemix/bluemix-go/trace"
)

func main() {
	var org string
	flag.StringVar(&org, "org", "", "Bluemix Organization")

	var space string
	flag.StringVar(&space, "space", "", "Bluemix Space")

	var host string
	flag.StringVar(&host, "host", "myexample", "Bluemix Host")

	var path string
	flag.StringVar(&path, "path", "/mypath", "Bluemix Path")

	var newHost string
	flag.StringVar(&newHost, "new_host", "myexample1", "Bluemix Path")

	flag.Parse()

	if org == "" || space == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	client, err := cfv2.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	orgAPI := client.Organizations()
	myorg, err := orgAPI.FindByName(org)
	if err != nil {
		log.Fatal(err)
	}

	spaceAPI := client.Spaces()
	myspace, err := spaceAPI.FindByNameInOrg(myorg.GUID, space)
	if err != nil {
		log.Fatal(err)
	}

	sharedDomainAPI := client.SharedDomains()

	sd, err := sharedDomainAPI.FindByName("mybluemix.net")
	if err != nil {
		log.Fatal(err)
	}

	routesAPI := client.Routes()

	payload := cfv2.RouteRequest{
		Host:       host,
		SpaceGUID:  myspace.GUID,
		DomainGUID: sd.GUID,
		Path:       path,
	}
	r, err := routesAPI.Create(payload)
	if err != nil {
		log.Fatal(err)
	}

	payload.Host = newHost

	updatedRoute, err := routesAPI.Update(r.Metadata.GUID, payload)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*updatedRoute)

	err = routesAPI.Delete(r.Metadata.GUID, true)
	if err != nil {
		log.Fatal(err)
	}

}
