package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/session"
	"github.com/IBM-Bluemix/bluemix-go/trace"
)

func main() {

	var path string
	flag.StringVar(&path, "path", "", "Bluemix path for application")

	var org string
	flag.StringVar(&org, "org", "", "Bluemix Organization")

	var name string
	flag.StringVar(&name, "name", "", "Bluemix app name")

	var space string
	flag.StringVar(&space, "space", "", "Bluemix Space")

	var routeName string
	flag.StringVar(&routeName, "routeName", "", "Bluemix app route")

	var buildpack string
	flag.StringVar(&buildpack, "buildpack", "", "Bluemix buildpack")

	var instance int
	flag.IntVar(&instance, "instance", 2, "Bluemix App Instance")

	var memory int
	flag.IntVar(&memory, "memory", 128, "Bluemix app memory")

	var diskQuota int
	flag.IntVar(&diskQuota, "diskQuota", 512, "Bluemix app diskquota")

	var async bool
	flag.BoolVar(&async, "async", false, "For asynchronous and synchronous request. Defaults to false")

	flag.Parse()

	if name == "" || space == "" || buildpack == "" || org == "" || path == "" || routeName == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	client, err := cfv2.New(sess)

	if err != nil {
		log.Fatal(err)
	}
	orgAPI := client.Organizations()
	myorg, err := orgAPI.FindByName(org)

	if err != nil {
		log.Fatal(err)
	}

	spaceAPI := client.Spaces()
	myspace, err := spaceAPI.FindByNameInOrg(myorg.GUID, space)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(myorg.GUID, myspace.GUID)

	appAPI := client.Apps()
	_, err = appAPI.Exists(myspace.GUID, name)

	if err != nil {
		log.Fatal(err)
	}

	var appPayload = &cfv2.AppCreateRequest{
		Name:      name,
		SpaceGUID: myspace.GUID,
		BuildPack: buildpack,
		Instances: instance,
		Memory:    memory,
		DiskQuota: diskQuota,
	}

	newapp, err := appAPI.Create(appPayload)
	fmt.Println(newapp)
	if err != nil {
		log.Fatal(err)
	}

	routeAPI := client.Routes()
	domain, err := routeAPI.GetSharedDomains("mybluemix.net")
	fmt.Println(domain)
	if err != nil {
		log.Fatal(err)
	}

	route, err := routeAPI.Find(routeName, domain.GUID)
	fmt.Println(route)
	if err != nil {
		log.Fatal(err)
	}

	if len(route) == 0 {
		newroute, err := routeAPI.Create(routeName, domain.GUID, myspace.GUID)
		fmt.Println(newroute)
		if err != nil {
			log.Fatal(err)
		}
		bindRoute, err := appAPI.BindRoute(newapp.Metadata.GUID, newroute.Metadata.GUID)
		fmt.Println(bindRoute)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		bindRoute, err := appAPI.BindRoute(newapp.Metadata.GUID, route[0].GUID)
		fmt.Println(bindRoute)
		if err != nil {
			log.Fatal(err)
		}

	}

	_, err = appAPI.Upload(newapp.Metadata.GUID, path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = appAPI.Start(newapp.Metadata.GUID, async)
	if err != nil {
		log.Fatal(err)
	}

	apps, err := appAPI.Get(newapp.Metadata.GUID)
	fmt.Println(apps)
	if err != nil {
		log.Fatal(err)
	}
	var appUpdatePayload = &cfv2.AppCreateRequest{
		Name:      "testappupdate",
		SpaceGUID: myspace.GUID,
	}

	updateapp, err := appAPI.Update(newapp.Metadata.GUID, appUpdatePayload)
	fmt.Println(updateapp)
	if err != nil {
		log.Fatal(err)
	}

	appInstances, err := appAPI.Instances(newapp.Metadata.GUID)
	fmt.Println(appInstances)
	if err != nil {
		log.Fatal(err)
	}

	err = appAPI.Delete(newapp.Metadata.GUID)
	if err != nil {
		log.Fatal(err)
	}

}
