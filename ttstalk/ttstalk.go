package ttstalk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"s/utils"
)

type Result struct {
	Username  string
	Nickname  string
	ID        string
	SecUID    string
	Bio       string
	BioLink   string
	Verified  bool
	Private   bool
	Seller    bool
	Avatar    string
	Followers int
	Following int
	Likes     int
	Videos    int
	Friends   int
	Created   string
}

func Get(handle string) (*Result, error) {
	clean := strings.TrimPrefix(strings.TrimSpace(handle), "@")
	if clean == "" {
		return nil, fmt.Errorf("username required")
	}

	u := "https://www.tiktok.com/@" + url.QueryEscape(clean)

	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 14; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Accept", "text/html")

	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	re := regexp.MustCompile(`<script id="__UNIVERSAL_DATA_FOR_REHYDRATION__"[^>]*>([\s\S]*?)</script>`)
	m := re.FindStringSubmatch(html)
	if len(m) < 2 {
		return nil, fmt.Errorf("data not found")
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(m[1]), &raw); err != nil {
		return nil, err
	}

	r := &Result{Username: clean}

	defScope, _ := raw["__DEFAULT_SCOPE__"].(map[string]interface{})
	if defScope == nil {
		defScope, _ = raw["DEFAULT_SCOPE"].(map[string]interface{})
	}

	userDetail, _ := defScope["webapp.user-detail"].(map[string]interface{})
	if userDetail == nil {
		userDetail, _ = defScope["webapp.reflow.profile.initial"].(map[string]interface{})
	}

	userInfo, _ := userDetail["userInfo"].(map[string]interface{})

	user, _ := userInfo["user"].(map[string]interface{})
	stats, _ := userInfo["stats"].(map[string]interface{})

	if user == nil {
		usersModule, _ := raw["UserModule"].(map[string]interface{})
		users, _ := usersModule["users"].(map[string]interface{})
		statsModule, _ := usersModule["stats"].(map[string]interface{})

		for k, v := range users {
			user, _ = v.(map[string]interface{})
			stats, _ = statsModule[k].(map[string]interface{})
			break
		}
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	r.Nickname = str(user["nickname"])
	r.ID = str(user["id"])
	r.SecUID = str(user["secUid"])
	r.Bio = str(user["signature"])
	r.Verified = user["verified"] == true
	r.Private = user["privateAccount"] == true
	r.Seller = user["ttSeller"] == true
	r.Avatar = str(user["avatarLarger"])
	if r.Avatar == "" {
		r.Avatar = str(user["avatarMedium"])
	}
	if r.Avatar == "" {
		r.Avatar = str(user["avatarThumb"])
	}

	if bioLink, ok := user["bioLink"].(map[string]interface{}); ok {
		r.BioLink = str(bioLink["link"])
	}

	if ct, ok := user["createTime"].(float64); ok {
		r.Created = fmt.Sprintf("%d", int64(ct))
	}

	if stats != nil {
		r.Followers = intNum(stats["followerCount"])
		r.Following = intNum(stats["followingCount"])
		r.Likes = intNum(stats["heartCount"])
		if r.Likes == 0 {
			r.Likes = intNum(stats["heart"])
		}
		r.Videos = intNum(stats["videoCount"])
		r.Friends = intNum(stats["friendCount"])
	}

	return r, nil
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ TIKTOK STALK ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Username: ") + utils.Wht("@"+r.Username))
	if r.Nickname != "" {
		fmt.Println(utils.Gry("Nickname: ") + utils.Wht(r.Nickname))
	}
	if r.ID != "" {
		fmt.Println(utils.Gry("ID: ") + utils.Wht(r.ID))
	}
	if r.SecUID != "" {
		fmt.Println(utils.Gry("SecUID: ") + utils.Wht(r.SecUID))
	}
	fmt.Println(utils.Div())
	if r.Bio != "" {
		fmt.Println(utils.Gry("Bio: ") + utils.Wht(r.Bio))
	}
	if r.BioLink != "" {
		fmt.Println(utils.Gry("Bio Link: ") + utils.Wht(r.BioLink))
	}
	fmt.Println(utils.Div())
	if r.Verified {
		fmt.Println(utils.Gry("Verified: ") + utils.Grn("YES"))
	} else {
		fmt.Println(utils.Gry("Verified: ") + utils.Red("NO"))
	}
	if r.Private {
		fmt.Println(utils.Gry("Private: ") + utils.Red("YES"))
	} else {
		fmt.Println(utils.Gry("Private: ") + utils.Grn("NO"))
	}
	if r.Seller {
		fmt.Println(utils.Gry("Seller: ") + utils.Grn("YES"))
	} else {
		fmt.Println(utils.Gry("Seller: ") + utils.Red("NO"))
	}
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Followers: ") + utils.Wht(fmt.Sprintf("%d", r.Followers)))
	fmt.Println(utils.Gry("Following: ") + utils.Wht(fmt.Sprintf("%d", r.Following)))
	fmt.Println(utils.Gry("Likes: ") + utils.Wht(fmt.Sprintf("%d", r.Likes)))
	fmt.Println(utils.Gry("Videos: ") + utils.Wht(fmt.Sprintf("%d", r.Videos)))
	fmt.Println(utils.Gry("Friends: ") + utils.Wht(fmt.Sprintf("%d", r.Friends)))
	if r.Created != "" {
		fmt.Println(utils.Gry("Created: ") + utils.Wht(r.Created))
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
	}
	return 0
}
