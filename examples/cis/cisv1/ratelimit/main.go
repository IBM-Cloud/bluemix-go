package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Cloud/bluemix-go/api/cis/cisv1"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var cis_id string
	flag.StringVar(&cis_id, "cis_id", "", "CRN of the CIS service instance")
	var zone_id string
	flag.StringVar(&zone_id, "zone_id", "", "zone_id for zone")

	flag.Parse()
	if zone_id == "" || cis_id == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	cisClient, err := cisv1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	rateLimitAPI := cisClient.RateLimit()

	log.Println(">>>>>>>>>  RateLimit create")

	params := cisv1.RateLimitRecord{
		Description: "test",
		Threshold:   60,
		Period:      600,
		Correlate: &cisv1.RateLimitCorrelate{
			By: "nat",
		},
		Action: cisv1.RateLimitAction{
			Mode:    "simulate",
			Timeout: 40000,
			Response: &cisv1.ActionResponse{
				ContentType: "text/plain",
				Body:        "custom response body",
			},
		},
		Match: cisv1.RateLimitMatch{
			Request: cisv1.MatchRequest{
				URL: "https",
				Methods: []string{
					"GET", "POST", "PUT", "DELETE",
				},
				Schemes: []string{
					"HTTP", "HTTPS",
				},
			},
			Response: cisv1.MatchResponse{
				Statuses: []int{
					200, 201, 202,
				},
			},
		},
	}

	myRateLimitPtr, err := rateLimitAPI.CreateRateLimit(cis_id, zone_id, params)
	if err != nil {
		log.Fatal(err)
	}

	myRateLimit := *myRateLimitPtr
	ruleID := myRateLimit.ID

	log.Println(">>>>>>>>>  RateLimit Read")
	myRateLimitPtr, err = rateLimitAPI.GetRateLimit(cis_id, zone_id, ruleID)
	if err != nil {
		log.Fatal(err)
	}
	myRateLimit = *myRateLimitPtr
	log.Printf("My Rate Limit %v", myRateLimit)

	log.Println(">>>>>>>>>  RateLimit Update")
	_, err = rateLimitAPI.UpdateRateLimit(cis_id, zone_id, ruleID, params)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(">>>>>>>>>  RateLimit Delete")
	err = rateLimitAPI.DeleteRateLimit(cis_id, zone_id, ruleID)
	if err != nil {
		log.Fatal(err)
	}
}
