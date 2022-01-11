package kubernetes

import (
	"github.com/herlon214/iptv-recording/recording"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
)

var KubernetesCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Starts a metacontroller to manage recording in a kubernetes level dealing with the CRs",
	Long:  "Checkout how metacontroller works at https://metacontroller.github.io/",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	log.Println("Starting controller at port 8080...")
	http.HandleFunc("/sync", syncHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Controller struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              recording.Item   `json:"spec"`
	Status            ControllerStatus `json:"status"`
}

type ControllerStatus struct {
	Replicas  int `json:"replicas"`
	Succeeded int `json:"succeeded"`
}

type SyncRequest struct {
	Parent   Controller          `json:"parent"`
	Children SyncRequestChildren `json:"children"`
}

type SyncRequestChildren struct {
	Pods map[string]*v1.Pod `json:"Pod.v1"`
}

type SyncResponse struct {
	Status   ControllerStatus `json:"status"`
	Children []runtime.Object `json:"children"`
}

func sync(request *SyncRequest) (*SyncResponse, error) {
	response := &SyncResponse{}

	// Compute status based on latest observed state.
	for _, pod := range request.Children.Pods {
		response.Status.Replicas++
		if pod.Status.Phase == v1.PodSucceeded {
			response.Status.Succeeded++
		}
	}

	// Generate desired children.
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: request.Parent.Name,
		},
		Spec: v1.PodSpec{
			RestartPolicy: v1.RestartPolicyOnFailure,
			Containers: []v1.Container{
				{
					Name:    "hello",
					Image:   "busybox",
					Command: []string{"echo", request.Parent.Spec.Name},
				},
			},
		},
	}
	response.Children = append(response.Children, pod)

	return response, nil
}

func syncHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	request := &SyncRequest{}
	if err := json.Unmarshal(body, request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := sync(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err = json.Marshal(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
