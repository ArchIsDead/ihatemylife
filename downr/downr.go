package downr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"s/utils"
)

type Result struct {
	URL      string
	Title    string
	Type     string
	Media    []Media
}

type Media struct {
	Quality string
	URL     string
	Size    string
}

func Download(u string) (*Result, error) {
	if !strings.HasPrefix(u, "https://") {
		return nil, errors.New("invalid url")
	}

	req, _ := http.NewRequest("GET", "https://downr.org/.netlify/functions/analytics", nil)
	req.Header.Set("Referer", "https://downr.org/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	cookie := resp.Header.Get("Set-Cookie")
	resp.Body.Close()

	body := map[string]string{"url": u}
	b, _ := json.Marshal(body)
	req2, _ := http.NewRequest("POST", "https://downr.org/.netlify/functions/nyt", bytes.NewReader(b))
	req2.Header.Set("Accept", "*/*")
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Cookie", cookie)
	req2.Header.Set("Origin", "https://downr.org")
	req2.Header.Set("Referer", "https://downr.org/")
	req2.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	resp2, err := utils.Hc.Do(req2)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&raw); err != nil {
		return nil, err
	}

	r := &Result{URL: u}

	if data, ok := raw["data"].(map[string]interface{}); ok {
		r.Title = str(data["title"])
		r.Type = str(data["type"])

		if medias, ok := data["media"].([]interface{}); ok {
			for _, m := range medias {
				mm, _ := m.(map[string]interface{})
				r.Media = append(r.Media, Media{
					Quality: str(mm["quality"]),
					URL:     str(mm["url"]),
					Size:    str(mm["size"]),
				})
			}
		}
	}

	return r, nil
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ DOWNLOADER ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("URL: ") + utils.Wht(r.URL))
	if r.Title != "" {
		fmt.Println(utils.Gry("Title: ") + utils.Wht(r.Title))
	}
	if r.Type != "" {
		fmt.Println(utils.Gry("Type: ") + utils.Wht(r.Type))
	}
	fmt.Println(utils.Div())

	if len(r.Media) == 0 {
		fmt.Println(utils.Gry("No media found"))
		return
	}

	for i, m := range r.Media {
		fmt.Println(utils.Gry(fmt.Sprintf("[%d] %s", i+1, m.Quality)))
		if m.Size != "" {
			fmt.Println(utils.Gry("Size: ") + utils.Wht(m.Size))
		}
		fmt.Println(utils.Gry("URL: ") + utils.Wht(m.URL))
		fmt.Println()
	}
}

func str(x interface{}) string {
	if x == nil {
		return ""
	}
	s, _ := x.(string)
	return s
}
