package main

import (
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"

	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	target := v2.LoggingTargetHeader{}

	loggingClient, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	loggingAPI := loggingClient.Logging()

	var loggingInfo = v2.LoggingDeleteRequest{
		Cluster:  "DragonBoat-cluster",
		Instance: "bb8011cf-a42c-4637-867d-585ca27eac8d",
	}

	err1 := loggingAPI.DeleteLoggingConfig(loggingInfo, target)
	if err1 != nil {
		log.Fatal(err)
	}
	fmt.Println("out=", err1)

}
