package recording

import (
	"time"

	"github.com/herlon214/iptv-recording/process"

	cron "github.com/robfig/cron/v3"
)

type Item struct {
	Name     string        `yaml:"name"`
	URL      string        `yaml:"url"`
	FileName string        `yaml:"fileName"`
	Schedule string        `yaml:"schedule"` // Cron
	Duration time.Duration `yaml:"duration"`
	Folder   string        `yaml:"folder"`

	process *process.Recording
}

func (i *Item) ShouldRun() (bool, error) {
	schedule, err := cron.ParseStandard(i.Schedule)
	if err != nil {
		return false, err
	}

	next := schedule.Next(time.Now().Add(i.Duration * -1))

	if time.Now().After(next) && time.Now().Before(next.Add(i.Duration)) {
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
