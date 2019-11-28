package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	v2 "github.com/IBM-Cloud/bluemix-go/api/iam/iamv2"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {

	const (
		Member = "MEMBER"
	)
	var userEmail string
	flag.StringVar(&userEmail, "userEmail", "", "Email of the user to be invited")

	var accountID string
	flag.StringVar(&accountID, "accountID", "", "Account ID of the inviter")

	trace.Logger = trace.NewLogger("true")
	c := new(bluemix.Config)
	flag.Parse()

	if userEmail == "" || accountID == "" {
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
	users := make([]models.User, 0)
	users = append(users, models.User{Email: userEmail, AccountRole: Member})
	payload := models.UserInvite{
		Users: users,
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
