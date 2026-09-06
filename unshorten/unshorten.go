package unshorten

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"s/utils"
)

func Do(u string) (string, error) {
	if !strings.HasPrefix(u, "https://") {
		return "", fmt.Errorf("url required")
	}

	req, _ := http.NewRequest("GET", "https://unshorten.it/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return "", err
	}
	body, _ := io.ReadAll(resp.Body)
	cookie := resp.Header.Get("Set-Cookie")
	resp.Body.Close()

	re := regexp.MustCompile(`name="csrfmiddlewaretoken" value="([^"]+)"`)
	m := re.FindStringSubmatch(string(body))
	if len(m) < 2 {
		return "", fmt.Errorf("csrf token not found")
	}
	token := m[1]

	form := url.Values{
		"short-url":             {u},
		"csrfmiddlewaretoken":   {token},
	}

	req2, _ := http.NewRequest("POST", "https://unshorten.it/main/get_long_url", strings.NewReader(form.Encode()))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req2.Header.Set("Cookie", cookie)
	req2.Header.Set("Origin", "https://unshorten.it")
	req2.Header.Set("Referer", "https://unshorten.it/")
	req2.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	req2.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp2, err := utils.Hc.Do(req2)
	if err != nil {
		return "", err
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)
	var raw map[string]interface{}
	if err := json.Unmarshal(body2, &raw); err != nil {
		return "", err
	}
	longURL, _ := raw["long_url"].(string)
	if longURL == "" {
		return "", fmt.Errorf("no result")
	}
	return longURL, nil
}
