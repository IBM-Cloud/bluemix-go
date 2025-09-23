package main

import (
	"flag"
	"log"
	"os"

	"github.com/Mavrickk3/bluemix-go"
	"github.com/Mavrickk3/bluemix-go/api/hpcs"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"
)

func main() {

	var instanceID string
	flag.StringVar(&instanceID, "instanceID", "", "instance ID")

	flag.Parse()
	if instanceID == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")

	c := new(bluemix.Config)

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}
	// sess.Config.Region = "us-east"
	hpcsClient, err := hpcs.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	hsAPI := hpcsClient.Endpoint()

	resp, err := hsAPI.GetAPIEndpoint(instanceID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("\nresp=", resp)

}
