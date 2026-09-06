package whatsmyname

import (
	"bytes"
	"encoding/json"
	"net/http"

	"s/utils"
)

type Result struct {
	Url string `json:"url"`
	Src string `json:"source"`
	Ex  bool   `json:"exists"`
	Cat string `json:"category"`
}

func Scan(n, m, rc string) (map[string]interface{}, error) {
	bd := map[string]interface{}{
		"source": n,
		"type":   "name",
		"rescan": rc == "true",
	}
	b, _ := json.Marshal(bd)
	req, _ := http.NewRequest("POST", "https://api.whatsmyname.io/discoverprofile", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if raw["error"] != nil {
		return raw, nil
	}
	rr, _ := raw["result"].([]interface{})
	var rs []Result
	for _, x := range rr {
		xx, _ := x.(map[string]interface{})
		rs = append(rs, Result{
			Url: str(xx["url"]),
			Src: str(xx["source"]),
			Ex:  xx["isExist"] == true,
			Cat: str(xx["category"]),
		})
	}
	ec := 0
	for _, r := range rs {
		if r.Ex {
			ec++
		}
	}
	out := map[string]interface{}{
		"username": n,
		"rescan":   rc == "true",
		"total":    len(rs),
		"exists":   ec,
		"not":      len(rs) - ec,
		"results":  rs,
	}
	if m == "exist" {
		var f []Result
		for _, r := range rs {
			if r.Ex {
				f = append(f, r)
			}
		}
		out = map[string]interface{}{"username": n, "total": len(f), "found": f}
	}
	if m == "notexist" {
		var nf []Result
		for _, r := range rs {
			if !r.Ex {
				nf = append(nf, r)
			}
		}
		out = map[string]interface{}{"username": n, "total": len(nf), "not_found": nf}
	}
	return out, nil
}

func str(x interface{}) string {
	if x == nil {
		return ""
	}
	s, _ := x.(string)
	return s
}
