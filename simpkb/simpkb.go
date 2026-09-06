package simpkb

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"s/utils"
)

const baseURL = "https://app.simpkb.id"

func fetchJSON(path string) (map[string]interface{}, error) {
	bd, err := utils.F(baseURL+path, map[string]string{
		"User-Agent":       "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36",
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          baseURL + "/akun/ptk",
	})
	if err != nil {
		return nil, err
	}
	var o map[string]interface{}
	if err := json.Unmarshal(bd, &o); err != nil {
		return nil, err
	}
	return o, nil
}

func Provinsi() (map[string]interface{}, error) {
	return fetchJSON("/asset/js/configs/propinsi.json?version=202309261116")
}

func Kota() (map[string]interface{}, error) {
	return fetchJSON("/asset/js/configs/kota.json?version=202309261116")
}

func Cari(keyword, prov, kab, paspor, dapodik, page string) (map[string]interface{}, error) {
	if keyword == "" {
		return nil, errors.New("keyword empty")
	}
	if page == "" {
		page = "1"
	}
	fm := url.Values{
		"k_propinsi": {prov},
		"k_kota":     {kab},
		"is_paspor":  {paspor},
		"is_dapodik": {dapodik},
		"keyword":    {keyword},
		"page":       {page},
	}
	req, _ := http.NewRequest("POST", baseURL+"/akun/ptk-solr", strings.NewReader(fm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", baseURL+"/akun/ptk")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var o map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return nil, err
	}
	return o, nil
}
