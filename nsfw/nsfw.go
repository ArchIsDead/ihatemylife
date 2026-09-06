package nsfw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"s/utils"
)

type Result struct {
	Label      string
	LabelID    string
	Confidence float64
}

func Check(input string) (*Result, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "image.jpg")

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		resp, err := utils.Hc.Get(input)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		io.Copy(part, resp.Body)
	} else {
		file, err := os.Open(input)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		io.Copy(part, file)
	}

	writer.Close()

	req, _ := http.NewRequest("POST", "https://www.nyckel.com/v1/functions/o2f0jzcdyut2qxhu/invoke", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return &Result{
		Label:      str(raw["labelName"]),
		LabelID:    str(raw["labelId"]),
		Confidence: num(raw["confidence"]),
	}, nil
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ NSFW CHECK ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Label: ") + utils.Wht(r.Label))
	fmt.Println(utils.Gry("Confidence: ") + utils.Wht(fmt.Sprintf("%.2f%%", r.Confidence*100)))
}

func str(x interface{}) string {
	if x == nil {
		return ""
	}
	s, _ := x.(string)
	return s
}

func num(x interface{}) float64 {
	if x == nil {
		return 0
	}
	f, _ := x.(float64)
	return f
}
