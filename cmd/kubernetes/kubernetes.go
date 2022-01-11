package kubernetes

import (
	"fmt"
	"github.com/herlon214/iptv-recording/recording"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

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

	// Check if should be running
	shouldRun, err := request.Parent.Spec.ShouldRun()
	if err != nil {
		log.Printf("Error parsing cron: %s", err.Error())
		return nil, err
	}

	// Generate desired children.
	if shouldRun {
		// Parse date in the filename
		output := fmt.Sprintf("/output/%s/%s", request.Parent.Spec.Folder, request.Parent.Spec.FileName)
		output = strings.Replace(output, "$date", time.Now().Format("2006-01-02.1504"), -1)

		pod := &v1.Pod{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: fmt.Sprintf("recording-%s", request.Parent.Name),
			},
			Spec: v1.PodSpec{
				RestartPolicy: v1.RestartPolicyOnFailure,
				Containers: []v1.Container{
					{
						Name:    "recording",
						Image:   "ghcr.io/herlon214/iptv-recording:v0.3.4",
						Command: []string{"ffmpeg", "-reconnect", "1", "-reconnect_delay_max", "5", "-i", request.Parent.Spec.URL, "-map", "0", "-codec:", "copy", "-f", "mpegts", fmt.Sprintf("%s.mp4", output)},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "output",
								MountPath: "/output",
							},
						},
					},
				},
				Volumes: []v1.Volume{
					{
						Name: "output",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{
								Path: request.Parent.Spec.HostPath,
							},
						},
					},
				},
			},
		}

		response.Children = append(response.Children, pod)
	}

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
