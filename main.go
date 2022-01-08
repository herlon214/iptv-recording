package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/herlon214/iptv-recording/process"
	"github.com/herlon214/iptv-recording/recording"

	"gopkg.in/yaml.v2"
)

var Running = make([]string, 0)

func main() {
	data, err := ioutil.ReadFile("recording.yaml")
	if err != nil {
		panic(err)
	}

	var config recording.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	log.Printf("Found %d items to record \n", len(config.Items))
	log.Println("---------------------")

	for {
		for _, item := range config.Items {
			shouldRun, err := item.ShouldRun()
			if err != nil {
				log.Printf("Error parsing cron: %s", err.Error())
				continue
			}

			if !shouldRun {
				if isRunning(item.Name) {
					log.Println("---> Stop recording")
					if err := item.Process().Stop(); err != nil {
						log.Println("Failed to stop recording:", err.Error())
						continue
					}

					item.SetProcess(nil)
					removeRunning(item.Name)
				}

				continue
			}

			log.Printf("%s [%s] -> [%s] > live? %t", item.Name, item.Schedule, item.Duration, shouldRun)

			// Check if it's running already
			if !isRunning(item.Name) {
				proc, err := process.New(item.URL, fmt.Sprintf("%s/%s", item.Folder, item.FileName))
				if err != nil {
					log.Printf("Error creating process: %s", err.Error())
					continue
				}

				// Set current process
				item.SetProcess(proc)

				log.Println("-_-> Start recording", item.Name)
				if err := proc.Start(); err != nil {
					log.Println("Failed to start recording:", err.Error())
					continue
				}

				// Set running
				Running = append(Running, item.Name)
			}
		}

		time.Sleep(time.Second)
	}

}

func isRunning(name string) bool {
	for _, n := range Running {
		if n == name {
			return true
		}
	}

	return false
}

func removeRunning(name string) {
	current := make([]string, 0)

	for _, n := range Running {
		if n != name {
			current = append(current, n)
		}
	}

	Running = current
}
