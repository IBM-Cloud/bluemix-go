package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/iam/iamv1"
	"github.com/IBM-Cloud/bluemix-go/api/iampap/iampapv1"
	v2 "github.com/IBM-Cloud/bluemix-go/api/usermanagement/usermanagementv2"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
	"github.com/IBM-Cloud/bluemix-go/utils"
)

func main() {

	const (
		Member = "MEMBER"
	)
	var userEmail string
	flag.StringVar(&userEmail, "userEmail", "", "Email of the user to be invited")

	var accountID string
	flag.StringVar(&accountID, "accountID", "", "Account ID of the inviter")

	var service string
	flag.StringVar(&service, "service", "", "Bluemix service name")

	var roles string
	flag.StringVar(&roles, "roles", "", "Comma seperated list of roles")

	var serviceInstance string
	flag.StringVar(&serviceInstance, "serviceInstance", "", "Bluemix service instance name")

	var region string
	flag.StringVar(&region, "region", "", "Bluemix region")

	var resourceType string
	flag.StringVar(&resourceType, "resourceType", "", "Bluemix resource type")

	var resource string
	flag.StringVar(&resource, "resource", "", "Bluemix resource")

	var resourceGroupID string
	flag.StringVar(&resourceGroupID, "resourceGroupID", "", "Bluemix resource group ")

	var serviceType string
	flag.StringVar(&serviceType, "serviceType", "", "service type")

	trace.Logger = trace.NewLogger("true")
	c := new(bluemix.Config)
	flag.Parse()

	if userEmail == "" || accountID == "" || roles == "" {
		flag.Usage()
		os.Exit(1)
	}

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}
	userManagementHandler, err := v2.New(sess)
	if err != nil {
		log.Println("failed to initilize userManagementHandler", err)
	}
	userInvite := userManagementHandler.UserInvite()
	users := make([]v2.User, 0)
	users = append(users, v2.User{Email: userEmail, AccountRole: Member})

	// user roles
	iamClient, err := iamv1.New(sess)
	if err != nil {
		log.Fatal(err)
	}

	serviceRolesAPI := iamClient.ServiceRoles()

	var definedRoles []models.PolicyRole

	if service == "" {
		definedRoles, err = serviceRolesAPI.ListSystemDefinedRoles()
	} else {
		definedRoles, err = serviceRolesAPI.ListServiceRoles(service)
	}

	filterRoles, err := utils.GetRolesFromRoleNames(strings.Split(roles, ","), definedRoles)

	if err != nil {
		log.Fatal(err)
	}

	policyResource := iampapv1.Resource{}

	if service != "" {
		policyResource.SetServiceName(service)
	}

	if serviceInstance != "" {
		policyResource.SetServiceInstance(serviceInstance)
	}

	if region != "" {
		policyResource.SetRegion(region)
	}

	if resourceType != "" {
		policyResource.SetResourceType(resourceType)
	}

	if resource != "" {
		policyResource.SetResource(resource)
	}

	if resourceGroupID != "" {
		policyResource.SetResourceGroupID(resourceGroupID)
	}

	switch serviceType {
	case "service":
		fallthrough
	case "platform_service":
		policyResource.SetServiceType(serviceType)
	}

	if len(policyResource.Attributes) == 0 {
		policyResource.SetServiceType("service")
	}
	policyResource.SetAccountID(accountID)

	policy := v2.UserPolicy{Roles: iampapv1.ConvertRoleModels(filterRoles), Type: "access", Resources: []iampapv1.Resource{policyResource}}
	var Policies = []v2.UserPolicy{policy}

	payload := v2.UserInvite{
		Users:     users,
		IAMPolicy: Policies,
	}

	out, err := userInvite.InviteUsers(accountID, payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Invited User=", out)

	usersList, errList := userInvite.GetUsers(accountID)
	if errList != nil {
		log.Fatal(errList)
	}
	fmt.Println("List Of Users=", usersList)
	var UserIAMID string
	for _, u := range usersList.Resources {
		if u.Email == userEmail {
			UserIAMID = u.IamID
			break
		}
	}

	profile, errProf := userInvite.GetUserProfile(accountID, UserIAMID)
	if errProf != nil {
		log.Fatal(errProf)
	}
	fmt.Println("UserProfile=", profile)

	errRemove := userInvite.RemoveUsers(accountID, UserIAMID)
	if errRemove != nil {
		log.Fatal(errRemove)
	}
}
