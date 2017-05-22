package main

import (
	"flag"
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
	log.Println(myorg.GUID, myspace.GUID)

	quotaAPI := client.SpaceQuotas()

	createRequest := cfv2.SpaceQuotaCreateRequest{
		Name:                    "test1",
		OrgGUID:                 myorg.GUID,
		MemoryLimitInMB:         1024,
		InstanceMemoryLimitInMB: 1024,
		RoutesLimit:             50,
		ServicesLimit:           150,
		NonBasicServicesAllowed: false,
	}

	myquota, err := quotaAPI.Create(createRequest)
	if err != nil {
		log.Fatal(err)
	}

	newspace, err := spaceAPI.Create("test", myorg.GUID, myquota.Metadata.GUID)
	if err != nil {
		log.Fatal(err)
	}

	newspace, err = spaceAPI.Get(newspace.Metadata.GUID)
	if err != nil {
		log.Fatal(err)
	}

	newspace, err = spaceAPI.Update("testupdate", newspace.Metadata.GUID)
	if err != nil {
		log.Fatal(err)
	}

	err = spaceAPI.Delete(newspace.Metadata.GUID)
	if err != nil {
		log.Fatal(err)
	}

	err = quotaAPI.Delete(myquota.Metadata.GUID)
	if err != nil {
		log.Fatal(err)
	}
}
