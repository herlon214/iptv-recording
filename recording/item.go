package recording

import (
	"time"

	"github.com/herlon214/iptv-recording/process"

	cron "github.com/robfig/cron/v3"
)

type Item struct {
	Name     string `yaml:"name" json:"name"`
	URL      string `yaml:"url" json:"url"`
	FileName string `yaml:"fileName" json:"fileName"`
	Schedule string `yaml:"schedule" json:"schedule"` // Cron
	Duration string `yaml:"duration" json:"duration"`
	Folder   string `yaml:"folder" json:"folder"`
	HostPath string `json:"hostPath"`

	process *process.Recording
}

func (i *Item) ShouldRun() (bool, error) {
	// Parse cron
	schedule, err := cron.ParseStandard(i.Schedule)
	if err != nil {
		return false, err
	}

	// Parse duration
	dur, err := time.ParseDuration(i.Duration)
	if err != nil {
		return false, err
	}

	next := schedule.Next(time.Now().Add(dur * -1))

	if time.Now().After(next) && time.Now().Before(next.Add(dur)) {
		return true, nil
	}

	return false, nil
}

func (i *Item) SetProcess(proc *process.Recording) {
	i.process = proc
}

func (i *Item) Process() *process.Recording {
	return i.process
}
