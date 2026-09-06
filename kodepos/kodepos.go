package kodepos

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"s/utils"
)

type Result struct {
	No string `json:"no"`
	Kp string `json:"postal_code"`
	D  string `json:"village"`
	Kc string `json:"district"`
	Kk string `json:"city_regency"`
	P  string `json:"province"`
}

func extractRows(h string) []Result {
	var rs []Result
	re := regexp.MustCompile(`(?s)<tr>\s*<td>(.*?)</td>\s*<td>(.*?)</td>\s*<td>(.*?)</td>\s*<td>(.*?)</td>\s*<td>(.*?)</td>\s*<td>(.*?)</td>\s*</tr>`)
	ms := re.FindAllStringSubmatch(h, -1)
	for _, m := range ms {
		rs = append(rs, Result{
			No: strings.TrimSpace(m[1]),
			Kp: strings.TrimSpace(m[2]),
			D:  strings.TrimSpace(m[3]),
			Kc: strings.TrimSpace(m[4]),
			Kk: strings.TrimSpace(m[5]),
			P:  strings.TrimSpace(m[6]),
		})
	}
	return rs
}

func Cari(k string, pg int) ([]Result, int, error) {
	fm := url.Values{"kodepos": {k}}
	req, _ := http.NewRequest("POST", "https://kodepos.posindonesia.co.id/CariKodepos", strings.NewReader(fm.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Origin", "https://kodepos.posindonesia.co.id")
	req.Header.Set("Referer", "https://kodepos.posindonesia.co.id/CariKodepos")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	bd, _ := io.ReadAll(resp.Body)
	rs := extractRows(string(bd))
	total := len(rs)
	per := 25
	tp := (total + per - 1) / per
	start := (pg - 1) * per
	end := start + per
	if end > total {
		end = total
	}
	if start >= total {
		return []Result{}, tp, nil
	}
	return rs[start:end], tp, nil
}
