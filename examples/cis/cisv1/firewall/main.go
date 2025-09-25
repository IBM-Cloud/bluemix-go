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
	var firewallType string
	flag.StringVar(&firewallType, "firewallType", "", "firewallType for zone")

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

	log.Println(">>>>>>>>>  Firewall create")

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

	myFirewallPtr, err := firewallAPI.CreateFirewall(cis_id, zone_id, firewallType, params)

	if err != nil {
		log.Fatal(err)
	}

	myFirewall := *myFirewallPtr
	FirewallID := myFirewall.ID
	log.Println("Firewall create :", myFirewall)

	log.Println(">>>>>>>>>  Firewall read")
	myFirewallPtr, err = firewallAPI.GetFirewall(cis_id, zone_id, firewallType, FirewallID)

	if err != nil {
		log.Fatal(err)
	}

	myFirewall = *myFirewallPtr

	log.Println("Firewall Details by ID:", myFirewall)

	log.Println(">>>>>>>>>  Firewall delete")
	err = firewallAPI.DeleteFirewall(cis_id, zone_id, firewallType, FirewallID)
	if err != nil {
		log.Fatal(err)
	}
}
