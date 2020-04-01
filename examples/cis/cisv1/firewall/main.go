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

	var cis_id = "crn:v1:bluemix:public:internet-svcs:global:a/4448261269a14562b839e0a3019ed980:7e8666d0-67ce-4fc1-a17f-f80f5459c5c2::"
	//flag.StringVar(&cis_id, "cis_id", "", "CRN of the CIS service instance")
	var zone_id = "5e7dd209f43b1d247ea3863c7e465db0"
	//flag.StringVar(&zone_id, "zone_id", "", "zone_id for zone")
	var firewallType = "lockdowns"
	//flag.StringVar(&firewallType, "firewallType", "", "firewallType for zone")

	flag.Parse()

	if zone_id == "" || cis_id == "" || firewallType == "" {
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
	firewallAPI := cisClient.Firewall()

	log.Println(">>>>>>>>>  Dns create")

	params := cisv1.FirewallBody{
		Paused:         false,
		Description:    "testlockdown1",
		Urls:           []string{"www.cis-terraform.com"},
		Configurations: []cisv1.Configuration{},
	}
	configuration := cisv1.Configuration{
		Target: "ip",
		Value:  "127.0.0.2",
	}
	params.Configurations = append(params.Configurations, configuration)

	log.Println(params)

	// params.Configurations[0].Target = "ip"
	// params.Configurations[0].Value = "127.0.0.2"

	myFirewallPtr, err := firewallAPI.CreateFirewall(cis_id, zone_id, firewallType, params)

	if err != nil {
		log.Fatal(err)
	}

	myFirewall := *myFirewallPtr
	FirewallID := myFirewall.ID
	log.Println("Dns create :", myFirewall)

	log.Println(">>>>>>>>>  Dns read")
	myFirewallPtr, err = firewallAPI.GetFirewall(cis_id, zone_id, firewallType, FirewallID)

	if err != nil {
		log.Fatal(err)
	}

	myFirewall = *myFirewallPtr

	log.Println("Dns Details by ID:", myFirewall)

	log.Println(">>>>>>>>>  Dns delete")
	err = firewallAPI.DeleteFirewall(cis_id, zone_id, firewallType, FirewallID)
	if err != nil {
		log.Fatal(err)
	}
}
