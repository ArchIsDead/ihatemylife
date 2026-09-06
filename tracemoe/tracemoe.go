package tracemoe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"s/utils"
)

type Result struct {
	Total int
	Data  []Item
}

type Item struct {
	Title    string
	Episode  int
	Similarity float64
	Video    string
	Image    string
	Anilist  int
	MyAnimeList int
}

func Search(path string) (*Result, error) {
	if path == "" {
		return nil, errors.New("image path required")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("image", "image.jpg")
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest("POST", "https://api.trace.moe/search", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Origin", "https://trace.moe")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36")
	resp, err := utils.Hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	r := &Result{}
	data, _ := raw["result"].([]interface{})
	r.Total = len(data)

	for _, d := range data {
		dd, _ := d.(map[string]interface{})
		item := Item{
			Title:     str(dd["filename"]),
			Episode:   int(num(dd["episode"])),
			Similarity: num(dd["similarity"]),
			Video:     str(dd["video"]),
			Image:     str(dd["image"]),
			Anilist:   int(num(dd["anilist"])),
			MyAnimeList: int(num(dd["mal_id"])),
		}
		r.Data = append(r.Data, item)
	}

	return r, nil
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ TRACE MOE ]")))
	fmt.Println(utils.Gry(fmt.Sprintf("Results: %d", r.Total)))
	fmt.Println(utils.Div())

	if r.Total == 0 {
		fmt.Println(utils.Gry("No matches"))
		return
	}

	for i, d := range r.Data {
		fmt.Println(utils.Gry(fmt.Sprintf("[%d] %s", i+1, d.Title)))
		fmt.Println(utils.Gry("Episode: ") + utils.Wht(fmt.Sprintf("%d", d.Episode)))
		fmt.Println(utils.Gry("Similarity: ") + utils.Wht(fmt.Sprintf("%.2f%%", d.Similarity*100)))
		fmt.Println(utils.Gry("Anilist ID: ") + utils.Wht(fmt.Sprintf("%d", d.Anilist)))
		fmt.Println(utils.Gry("MAL ID: ") + utils.Wht(fmt.Sprintf("%d", d.MyAnimeList)))
		if d.Video != "" {
			fmt.Println(utils.Gry("Video: ") + utils.Wht(d.Video))
		}
		if d.Image != "" {
			fmt.Println(utils.Gry("Image: ") + utils.Wht(d.Image))
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

func num(x interface{}) float64 {
	if x == nil {
		return 0
	}
	f, _ := x.(float64)
	return f
}
