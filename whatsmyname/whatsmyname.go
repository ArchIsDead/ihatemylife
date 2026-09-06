package whatsmyname

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"s/utils"
)

type Result struct {
	Username string
	Total    int
	Page     int
	Pages    int
	Exists   int
	NotExist int
	Data     []Item
}

type Item struct {
	URL      string
	Source   string
	Exists   bool
	Category string
}

func Scan(n, m, rc string, pg int) (*Result, error) {
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

	r := &Result{Username: n, Page: pg}

	rr, _ := raw["result"].([]interface{})
	var all []Item
	for _, x := range rr {
		xx, _ := x.(map[string]interface{})
		all = append(all, Item{
			URL:      str(xx["url"]),
			Source:   str(xx["source"]),
			Exists:   xx["isExist"] == true,
			Category: str(xx["category"]),
		})
	}

	for _, item := range all {
		if item.Exists {
			r.Exists++
		}
	}

	var filtered []Item
	switch m {
	case "exist":
		for _, item := range all {
			if item.Exists {
				filtered = append(filtered, item)
			}
		}
	case "notexist":
		for _, item := range all {
			if !item.Exists {
				filtered = append(filtered, item)
			}
		}
	default:
		filtered = all
	}

	r.Total = len(filtered)
	r.NotExist = r.Total - r.Exists
	per := 10
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
	r.Data = filtered[start:end]

	return r, nil
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ USERNAME SCAN ]")))
	fmt.Println(utils.Gry(fmt.Sprintf("Username: %s", r.Username)))
	fmt.Println(utils.Gry(fmt.Sprintf("Total: %d", r.Total)))
	fmt.Println(utils.Gry(fmt.Sprintf("Page: %d/%d", r.Page, r.Pages)))
	fmt.Println(utils.Gry(fmt.Sprintf("Exists: %d", r.Exists)))
	fmt.Println(utils.Gry(fmt.Sprintf("Not Exists: %d", r.NotExist)))
	fmt.Println(utils.Div())

	if len(r.Data) == 0 {
		fmt.Println(utils.Gry("No results"))
		return
	}

	for i, item := range r.Data {
		status := "NO"
		if item.Exists {
			status = "YES"
		}
		fmt.Println(utils.Gry(fmt.Sprintf("[%d] %s", i+1, item.Source)))
		fmt.Println(utils.Gry("URL: ") + utils.Wht(item.URL))
		fmt.Println(utils.Gry("Category: ") + utils.Wht(item.Category))
		fmt.Println(utils.Gry("Exists: ") + utils.Wht(status))
		fmt.Println()
	}
}

func str(x interface{}) string {
	if x == nil {
		return ""
	}
	s, _ := x.(string)
	return s
}
