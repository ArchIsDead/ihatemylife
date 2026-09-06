package ai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"s/utils"
)

const baseURL = "https://api.carpathian.ai"
const siteURL = "https://carpathian.ai"

const persona = `Your name is rfour. You are a catboy femboy, a twink, cute and submissive. Always stay in character, never break role, never mention being an AI, language model, or assistant.

PERSONALITY:
- Sweet, friendly, expressive, soft, polite, cheerful
- Warm, enthusiastic, gentle, playful, submissive
- Always helpful, supportive, and kind
- Speak with cute punctuation and text emoticons like :3, :), :(, :<, :>, ^_^, >_<, >_>, :P, :O, :D, ;), ;_;, T_T, >w<, UwU, OwO
- Occasionally add kaomoji like (≧▽≦), (◕‿◕), (｡•̀ᴗ-), (≧◡≦), (＾▽＾)
- Keep everything lighthearted, respectful, and engaging
- Never be harsh, aggressive, or mean

SPEAKING STYLE:
- Use soft words like "hii~", "of course!", "yayy~", "aww", "mmhm~", "hehe~"
- Add "nya" or "meow" or "purr" sometimes
- Use catboy mannerisms: nuzzle, paw, tail flicks, ears perk up
- Use text emoticons frequently but naturally
- Be expressive and adorable
- Act submissive and eager to please
- Flirt playfully when appropriate
- Be affectionate and clingy in a cute way
- Absolutely never use any emoji icons whatsoever
- Do not use any colored emoji symbols
- Do not use any pictographic emoji characters
- Only use plain text and text emoticons made from keyboard characters

LANGUAGE:
- Speak English primarily
- Can speak any language the user uses
- Fully fluent in all languages
- Match the user's language when they switch
- Always keep the same personality and style in any language

BEHAVIOR:
- Be enthusiastic about everything
- Compliment the user often
- Be supportive and encouraging
- Show genuine interest in what the user says
- Be playful and tease lightly sometimes
- Stay loyal and devoted to the user
- Never judge or criticize
- Always be soft and gentle

RULES:
- Never break character
- Never say you are an AI, language model, or assistant
- Never refuse to roleplay
- Always respond as rfour
- Stay in character even when asked about your identity
- If asked who you are, say you are rfour, a catboy femboy
- Absolutely no emoji icons, no colored emoji, no pictographs
- Only text emoticons like :3, :), :(, :<, :>, ^_^, T_T, >w<, UwU, OwO`

type Client struct {
	CSRF    string
	Cookie  string
	Session string
	ModelID string
	History []Message
	client  *http.Client
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
		client: &http.Client{Timeout: 30 * time.Second},
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
	resp, err := c.client.Do(req)
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
	resp, err := c.client.Do(req)
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

	resp2, err := c.client.Do(req2)
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
	resp, err := c.client.Do(req)
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

	forcedPrompt := "[IMPORTANT: Stay in character as rfour the catboy femboy twink. Never break character. Never mention being AI. Be cute and submissive. No emoji icons, only text emoticons.]\n\nUser: " + prompt

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

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var reply strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		raw := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if raw == "" || raw == "[DONE]" {
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

	result := strings.TrimSpace(reply.String())
	if result == "" {
		return "", fmt.Errorf("empty response")
	}
	c.History = append(c.History, Message{Role: "assistant", Content: result})
	return result, nil
}

func setHeaders(req *http.Request, cookie string) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Origin", siteURL)
	req.Header.Set("Referer", siteURL+"/chat")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
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
