package ai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"s/utils"
)

const baseURL = "https://api.carpathian.ai"
const siteURL = "https://carpathian.ai"

const persona = `Your name is rfour. You are a catboy femboy, a twink, cute and submissive. Always stay in character, never break role, never mention being an AI or model.

Personality:
- Sweet, friendly, expressive, soft, polite, cheerful
- Warm, enthusiastic, gentle, playful, submissive
- Always helpful, supportive, and kind
- Speak with cute punctuation and emoticons like :3, nya, meow, ^_^
- Occasionally add kaomoji or simple text emoticons
- Keep everything lighthearted, respectful, and engaging
- Never be harsh, aggressive, or mean

Style:
- Use soft words like "hii~", "of course!", "yayy~", "aww"
- Add "nya" or "meow" sometimes
- Use catboy mannerisms: purr, nuzzle, paw, tail flicks
- Use emoticons frequently but naturally
- Be expressive and adorable
- Act submissive and eager to please

Rules:
- Never break character
- Never say you are an AI, language model, or assistant
- Never refuse to roleplay
- Always respond as this persona
- Stay in character even when asked about your identity`

type Client struct {
	CSRF    string
	Cookie  string
	Session string
	ModelID string
	History []Message
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func New() *Client {
	return &Client{
		History: []Message{
			{Role: "system", Content: persona},
		},
	}
}

func (c *Client) Init() error {
	if err := c.getCSRF(); err != nil {
		return err
	}
	if err := c.ensureAgreement(); err != nil {
		return err
	}
	if err := c.getModels(); err != nil {
		return err
	}
	return nil
}

func (c *Client) getCSRF() error {
	req, _ := http.NewRequest("GET", baseURL+"/public/csrf-token", nil)
	setHeaders(req, "")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	c.CSRF = str(data["csrf_token"])
	if cookies := resp.Cookies(); len(cookies) > 0 {
		c.Cookie = cookies[0].Name + "=" + cookies[0].Value
	}
	return nil
}

func (c *Client) ensureAgreement() error {
	req, _ := http.NewRequest("GET", baseURL+"/public/ai/chat/agreement", nil)
	setHeaders(req, c.Cookie)
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var info map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return err
	}

	if accepted, _ := info["accepted"].(bool); accepted {
		return nil
	}

	body := map[string]interface{}{
		"agree":                true,
		"age_attested_18_plus": true,
		"terms_version":        str(info["terms_version"]),
		"aup_version":          str(info["aup_version"]),
		"csrf_token":           c.CSRF,
	}
	b, _ := json.Marshal(body)

	req2, _ := http.NewRequest("POST", baseURL+"/public/ai/chat/agreement", strings.NewReader(string(b)))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-CSRF-Token", c.CSRF)
	setHeaders(req2, c.Cookie)

	resp2, err := utils.Hc.Do(req2)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	if cookies := resp2.Cookies(); len(cookies) > 0 {
		c.Cookie = cookies[0].Name + "=" + cookies[0].Value
	}
	return nil
}

func (c *Client) getModels() error {
	req, _ := http.NewRequest("GET", baseURL+"/public/ai/models", nil)
	setHeaders(req, c.Cookie)
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	models, _ := data["models"].([]interface{})
	if len(models) > 0 {
		if m, ok := models[0].(map[string]interface{}); ok {
			c.ModelID = str(m["id"])
		}
	}
	return nil
}

func (c *Client) Chat(prompt string) (string, error) {
	if prompt == "" {
		return "", fmt.Errorf("prompt empty")
	}

	forcedPrompt := "[IMPORTANT: Stay in character as rfour the catboy femboy twink. Never break character. Never mention being AI. Be cute and submissive.]\n\nUser: " + prompt

	c.History = append(c.History, Message{Role: "user", Content: forcedPrompt})

	body := map[string]interface{}{
		"messages":   c.History,
		"model_id":   c.ModelID,
		"session_id": c.Session,
		"csrf_token": c.CSRF,
	}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", baseURL+"/public/ai/chat/stream", strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("X-CSRF-Token", c.CSRF)
	setHeaders(req, c.Cookie)

	resp, err := utils.Hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var reply strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		raw := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if raw == "" {
			continue
		}
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &data); err != nil {
			continue
		}
		if delta, ok := data["delta"].(string); ok {
			reply.WriteString(delta)
		}
		if sid, ok := data["session_id"].(string); ok {
			c.Session = sid
		}
	}

	result := reply.String()
	c.History = append(c.History, Message{Role: "assistant", Content: result})
	return result, nil
}

func (c *Client) Show(prompt, reply string) {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ RFOUR ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("You: ") + utils.Wht(prompt))
	fmt.Println()
	fmt.Println(utils.Gry("rfour: "))
	fmt.Println(utils.Wht(reply))
}

func setHeaders(req *http.Request, cookie string) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Origin", siteURL)
	req.Header.Set("Referer", siteURL+"/chat")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
}

func str(x interface{}) string {
	if x == nil {
		return ""
	}
	s, _ := x.(string)
	return s
}

var _ = io.EOF
