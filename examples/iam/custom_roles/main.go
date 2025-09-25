package main

import (
	"flag"
	"log"
	"os"

	"github.com/IBM-Cloud/bluemix-go/api/iampap/iampapv2"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	var acID string
	flag.StringVar(&acID, "acID", "", "Account ID")

	flag.Parse()
	if acID == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")
	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	roleClient, err := iampapv2.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(acID)

	rAPI := roleClient.IAMRoles()

	rolereq := iampapv2.CreateRoleRequest{
		Name:        "Testrole44",
		DisplayName: "Testrole44 disp",
		Description: "Custom role for example",
		ServiceName: "kms",
		Actions:     []string{"kms.policies.write"},
	}
	rolereq.AccountID = acID

	updatereq := iampapv2.UpdateRoleRequest{
		DisplayName: "Example Role updated3",
		Description: "aaa",
	}
	resp, err := rAPI.Create(rolereq)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("\nresp=", resp)

	listres, err := rAPI.ListAll(iampapv2.RoleQuery{AccountID: acID, ServiceName: "kms"})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("\nlistres=", listres)

	getresp, etag, err := rAPI.Get(resp.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("\ngetresp=", getresp)

	upres, err := rAPI.Update(updatereq, getresp.ID, etag)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("\nupresp=", upres)

	err1 := rAPI.Delete(getresp.ID)
	if err != nil {
		log.Fatal(err1)
	}
	log.Println("deleted")

}
