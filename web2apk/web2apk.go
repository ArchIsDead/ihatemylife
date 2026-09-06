package web2apk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"s/utils"
)

const baseURL = "https://rfweb2apk.rfdevv.com"

type Client struct {
	Token string
}

type Permission struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Safe bool   `json:"safe"`
}

var Permissions = []Permission{
	{ID: "1", Name: "android.permission.INTERNET", Desc: "Internet Access", Safe: true},
	{ID: "2", Name: "android.permission.ACCESS_NETWORK_STATE", Desc: "Network State", Safe: true},
	{ID: "3", Name: "android.permission.ACCESS_WIFI_STATE", Desc: "WiFi Info", Safe: true},
	{ID: "4", Name: "android.permission.VIBRATE", Desc: "Vibrate", Safe: true},
	{ID: "5", Name: "android.permission.WAKE_LOCK", Desc: "Keep Screen On", Safe: true},
	{ID: "6", Name: "android.permission.FOREGROUND_SERVICE", Desc: "Foreground Service", Safe: true},
	{ID: "7", Name: "android.permission.POST_NOTIFICATIONS", Desc: "Push Notifications", Safe: true},
	{ID: "8", Name: "android.permission.RECEIVE_BOOT_COMPLETED", Desc: "Auto Start", Safe: true},
	{ID: "9", Name: "android.permission.CAMERA", Desc: "Camera Access", Safe: false},
	{ID: "10", Name: "android.permission.RECORD_AUDIO", Desc: "Record Audio", Safe: false},
	{ID: "11", Name: "android.permission.READ_EXTERNAL_STORAGE", Desc: "Read Storage", Safe: false},
	{ID: "12", Name: "android.permission.WRITE_EXTERNAL_STORAGE", Desc: "Write Storage", Safe: false},
	{ID: "13", Name: "android.permission.ACCESS_FINE_LOCATION", Desc: "GPS Location", Safe: false},
	{ID: "14", Name: "android.permission.ACCESS_COARSE_LOCATION", Desc: "Network Location", Safe: false},
	{ID: "15", Name: "android.permission.READ_CONTACTS", Desc: "Read Contacts", Safe: false},
	{ID: "16", Name: "android.permission.WRITE_CONTACTS", Desc: "Write Contacts", Safe: false},
	{ID: "17", Name: "android.permission.READ_SMS", Desc: "Read SMS", Safe: false},
	{ID: "18", Name: "android.permission.SEND_SMS", Desc: "Send SMS", Safe: false},
	{ID: "19", Name: "android.permission.RECEIVE_SMS", Desc: "Receive SMS", Safe: false},
	{ID: "20", Name: "android.permission.CALL_PHONE", Desc: "Call Phone", Safe: false},
	{ID: "21", Name: "android.permission.READ_CALL_LOG", Desc: "Read Call Log", Safe: false},
	{ID: "22", Name: "android.permission.BLUETOOTH", Desc: "Bluetooth", Safe: false},
	{ID: "23", Name: "android.permission.BLUETOOTH_CONNECT", Desc: "Bluetooth Connect", Safe: false},
	{ID: "24", Name: "android.permission.NFC", Desc: "NFC", Safe: false},
	{ID: "25", Name: "android.permission.FLASHLIGHT", Desc: "Flashlight", Safe: true},
	{ID: "26", Name: "android.permission.USE_FINGERPRINT", Desc: "Fingerprint", Safe: false},
	{ID: "27", Name: "android.permission.USE_BIOMETRIC", Desc: "Biometric", Safe: false},
}

var PermissionPresets = map[string]struct {
	Name  string
	Perms []string
}{
	"1": {Name: "Standard / Recommended (Safe)", Perms: []string{
		"android.permission.INTERNET",
		"android.permission.ACCESS_NETWORK_STATE",
		"android.permission.ACCESS_WIFI_STATE",
		"android.permission.VIBRATE",
		"android.permission.WAKE_LOCK",
		"android.permission.POST_NOTIFICATIONS",
	}},
	"2": {Name: "All Permissions (27)", Perms: func() []string {
		var all []string
		for _, p := range Permissions {
			all = append(all, p.Name)
		}
		return all
	}()},
	"3": {Name: "Media, Camera & Storage", Perms: []string{
		"android.permission.INTERNET",
		"android.permission.ACCESS_NETWORK_STATE",
		"android.permission.CAMERA",
		"android.permission.RECORD_AUDIO",
		"android.permission.READ_EXTERNAL_STORAGE",
		"android.permission.WRITE_EXTERNAL_STORAGE",
		"android.permission.FLASHLIGHT",
	}},
	"4": {Name: "Location & Maps (GPS)", Perms: []string{
		"android.permission.INTERNET",
		"android.permission.ACCESS_NETWORK_STATE",
		"android.permission.ACCESS_FINE_LOCATION",
		"android.permission.ACCESS_COARSE_LOCATION",
	}},
	"5": {Name: "Social & Communication", Perms: []string{
		"android.permission.INTERNET",
		"android.permission.ACCESS_NETWORK_STATE",
		"android.permission.READ_CONTACTS",
		"android.permission.WRITE_CONTACTS",
		"android.permission.READ_SMS",
		"android.permission.SEND_SMS",
		"android.permission.RECEIVE_SMS",
		"android.permission.CALL_PHONE",
	}},
	"6": {Name: "Minimal (Internet Only)", Perms: []string{
		"android.permission.INTERNET",
		"android.permission.ACCESS_NETWORK_STATE",
	}},
}

type BuildRequest struct {
	AppName     string
	PackageName string
	URL         string
	HTMLPath    string
	ZIPPath     string
	VersionName string
	VersionCode string
	Orientation string
	SplashType  string
	Perms       []string
	IconPath    string
	SplashPath  string
	SplashHTML  string
	SplashVideo string
}

type BuildResult struct {
	Success     bool
	DownloadURL string
	Size        int64
	BuildID     string
	FileName    string
	Error       string
}

func (c *Client) headers() map[string]string {
	h := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
	if c.Token != "" {
		h["Authorization"] = "Bearer " + c.Token
	}
	return h
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func addFile(w *multipart.Writer, field, pathOrURL, filename string) error {
	if pathOrURL == "" {
		return nil
	}

	var data io.Reader
	var name string

	if isURL(pathOrURL) {
		resp, err := http.Get(pathOrURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		data = resp.Body
		name = filename
	} else {
		file, err := os.Open(pathOrURL)
		if err != nil {
			return err
		}
		defer file.Close()
		data = file
		name = filepath.Base(pathOrURL)
	}

	part, err := w.CreateFormFile(field, name)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, data)
	return err
}

func (c *Client) Build(reqData BuildRequest) (*BuildResult, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	w.WriteField("appName", reqData.AppName)
	w.WriteField("packageName", reqData.PackageName)
	w.WriteField("versionName", reqData.VersionName)
	w.WriteField("versionCode", reqData.VersionCode)
	w.WriteField("orientation", reqData.Orientation)
	w.WriteField("splashType", reqData.SplashType)

	if reqData.URL != "" {
		w.WriteField("url", reqData.URL)
	}

	if reqData.Perms != nil && len(reqData.Perms) > 0 {
		pb, _ := json.Marshal(reqData.Perms)
		w.WriteField("permissions", string(pb))
	}

	addFile(w, "html", reqData.HTMLPath, "index.html")
	addFile(w, "zip", reqData.ZIPPath, "website.zip")
	addFile(w, "icon", reqData.IconPath, "icon.png")

	if reqData.SplashType == "image" {
		addFile(w, "splash", reqData.SplashPath, "splash.png")
	} else if reqData.SplashType == "html" {
		addFile(w, "splashHtml", reqData.SplashHTML, "splash.html")
	} else if reqData.SplashType == "video" {
		addFile(w, "splashVideo", reqData.SplashVideo, "splash.mp4")
	}

	w.Close()

	req, _ := http.NewRequest("POST", baseURL+"/api/apk/build", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	for k, v := range c.headers() {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&raw)

	result := &BuildResult{}
	if s, ok := raw["success"].(bool); ok {
		result.Success = s
	}
	if u, ok := raw["downloadUrl"].(string); ok {
		result.DownloadURL = u
	}
	if sz, ok := raw["size"].(float64); ok {
		result.Size = int64(sz)
	}
	if id, ok := raw["buildId"].(string); ok {
		result.BuildID = id
	}
	if fn, ok := raw["fileName"].(string); ok {
		result.FileName = fn
	}
	if e, ok := raw["error"].(string); ok {
		result.Error = e
	}

	return result, nil
}

func (r *BuildResult) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ WEB2APK BUILD ]")))
	fmt.Println(utils.Div())
	if !r.Success {
		fmt.Println(utils.Red("Error: " + r.Error))
		return
	}
	fmt.Println(utils.Gry("Status: ") + utils.Wht("Success"))
	if r.FileName != "" {
		fmt.Println(utils.Gry("File: ") + utils.Wht(r.FileName))
	}
	if r.BuildID != "" {
		fmt.Println(utils.Gry("Build ID: ") + utils.Wht(r.BuildID))
	}
	if r.Size > 0 {
		fmt.Println(utils.Gry("Size: ") + utils.Wht(fmt.Sprintf("%.2f MB", float64(r.Size)/1024/1024)))
	}
	if r.DownloadURL != "" {
		full := r.DownloadURL
		if !strings.HasPrefix(full, "http") {
			full = baseURL + "/" + strings.TrimPrefix(full, "/")
		}
		fmt.Println(utils.Gry("Download: ") + utils.Wht(full))
	}
}
