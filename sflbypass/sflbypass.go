package sflbypass

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
	OK          bool     `json:"ok"`
	Destination string   `json:"destination,omitempty"`
	Hops        []string `json:"hops,omitempty"`
	Error       string   `json:"error,omitempty"`
	Elapsed     float64  `json:"elapsed_sec"`
}

func Do(rawURL string) *Result {
	t0 := time.Now()
	cur := rawURL
	if !strings.HasPrefix(cur, "http://") && !strings.HasPrefix(cur, "https://") {
		cur = "https://" + cur
	}

	visited := make(map[string]bool)
	var hops []string

	for i := 0; i < 15; i++ {
		if visited[cur] {
			break
		}
		visited[cur] = true
		hops = append(hops, cur)

		jar := NewJar()
		resp, html, err := fetch(cur, jar)
		if err != nil {
			return &Result{OK: false, Error: err.Error(), Hops: hops, Elapsed: elapsed(t0)}
		}

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
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

		re := regexp.MustCompile(`(?i)<meta[^>]+http-equiv=["']refresh["'][^>]+content=["'][^"']*url=([^"']+)["']`)
		m := re.FindStringSubmatch(html)
		if len(m) > 1 {
			cur = resolveURL(strings.TrimSpace(m[1]), cur)
			continue
		}

		break
	}

	return &Result{
		OK:          true,
		Destination: cur,
		Hops:        hops,
		Elapsed:     elapsed(t0),
	}
}

func fetch(u string, jar *Jar) (*http.Response, string, error) {
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	if c := jar.Header(); c != "" {
		req.Header.Set("Cookie", c)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	jar.Ingest(resp)

	body, _ := io.ReadAll(resp.Body)
	return resp, string(body), nil
}

func bypass(pageURL, html string, jar *Jar) (string, error) {
	formRe := regexp.MustCompile(`(?i)<form[^>]+action=["']([^"']+)["']`)
	rayRe := regexp.MustCompile(`(?i)name=["']ray_id["']\s+value=["']([^"']+)["']`)
	aliasRe := regexp.MustCompile(`(?i)name=["']alias["']\s+value=["']([^"']+)["']`)

	formM := formRe.FindStringSubmatch(html)
	rayM := rayRe.FindStringSubmatch(html)
	aliasM := aliasRe.FindStringSubmatch(html)

	if len(formM) < 2 || len(rayM) < 2 || len(aliasM) < 2 {
		return "", fmt.Errorf("missing form data")
	}

	form := formM[1]
	rayID := rayM[1]
	alias := aliasM[1]

	origin := resolveURL(form, pageURL)
	idx := strings.Index(origin[8:], "/")
	if idx == -1 {
 	   idx = len(origin[8:])
	}
	origin = origin[:8+idx]
	redirectURL := fmt.Sprintf("%s/redirect.php?ray_id=%s&alias=%s", origin, url.QueryEscape(rayID), url.QueryEscape(alias))

	resp, _, err := fetch(redirectURL, jar)
	if err != nil {
		return "", err
	}
	pageAfterRedirect := resp.Request.URL.String()

	xsrf := jar.Get("XSRF-TOKEN")
	if xsrf == "" {
		return "", fmt.Errorf("missing xsrf")
	}

	apiHeaders := map[string]string{
		"Origin":           origin,
		"Referer":          pageAfterRedirect,
		"Content-Type":     "application/json",
		"Accept":           "application/json, */*",
		"X-Requested-With": "XMLHttpRequest",
	}

	sessionURL := origin + "/api/session"
	token := makeToken(xsrf)
	s1, _ := postJSON(sessionURL, map[string]interface{}{"_token": token}, apiHeaders, jar)

	step := 1
	if s1 != nil {
		if s, ok := s1["step"].(float64); ok {
			step = int(s)
		}
	}

	ref := pageAfterRedirect

	if step == 1 {
		time.Sleep(300 * time.Millisecond)
		verifyURL := origin + "/api/verify"
		vd, _ := postJSON(verifyURL, map[string]interface{}{"_a": 0, "captcha": nil, "passcode": nil}, apiHeaders, jar)

		target := "/redirect.php"
		if vd != nil {
			if t, ok := vd["target"].(string); ok && t != "" {
				target = t
			}
		}
		if strings.HasPrefix(target, "/") {
			target = origin + target
		}

		resp2, _, err := fetch(target, jar)
		if err == nil {
			ref = resp2.Request.URL.String()
		}

		apiHeaders["Referer"] = ref
		newXsrf := jar.Get("XSRF-TOKEN")
		if newXsrf == "" {
			newXsrf = xsrf
		}
		postJSON(sessionURL, map[string]interface{}{"_token": makeToken(newXsrf)}, apiHeaders, jar)
	}

	time.Sleep(500 * time.Millisecond)
	key := 500
	size := fmt.Sprintf("%d.%d", (1920+key)*2, (1080+key)*2)
	apiHeaders["Referer"] = ref
	gd, _ := postJSON(origin+"/api/go", map[string]interface{}{"key": key, "size": size}, apiHeaders, jar)

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

	_, body, _ := fetch(readyURL, jar)
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
		Timeout: 15 * time.Second,
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

func elapsed(t time.Time) float64 {
	return float64(time.Since(t).Milliseconds()) / 1000.0
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ SFL BYPASS ]")))
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
