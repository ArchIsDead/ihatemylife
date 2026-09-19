package bypass

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"s/utils"
)

const ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
const fpRaw = "webgl:ANGLE (Intel, Intel(R) UHD Graphics Direct3D11 vs_5_0 ps_5_0, D3D11)||audio:12.345678901234||canvas:abcdef123||fonts:20/25||system:8CPU,8GB,Win32,0,1,0||env:Asia/Jakarta,1920,1080,24,en-US||network:4g,10,50"

type Jar struct {
	cookies map[string]string
}

func NewJar() *Jar {
	return &Jar{cookies: make(map[string]string)}
}

func (j *Jar) Ingest(resp *http.Response) {
	for _, c := range resp.Cookies() {
		j.cookies[c.Name] = c.Value
	}
}

func (j *Jar) Header() string {
	var parts []string
	for k, v := range j.cookies {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, "; ")
}

func (j *Jar) Get(k string) string {
	return j.cookies[k]
}

type Result struct {
	OK          bool
	Destination string
	Hops        []string
	Error       string
	Elapsed     float64
}

func rawGet(u string, jar *Jar, followRedirect bool) (*http.Response, string, error) {
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	if jar != nil {
		if c := jar.Header(); c != "" {
			req.Header.Set("Cookie", c)
		}
	}

	client := &http.Client{Timeout: 20 * time.Second}
	if !followRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if jar != nil {
		jar.Ingest(resp)
	}

	body, _ := io.ReadAll(resp.Body)
	return resp, string(body), nil
}

func postJSON(u string, payload map[string]interface{}, headers map[string]string, jar *Jar) (map[string]interface{}, error) {
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", u, strings.NewReader(string(b)))
	req.Header.Set("User-Agent", ua)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if c := jar.Header(); c != "" {
		req.Header.Set("Cookie", c)
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jar.Ingest(resp)

	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func makeToken(xsrf string) string {
	fp := sha256.Sum256([]byte(fpRaw))
	fpHex := fmt.Sprintf("%x", fp)
	suffix := "#" + base64.StdEncoding.EncodeToString([]byte(fpHex))

	decoded, _ := url.QueryUnescape(xsrf)
	maxLen := 128 - len(suffix)
	if len(decoded) > maxLen {
		decoded = decoded[:maxLen]
	}
	return decoded + suffix
}

func resolveURL(ref, base string) string {
	u, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	b, err := url.Parse(base)
	if err != nil {
		return ref
	}
	return b.ResolveReference(u).String()
}

func extractOrigin(formAction, pageURL string) string {
	full := resolveURL(formAction, pageURL)
	u, err := url.Parse(full)
	if err != nil {
		return ""
	}
	return u.Scheme + "://" + u.Host
}

func bypass(cur, html string, jar *Jar) (string, error) {
	formRe := regexp.MustCompile(`(?i)<form[^>]+action=["']([^"']+)["']`)
	rayRe := regexp.MustCompile(`(?i)name=["']ray_id["']\s+value=["']([^"']+)["']`)
	aliasRe := regexp.MustCompile(`(?i)name=["']alias["']\s+value=["']([^"']+)["']`)

	formM := formRe.FindStringSubmatch(html)
	rayM := rayRe.FindStringSubmatch(html)
	aliasM := aliasRe.FindStringSubmatch(html)

	if len(formM) < 2 || len(rayM) < 2 || len(aliasM) < 2 {
		return "", fmt.Errorf("missing form data")
	}

	origin := extractOrigin(formM[1], cur)
	if origin == "" {
		return "", fmt.Errorf("no origin")
	}

	redirectURL := fmt.Sprintf("%s/redirect.php?ray_id=%s&alias=%s", origin, url.QueryEscape(rayM[1]), url.QueryEscape(aliasM[1]))

	jar2 := NewJar()
	_, _, err := rawGet(redirectURL, jar2, true)
	if err != nil {
		return "", err
	}

	xsrf := jar2.Get("XSRF-TOKEN")
	if xsrf == "" {
		return "", fmt.Errorf("missing xsrf")
	}

	apiHeaders := map[string]string{
		"Origin":           origin,
		"Referer":          redirectURL,
		"Content-Type":     "application/json",
		"Accept":           "application/json, */*",
		"X-Requested-With": "XMLHttpRequest",
	}

	sessionURL := origin + "/api/session"
	_, _ = postJSON(sessionURL, map[string]interface{}{"_token": makeToken(xsrf)}, apiHeaders, jar2)

	time.Sleep(300 * time.Millisecond)

	verifyURL := origin + "/api/verify"
	vd, _ := postJSON(verifyURL, map[string]interface{}{"_a": 0, "captcha": nil, "passcode": nil}, apiHeaders, jar2)

	target := "/redirect.php"
	if vd != nil {
		if t, ok := vd["target"].(string); ok && t != "" {
			target = t
		}
	}
	if strings.HasPrefix(target, "/") {
		target = origin + target
	}

	_, _, _ = rawGet(target, jar2, true)

	xsrf2 := jar2.Get("XSRF-TOKEN")
	if xsrf2 == "" {
		xsrf2 = xsrf
	}
	_, _ = postJSON(sessionURL, map[string]interface{}{"_token": makeToken(xsrf2)}, apiHeaders, jar2)

	time.Sleep(500 * time.Millisecond)

	key := 500
	size := fmt.Sprintf("%d.%d", (1920+key)*2, (1080+key)*2)
	gd, _ := postJSON(origin+"/api/go", map[string]interface{}{"key": key, "size": size}, apiHeaders, jar2)

	readyURL := ""
	if gd != nil {
		if u, ok := gd["url"].(string); ok {
			readyURL = u
		}
	}
	if readyURL == "" {
		return "", fmt.Errorf("no url")
	}

	if !strings.HasPrefix(readyURL, "http://") && !strings.HasPrefix(readyURL, "https://") {
		readyURL = resolveURL(readyURL, origin)
	}

	_, body, _ := rawGet(readyURL, jar2, true)
	destRe := regexp.MustCompile(`(?i)window\.location\.href\s*=\s*["']([^"']+)["']`)
	dm := destRe.FindStringSubmatch(body)
	if len(dm) > 1 {
		dest := dm[1]
		dest = strings.ReplaceAll(dest, `\/`, "/")
		dest = strings.ReplaceAll(dest, `\u0026`, "&")
		return dest, nil
	}

	return readyURL, nil
}

func DoSFL(rawURL string) *Result {
	t0 := time.Now()
	cur := rawURL
	if !strings.HasPrefix(cur, "http://") && !strings.HasPrefix(cur, "https://") {
		cur = "https://" + cur
	}

	visited := make(map[string]bool)
	var hops []string

	for i := 0; i < 20; i++ {
		if visited[cur] {
			break
		}
		visited[cur] = true
		hops = append(hops, cur)

		jar := NewJar()
		resp, html, err := rawGet(cur, jar, true)
		if err != nil {
			break
		}

		finalURL := ""
		if resp.Request != nil && resp.Request.URL != nil {
			finalURL = resp.Request.URL.String()
		}

		if finalURL != "" && finalURL != cur {
			cur = finalURL
			continue
		}

		status := resp.StatusCode
		if status >= 300 && status < 400 {
			loc := resp.Header.Get("Location")
			if loc != "" {
				cur = resolveURL(loc, cur)
				continue
			}
		}

		if strings.Contains(html, `name="ray_id"`) || strings.Contains(html, "redirect.php") {
			dest, err := bypass(cur, html, jar)
			if err == nil && dest != "" && dest != cur {
				cur = dest
				continue
			}
		}

		metaRe := regexp.MustCompile(`(?i)<meta[^>]+http-equiv=["']refresh["'][^>]+content=["'][^"']*url=([^"']+)["']`)
		if m := metaRe.FindStringSubmatch(html); len(m) > 1 {
			cur = resolveURL(strings.TrimSpace(m[1]), cur)
			continue
		}

		jsRe := regexp.MustCompile(`(?i)window\.location(?:\.href)?\s*=\s*["']([^"']+)["']`)
		if m := jsRe.FindStringSubmatch(html); len(m) > 1 {
			cur = resolveURL(strings.TrimSpace(m[1]), cur)
			continue
		}

		aRe := regexp.MustCompile(`(?i)<a[^>]+href=["']([^"']+)["'][^>]*>\s*(?:continue|lanjut|klik|click|here|di sini)`)
		if m := aRe.FindStringSubmatch(html); len(m) > 1 {
			cur = resolveURL(strings.TrimSpace(m[1]), cur)
			continue
		}

		break
	}

	if cur == rawURL || cur == "https://"+rawURL {
		apiRes, err := DoBypassTools(rawURL)
		if err == nil && apiRes != "" {
			cur = apiRes
		}
	}

	return &Result{
		OK:          true,
		Destination: cur,
		Hops:        hops,
		Elapsed:     elapsed(t0),
	}
}

func DoBypassTools(target string) (string, error) {
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
		return "", fmt.Errorf("bypasstools failed")
	}

	return raw.Result, nil
}

func elapsed(t time.Time) float64 {
	return float64(time.Since(t).Milliseconds()) / 1000.0
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ BYPASS SFL ]")))
	fmt.Println(utils.Div())

	if !r.OK {
		fmt.Println(utils.Red("Error: " + r.Error))
		return
	}

	fmt.Println(utils.Gry("Destination: ") + utils.Wht(r.Destination))
	fmt.Println(utils.Gry(fmt.Sprintf("Hops: %d", len(r.Hops))))
	for i, h := range r.Hops {
		fmt.Println(utils.Gry(fmt.Sprintf("  [%d] %s", i+1, h)))
	}
	fmt.Println(utils.Gry(fmt.Sprintf("Elapsed: %.2fs", r.Elapsed)))
}
