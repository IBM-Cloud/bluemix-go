package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
)

func main() {

	var workerID string
	flag.StringVar(&workerID, "workerID", "", "worker ID")
	flag.Parse()

	if workerID == "" {
		flag.Usage()
		os.Exit(1)
	}
	c := new(bluemix.Config)

	trace.Logger = trace.NewLogger("true")

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	target := v1.ClusterTargetHeader{}

	clusterClient, err := v1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	workersAPI := clusterClient.Workers()

	params := v1.WorkerUpdateParam{
		Action: "reload",
	}

	err = workersAPI.Update("", workerID, params, target)
	if err != nil {
		fmt.Println("error: ", err.Error())
	} else {
		fmt.Println("ok")
	}
}
