package main

import (
	"fmt"
	"log"

	bluemix "github.com/IBM-Cloud/bluemix-go"

	"github.com/IBM-Cloud/bluemix-go/session"

	sch "github.com/IBM-Cloud/bluemix-go/api/schematics"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {
	c := new(bluemix.Config)

	workspaceID := "myworkspaceptab-bc54176d-bbcb-42"
	templateid := "29291199-ca08-46"
	trace.Logger = trace.NewLogger("true")

	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}
	schClient, err := sch.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	schAPI := schClient.Workspaces()
	//Get the workspace
	works, err := schAPI.GetWorkspaceByID(workspaceID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nThe workspace info= ", works)

	state, err := schAPI.GetStateStore(workspaceID, templateid)
	if err != nil {
		log.Fatal(err)
	}
	statestr := fmt.Sprintf("%v", state)
	fmt.Println("The state info= ", statestr)

	out, err := schAPI.GetOutputValues(workspaceID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("The output info= ", out)
}
