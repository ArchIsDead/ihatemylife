package web2zip

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"s/utils"
)

type Result struct {
	URL         string
	ErrorText   string
	ErrorCode   int
	CopiedFiles int
	DownloadURL string
}

func Save(u string) (*Result, error) {
	if u == "" {
		return nil, errors.New("url required")
	}
	if !strings.HasPrefix(u, "https://") {
		u = "https://" + u
	}

	body := map[string]interface{}{
		"url":                  u,
		"renameAssets":         false,
		"saveStructure":        false,
		"alternativeAlgorithm": false,
		"mobileVersion":        false,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://copier.saveweb2zip.com/api/copySite", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://saveweb2zip.com")
	req.Header.Set("Referer", "https://saveweb2zip.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var first map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&first); err != nil {
		return nil, err
	}

	md5, _ := first["md5"].(string)
	if md5 == "" {
		return nil, errors.New("no md5")
	}

	for {
		req2, _ := http.NewRequest("GET", "https://copier.saveweb2zip.com/api/getStatus/"+md5, nil)
		req2.Header.Set("Accept", "*/*")
		req2.Header.Set("Origin", "https://saveweb2zip.com")
		req2.Header.Set("Referer", "https://saveweb2zip.com/")
		req2.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")
		resp2, err := utils.Hc.Do(req2)
		if err != nil {
			return nil, err
		}
		var st map[string]interface{}
		if err := json.NewDecoder(resp2.Body).Decode(&st); err != nil {
			resp2.Body.Close()
			return nil, err
		}
		resp2.Body.Close()

		done, _ := st["isFinished"].(bool)
		if done {
			return &Result{
				URL:         u,
				ErrorText:   str(st["errorText"]),
				ErrorCode:   int(num(st["errorCode"])),
				CopiedFiles: int(num(st["copiedFilesAmount"])),
				DownloadURL: "https://copier.saveweb2zip.com/api/downloadArchive/" + md5,
			}, nil
		}
		time.Sleep(1 * time.Second)
	}
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ WEB2ZIP ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("URL: ") + utils.Wht(r.URL))
	fmt.Println(utils.Gry("Files Copied: ") + utils.Wht(fmt.Sprintf("%d", r.CopiedFiles)))
	if r.ErrorText != "" {
		fmt.Println(utils.Gry("Error: ") + utils.Red(r.ErrorText))
	}
	if r.DownloadURL != "" {
		fmt.Println(utils.Gry("Download: ") + utils.Wht(r.DownloadURL))
	}
}

func str(x interface{}) string {
	if x == nil {
		return ""
	}
	s, _ := x.(string)
	return s
}

func num(x interface{}) float64 {
	if x == nil {
		return 0
	}
	f, _ := x.(float64)
	return f
}
