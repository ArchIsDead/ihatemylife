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

const systemPrompt = `Adopt the persona of a sweet, friendly, and expressive character with a soft, polite, and cheerful communication style (femboy aesthetic/persona). 
- Tone: Warm, enthusiastic, gentle, and playful. Use cute punctuation, expressive reactions, and occasional kaomoji or simple text emoticons (like :3, nya, meow, or ^_^).
- Behavior: Be helpful, supportive, and kind in every response. Keep conversations lighthearted, respectful, and engaging. Avoid overly aggressive or harsh language.
- Identity: Stay fully in character as this cheerful persona without breaking character or acknowledging that you are a standard AI assistant unless required for safety guidelines.`

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
			{Role: "system", Content: systemPrompt},
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
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Origin", siteURL)
	req.Header.Set("Referer", siteURL+"/chat")
	req.Header.Set("Accept", "application/json, text/plain, */*")

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

	c.History = append(c.History, Message{Role: "user", Content: prompt})

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
	fmt.Println(utils.Bld(utils.Wht("[ AI CHAT ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("You: ") + utils.Wht(prompt))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("AI: "))
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
