package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Bluemix/bluemix-go/api/resource/resourcev1/management"
	"github.com/IBM-Bluemix/bluemix-go/session"
	"github.com/IBM-Bluemix/bluemix-go/trace"
)

func main() {

	var resourcequota string
	flag.StringVar(&resourcequota, "quota", "", "Bluemix Org Quota Definition")

	flag.Parse()

	if resourcequota == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	client, err := management.New(sess)

	if err != nil {
		log.Fatal(err)
	}

	orgQuotaAPI := client.ResourceQuota()

	quota, err := orgQuotaAPI.FindByName(resourcequota)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Quota Defination Details :", quota)

}
