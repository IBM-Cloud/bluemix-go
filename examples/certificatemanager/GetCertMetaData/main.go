package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Mavrickk3/bluemix-go"
	v "github.com/Mavrickk3/bluemix-go/api/certificatemanager"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"
)

func main() {

	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")
	var CertID string
	flag.StringVar(&CertID, "CertID", "", "Id of certificate")
	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	certClient, err := v.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	certificateAPI := certClient.Certificate()

	out, err := certificateAPI.GetMetaData(CertID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("out=", out)
}
