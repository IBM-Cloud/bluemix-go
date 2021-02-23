package main

import (
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/session"

	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	target := v2.LoggingTargetHeader{}

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	loggingClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	loggingAPI := loggingClient.Logging()

	out, err := loggingAPI.GetLoggingConfig("DragonBoat-cluster", "bb8011cf-a42c-4637-867d-585ca27eac8d", target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("out=", out)
}
