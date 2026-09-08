package pinstalk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"s/utils"
)

type Board struct {
	Name      string
	PinCount  int
	Thumbnail string
}

type Result struct {
	Username    string
	DisplayName string
	Bio         string
	Followers   int
	Following   int
	Verified    bool
	Website     string
	Avatar      string
	TotalBoards int
	TotalPins   int
	ProfileScore int
	Boards      []Board
}

func Get(username string) (*Result, error) {
	clean := strings.TrimPrefix(strings.TrimSpace(username), "@")
	if clean == "" {
		return nil, fmt.Errorf("username required")
	}

	u := "https://pinout.in/api/profile-analyze?username=" + url.QueryEscape(clean)

	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", "https://pinout.in/pinterest-profile-analyzer/"+clean)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Mobile Safari/537.36")

	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw struct {
		Status bool `json:"status"`
		Data   struct {
			Boards []struct {
				Name      string `json:"name"`
				PinCount  int    `json:"pinCount"`
				Thumbnail string `json:"thumbnail"`
			} `json:"boards"`
			Profile struct {
				Avatar      string `json:"avatar"`
				Bio         string `json:"bio"`
				DisplayName string `json:"displayName"`
				Followers   int    `json:"followers"`
				Following   int    `json:"following"`
				Username    string `json:"username"`
				Verified    bool   `json:"verified"`
				Website     string `json:"website"`
			} `json:"profile"`
			Stats struct {
				AvgPinsPerBoard int `json:"avgPinsPerBoard"`
				ProfileScore    int `json:"profileScore"`
				TotalBoards     int `json:"totalBoards"`
				TotalPinsEstimate int `json:"totalPinsEstimate"`
			} `json:"stats"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	if !raw.Status {
		return nil, fmt.Errorf("failed to fetch")
	}

	r := &Result{
		Username:     raw.Data.Profile.Username,
		DisplayName:  raw.Data.Profile.DisplayName,
		Bio:          raw.Data.Profile.Bio,
		Followers:    raw.Data.Profile.Followers,
		Following:    raw.Data.Profile.Following,
		Verified:     raw.Data.Profile.Verified,
		Website:      raw.Data.Profile.Website,
		Avatar:       raw.Data.Profile.Avatar,
		TotalBoards:  raw.Data.Stats.TotalBoards,
		TotalPins:    raw.Data.Stats.TotalPinsEstimate,
		ProfileScore: raw.Data.Stats.ProfileScore,
	}

	for _, b := range raw.Data.Boards {
		r.Boards = append(r.Boards, Board{
			Name:      b.Name,
			PinCount:  b.PinCount,
			Thumbnail: b.Thumbnail,
		})
	}

	return r, nil
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ PINTEREST STALK ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Username: ") + utils.Wht(r.Username))
	if r.DisplayName != "" {
		fmt.Println(utils.Gry("Display Name: ") + utils.Wht(r.DisplayName))
	}
	fmt.Println(utils.Div())
	if r.Bio != "" {
		fmt.Println(utils.Gry("Bio: ") + utils.Wht(r.Bio))
	}
	if r.Website != "" {
		fmt.Println(utils.Gry("Website: ") + utils.Wht(r.Website))
	}
	fmt.Println(utils.Div())
	if r.Verified {
		fmt.Println(utils.Gry("Verified: ") + utils.Grn("YES"))
	} else {
		fmt.Println(utils.Gry("Verified: ") + utils.Red("NO"))
	}
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Followers: ") + utils.Wht(fmt.Sprintf("%d", r.Followers)))
	fmt.Println(utils.Gry("Following: ") + utils.Wht(fmt.Sprintf("%d", r.Following)))
	fmt.Println(utils.Gry("Boards: ") + utils.Wht(fmt.Sprintf("%d", r.TotalBoards)))
	fmt.Println(utils.Gry("Pins: ") + utils.Wht(fmt.Sprintf("%d", r.TotalPins)))
	fmt.Println(utils.Gry("Score: ") + utils.Wht(fmt.Sprintf("%d", r.ProfileScore)))
	if r.Avatar != "" {
		fmt.Println(utils.Gry("Avatar: ") + utils.Wht(r.Avatar))
	}
	if len(r.Boards) > 0 {
		fmt.Println(utils.Div())
		fmt.Println(utils.Bld(utils.Wht("[ BOARDS ]")))
		fmt.Println(utils.Div())
		for i, b := range r.Boards {
			fmt.Println(utils.Gry(fmt.Sprintf("[%d] %s", i+1, b.Name)))
			fmt.Println(utils.Gry("    Pins: ") + utils.Wht(fmt.Sprintf("%d", b.PinCount)))
			if b.Thumbnail != "" {
				fmt.Println(utils.Gry("    Thumbnail: ") + utils.Wht(b.Thumbnail))
			}
		}
	}
}
