package decoder

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
	URLDecode string
	Hex       string
}

func Do(input string) *Result {
	r := &Result{Input: input}

	if b, err := base64.StdEncoding.DecodeString(input); err == nil {
		r.Base64 = string(b)
	} else {
		r.Base64 = "invalid"
	}

	if u, err := url.QueryUnescape(input); err == nil {
		r.URLDecode = u
	} else {
		r.URLDecode = "invalid"
	}

	if h, err := hex.DecodeString(input); err == nil {
		r.Hex = string(h)
	} else {
		r.Hex = "invalid"
	}

	return r
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ DECODER ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Input: ") + utils.Wht(r.Input))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Base64: ") + utils.Wht(r.Base64))
	fmt.Println(utils.Gry("URL: ") + utils.Wht(r.URLDecode))
	fmt.Println(utils.Gry("Hex: ") + utils.Wht(r.Hex))
}
