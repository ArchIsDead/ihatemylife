package encoder

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"

	"s/utils"
)

type Result struct {
	Input     string
	Base64    string
	URLEncode string
	Hex       string
}

func Do(input string) *Result {
	return &Result{
		Input:     input,
		Base64:    base64.StdEncoding.EncodeToString([]byte(input)),
		URLEncode: url.QueryEscape(input),
		Hex:       hex.EncodeToString([]byte(input)),
	}
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ ENCODER ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Input: ") + utils.Wht(r.Input))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Base64: ") + utils.Wht(r.Base64))
	fmt.Println(utils.Gry("URL: ") + utils.Wht(r.URLEncode))
	fmt.Println(utils.Gry("Hex: ") + utils.Wht(r.Hex))
}
