package process

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type Recording struct {
	cmd *exec.Cmd

	stdin io.WriteCloser
}

func New(URL string, output string) (*Recording, error) {
	// Apply $date var
	output = strings.Replace(output, "$date", time.Now().Format("2006-01-02 15:04"), -1)

	cmd := exec.Command("ffmpeg", "-reconnect", "1", "-reconnect_delay_max", "5", "-i", URL, "-map", "0", "-codec:", "copy", "-f", "mpegts", fmt.Sprintf("%s.mp4", output))
	stdin, err := cmd.StdinPipe()
	if nil != err {
		return nil, err
	}

	return &Recording{
		cmd:   cmd,
		stdin: stdin,
	}, nil
}

func (rp *Recording) Start() error {
	return rp.cmd.Start()
}

func (rp *Recording) Stop() error {
	return rp.cmd.Process.Kill()
}
