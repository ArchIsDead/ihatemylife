package songfinder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"s/utils"
)

type Track struct {
	Title       string
	Artist      string
	Album       string
	ReleaseDate string
	SpotifyURL  string
	AppleURL    string
}

func Identify(input string) (*Track, error) {
	if input == "" {
		return nil, fmt.Errorf("audio path or url required")
	}

	var data io.Reader
	var filename string

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		resp, err := http.Get(input)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		data = resp.Body
		filename = "sample.mp3"
	} else {
		file, err := os.Open(input)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		data = file
		filename = filepath.Base(input)
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, _ := w.CreateFormFile("file", filename)
	io.Copy(part, data)
	w.WriteField("api_token", "test")
	w.WriteField("return", "spotify,apple_music")
	w.Close()

	client := &http.Client{Timeout: 60 * time.Second}
	req, _ := http.NewRequest("POST", "https://api.audd.io/", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	if raw["status"] != "success" {
		return nil, fmt.Errorf("track not found")
	}

	result, ok := raw["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("track not found")
	}

	track := &Track{
		Title:       str(result["title"]),
		Artist:      str(result["artist"]),
		Album:       str(result["album"]),
		ReleaseDate: str(result["release_date"]),
	}

	if spotify, ok := result["spotify"].(map[string]interface{}); ok {
		if ext, ok := spotify["external_urls"].(map[string]interface{}); ok {
			track.SpotifyURL = str(ext["spotify"])
		}
	}

	if apple, ok := result["apple_music"].(map[string]interface{}); ok {
		track.AppleURL = str(apple["url"])
	}

	return track, nil
}

func (t *Track) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ SONG IDENTIFIED ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Title: ") + utils.Wht(t.Title))
	if t.Artist != "" {
		fmt.Println(utils.Gry("Artist: ") + utils.Wht(t.Artist))
	}
	if t.Album != "" {
		fmt.Println(utils.Gry("Album: ") + utils.Wht(t.Album))
	}
	if t.ReleaseDate != "" {
		fmt.Println(utils.Gry("Release: ") + utils.Wht(t.ReleaseDate))
	}
	if t.SpotifyURL != "" {
		fmt.Println(utils.Gry("Spotify: ") + utils.Wht(t.SpotifyURL))
	}
	if t.AppleURL != "" {
		fmt.Println(utils.Gry("Apple Music: ") + utils.Wht(t.AppleURL))
	}
}

func str(x interface{}) string {
	if x == nil {
		return ""
	}
	s, _ := x.(string)
	return s
}
