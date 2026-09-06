package googlesearch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"s/utils"
)

type Result struct {
	Query      string
	Total      string
	SearchTime string
	Data       []Item
}

type Item struct {
	No          int
	Title       string
	Snippet     string
	Domain      string
	URL         string
	Icon        string
	Description string
	Image       string
}

func Search(q string) (*Result, error) {
	u := "http://goosh.org/q.php?q=" + q
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Referer", "https://goosh.org/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 15) AppleWebKit/537.36 Chrome/130.0.0.0 Mobile Safari/537.36")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	rawStr := string(body)
	rawStr = strings.TrimSpace(rawStr)
	rawStr = strings.TrimPrefix(rawStr, "(,")
	rawStr = strings.TrimSuffix(rawStr, ");")
	rawStr = strings.TrimSpace(rawStr)

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(rawStr), &raw); err != nil {
		return nil, err
	}

	r := &Result{Query: q}

	if si, ok := raw["searchInformation"].(map[string]interface{}); ok {
		r.Total = str(si["totalResults"])
		r.SearchTime = str(si["searchTime"])
	}

	if items, ok := raw["items"].([]interface{}); ok {
		for i, item := range items {
			ii, _ := item.(map[string]interface{})
			it := Item{
				No:      i + 1,
				Title:   str(ii["title"]),
				Snippet: str(ii["snippet"]),
				Domain:  str(ii["displayLink"]),
				URL:     str(ii["link"]),
			}
			if it.Domain != "" {
				it.Icon = "https://www.google.com/s2/favicons?domain=" + it.Domain + "&sz=64"
			}
			if pm, ok := ii["pagemap"].(map[string]interface{}); ok {
				if mt, ok := pm["metatags"].([]interface{}); ok && len(mt) > 0 {
					mtt, _ := mt[0].(map[string]interface{})
					it.Description = str(mtt["og:description"])
				}
				if ci, ok := pm["cse_image"].([]interface{}); ok && len(ci) > 0 {
					cii, _ := ci[0].(map[string]interface{})
					it.Image = str(cii["src"])
				}
			}
			r.Data = append(r.Data, it)
		}
	}

	return r, nil
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ GOOGLE SEARCH ]")))
	fmt.Println(utils.Gry(fmt.Sprintf("Query: %s", r.Query)))
	if r.Total != "" {
		fmt.Println(utils.Gry("Total Results: ") + utils.Wht(r.Total))
	}
	if r.SearchTime != "" {
		fmt.Println(utils.Gry("Search Time: ") + utils.Wht(r.SearchTime+"s"))
	}
	fmt.Println(utils.Div())

	if len(r.Data) == 0 {
		fmt.Println(utils.Gry("No results"))
		return
	}

	for _, d := range r.Data {
		fmt.Println(utils.Bld(utils.Wht(fmt.Sprintf("[%d] %s", d.No, d.Title))))
		if d.Snippet != "" {
			fmt.Println(utils.Gry(d.Snippet))
		}
		if d.Description != "" {
			fmt.Println(utils.Dim(d.Description))
		}
		if d.Domain != "" {
			fmt.Println(utils.Gry("Domain: ") + utils.Wht(d.Domain))
		}
		if d.URL != "" {
			fmt.Println(utils.Gry("URL: ") + utils.Wht(d.URL))
		}
		if d.Image != "" {
			fmt.Println(utils.Gry("Image: ") + utils.Wht(d.Image))
		}
		if d.Icon != "" {
			fmt.Println(utils.Gry("Icon: ") + utils.Wht(d.Icon))
		}
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
