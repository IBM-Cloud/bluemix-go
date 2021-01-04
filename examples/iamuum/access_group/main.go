package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Cloud/bluemix-go/models"

	"github.com/IBM-Cloud/bluemix-go/api/account/accountv2"
	"github.com/IBM-Cloud/bluemix-go/api/iamuum/iamuumv2"
	"github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {
	var org string
	flag.StringVar(&org, "org", "umaraliorg", "Bluemix Organization")

	var accessGroup string
	flag.StringVar(&accessGroup, "accessGroup", "test", "Bluemix access group name")

	var user string
	flag.StringVar(&user, "user", "umarali.nagoor@in.ibm.com", "IBM-id or email id of the user")

	flag.Parse()
	if org == "" || accessGroup == "" || user == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	client, err := mccpv2.New(sess)

	if err != nil {
		log.Fatal(err)
	}
	orgAPI := client.Organizations()
	myorg, err := orgAPI.FindByName(org, sess.Config.Region)

	if err != nil {
		log.Fatal(err)
	}

	accClient, err := accountv2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	accountAPI := accClient.Accounts()
	myAccount, err := accountAPI.FindByOrg(myorg.GUID, sess.Config.Region)
	if err != nil {
		log.Fatal(err)
	}

	iamuumClient, err := iamuumv2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	accessGroupAPI := iamuumClient.AccessGroup()

	data := models.AccessGroupV2{
		AccessGroup: models.AccessGroup{
			Name: accessGroup,
		},
	}
	agID, err := accessGroupAPI.Create(data, myAccount.GUID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("create", agID)

	agID, _, err = accessGroupAPI.Get(agID.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("get", agID)

	err = accessGroupAPI.Delete(agID.ID, false)
	if err != nil {
		log.Fatal(err)
	}

	accGroups, err1 := accessGroupAPI.List(myAccount.GUID, user)
	if err1 != nil {
		log.Fatal(err1)
	}
	log.Println("ListAccGroups=> ", accGroups)

}
