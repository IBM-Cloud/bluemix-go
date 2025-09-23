package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Mavrickk3/bluemix-go"
	v "github.com/Mavrickk3/bluemix-go/api/certificatemanager"
	"github.com/Mavrickk3/bluemix-go/models"
	"github.com/Mavrickk3/bluemix-go/session"
	"github.com/Mavrickk3/bluemix-go/trace"
)

func main() {

	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	var CertID string
	flag.StringVar(&CertID, "CertID", "", "Id of Certificate")
	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	renewdata := models.CertificateRenewData{
		RotateKeys: false,
	}

	certClient, err := v.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	certificateAPI := certClient.Certificate()

	out, err := certificateAPI.RenewCertificate(CertID, renewdata)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("out=", out)
}
