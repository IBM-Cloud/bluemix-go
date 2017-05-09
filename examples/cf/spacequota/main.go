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

	flag.Parse()

	if org == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	client, err := cfv2.New(sess)

	orgAPI := client.Organizations()
	myorg, err := orgAPI.FindByName(org)

	if err != nil {
		log.Fatal(err)
	}

	quotaAPI := client.SpaceQuotas()

	myquota, err := quotaAPI.Create("test1", myorg.GUID, 1024, 1024, 50, 150, false)

	if err != nil {
		log.Fatal(err)
	}

	myquota, err = quotaAPI.Get(myquota.Metadata.GUID)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(myquota.Metadata.GUID)

	myquota, err = quotaAPI.Update("testnew", myquota.Metadata.GUID, myorg.GUID, 1024, 1024, 50, 150, false)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(myquota.Metadata.GUID)

	quota, err := quotaAPI.FindByName("testnew", myorg.GUID)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(quota)

	err = quotaAPI.Delete(myquota.Metadata.GUID)

	if err != nil {
		log.Fatal(err)
	}
}
