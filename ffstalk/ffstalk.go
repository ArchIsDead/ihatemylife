package ffstalk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"s/utils"
)

type Result struct {
	UID       string
	Nickname  string
	Level     int
	Exp       int
	Likes     int
	Rank      int
	CSRank    int
	Region    string
	Banned    bool
	BanPeriod int
	BanStatus string
	LastLogin string
	CreateAt  string
	Avatar    string
	Character string
	CharIcon  string
}

func Get(uid string) (*Result, error) {
	if uid == "" {
		return nil, fmt.Errorf("uid required")
	}

	u := "https://freefire.my.id/api/ff?uid=" + uid

	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Referer", "https://freefire.my.id/stalk/"+uid)
	req.Header.Set("Accept", "application/json")

	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	r := &Result{UID: uid}

	if player, ok := raw["player"].(map[string]interface{}); ok {
		r.Nickname = str(player["nickname"])
		r.Level = intNum(player["level"])
		r.Exp = intNum(player["exp"])
		r.Likes = intNum(player["liked"])
		r.Rank = intNum(player["rank"])
		r.CSRank = intNum(player["csRank"])
		r.Region = str(player["region"])
		r.LastLogin = str(player["lastLoginAt"])
		r.CreateAt = str(player["createAt"])
		r.Avatar = str(player["avatarUrl"])

		if char, ok := player["equippedCharacter"].(map[string]interface{}); ok {
			r.Character = str(char["name"])
			r.CharIcon = str(char["icon"])
		}
	}

	if ban, ok := raw["ban"].(map[string]interface{}); ok {
		r.Banned = ban["isBanned"] == true
		r.BanPeriod = intNum(ban["banPeriod"])
		r.BanStatus = str(ban["status"])
	}

	return r, nil
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ FREE FIRE STALK ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("UID: ") + utils.Wht(r.UID))
	if r.Nickname != "" {
		fmt.Println(utils.Gry("Nickname: ") + utils.Wht(r.Nickname))
	}
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Level: ") + utils.Wht(fmt.Sprintf("%d", r.Level)))
	fmt.Println(utils.Gry("Exp: ") + utils.Wht(fmt.Sprintf("%d", r.Exp)))
	fmt.Println(utils.Gry("Likes: ") + utils.Wht(fmt.Sprintf("%d", r.Likes)))
	fmt.Println(utils.Gry("Rank: ") + utils.Wht(fmt.Sprintf("%d", r.Rank)))
	fmt.Println(utils.Gry("CS Rank: ") + utils.Wht(fmt.Sprintf("%d", r.CSRank)))
	if r.Region != "" {
		fmt.Println(utils.Gry("Region: ") + utils.Wht(r.Region))
	}
	fmt.Println(utils.Div())
	if r.Banned {
		fmt.Println(utils.Gry("Banned: ") + utils.Red("YES"))
		fmt.Println(utils.Gry("Ban Period: ") + utils.Wht(fmt.Sprintf("%d days", r.BanPeriod)))
		fmt.Println(utils.Gry("Ban Status: ") + utils.Wht(r.BanStatus))
	} else {
		fmt.Println(utils.Gry("Banned: ") + utils.Grn("NO"))
	}
	if r.LastLogin != "" && r.LastLogin != "0" {
		fmt.Println(utils.Gry("Last Login: ") + utils.Wht(r.LastLogin))
	}
	if r.CreateAt != "" && r.CreateAt != "0" {
		fmt.Println(utils.Gry("Created: ") + utils.Wht(r.CreateAt))
	}
	if r.Character != "" {
		fmt.Println(utils.Div())
		fmt.Println(utils.Gry("Character: ") + utils.Wht(r.Character))
	}
	if r.CharIcon != "" {
		fmt.Println(utils.Gry("Character Icon: ") + utils.Wht(r.CharIcon))
	}
	if r.Avatar != "" {
		fmt.Println(utils.Gry("Avatar: ") + utils.Wht(r.Avatar))
	}
}

func str(x interface{}) string {
	if x == nil {
		return ""
	}
	s, _ := x.(string)
	return s
}

func intNum(x interface{}) int {
	if x == nil {
		return 0
	}
	switch v := x.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		n, _ := strconv.Atoi(v)
		return n
	}
	return 0
}
