package cfv2

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//ErrCodeAppAlreadyExist ...
var ErrCodeAppAlreadyExist = "AppAlreadyExist"

const TimeOutError = "Timeout"

//AppCreateRequest ...
type AppCreateRequest struct {
	Name                     string `json:"name"`
	Memory                   int    `json:"memory,omitempty"`
	Instances                int    `json:"instances,omitempty"`
	DiskQuota                int    `json:"disk_quota,omitempty"`
	SpaceGUID                string `json:"space_guid"`
	StackGUID                string `json:"stack_guid,omitempty"`
	State                    string `json:"state,omitempty"`
	DetectedStartCommand     string `json:"detected_start_command,omitempty"`
	Command                  string `json:"command,omitempty"`
	BuildPack                string `json:"buildpack,omitempty"`
	HealthCheckType          string `json:"health_check_type,omitempty"`
	HealthCheckTimeout       int    `json:"health_check_type,omitempty"`
	Diego                    bool   `json:"diego,omitempty"`
	EnableSSH                bool   `json:"enable_ssh,omitempty"`
	DockerImage              string `json:"docker_image,omitempty"`
	StagingFailedReason      string `json:"staging_failed_reason,omitempty"`
	StagingFailedDescription string `json:"staging_failed_description,omitempty"`
	Ports                    []int  `json:"ports,omitempty"`
}

//AppsStateUpdateRequest ...
type AppsStateUpdateRequest struct {
	State string `json:"state"`
}

//Metadata ...
type AppMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//AppEntity ...
type AppEntity struct {
	Name                     string `json:"name"`
	SpaceGUID                string `json:"space_guid"`
	StackGUID                string `json:"stack_guid"`
	State                    string `json:"state"`
	PackageState             string `json:"package_state"`
	Memory                   int    `json:"memory"`
	Instances                int    `json:"instances"`
	DiskQuota                int    `json:"disk_quota"`
	Version                  string `json:"version"`
	Command                  string `json:"command"`
	Console                  bool   `json:"console"`
	Debug                    string `json:"debug"`
	StagingTaskID            string `json:"staging_task_id"`
	HealthCheckType          string `json:"health_check_type"`
	HealthCheckTimeout       string `json:"health_check_timeout"`
	StagingFailedReason      string `json:"staging_failed_reason"`
	StagingFailedDescription string `json:"staging_failed_description"`
	Diego                    bool   `json:"diego"`
	DockerImage              string `json:"docker_image"`
	EnableSSH                bool   `json:"enable_ssh"`
	Ports                    []int  `json:"ports"`
}

//AppResource ...
type AppResource struct {
	Resource
	Entity AppEntity
}

type AppFields struct {
	Metadata AppMetadata
	Entity   AppEntity
}

type AppSummaryFields struct {
	GUID             string `json:"guid"`
	Name             string `json:"name"`
	State            string `json:"state"`
	PackageState     string `json:"package_state"`
	RunningInstances int    `json:"running_instances"`
}

type AppStats struct {
	State string `json:"state"`
}

type statTime struct {
	time.Time
}

//ToFields ..
func (resource AppResource) ToFields() App {
	entity := resource.Entity

	return App{
		GUID:      resource.Metadata.GUID,
		Name:      entity.Name,
		SpaceGUID: entity.SpaceGUID,
		StackGUID: entity.StackGUID,
	}
}

//App model
type App struct {
	GUID      string
	Name      string
	SpaceGUID string
	StackGUID string
}

//Apps ...
type Apps interface {
	Create(appPayload *AppCreateRequest) (*AppFields, error)
	List() ([]App, error)
	Get(appGUID string) (*AppFields, error)
	Update(appGUID string, appPayload *AppCreateRequest) (*AppFields, error)
	Delete(appGUID string) error
	Exists(spaceGUID string, name string) ([]App, error)
	BindRoute(appGUID, routeGUID string) (*AppFields, error)
	Start(appGUID string, async bool) (*AppFields, error)
	Upload(path string, name string) (*AppFields, error)
	Summary(appGUID string) (*AppSummaryFields, error)
	Stat(appGUID string) (map[string]AppStats, error)
	CheckAppStatus(waitForThisState, appGUID string) (bool, error)
	CheckInstanceStatus(waitForThisState, appGUID string) (bool, error)
	Instances(appGUID string) (map[string]AppStats, error)
}

type app struct {
	client *client.Client
}

func newAppAPI(c *client.Client) Apps {
	return &app{
		client: c,
	}
}

func (r *app) Exists(spaceGUID string, name string) ([]App, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s/apps", spaceGUID)
	req := rest.GetRequest(rawURL).Query("q", "name:"+name)
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	apps, err := r.listAppWithPath(path)
	if err != nil {
		return nil, err
	}
	if len(apps) != 0 {
		return nil, bmxerror.New(ErrCodeAppAlreadyExist,
			fmt.Sprintf("Given app:  %q already exist in given space: %q", name, spaceGUID))

	}
	return apps, nil
}

func (r *app) Create(appPayload *AppCreateRequest) (*AppFields, error) {

	rawURL := "/v2/apps?async=true"
	appFields := AppFields{}
	_, err := r.client.Post(rawURL, appPayload, &appFields)
	if err != nil {
		return nil, err
	}
	return &appFields, nil
}

func (r *app) BindRoute(appGUID, routeGUID string) (*AppFields, error) {

	rawURL := fmt.Sprintf("/v2/apps/%s/routes/%s", appGUID, routeGUID)
	appFields := AppFields{}
	_, err := r.client.Put(rawURL, nil, &appFields)
	if err != nil {
		return nil, err
	}
	return &appFields, nil
}

func (r *app) listAppWithPath(path string) ([]App, error) {
	var apps []App
	_, err := r.client.GetPaginated(path, AppResource{}, func(resource interface{}) bool {
		if appResource, ok := resource.(AppResource); ok {
			apps = append(apps, appResource.ToFields())
			return true
		}
		return false
	})
	return apps, err
}

func (r *app) newfileUploadRequest(rawURL string, params map[string]string, paramName, path string) (*AppFields, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	h := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	appFields := AppFields{}
	_, err = r.client.Put(rawURL, body, &appFields, h)
	if err != nil {
		return nil, err
	}

	return &appFields, err

}

func (r *app) Upload(appGUID string, zippath string) (*AppFields, error) {

	extraParams := map[string]string{

		"resources": "[]",
	}
	uploadResp, err := r.newfileUploadRequest("/v2/apps/"+appGUID+"/bits?async=false", extraParams, "application", zippath)
	if err != nil {
		log.Fatal(err)
	}
	return uploadResp, err
}

func (r *app) Start(appGUID string, async bool) (*AppFields, error) {
	payload := AppsStateUpdateRequest{
		State: "STARTED",
	}
	rawURL := fmt.Sprintf("/v2/apps/%s", appGUID)
	appFields := AppFields{}
	_, err := r.client.Put(rawURL, payload, &appFields)
	if err != nil {
		return nil, err
	}
	if !async {
		isFinished, err := r.CheckAppStatus("STAGED", appGUID)
		if isFinished {
			fmt.Println("APP is staged.")
			// Check instance is in running state
			isRunning, err := r.CheckInstanceStatus("RUNNING", appGUID)
			if isRunning {
				fmt.Println("All Instance  is Running.")
				return &appFields, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &appFields, nil
}

func (r *app) Get(appGUID string) (*AppFields, error) {

	rawURL := fmt.Sprintf("/v2/apps/%s", appGUID)
	appFields := AppFields{}
	_, err := r.client.Get(rawURL, &appFields, nil)
	if err != nil {
		return nil, err
	}
	return &appFields, nil
}

func (r *app) Summary(appGUID string) (*AppSummaryFields, error) {

	rawURL := fmt.Sprintf("/v2/apps/%s/summary", appGUID)
	appFields := AppSummaryFields{}
	_, err := r.client.Get(rawURL, &appFields, nil)
	if err != nil {
		return nil, err
	}
	return &appFields, nil
}

func (r *app) Stat(appGUID string) (map[string]AppStats, error) {

	rawURL := fmt.Sprintf("/v2/apps/%s/stats", appGUID)
	appStats := map[string]AppStats{}
	_, err := r.client.Get(rawURL, &appStats, nil)
	if err != nil {
		return nil, err
	}
	return appStats, nil
}

func (r *app) Instances(appGUID string) (map[string]AppStats, error) {

	rawURL := fmt.Sprintf("/v2/apps/%s/instances", appGUID)
	appInstances := map[string]AppStats{}
	_, err := r.client.Get(rawURL, &appInstances, nil)
	if err != nil {
		return nil, err
	}
	return appInstances, nil
}

func (r *app) List() ([]App, error) {
	rawURL := "v2/apps"
	req := rest.GetRequest(rawURL)
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	apps, err := r.listAppWithPath(path)
	if err != nil {
		return nil, err
	}
	return apps, nil

}

func (r *app) Update(appGUID string, appPayload *AppCreateRequest) (*AppFields, error) {

	rawURL := fmt.Sprintf("/v2/apps/%s", appGUID)
	appFields := AppFields{}
	_, err := r.client.Put(rawURL, appPayload, &appFields)
	if err != nil {
		return nil, err
	}
	return &appFields, nil
}

func (r *app) Delete(appGUID string) error {
	rawURL := fmt.Sprintf("/v2/apps/%s", appGUID)
	_, err := r.client.Delete(rawURL)
	return err
}

func (r *app) CheckAppStatus(waitForThisState, appGUID string) (bool, error) {

	timeout := time.After(10 * time.Minute)
	tick := time.Tick(5 * time.Second)

	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			return false, bmxerror.New(TimeOutError, "Time out while waiting for the app to start")
		// Got a tick, we should check the App package status
		case <-tick:
			apps, err := r.Get(appGUID)
			if err != nil {
				return false, err
			}
			fmt.Println("apps.Entity.PackageState  ===>>> ", apps.Entity.PackageState)
			if apps.Entity.PackageState == waitForThisState {
				return true, nil

			}
		}
	}

}

func (r *app) CheckInstanceStatus(waitForThisState, appGUID string) (bool, error) {

	timeout := time.After(10 * time.Minute)
	tick := time.Tick(5 * time.Second)

	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			return false, bmxerror.New(TimeOutError, "Time out while waiting for the instance to start")
		// Got a tick, we should check the instance status
		case <-tick:
			appStat, err := r.Stat(appGUID)
			if err != nil {
				return false, err
			}
			running := 0
			for k, v := range appStat {
				fmt.Printf("Instance[%s] State is %s", k, v)
				if v.State == waitForThisState {
					running = running + 1
				}
			}
			if running == len(appStat) {
				return true, nil
			}

		}
	}

}
