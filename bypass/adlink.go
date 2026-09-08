package bypass

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"s/utils"
)

func BypassTools(target string) (string, error) {
	encoded := url.QueryEscape(target)
	u := "https://r4-api.my.id/bypasstools?url=" + encoded

	client := &http.Client{Timeout: 90 * time.Second}
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw struct {
		Status bool   `json:"status"`
		Result string `json:"result"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return "", err
	}

	if !raw.Status || raw.Result == "" {
		return "", fmt.Errorf("bypass failed")
	}

	return raw.Result, nil
}

func BypassLink(target string) (string, error) {
	encoded := url.QueryEscape(target)
	u := "https://r4-api.my.id/bypasslink?url=" + encoded

	client := &http.Client{Timeout: 60 * time.Second}
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw struct {
		Success bool   `json:"success"`
		Direct  string `json:"direct"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return "", err
	}

	if !raw.Success || raw.Direct == "" {
		return "", fmt.Errorf("bypass failed")
	}

	return raw.Direct, nil
}

func Bypass(target string, method string) (string, error) {
	switch method {
	case "1":
		return BypassTools(target)
	case "2":
		return BypassLink(target)
	default:
		return BypassTools(target)
	}
}

func Show(target, result, method string) {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ BYPASS RESULT ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Method: ") + utils.Wht(method))
	fmt.Println(utils.Gry("Original: ") + utils.Wht(target))
	fmt.Println(utils.Gry("Direct: ") + utils.Wht(result))
}
