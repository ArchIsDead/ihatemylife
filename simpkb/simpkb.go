package simpkb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"s/utils"
)

const baseURL = "https://app.simpkb.id"

type Result struct {
	Keyword  string
	Total    int
	Data     []Person
}

type Person struct {
	Nama      string
	Nuptk     string
	Provinsi  string
	Kota      string
	Kecamatan string
	Sekolah   string
	Status    string
}

func Cari(keyword, prov, kab, paspor, dapodik, page string) (*Result, error) {
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

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	result := &Result{Keyword: keyword}

	data, _ := raw["data"].([]interface{})
	if len(data) == 0 {
		return result, nil
	}

	for _, d := range data {
		dd, _ := d.(map[string]interface{})
		if content, ok := dd["content"].([]interface{}); ok {
			for _, c := range content {
				cc, _ := c.(map[string]interface{})
				p := Person{
					Nama:      str(cc["nama"]),
					Nuptk:     str(cc["nuptk"]),
					Provinsi:  str(cc["provinsi"]),
					Kota:      str(cc["kota"]),
					Kecamatan: str(cc["kecamatan"]),
					Sekolah:   str(cc["sekolah"]),
					Status:    str(cc["status"]),
				}
				result.Data = append(result.Data, p)
			}
		} else {
			p := Person{
				Nama:      str(dd["nama"]),
				Nuptk:     str(dd["nuptk"]),
				Provinsi:  str(dd["provinsi"]),
				Kota:      str(dd["kota"]),
				Kecamatan: str(dd["kecamatan"]),
				Sekolah:   str(dd["sekolah"]),
				Status:    str(dd["status"]),
			}
			result.Data = append(result.Data, p)
		}
	}

	result.Total = len(result.Data)
	return result, nil
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ NUPTK SEARCH ]")))
	fmt.Println(utils.Gry(fmt.Sprintf("Keyword: %s", r.Keyword)))
	fmt.Println(utils.Gry(fmt.Sprintf("Total: %d", r.Total)))
	fmt.Println(utils.Div())

	if r.Total == 0 {
		fmt.Println(utils.Gry("No results"))
		return
	}

	for i, p := range r.Data {
		fmt.Println(utils.Gry(fmt.Sprintf("[%d] %s", i+1, p.Nama)))
		if p.Nuptk != "-" {
			fmt.Println(utils.Gry("NUPTK: ") + utils.Wht(p.Nuptk))
		}
		if p.Sekolah != "-" {
			fmt.Println(utils.Gry("School: ") + utils.Wht(p.Sekolah))
		}
		if p.Kecamatan != "-" {
			fmt.Println(utils.Gry("District: ") + utils.Wht(p.Kecamatan))
		}
		if p.Kota != "-" {
			fmt.Println(utils.Gry("City: ") + utils.Wht(p.Kota))
		}
		if p.Provinsi != "-" {
			fmt.Println(utils.Gry("Province: ") + utils.Wht(p.Provinsi))
		}
		if p.Status != "-" {
			fmt.Println(utils.Gry("Status: ") + utils.Wht(p.Status))
		}
		fmt.Println()
	}
}

func str(x interface{}) string {
	if x == nil {
		return "-"
	}
	s, _ := x.(string)
	if s == "" {
		return "-"
	}
	return s
}
