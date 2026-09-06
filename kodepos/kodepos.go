package kodepos

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"s/utils"
)

type Result struct {
	Query string
	Total int
	Page  int
	Pages int
	Data  []Item
}

type Item struct {
	No          string
	PostalCode  string
	Village     string
	District    string
	CityRegency string
	Province    string
}

func Cari(k string, pg int) (*Result, error) {
	fm := url.Values{"kodepos": {k}}
	req, _ := http.NewRequest("POST", "https://kodepos.posindonesia.co.id/CariKodepos", strings.NewReader(fm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Origin", "https://kodepos.posindonesia.co.id")
	req.Header.Set("Referer", "https://kodepos.posindonesia.co.id/CariKodepos")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bd, _ := io.ReadAll(resp.Body)
	all := extractRows(string(bd))

	r := &Result{Query: k, Page: pg}
	r.Total = len(all)
	per := 25
	r.Pages = (r.Total + per - 1) / per
	start := (pg - 1) * per
	end := start + per
	if end > r.Total {
		end = r.Total
	}
	if start >= r.Total {
		start = 0
		end = 0
	}
	r.Data = all[start:end]

	return r, nil
}

func extractRows(h string) []Item {
	var rs []Item
	re := regexp.MustCompile(`(?s)<tr>\s*<td>(.*?)</td>\s*<td>(.*?)</td>\s*<td>(.*?)</td>\s*<td>(.*?)</td>\s*<td>(.*?)</td>\s*<td>(.*?)</td>\s*</tr>`)
	ms := re.FindAllStringSubmatch(h, -1)
	for _, m := range ms {
		rs = append(rs, Item{
			No:          strings.TrimSpace(m[1]),
			PostalCode:  strings.TrimSpace(m[2]),
			Village:     strings.TrimSpace(m[3]),
			District:    strings.TrimSpace(m[4]),
			CityRegency: strings.TrimSpace(m[5]),
			Province:    strings.TrimSpace(m[6]),
		})
	}
	return rs
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ POSTAL CODE ]")))
	fmt.Println(utils.Gry(fmt.Sprintf("Query: %s", r.Query)))
	fmt.Println(utils.Gry(fmt.Sprintf("Total: %d", r.Total)))
	fmt.Println(utils.Gry(fmt.Sprintf("Page: %d/%d", r.Page, r.Pages)))
	fmt.Println(utils.Div())

	if len(r.Data) == 0 {
		fmt.Println(utils.Gry("No results"))
		return
	}

	for _, d := range r.Data {
		fmt.Println(utils.Bld(utils.Wht(d.PostalCode + " - " + d.Village)))
		fmt.Println(utils.Gry("District: ") + utils.Wht(d.District))
		fmt.Println(utils.Gry("City/Regency: ") + utils.Wht(d.CityRegency))
		fmt.Println(utils.Gry("Province: ") + utils.Wht(d.Province))
		fmt.Println()
	}
}
