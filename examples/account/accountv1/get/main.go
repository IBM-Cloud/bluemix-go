package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/account/accountv1"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var accountId string
	flag.StringVar(&accountId, "accountId", "", "Bluemix Account ID")

	if accountId == "" {
		flag.Usage()
		os.Exit(1)
	}

	c := new(bluemix.Config)
	trace.Logger = trace.NewLogger("true")

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	accClient1, err := accountv1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	accountAPIV1 := accClient1.Accounts()
	//Get list of users under account
	account, err := accountAPIV1.Get(accountId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Accounts", account)
}
