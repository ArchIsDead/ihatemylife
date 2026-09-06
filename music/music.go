package music

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var cmd *exec.Cmd

func getCacheDir() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, "suic1.de", "music")
	os.MkdirAll(dir, 0755)
	return dir
}

func getCachedPath(url string) string {
	name := ""
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		name = parts[len(parts)-1]
	}
	if name == "" {
		name = "music.mp3"
	}
	return filepath.Join(getCacheDir(), name)
}

func download(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func Play(path string) {
	Stop()

	actualPath := path
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		cached := getCachedPath(path)
		if _, err := os.Stat(cached); os.IsNotExist(err) {
			fmt.Println("Downloading music...")
			if err := download(path, cached); err != nil {
				fmt.Println("Download failed:", err)
				return
			}
		}
		actualPath = cached
	}

	if _, err := os.Stat(actualPath); os.IsNotExist(err) {
		fmt.Println("Music file not found:", actualPath)
		return
	}

	cmd = exec.Command("mpv", "--no-video", "--loop=inf", actualPath)
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
