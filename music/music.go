package music

import (
	"os/exec"
)

var cmd *exec.Cmd

func Play(path string) {
	Stop()
	cmd = exec.Command("mpv", "--no-video", "--loop=inf", path)
	cmd.Start()
}

func Stop() {
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
		cmd = nil
	}
}

func IsPlaying() bool {
	return cmd != nil && cmd.Process != nil
}
