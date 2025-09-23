package main

import (
	"fmt"
	"log"

	bluemix "github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/api/account/accountv1"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"
)

func main() {

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
	accounts, err := accountAPIV1.List()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total accounts", len(accounts))

	for _, acc := range accounts {
		fmt.Println(acc.Guid)
	}
}
