package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/trace"
	"github.com/IBM-Bluemix/bluemix-go/api/account/accountv2"
	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/session"
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

	client, err := cfv2.NewClient(sess)

	if err != nil {
		log.Fatal(err)
	}

	orgAPI := client.Organizations()
	myorg, err := orgAPI.FindByName(org)

	if err != nil {
		log.Fatal(err)
	}
	accClient, err := accountv2.NewClient(sess)
	if err != nil {
		log.Fatal(err)
	}
	accountAPI := accClient.Accounts()
	myAccount, err := accountAPI.FindByOrg(myorg.GUID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(myAccount.Name, myAccount.CountryCode, myAccount.OwnerUserID)

}
