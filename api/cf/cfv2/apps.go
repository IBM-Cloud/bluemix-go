package cfv2

import (
	"fmt"
	"os"
	"time"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/helpers"
	"github.com/IBM-Bluemix/bluemix-go/rest"
	"github.com/IBM-Bluemix/bluemix-go/trace"
)

//AppState ...
type AppState struct {
	PackageState  string
	InstanceState string
}

const (
	//ErrCodeAppDoesnotExist ...
	ErrCodeAppDoesnotExist = "AppADoesnotExist"

	//AppRunningState ...
	AppRunningState = "RUNNING"

	//AppStartedState ...
	AppStartedState = "STARTED"

	//AppStagedState ...
	AppStagedState = "STAGED"

	//AppPendingState ...
	AppPendingState = "PENDING"

	//AppStoppedState ...
	AppStoppedState = "STOPPED"

	//AppFailedState ...
	AppFailedState = "FAILED"

	//AppUnKnownState ...
	AppUnKnownState = "UNKNOWN"

	//DefaultRetryDelayForStatusCheck ...
	DefaultRetryDelayForStatusCheck = 10 * time.Second
)

//AppRequest ...
type AppRequest struct {
	Name                     *string                 `json:"name,omitempty"`
	Memory                   int                     `json:"memory,omitempty"`
	Instances                int                     `json:"instances,omitempty"`
	DiskQuota                int                     `json:"disk_quota,omitempty"`
	SpaceGUID                *string                 `json:"space_guid,omitempty"`
	StackGUID                *string                 `json:"stack_guid,omitempty"`
	State                    *string                 `json:"state,omitempty"`
	DetectedStartCommand     *string                 `json:"detected_start_command,omitempty"`
	Command                  *string                 `json:"command,omitempty"`
	BuildPack                *string                 `json:"buildpack,omitempty"`
	HealthCheckType          *string                 `json:"health_check_type,omitempty"`
	HealthCheckTimeout       int                     `json:"health_check_timeout,omitempty"`
	Diego                    bool                    `json:"diego,omitempty"`
	EnableSSH                bool                    `json:"enable_ssh,omitempty"`
	DockerImage              *string                 `json:"docker_image,omitempty"`
	StagingFailedReason      *string                 `json:"staging_failed_reason,omitempty"`
	StagingFailedDescription *string                 `json:"staging_failed_description,omitempty"`
	Ports                    *[]int                  `json:"ports,omitempty"`
	DockerCredentialsJSON    *map[string]interface{} `json:"docker_credentials_json,omitempty"`
	EnvironmentJSON          *map[string]interface{} `json:"environment_json,omitempty"`
}

//AppEntity ...
type AppEntity struct {
	Name                     string                 `json:"name"`
	SpaceGUID                string                 `json:"space_guid"`
	StackGUID                string                 `json:"stack_guid"`
	State                    string                 `json:"state"`
	PackageState             string                 `json:"package_state"`
	Memory                   int                    `json:"memory"`
	Instances                int                    `json:"instances"`
	DiskQuota                int                    `json:"disk_quota"`
	Version                  string                 `json:"version"`
	BuildPack                *string                `json:"buildpack"`
	Command                  *string                `json:"command"`
	Console                  bool                   `json:"console"`
	Debug                    *string                `json:"debug"`
	StagingTaskID            string                 `json:"staging_task_id"`
	HealthCheckType          string                 `json:"health_check_type"`
	HealthCheckTimeout       *int                   `json:"health_check_timeout"`
	StagingFailedReason      string                 `json:"staging_failed_reason"`
	StagingFailedDescription string                 `json:"staging_failed_description"`
	Diego                    bool                   `json:"diego"`
	DockerImage              *string                `json:"docker_image"`
	EnableSSH                bool                   `json:"enable_ssh"`
	Ports                    []int                  `json:"ports"`
	DockerCredentialsJSON    map[string]interface{} `json:"docker_credentials_json"`
	EnvironmentJSON          map[string]interface{} `json:"environment_json"`
}

//AppResource ...
type AppResource struct {
	Resource
	Entity AppEntity
}

//AppFields ...
type AppFields struct {
	Metadata Metadata
	Entity   AppEntity
}

//UploadBitsEntity ...
type UploadBitsEntity struct {
	GUID   string `json:"guid"`
	Status string `json:"status"`
}

//UploadBitFields ...
type UploadBitFields struct {
	Metadata Metadata
	Entity   UploadBitsEntity
}

//AppSummaryFields ...
type AppSummaryFields struct {
	GUID             string `json:"guid"`
	Name             string `json:"name"`
	State            string `json:"state"`
	PackageState     string `json:"package_state"`
	RunningInstances int    `json:"running_instances"`
}

//AppStats ...
type AppStats struct {
	State string `json:"state"`
}

//ToFields ..
func (resource AppResource) ToFields() App {
	entity := resource.Entity

	return App{
		GUID:               resource.Metadata.GUID,
		Name:               entity.Name,
		SpaceGUID:          entity.SpaceGUID,
		StackGUID:          entity.StackGUID,
		State:              entity.State,
		PackageState:       entity.PackageState,
		Memory:             entity.Memory,
		Instances:          entity.Instances,
		DiskQuota:          entity.DiskQuota,
		Version:            entity.Version,
		BuildPack:          entity.BuildPack,
		Command:            entity.Command,
		Console:            entity.Console,
		Debug:              entity.Debug,
		StagingTaskID:      entity.StagingTaskID,
		HealthCheckType:    entity.HealthCheckType,
		HealthCheckTimeout: entity.HealthCheckTimeout,
		Diego:              entity.Diego,
		DockerImage:        entity.DockerImage,
		EnableSSH:          entity.EnableSSH,
		Ports:              entity.Ports,
		DockerCredentialsJSON: entity.DockerCredentialsJSON,
		EnvironmentJSON:       entity.EnvironmentJSON,
	}
}

//App model
type App struct {
	Name                  string
	SpaceGUID             string
	GUID                  string
	StackGUID             string
	State                 string
	PackageState          string
	Memory                int
	Instances             int
	DiskQuota             int
	Version               string
	BuildPack             *string
	Command               *string
	Console               bool
	Debug                 *string
	StagingTaskID         string
	HealthCheckType       string
	HealthCheckTimeout    *int
	Diego                 bool
	DockerImage           *string
	EnableSSH             bool
	Ports                 []int
	DockerCredentialsJSON map[string]interface{}
	EnvironmentJSON       map[string]interface{}
}

//Apps ...
type Apps interface {
	Create(appPayload AppRequest) (*AppFields, error)
	List() ([]App, error)
	Get(appGUID string) (*AppFields, error)
	Update(appGUID string, appPayload AppRequest) (*AppFields, error)
	Delete(appGUID string, async bool, recursive bool) error
	FindByName(spaceGUID, name string) (*App, error)
	Start(appGUID string, timeout time.Duration) (*AppState, error)
	Upload(path string, name string) (*UploadBitFields, error)
	Summary(appGUID string) (*AppSummaryFields, error)
	Stat(appGUID string) (map[string]AppStats, error)
	WaitForAppStatus(waitForThisState, appGUID string, timeout time.Duration) (string, error)
	WaitForInstanceStatus(waitForThisState, appGUID string, timeout time.Duration) (string, error)
	Instances(appGUID string) (map[string]AppStats, error)
	Restage(appGUID string, timeout time.Duration) (*AppState, error)
	WaitForStatus(appGUID string, maxWaitTime time.Duration) (*AppState, error)

	//Routes related
	BindRoute(appGUID, routeGUID string) (*AppFields, error)
	ListRoutes(appGUID string) ([]Route, error)
	UnBindRoute(appGUID, routeGUID string) error

	//Service bindings
	ListServiceBindings(appGUID string) ([]ServiceBinding, error)
	DeleteServiceBindings(appGUID string, bindingGUIDs ...string) error
}

type app struct {
	client *client.Client
}

func newAppAPI(c *client.Client) Apps {
	return &app{
		client: c,
	}
}

func (r *app) FindByName(spaceGUID string, name string) (*App, error) {
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
	if len(apps) == 0 {
		return nil, bmxerror.New(ErrCodeAppDoesnotExist,
			fmt.Sprintf("Given app:  %q doesn't exist in given space: %q", name, spaceGUID))

	}
	return &apps[0], nil
}

func (r *app) Create(appPayload AppRequest) (*AppFields, error) {
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

func (r *app) ListRoutes(appGUID string) ([]Route, error) {
	rawURL := fmt.Sprintf("/v2/apps/%s/routes", appGUID)
	req := rest.GetRequest(rawURL)
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	route, err := listRouteWithPath(r.client, path)
	if err != nil {
		return nil, err
	}
	return route, nil
}

func (r *app) UnBindRoute(appGUID, routeGUID string) error {
	rawURL := fmt.Sprintf("/v2/apps/%s/routes/%s", appGUID, routeGUID)
	_, err := r.client.Delete(rawURL)
	return err
}

func (r *app) DeleteServiceBindings(appGUID string, sbGUIDs ...string) error {
	for _, g := range sbGUIDs {
		rawURL := fmt.Sprintf("/v2/apps/%s/service_bindings/%s", appGUID, g)
		_, err := r.client.Delete(rawURL)
		return err
	}
	return nil
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

func (r *app) Upload(appGUID string, zipPath string) (*UploadBitFields, error) {
	req := rest.PutRequest(r.client.URL("/v2/apps/"+appGUID+"/bits")).Query("async", "false")
	file, err := os.Open(zipPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	f := rest.File{
		Name:    file.Name(),
		Content: file,
	}
	req.File("application", f)
	req.Field("resources", "[]")
	uploadBitResponse := &UploadBitFields{}
	_, err = r.client.SendRequest(req, uploadBitResponse)
	return uploadBitResponse, err
}

func (r *app) Start(appGUID string, maxWaitTime time.Duration) (*AppState, error) {
	payload := AppRequest{
		State: helpers.String(AppStartedState),
	}
	rawURL := fmt.Sprintf("/v2/apps/%s", appGUID)
	appFields := AppFields{}
	_, err := r.client.Put(rawURL, payload, &appFields)
	if err != nil {
		return nil, err
	}
	appState := &AppState{
		PackageState:  AppPendingState,
		InstanceState: AppUnKnownState,
	}
	if maxWaitTime == 0 {
		appState.PackageState = appFields.Entity.PackageState
		appState.InstanceState = appFields.Entity.State
		return appState, nil
	}
	return r.WaitForStatus(appGUID, maxWaitTime)

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

func (r *app) Update(appGUID string, appPayload AppRequest) (*AppFields, error) {
	rawURL := fmt.Sprintf("/v2/apps/%s", appGUID)
	appFields := AppFields{}
	_, err := r.client.Put(rawURL, appPayload, &appFields)
	if err != nil {
		return nil, err
	}
	return &appFields, nil
}

func (r *app) Delete(appGUID string, async bool, recursive bool) error {
	req := rest.GetRequest(fmt.Sprintf("/v2/apps/%s", appGUID))
	if async {
		req.Query("async", "true")
	}
	if recursive {
		req.Query("recursive", "true")
	}
	httpReq, err := req.Build()
	if err != nil {
		return err
	}
	path := httpReq.URL.String()
	_, err = r.client.Delete(path)
	return err
}

func (r *app) Restage(appGUID string, maxWaitTime time.Duration) (*AppState, error) {
	rawURL := fmt.Sprintf("/v2/apps/%s/restage", appGUID)
	appFields := AppFields{}
	_, err := r.client.Post(rawURL, nil, &appFields)
	if err != nil {
		return nil, err
	}
	appState := &AppState{
		PackageState:  AppPendingState,
		InstanceState: AppUnKnownState,
	}
	if maxWaitTime == 0 {
		appState.PackageState = appFields.Entity.PackageState
		appState.InstanceState = appFields.Entity.State
		return appState, nil
	}
	return r.WaitForStatus(appGUID, maxWaitTime)

}

func (r *app) WaitForAppStatus(waitForThisState, appGUID string, maxWaitTime time.Duration) (string, error) {
	timeout := time.After(maxWaitTime)
	tick := time.Tick(DefaultRetryDelayForStatusCheck)
	status := AppPendingState
	for {
		select {
		case <-timeout:
			trace.Logger.Printf("Timed out while checking the app status for %q.  Waited for %q for the state to be %q", appGUID, maxWaitTime, waitForThisState)
			return status, nil
		case <-tick:
			appFields, err := r.Get(appGUID)
			if err != nil {
				return "", err
			}
			status = appFields.Entity.PackageState
			trace.Logger.Println("apps.Entity.PackageState  ===>>> ", status)
			if status == waitForThisState || status == AppFailedState {
				return status, nil
			}
		}
	}
}

func (r *app) WaitForInstanceStatus(waitForThisState, appGUID string, maxWaitTime time.Duration) (string, error) {
	timeout := time.After(maxWaitTime)
	tick := time.Tick(DefaultRetryDelayForStatusCheck)
	status := AppStartedState
	for {
		select {
		case <-timeout:
			trace.Logger.Printf("Timed out while checking the app status for %q. Waited for %q for the state to be %q", appGUID, maxWaitTime, waitForThisState)
			return status, nil
		case <-tick:
			appStat, err := r.Stat(appGUID)
			if err != nil {
				return status, err
			}
			stateCount := 0
			for k, v := range appStat {
				fmt.Printf("Instance[%s] State is %s", k, v)
				if v.State == waitForThisState {
					stateCount++
				}
			}
			if stateCount == len(appStat) {
				return waitForThisState, nil
			}

		}
	}

}

func (r *app) WaitForStatus(appGUID string, maxWaitTime time.Duration) (*AppState, error) {
	appState := &AppState{
		PackageState:  AppPendingState,
		InstanceState: AppUnKnownState,
	}
	status, err := r.WaitForAppStatus(AppStagedState, appGUID, maxWaitTime/2)
	appState.PackageState = status
	if err != nil || status == AppFailedState {
		return appState, err
	}
	status, err = r.WaitForInstanceStatus(AppRunningState, appGUID, maxWaitTime/2)
	appState.InstanceState = status
	return appState, nil
}

//TODO pull the wait logic in a auxiliary function which can be used by all

func (r *app) ListServiceBindings(appGUID string) ([]ServiceBinding, error) {
	rawURL := fmt.Sprintf("/v2/apps/%s/service_bindings", appGUID)
	req := rest.GetRequest(rawURL)
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	sb, err := listServiceBindingWithPath(r.client, path)
	if err != nil {
		return nil, err
	}
	return sb, nil
}
