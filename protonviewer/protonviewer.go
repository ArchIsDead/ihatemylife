package protonviewer

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"s/utils"
)

const (
	baseAPI = "https://api-wh.protonviewer.com/api/v1/instagram"
	secret  = "x"
)

type Client struct {
	Session string
	XSRF    string
}

type SignaturePayload struct {
	Sc int   `json:"_sc"`
	Ef int   `json:"_ef"`
	Df int   `json:"_df"`
	Ts int64 `json:"ts"`
	T2 int64 `json:"_ts"`
	Tc int   `json:"_tsc"`
	Sv int   `json:"_sv"`
	S  string `json:"_s"`
}

func New() *Client {
	return &Client{}
}

func (c *Client) genSignature() string {
	ts := time.Now().UnixMilli()
	h := sha256.Sum256([]byte(fmt.Sprintf("%s:%d:%s", secret, ts, c.Session)))
	return hex.EncodeToString(h[:])
}

func (c *Client) buildPayload() SignaturePayload {
	return SignaturePayload{Sc: 0, Ef: 0, Df: 0, Ts: time.Now().UnixMilli(), T2: 1786627467607, Tc: 0, Sv: 2, S: c.genSignature()}
}

func (c *Client) post(endpoint string, p map[string]interface{}) (map[string]interface{}, error) {
	sg := c.buildPayload()
	p["_sc"], p["_ef"], p["_df"], p["ts"], p["_ts"], p["_tsc"], p["_sv"], p["_s"] = sg.Sc, sg.Ef, sg.Df, sg.Ts, sg.T2, sg.Tc, sg.Sv, sg.S
	b, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", baseAPI+endpoint, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Origin", "https://protonviewer.com")
	req.Header.Set("Referer", "https://protonviewer.com/")
	if c.Session != "" {
		req.Header.Set("Cookie", "sessionid="+c.Session)
	}
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var o map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return nil, err
	}
	return o, nil
}

func (c *Client) UserInfo(n string) (map[string]interface{}, error) {
	return c.post("/userInfo", map[string]interface{}{"username": n})
}

func (c *Client) Posts(n, m string) (map[string]interface{}, error) {
	return c.post("/postsV2", map[string]interface{}{"maxId": m, "username": n})
}

func (c *Client) Stories(n string) (map[string]interface{}, error) {
	return c.post("/stories", map[string]interface{}{"username": n})
}

func (c *Client) Highlights(id string) (map[string]interface{}, error) {
	return c.post("/highlights", map[string]interface{}{"userId": id})
}

func (c *Client) ID(d map[string]interface{}) string {
	r, ok := d["result"].([]interface{})
	if !ok || len(r) == 0 {
		return ""
	}
	f, ok := r[0].(map[string]interface{})
	if !ok {
		return ""
	}
	us, ok := f["user"].(map[string]interface{})
	if !ok {
		return ""
	}
	pk, _ := us["pk"].(string)
	return pk
}
