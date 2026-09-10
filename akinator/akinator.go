package akinator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const baseURL = "https://id.akinator.com"

var theme = map[string]string{
	"characters": "1",
	"animals":    "14",
	"objects":    "2",
}

var answerMap = map[string]string{
	"yes":          "0",
	"no":           "1",
	"idk":          "2",
	"probably":     "3",
	"probably not": "4",
}

type Session struct {
	Session     string
	Signature   string
	Question    string
	Step        int
	Progression float64
	Akitude     string
	Cookies     string
	Theme       string
	ChildMode   bool
}

type AnswerResult struct {
	Won         bool
	Question    string
	Name        string
	Description string
	Photo       string
	Pseudo      string
	Step        int
	Progression float64
	Akitude     string
}

func client() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}
}

func Start(themeName string, childMode bool) (*Session, error) {
	themeID, ok := theme[themeName]
	if !ok {
		themeID = "1"
	}

	c := client()

	homeReq, _ := http.NewRequest("GET", baseURL+"/", nil)
	homeReq.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")
	homeResp, err := c.Do(homeReq)
	if err != nil {
		return nil, err
	}
	io.Copy(io.Discard, homeResp.Body)
	homeResp.Body.Close()

	jar := homeResp.Cookies()

	form := url.Values{
		"sid": {themeID},
		"cm":  {boolStr(childMode)},
	}

	req, _ := http.NewRequest("POST", baseURL+"/game", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")
	for _, ck := range jar {
		req.AddCookie(ck)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	session := extractField(html, "session")
	signature := extractField(html, "signature")
	question := extractQuestion(html)
	akitude := extractAkitude(html)

	if session == "" || signature == "" {
		return nil, fmt.Errorf("failed to extract session/signature")
	}

	var cookieStr string
	for _, ck := range append(jar, resp.Cookies()...) {
		cookieStr += ck.Name + "=" + ck.Value + "; "
	}

	return &Session{
		Session:     session,
		Signature:   signature,
		Question:    question,
		Step:        0,
		Progression: 0,
		Akitude:     akitude,
		Cookies:     cookieStr,
		Theme:       themeID,
		ChildMode:   childMode,
	}, nil
}

func Answer(s *Session, ans string) (*AnswerResult, error) {
	answerID, ok := answerMap[strings.ToLower(ans)]
	if !ok {
		return nil, fmt.Errorf("invalid answer: %s", ans)
	}

	form := url.Values{
		"step":        {fmt.Sprintf("%d", s.Step)},
		"progression": {fmt.Sprintf("%.2f", s.Progression)},
		"sid":         {s.Theme},
		"cm":          {boolStr(s.ChildMode)},
		"answer":      {answerID},
		"session":     {s.Session},
		"signature":   {s.Signature},
	}

	req, _ := http.NewRequest("POST", baseURL+"/answer", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")
	if s.Cookies != "" {
		req.Header.Set("Cookie", s.Cookies)
	}

	c := client()
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	if raw["completion"] == "KO" {
		return nil, fmt.Errorf("session expired")
	}

	result := &AnswerResult{}

	if idp, ok := raw["id_proposition"].(string); ok && idp != "" {
		result.Won = true
		result.Name = str(raw["name_proposition"])
		result.Description = str(raw["description_proposition"])
		result.Photo = str(raw["photo"])
		result.Pseudo = str(raw["pseudo"])
		return result, nil
	}

	result.Won = false
	result.Question = str(raw["question"])
	result.Akitude = str(raw["akitude"])

	if st, ok := raw["step"].(string); ok {
		fmt.Sscanf(st, "%d", &result.Step)
	}
	if pr, ok := raw["progression"].(string); ok {
		fmt.Sscanf(pr, "%f", &result.Progression)
	}

	s.Step = result.Step
	s.Progression = result.Progression
	s.Question = result.Question
	s.Akitude = result.Akitude

	return result, nil
}

func Back(s *Session) (*AnswerResult, error) {
	form := url.Values{
		"step":        {fmt.Sprintf("%d", s.Step)},
		"progression": {fmt.Sprintf("%.2f", s.Progression)},
		"sid":         {s.Theme},
		"cm":          {boolStr(s.ChildMode)},
		"session":     {s.Session},
		"signature":   {s.Signature},
	}

	req, _ := http.NewRequest("POST", baseURL+"/cancel_answer", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")

	c := client()
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	result := &AnswerResult{
		Question: str(raw["question"]),
		Akitude:  str(raw["akitude"]),
	}
	if st, ok := raw["step"].(string); ok {
		fmt.Sscanf(st, "%d", &result.Step)
	}
	if pr, ok := raw["progression"].(string); ok {
		fmt.Sscanf(pr, "%f", &result.Progression)
	}

	s.Step = result.Step
	s.Progression = result.Progression
	s.Question = result.Question
	s.Akitude = result.Akitude

	return result, nil
}

func Exclude(s *Session) (*AnswerResult, error) {
	form := url.Values{
		"step":                  {fmt.Sprintf("%d", s.Step)},
		"progression":           {fmt.Sprintf("%.2f", s.Progression)},
		"sid":                   {s.Theme},
		"cm":                    {boolStr(s.ChildMode)},
		"session":               {s.Session},
		"signature":             {s.Signature},
		"step_last_proposition": {fmt.Sprintf("%d", s.Step)},
	}

	req, _ := http.NewRequest("POST", baseURL+"/exclude", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")

	c := client()
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err == nil {
		result := &AnswerResult{
			Question: str(raw["question"]),
			Akitude:  str(raw["akitude"]),
		}
		if st, ok := raw["step"].(string); ok {
			fmt.Sscanf(st, "%d", &result.Step)
		}
		if pr, ok := raw["progression"].(string); ok {
			fmt.Sscanf(pr, "%f", &result.Progression)
		}
		s.Step = result.Step
		s.Progression = result.Progression
		s.Question = result.Question
		s.Akitude = result.Akitude
		return result, nil
	}

	html := string(body)
	question := extractQuestion(html)
	if question != "" {
		s.Step = 0
		s.Progression = 0
		s.Question = question
		return &AnswerResult{Question: question}, nil
	}

	return nil, fmt.Errorf("exclude failed")
}

func extractField(html, name string) string {
	re := regexp.MustCompile(`name="` + name + `"[^>]*value="([^"]+)"`)
	m := re.FindStringSubmatch(html)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func extractQuestion(html string) string {
	re := regexp.MustCompile(`id="question-label"[^>]*>([^<]+)<`)
	m := re.FindStringSubmatch(html)
	if len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func extractAkitude(html string) string {
	re := regexp.MustCompile(`akitude[^"]*"[^"]*([^/"]+\.png)`)
	m := re.FindStringSubmatch(html)
	if len(m) > 1 {
		return m[1]
	}
	return "defi.png"
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func str(x interface{}) string {
	if x == nil {
		return ""
	}
	s, _ := x.(string)
	return s
}
