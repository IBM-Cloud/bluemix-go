package main

import (
	"fmt"
	"log"
	"time"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
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

	clusterID := "ck4ufagd0e3bch0a3e8g"
	albID := "public-crck4ufagd0e3bch0a3e8g-alb1"

	updateReq := v2.UpdateALBReq{
		ClusterID: clusterID,
		//ALBBuild:  "1.8.1_5365_iks",
		ALBBuild: "1.5.1_5367_iks",
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
		AutoUpdate: &autoUpdateOff,
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
