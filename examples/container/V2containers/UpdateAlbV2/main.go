package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
)

func main() {

	var clusterID, albID, albVersion string
	flag.StringVar(&clusterID, "clusterNameOrID", "", "cluster name or ID")
	flag.StringVar(&albID, "albID", "", "ALB ID")
	flag.StringVar(&albVersion, "albVersion", "", "target ALB build version")
	flag.Parse()

	trace.Logger = trace.NewLogger("true")

	if clusterID == "" || albID == "" || albVersion == "" {
		flag.Usage()
		os.Exit(1)
	}

	c := new(bluemix.Config)

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	updateReq := v2.UpdateALBReq{
		ClusterID: clusterID,
		ALBBuild:  albVersion,
		ALBList: []string{
			albID,
		},
	}

	targetV2 := v2.ClusterTargetHeader{}

	clusterClientV2, err := v2.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	targetV1 := v1.ClusterTargetHeader{}

	clusterClientV1, err := v1.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	albAPIV2 := clusterClientV2.Albs()
	albAPIV1 := clusterClientV1.Albs()

	autoUpdateOff := false
	albUpdatePolicy := v1.ALBUpdatePolicy{
		AutoUpdate: autoUpdateOff,
	}

	fmt.Println("turn off auto update")
	updatePolicyErr := albAPIV1.ChangeALBUpdatePolicy(clusterID, albUpdatePolicy, targetV1)
	if updatePolicyErr != nil {
		log.Fatal(updatePolicyErr)
	}

	albConf, getErr := albAPIV2.GetAlb(albID, targetV2)
	if err != nil {
		log.Fatal(getErr)
	}
	fmt.Println("before update version: ", albConf.AlbBuild)

	updateErr := albAPIV2.UpdateAlb(updateReq, targetV2)
	fmt.Println("updateErr: ", updateErr)

	fmt.Println("wait 60 sec...")
	time.Sleep(60 * time.Second)

	albConf, getErr = albAPIV2.GetAlb(albID, targetV2)
	if err != nil {
		log.Fatal(getErr)
	}
	fmt.Println("after update version: ", albConf.AlbBuild)
}
