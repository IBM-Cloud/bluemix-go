package k8sclusterv1

import (
	"fmt"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
)

//Worker ...
type Worker struct {
	Billing      string
	ErrorMessage string
	ID           string
	Isolation    string
	KubeVersion  string
	MachineType  string
	PrivateIP    string
	PrivateVlan  string
	PublicIP     string
	PublicVlan   string
	State        string
	Status       string
}

//WorkerParam ...
type WorkerParam struct {
	Action string
	Count  int
}

//Workers ...
type Workers interface {
	List(clusterName string, target *ClusterTargetHeader) ([]Worker, error)
	Get(clusterName string, target *ClusterTargetHeader) (Worker, error)
	Add(clusterName string, params WorkerParam, target *ClusterTargetHeader) error
	Delete(clusterName string, workerD string, target *ClusterTargetHeader) error
	Update(clusterName string, workerID string, params WorkerParam, target *ClusterTargetHeader) error
}

type worker struct {
	client *ClusterClient
	config *bluemix.Config
}

func newWorkerAPI(c *ClusterClient) Workers {
	return &worker{
		client: c,
		config: c.config,
	}
}

func (r *worker) Get(id string, target *ClusterTargetHeader) (Worker, error) {
	rawURL := fmt.Sprintf("/v1/workers/%s", id)
	worker := Worker{}
	_, err := r.client.get(rawURL, &worker, target)
	if err != nil {
		return worker, err
	}

	return worker, err
}

func (r *worker) Add(name string, params WorkerParam, target *ClusterTargetHeader) error {
	rawURL := fmt.Sprintf("/v1/clusters/%s/workers", name)
	_, err := r.client.post(rawURL, params, nil, target)
	return err
}

//DeleteWorker ...
func (r *worker) Delete(name string, workerID string, target *ClusterTargetHeader) error {
	rawURL := fmt.Sprintf("/v1/clusters/%s/workers/%s", name, workerID)
	_, err := r.client.delete(rawURL, target)
	return err
}

//UpdateWorker ...
func (r *worker) Update(name string, workerID string, params WorkerParam, target *ClusterTargetHeader) error {
	rawURL := fmt.Sprintf("/v1/clusters/%s/workers/%s", name, workerID)
	_, err := r.client.put(rawURL, params, nil, target)
	return err
}

func (r *worker) List(name string, target *ClusterTargetHeader) ([]Worker, error) {
	rawURL := fmt.Sprintf("/v1/clusters/%s/workers", name)
	workers := []Worker{}
	_, err := r.client.get(rawURL, &workers, target)
	if err != nil {
		return nil, err
	}
	return workers, err
}
