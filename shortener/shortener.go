package shortener

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"s/utils"
)

type Result struct {
	XGD  string
	Uto  string
	Walee string
	EJUZ string
	H1NU string
	Referis string
}

func Do(u, alias string) (*Result, error) {
	if !strings.HasPrefix(u, "https://") {
		return nil, fmt.Errorf("invalid url")
	}
	r := &Result{}

	if x, err := xgd(u, alias); err == nil {
		r.XGD = x
	}
	if x, err := uto(u); err == nil {
		r.Uto = x
	}
	if x, err := walee(u); err == nil {
		r.Walee = x
	}
	if x, err := ejuz(u, alias); err == nil {
		r.EJUZ = x
	}
	if x, err := h1nu(u, alias); err == nil {
		r.H1NU = x
	}
	if x, err := referis(u); err == nil {
		r.Referis = x
	}

	return r, nil
}

func xgd(u, alias string) (string, error) {
	req, _ := http.NewRequest("POST", "https://x.gd/api/V1/auth", nil)
	req.Header.Set("Origin", "https://x.gd")
	req.Header.Set("Referer", "https://x.gd/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var auth map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return "", err
	}
	result, _ := auth["result"].(map[string]interface{})
	s, _ := result["s"].(string)
	shift := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			shift = int(c - '0')
			break
		}
	}
	xacas := resp.Header.Get("xacas")
	decoded := caesarDecode(xacas, shift)
	form := url.Values{
		"url":        {u},
		"shortid":    {alias},
		"analytics":  {"1"},
		"filterbots": {"0"},
	}
	req2, _ := http.NewRequest("POST", "https://x.gd/api/V1/shorten", strings.NewReader(form.Encode()))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req2.Header.Set("Origin", "https://x.gd")
	req2.Header.Set("Referer", "https://x.gd/")
	req2.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	req2.Header.Set("xacas", decoded)
	resp2, err := utils.Hc.Do(req2)
	if err != nil {
		return "", err
	}
	defer resp2.Body.Close()
	var data map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&data); err != nil {
		return "", err
	}
	r, _ := data["result"].(map[string]interface{})
	xid, _ := r["xid"].(string)
	return "https://x.gd/" + xid, nil
}

func caesarDecode(s string, shift int) string {
	var out strings.Builder
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			out.WriteByte(byte((c-'a'-byte(shift)+26)%26 + 'a'))
		} else if c >= 'A' && c <= 'Z' {
			out.WriteByte(byte((c-'A'-byte(shift)+26)%26 + 'A'))
		} else {
			out.WriteRune(c)
		}
	}
	return out.String()
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
	if short, ok := data["shorturl"].(string); ok {
		return short, nil
	}
	if short, ok := data["short"].(string); ok {
		return short, nil
	}
	b2, _ := json.Marshal(data)
	return string(b2), nil
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
	b, _ := io.ReadAll(resp.Body)
	return string(b), nil
}

func ejuz(u, alias string) (string, error) {
	req, _ := http.NewRequest("GET", "https://ej.uz/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return "", err
	}
	body, _ := io.ReadAll(resp.Body)
	cookie := resp.Header.Get("Set-Cookie")
	resp.Body.Close()

	re := regexp.MustCompile(`name="csrf-token" content="([^"]+)"`)
	m := re.FindStringSubmatch(string(body))
	if len(m) < 2 {
		return "", fmt.Errorf("csrf not found")
	}
	csrf := m[1]

	form := url.Values{"smurl": {alias}, "url": {u}}
	req2, _ := http.NewRequest("POST", "https://ej.uz/", strings.NewReader(form.Encode()))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Set("Cookie", cookie)
	req2.Header.Set("csrf-token", csrf)
	req2.Header.Set("Origin", "https://ej.uz")
	req2.Header.Set("Referer", "https://ej.uz/")
	req2.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	resp2, err := utils.Hc.Do(req2)
	if err != nil {
		return "", err
	}
	defer resp2.Body.Close()
	var data map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&data); err != nil {
		return "", err
	}
	result, _ := data["result"].(map[string]interface{})
	smurl, _ := result["smurl"].(string)
	return "https://ej.uz/" + smurl, nil
}

func h1nu(u, alias string) (string, error) {
	form := url.Values{"url": {u}, "keyword": {alias}}
	req, _ := http.NewRequest("POST", "https://h1.nu/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://h1.nu")
	req.Header.Set("Referer", "https://h1.nu/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	re := regexp.MustCompile(`class="short-url" value="([^"]+)"`)
	m := re.FindStringSubmatch(string(body))
	if len(m) < 2 {
		return "", fmt.Errorf("no result")
	}
	return m[1], nil
}

func referis(u string) (string, error) {
	form := url.Values{"source": {"homepage"}, "url": {u}, "action": {"shorten"}}
	req, _ := http.NewRequest("POST", "https://refer.is/_root.data?index", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Origin", "https://refer.is")
	req.Header.Set("Referer", "https://refer.is/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15; SM-F958 Build/AP3A.240905.015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.86 Mobile Safari/537.36")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var data []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if len(data) >= 10 {
		return fmt.Sprintf("%v", data[9]), nil
	}
	return "", fmt.Errorf("no result")
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ SHORTENED URLS ]")))
	fmt.Println(utils.Div())
	if r.XGD != "" {
		fmt.Println(utils.Gry("X.GD: ") + utils.Wht(r.XGD))
	}
	if r.Uto != "" {
		fmt.Println(utils.Gry("U.TO: ") + utils.Wht(r.Uto))
	}
	if r.Walee != "" {
		fmt.Println(utils.Gry("WAL.EE: ") + utils.Wht(r.Walee))
	}
	if r.EJUZ != "" {
		fmt.Println(utils.Gry("EJ.UZ: ") + utils.Wht(r.EJUZ))
	}
	if r.H1NU != "" {
		fmt.Println(utils.Gry("H1.NU: ") + utils.Wht(r.H1NU))
	}
	if r.Referis != "" {
		fmt.Println(utils.Gry("REFER.IS: ") + utils.Wht(r.Referis))
	}
}
