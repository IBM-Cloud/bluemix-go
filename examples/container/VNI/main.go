package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	gohttp "net/http"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/api/container/graphql"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

func main() {
	var (
		apiKey    = flag.String("apikey", os.Getenv("IBMCLOUD_API_KEY"), "IBM Cloud API Key")
		region    = flag.String("region", "us-south", "IBM Cloud region")
		accountID = flag.String("account", os.Getenv("IBMCLOUD_ACCOUNT_ID"), "IBM Cloud Account ID")

		// Operation flags
		operation = flag.String("op", "list", "Operation: attach, detach, list, or get")

		// Resource identifiers
		clusterID  = flag.String("cluster", "", "Cluster ID or name")
		workerID   = flag.String("worker", "", "Worker ID")
		vniID      = flag.String("vni", "", "VNI ID (e.g., r006-xxx)")
		vlanID     = flag.Int("vlan", 100, "VLAN ID (1-500, for attach operation)")
		autoDelete = flag.Bool("auto-delete", false, "Auto-delete VNI when detaching")

		// Pagination
		first = flag.Int("first", 10, "Number of items to fetch (for list operation)")
		after = flag.String("after", "", "Cursor for pagination (for list operation)")
	)

	flag.Parse()

	if *apiKey == "" {
		log.Fatal("API key is required. Set IBMCLOUD_API_KEY environment variable or use -apikey flag")
	}

	if *accountID == "" {
		log.Fatal("Account ID is required. Set IBMCLOUD_ACCOUNT_ID environment variable or use -account flag")
	}

	// Enable trace for debugging
	trace.Logger = trace.NewLogger("true")

	// Create session
	sess, err := session.New(&bluemix.Config{
		BluemixAPIKey: *apiKey,
		Region:        *region,
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Create client configuration
	config := sess.Config.Copy()
	err = config.ValidateConfigForService(bluemix.ContainerService)
	if err != nil {
		log.Fatalf("Failed to validate config: %v", err)
	}

	if config.HTTPClient == nil {
		config.HTTPClient = http.NewHTTPClient(config)
	}

	tokenRefresher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"X-Original-User-Agent": []string{config.UserAgent},
			"User-Agent":            []string{http.UserAgent()},
		},
		HTTPClient: config.HTTPClient,
	})
	if err != nil {
		log.Fatalf("Failed to create token refresher: %v", err)
	}

	if config.IAMAccessToken == "" {
		err := authentication.PopulateTokens(tokenRefresher, config)
		if err != nil {
			log.Fatalf("Failed to populate tokens: %v", err)
		}
	}

	if config.Endpoint == nil {
		ep, err := config.EndpointLocator.ContainerEndpoint()
		if err != nil {
			log.Fatalf("Failed to get container endpoint: %v", err)
		}
		config.Endpoint = &ep
	}

	// Create client
	c := client.New(config, bluemix.ContainerService, tokenRefresher)

	// Create VNI API client
	vniAPI := graphql.NewVNIAPI(c)

	target := containerv1.ClusterTargetHeader{
		AccountID: *accountID,
	}

	// Execute operation
	switch *operation {
	case "attach":
		if err := attachVNI(vniAPI, target, *clusterID, *workerID, *vniID, *vlanID, *autoDelete); err != nil {
			log.Fatalf("Failed to attach VNI: %v", err)
		}

	case "detach":
		if err := detachVNI(vniAPI, target, *clusterID, *workerID, *vniID, *autoDelete); err != nil {
			log.Fatalf("Failed to detach VNI: %v", err)
		}

	case "list":
		if err := listVNIAttachments(vniAPI, target, *clusterID, *workerID, *first, *after); err != nil {
			log.Fatalf("Failed to list VNI attachments: %v", err)
		}

	case "get":
		if err := getVNIAttachment(vniAPI, target, *clusterID, *workerID, *vniID); err != nil {
			log.Fatalf("Failed to get VNI attachment: %v", err)
		}

	default:
		log.Fatalf("Unknown operation: %s. Use: attach, detach, list, or get", *operation)
	}
}

func attachVNI(api graphql.VNIs, target containerv1.ClusterTargetHeader, clusterID, workerID, vniID string, vlanID int, autoDelete bool) error {
	if vniID == "" {
		return fmt.Errorf("VNI ID is required for attach operation")
	}

	if clusterID == "" && workerID == "" {
		return fmt.Errorf("Either cluster ID or worker ID is required")
	}

	if clusterID != "" && workerID != "" {
		return fmt.Errorf("Specify either cluster ID or worker ID, not both")
	}

	input := graphql.AddVNIToBareMetalNodeInput{
		VirtualNetworkInterfaceID: vniID,
		VlanID:                    vlanID,
		AutoDelete:                autoDelete,
	}

	if clusterID != "" {
		input.Cluster = &clusterID
	} else {
		input.Node = &workerID
	}

	fmt.Printf("Attaching VNI %s to ", vniID)
	if clusterID != "" {
		fmt.Printf("cluster %s", clusterID)
	} else {
		fmt.Printf("worker %s", workerID)
	}
	fmt.Printf(" with VLAN ID %d...\n", vlanID)

	resp, err := api.AttachToBareMetalNode(input, target)
	if err != nil {
		return err
	}

	fmt.Println("\n✅ VNI attached successfully!")
	fmt.Printf("Worker ID: %s\n", resp.NetworkAttachment.AttachedTo.ID)
	fmt.Printf("VNI ID: %s\n", resp.NetworkAttachment.VirtualNetworkInterface.ExternalID)
	if resp.NetworkAttachment.VirtualNetworkInterface.Name != nil {
		fmt.Printf("VNI Name: %s\n", *resp.NetworkAttachment.VirtualNetworkInterface.Name)
	}
	if resp.NetworkAttachment.VlanID != nil {
		fmt.Printf("VLAN ID: %d\n", *resp.NetworkAttachment.VlanID)
	}
	if resp.NetworkAttachment.VirtualNetworkInterface.PrimaryIPAddress != nil {
		fmt.Printf("Primary IP: %s\n", *resp.NetworkAttachment.VirtualNetworkInterface.PrimaryIPAddress)
	}
	if resp.NetworkAttachment.VirtualNetworkInterface.MACAddress != nil {
		fmt.Printf("MAC Address: %s\n", *resp.NetworkAttachment.VirtualNetworkInterface.MACAddress)
	}

	return nil
}

func detachVNI(api graphql.VNIs, target containerv1.ClusterTargetHeader, clusterID, workerID, vniID string, autoDelete bool) error {
	if vniID == "" {
		return fmt.Errorf("VNI ID is required for detach operation")
	}

	if clusterID == "" && workerID == "" {
		return fmt.Errorf("Either cluster ID or worker ID is required")
	}

	input := graphql.RemoveVNIFromNodeInput{
		VirtualNetworkInterfaceID: vniID,
		AutoDelete:                autoDelete,
	}

	if clusterID != "" {
		input.Cluster = &clusterID
	} else {
		input.Node = &workerID
	}

	fmt.Printf("Detaching VNI %s from ", vniID)
	if clusterID != "" {
		fmt.Printf("cluster %s", clusterID)
	} else {
		fmt.Printf("worker %s", workerID)
	}
	fmt.Println("...")

	resp, err := api.DetachFromNode(input, target)
	if err != nil {
		return err
	}

	fmt.Println("\n✅ VNI detached successfully!")
	fmt.Printf("Cluster ID: %s\n", resp.Cluster.ID)
	fmt.Printf("Worker ID: %s\n", resp.Node.ID)
	fmt.Printf("VNI ID: %s\n", resp.VirtualNetworkInterface.ExternalID)
	if resp.VirtualNetworkInterface.Name != nil {
		fmt.Printf("VNI Name: %s\n", *resp.VirtualNetworkInterface.Name)
	}

	return nil
}

func listVNIAttachments(api graphql.VNIs, target containerv1.ClusterTargetHeader, clusterID, workerID string, first int, after string) error {
	nodeID := clusterID
	nodeType := "cluster"
	if workerID != "" {
		nodeID = workerID
		nodeType = "worker"
	}

	if nodeID == "" {
		return fmt.Errorf("Either cluster ID or worker ID is required for list operation")
	}

	input := graphql.ListVNIAttachmentsInput{
		NodeID: nodeID,
	}

	if first > 0 {
		input.First = &first
	}
	if after != "" {
		input.After = &after
	}

	fmt.Printf("Listing VNI attachments for %s %s...\n\n", nodeType, nodeID)

	resp, err := api.ListAttachments(input, target)
	if err != nil {
		return err
	}

	fmt.Printf("Attachable Type: %s\n", resp.AttachableType)
	fmt.Printf("Total Attachments: %d\n\n", len(resp.Connection.Edges))

	if len(resp.Connection.Edges) == 0 {
		fmt.Println("No VNI attachments found.")
		return nil
	}

	fmt.Println("VNI Attachments:")
	fmt.Println("================")
	for i, edge := range resp.Connection.Edges {
		attachment := edge.Node
		fmt.Printf("\n%d. Worker: %s\n", i+1, attachment.AttachedTo.ID)
		fmt.Printf("   VNI ID: %s\n", attachment.VirtualNetworkInterface.ExternalID)
		if attachment.VirtualNetworkInterface.Name != nil {
			fmt.Printf("   VNI Name: %s\n", *attachment.VirtualNetworkInterface.Name)
		}
		if attachment.VlanID != nil {
			fmt.Printf("   VLAN ID: %d\n", *attachment.VlanID)
		}
		if attachment.VirtualNetworkInterface.PrimaryIPAddress != nil {
			fmt.Printf("   Primary IP: %s\n", *attachment.VirtualNetworkInterface.PrimaryIPAddress)
		}
		if attachment.VirtualNetworkInterface.MACAddress != nil {
			fmt.Printf("   MAC Address: %s\n", *attachment.VirtualNetworkInterface.MACAddress)
		}
	}

	if resp.Connection.PageInfo.HasNextPage {
		fmt.Printf("\n📄 More results available. Use -after=%s to fetch next page\n", *resp.Connection.PageInfo.EndCursor)
	}

	return nil
}

func getVNIAttachment(api graphql.VNIs, target containerv1.ClusterTargetHeader, clusterID, workerID, vniID string) error {
	if clusterID == "" {
		return fmt.Errorf("Cluster ID is required for get operation")
	}

	// Get operation is essentially a filtered list
	input := graphql.ListVNIAttachmentsInput{
		NodeID: clusterID,
	}

	resp, err := api.ListAttachments(input, target)
	if err != nil {
		return err
	}

	// Filter by worker and VNI if specified
	var found *graphql.VNIAttachment
	for _, edge := range resp.Connection.Edges {
		attachment := edge.Node
		matchWorker := workerID == "" || attachment.AttachedTo.ID == workerID
		matchVNI := vniID == "" || attachment.VirtualNetworkInterface.ExternalID == vniID

		if matchWorker && matchVNI {
			found = &attachment
			break
		}
	}

	if found == nil {
		return fmt.Errorf("VNI attachment not found")
	}

	fmt.Println("✅ VNI Attachment Found:")
	fmt.Println("========================")
	fmt.Printf("Worker ID: %s\n", found.AttachedTo.ID)
	fmt.Printf("VNI ID: %s\n", found.VirtualNetworkInterface.ExternalID)
	if found.VirtualNetworkInterface.Name != nil {
		fmt.Printf("VNI Name: %s\n", *found.VirtualNetworkInterface.Name)
	}
	if found.VlanID != nil {
		fmt.Printf("VLAN ID: %d\n", *found.VlanID)
	}
	if found.VirtualNetworkInterface.PrimaryIPAddress != nil {
		fmt.Printf("Primary IP: %s\n", *found.VirtualNetworkInterface.PrimaryIPAddress)
	}
	if found.VirtualNetworkInterface.MACAddress != nil {
		fmt.Printf("MAC Address: %s\n", *found.VirtualNetworkInterface.MACAddress)
	}

	return nil
}
