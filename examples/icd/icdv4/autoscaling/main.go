package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Cloud/bluemix-go/api/icd/icdv4"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var icdID, groupID string
	flag.StringVar(&icdID, "icdID", "", "CRN of the IBM Cloud Database service instance")
	flag.Parse()
	flag.StringVar(&groupID, "icdID", "", "CRN of the IBM Cloud Database service instance")
	flag.Parse()

	if icdID == "" || groupID == "" {
		flag.Usage()
		os.Exit(1)
	}
	autoScaleType := "memory"
	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	icdClient, err := icdv4.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	autoScalingAPI := icdClient.AutoScaling()
	param := icdv4.AutoscalingSetGroup{}
	memoryReq := icdv4.ASGBody{
		Scalers: icdv4.ScalersBody{
			IO: &icdv4.IOBody{
				Enabled:      true,
				AbovePercent: 35,
				OverPeriod:   "15m",
			},
		},
		Rate: icdv4.RateBody{
			IncreasePercent:  10,
			LimitMBPerMember: 5000,
			PeriodSeconds:    900,
			Units:            "mb",
		},
	}
	if autoScaleType == "memory" {
		param.Autoscaling.Memory = &memoryReq
	}

	setgroup, err := autoScalingAPI.SetAutoScaling(icdID, groupID, param)
	group, err := autoScalingAPI.GetAutoScaling(icdID, groupID)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("AutoScaling :%+v", group)
	log.Printf("AutoScaling :%+v", setgroup)

}
