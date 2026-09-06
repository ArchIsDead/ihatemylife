package utils

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

var Hc = &http.Client{
	Timeout: 45 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func F(url string, h map[string]string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range h {
		req.Header.Set(k, v)
	}
	resp, err := Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
