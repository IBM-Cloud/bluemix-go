package main

import (
	"flag"
	"log"
	"os"
	"strings"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/account/accountv2"
	"github.com/IBM-Cloud/bluemix-go/api/iam/iamv1"
	"github.com/IBM-Cloud/bluemix-go/api/iampap/iampapv1"
	"github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
	"github.com/IBM-Cloud/bluemix-go/utils"

)

func main() {
	var org string
	flag.StringVar(&org, "org", "", "Bluemix Organization")
	
	var sourceServiceName string
	flag.StringVar(&sourceServiceName, "source_service_name", "", "Bluemix service name")
	
	var targetServiceName string
	flag.StringVar(&targetServiceName, "target_service_name", "", "Bluemix service name")

	var roles string
	flag.StringVar(&roles, "roles", "", "Comma seperated list of roles")

	var sourceServiceInstanceId string
	flag.StringVar(&sourceServiceInstanceId, "source_service_instance_id", "", "Bluemix source service instance id")
	
	var targetServiceInstanceId string
	flag.StringVar(&targetServiceInstanceId, "target_service_instance_id", "", "Bluemix target service instance id")
	
	var sourceResourceGroupId string
	flag.StringVar(&sourceResourceGroupId, "source_resource_group_id", "", "Bluemix source resource group id")
	
	var targetResourceGroupId string
	flag.StringVar(&targetResourceGroupId, "target_resource_group_id", "", "Bluemix target resource group id")
	
	var sourceResourceType string
	flag.StringVar(&sourceResourceType, "source_resource_type", "", "Source resource type")
	
	var targetResourceType string
	flag.StringVar(&targetResourceType, "target_resource_type", "", "Target resource type")

	trace.Logger = trace.NewLogger("true")
	c := new(bluemix.Config)
	flag.BoolVar(&c.Debug, "debug", false, "Show full trace if on")
	flag.Parse()

	if org == "" || sourceServiceName == "" || targetServiceName == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	client, err := mccpv2.New(sess)

	if err != nil {
		log.Fatal(err)
	}

	orgAPI := client.Organizations()
	myorg, err := orgAPI.FindByName(org, sess.Config.Region)

	if err != nil {
		log.Fatal(err)
	}

	accClient, err := accountv2.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	accountAPI := accClient.Accounts()
	myAccount, err := accountAPI.FindByOrg(myorg.GUID, sess.Config.Region)
	if err != nil {
		log.Fatal(err)
	}
	
	iamClient, err := iamv1.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	serviceRolesAPI := iamClient.ServiceRoles()
	
	var definedRoles []models.PolicyRole

	if sourceServiceName == "" {
		definedRoles, err = serviceRolesAPI.ListSystemDefinedRoles()
	} else {
		definedRoles, err = serviceRolesAPI.ListAuthorizationRoles(sourceServiceName, targetServiceName)
	}

	if err != nil {
		log.Fatal(err)
	}
	
	filterRoles, err := utils.GetRolesFromRoleNames(strings.Split(roles, ","), definedRoles)
	
	if err != nil {
		log.Fatal(err)
	}
	
	policy := iampapv1.Policy{
		Type: iampapv1.AuthorizationPolicyType,
	}
	
	policy.Roles = iampapv1.ConvertRoleModels(filterRoles)
	
	policy.Subjects = []iampapv1.Subject{
		{
			Attributes: []iampapv1.Attribute{
				{
					Name:  "accountId",
					Value: myAccount.GUID,
				},
				{
					Name:  "serviceName",
					Value: sourceServiceName,
				},
			},
		},
	}

	policy.Resources = []iampapv1.Resource{
		{
			Attributes: []iampapv1.Attribute{
				{
					Name:  "accountId",
					Value: myAccount.GUID,
				},
				{
					Name:  "serviceName",
					Value: targetServiceName,
				},
			},
		},
	}
	
	if sourceServiceInstanceId != "" {
		policy.Subjects[0].SetServiceInstance(sourceServiceInstanceId)
	}
	
	if targetServiceInstanceId != "" {
		policy.Resources[0].SetServiceInstance(targetServiceInstanceId)
	}
	
	if sourceResourceGroupId != "" {
		policy.Subjects[0].SetResourceGroupID(sourceResourceGroupId)
	}
	
	if targetResourceGroupId != "" {
		policy.Resources[0].SetResourceGroupID(targetResourceGroupId)
	}
	
	iampapClient, err := iampapv1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	
	authPolicy := iampapClient.V1Policy()
	
	createdAuthPolicy, err := authPolicy.Create(policy)

	if err != nil {
		log.Fatal(err)
	}
	
	log.Println(createdAuthPolicy)
	
	getPolicy, err := authPolicy.Get(createdAuthPolicy.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(getPolicy)

	err = authPolicy.Delete(createdAuthPolicy.ID)
	if err != nil {
		log.Fatal(err)
	}
		
}