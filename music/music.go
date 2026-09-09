package music

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var cmd *exec.Cmd
var volume = 100

func getCacheDir() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, "suic1.de", "music")
	os.MkdirAll(dir, 0755)
	return dir
}

func getCachedPath(url string) string {
	parts := strings.Split(url, "/")
	name := "music.mp3"
	if len(parts) > 0 && parts[len(parts)-1] != "" {
		name = parts[len(parts)-1]
	}
	return filepath.Join(getCacheDir(), name)
}

func download(url, dest string) error {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
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

func SetVolume(vol int) {
	if vol < 0 {
		vol = 0
	}
	if vol > 100 {
		vol = 100
	}
	volume = vol
}

func GetVolume() int {
	return volume
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

	volArg := strconv.Itoa(volume)
	cmd = exec.Command("mpv", "--no-video", "--loop=inf", "--volume="+volArg, actualPath)
	cmd.Stdout = nil
	cmd.Stderr = nil
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
