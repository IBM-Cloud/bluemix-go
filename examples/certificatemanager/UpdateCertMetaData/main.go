package main

import (
	"flag"
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

	updateMetadata := models.CertificateMetadataUpdate{
		Name:        "Kavya",
		Description: "lalala",
	}

	certClient, err := v.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	certificateAPI := certClient.Certificate()

	err2 := certificateAPI.UpdateCertificateMetaData(CertID, updateMetadata)
	if err != nil {
		log.Fatal(err2)
	}

}
