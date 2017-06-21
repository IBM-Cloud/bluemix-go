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

	var name string
	flag.StringVar(&name, "name", "myexample.net", "Private Domain Name")

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
	if err != nil {
		log.Fatal(err)
	}

	orgAPI := client.Organizations()
	region := sess.Config.Region
	myorg, err := orgAPI.FindByName(org, region)
	if err != nil {
		log.Fatal(err)
	}

	privateDomainAPI := client.PrivateDomains()

	payload := cfv2.PrivateDomainRequest{
		Name:    name,
		OrgGUID: myorg.GUID,
	}
	domain, err := privateDomainAPI.Create(payload)
	if err != nil {
		log.Fatal(err)
	}

	domain, err = privateDomainAPI.Get(domain.Metadata.GUID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(domain)

	err = privateDomainAPI.Delete(domain.Metadata.GUID, true)
	if err != nil {
		log.Fatal(err)
	}

}
