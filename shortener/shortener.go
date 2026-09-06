package shortener

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"s/utils"
)

type Result struct {
	Uto   string
	Walee string
}

func Do(u, service string) (*Result, error) {
	if !strings.HasPrefix(u, "https://") {
		return nil, fmt.Errorf("invalid url")
	}
	r := &Result{}

	switch service {
	case "all":
		if x, err := uto(u); err == nil {
			r.Uto = x
		}
		if x, err := walee(u); err == nil {
			r.Walee = x
		}
	case "uto":
		if x, err := uto(u); err == nil {
			r.Uto = x
		}
	case "walee":
		if x, err := walee(u); err == nil {
			r.Walee = x
		}
	}

	return r, nil
}

func uto(u string) (string, error) {
	body := map[string]string{"url": u}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://u.to/api/shorten/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://u.to")
	req.Header.Set("Referer", "https://u.to/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if short, ok := data["shortUrl"].(string); ok {
		return short, nil
	}
	return "", fmt.Errorf("no result")
}

func walee(u string) (string, error) {
	form := url.Values{"url": {u}}
	req, _ := http.NewRequest("POST", "https://wal.ee/proxy.php", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://wal.ee")
	req.Header.Set("Referer", "https://wal.ee/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if short, ok := data["shorturl"].(string); ok {
		return short, nil
	}
	return "", fmt.Errorf("no result")
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ SHORTENED URLS ]")))
	fmt.Println(utils.Div())
	if r.Uto != "" {
		fmt.Println(utils.Gry("U.TO: ") + utils.Wht(r.Uto))
	}
	if r.Walee != "" {
		fmt.Println(utils.Gry("WAL.EE: ") + utils.Wht(r.Walee))
	}
}
