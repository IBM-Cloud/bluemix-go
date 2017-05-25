package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/session"
	"github.com/IBM-Bluemix/bluemix-go/trace"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "myexample.net", "Shared Domain Name")

	flag.Parse()

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	client, err := cfv2.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	sharedDomainAPI := client.SharedDomains()

	payload := cfv2.SharedDomainRequest{
		Name: name,
	}
	domain, err := sharedDomainAPI.Create(payload)
	if err != nil {
		log.Fatal(err)
	}

	domain, err = sharedDomainAPI.Get(domain.Metadata.GUID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(domain)

	err = sharedDomainAPI.Delete(domain.Metadata.GUID, true)
	if err != nil {
		log.Fatal(err)
	}

}
