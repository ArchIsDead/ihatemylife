package dapo

import (
	"encoding/json"
	"regexp"

	"s/utils"
)

const baseURL = "https://dapo.kemendikdasmen.go.id"

var tk struct {
	T string
	U string
}

func token() (string, string, error) {
	if tk.T != "" {
		return tk.T, tk.U, nil
	}
	_, _ = utils.F(baseURL+"/", map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36",
	})
	bd, err := utils.F(baseURL+"/env.js", map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36",
	})
	if err != nil {
		return "", "", err
	}
	raw := string(bd)
	tk.T = regexp.MustCompile(`VITE_API_TOKEN["']?\s*:\s*["']([^"']+)["']`).FindStringSubmatch(raw)[1]
	tk.U = regexp.MustCompile(`VITE_STRAPI_URL["']?\s*:\s*["']([^"']+)["']`).FindStringSubmatch(raw)[1]
	return tk.T, tk.U, nil
}

func fetchAPI(p string, q map[string]string) (map[string]interface{}, error) {
	t, u := token()
	qq := ""
	for k, v := range q {
		if qq != "" {
			qq += "&"
		}
		qq += k + "=" + v
	}
	if qq != "" {
		p += "?" + qq
	}
	bd, err := utils.F(u+p, map[string]string{
		"Authorization": "Bearer " + t,
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

func ProgressProvinsi(j, s string) (map[string]interface{}, error) {
	q := map[string]string{}
	if j != "" {
		q["jenjang"] = j
	}
	if s != "" {
		q["status_sekolah"] = s
	}
	return fetchAPI("/api/progress-pengiriman/provinsi", q)
}

func ProgressKabupaten(k, j, s string) (map[string]interface{}, error) {
	return fetchAPI("/api/progress-pengiriman/kabupaten", map[string]string{
		"kode_provinsi":  k,
		"jenjang":        j,
		"status_sekolah": s,
	})
}

func ProgressKecamatan(k, j, s string) (map[string]interface{}, error) {
	return fetchAPI("/api/progress-pengiriman/kecamatan", map[string]string{
		"kode_kabupaten": k,
		"jenjang":         j,
		"status_sekolah":  s,
	})
}

func ProgressSekolah(k, j string) (map[string]interface{}, error) {
	return fetchAPI("/api/progress-pengiriman/kecamatan/school", map[string]string{
		"kode_kecamatan": k,
		"jenjang":         j,
	})
}

func CariSekolah(q string) (map[string]interface{}, error) {
	t, u := token()
	bd, err := utils.F(u+"/api/detail-sekolah/search?q="+q, map[string]string{
		"Authorization": "Bearer " + t,
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

func InfoSekolah(n string) (map[string]interface{}, error) {
	t, u := token()
	bd, err := utils.F(u+"/api/detail-sekolah?npsn="+n, map[string]string{
		"Authorization": "Bearer " + t,
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
