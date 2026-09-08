package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
)

func ShowBanner(b string) {
	fmt.Println(GW(b))
}

func Menu() {
	fmt.Println(Div())
	fmt.Println(Bld(Wht("[ MENU ]")))
	fmt.Println(Div())
	fmt.Println(Gry("[01]") + Wht(" Check PTK"))
	fmt.Println(Gry("[02]") + Wht(" DAPO Progress"))
	fmt.Println(Gry("[03]") + Wht(" Username Scan"))
	fmt.Println(Gry("[04]") + Wht(" Search School"))
	fmt.Println(Gry("[05]") + Wht(" School Info"))
	fmt.Println(Gry("[06]") + Wht(" IP Check"))
	fmt.Println(Gry("[07]") + Wht(" Current IP"))
	fmt.Println(Gry("[08]") + Wht(" Postal Code Search"))
	fmt.Println(Gry("[09]") + Wht(" NIK Parser"))
	fmt.Println(Gry("[10]") + Wht(" Web2Zip"))
	fmt.Println(Gry("[11]") + Wht(" Trace Moe"))
	fmt.Println(Gry("[12]") + Wht(" NSFW Check"))
	fmt.Println(Gry("[13]") + Wht(" Google Search"))
	fmt.Println(Gry("[14]") + Wht(" URL Shortener"))
	fmt.Println(Gry("[15]") + Wht(" Bypass Link"))
	fmt.Println(Gry("[16]") + Wht(" Hash Generator"))
	fmt.Println(Gry("[17]") + Wht(" Encoder"))
	fmt.Println(Gry("[18]") + Wht(" Decoder"))
	fmt.Println(Gry("[19]") + Wht(" Presets"))
	fmt.Println(Gry("[20]") + Wht(" Song Finder"))
	fmt.Println(Gry("[21]") + Wht(" TikTok Stalk"))
	fmt.Println(Gry("[22]") + Wht(" Pinterest Stalk"))
	fmt.Println(Gry("[23]") + Wht(" FF Stalk"))
	fmt.Println(Gry("[24]") + Wht(" Web2APK"))
	fmt.Println(Gry("[00]") + Red(" Exit"))
	fmt.Println(Div())
}

func PrintJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func Ask(r *bufio.Reader, l string) string {
	fmt.Print(Gry(l))
	v, _ := r.ReadString('\n')
	return strings.TrimSpace(v)
}

func Prompt(t string) string {
	return Gry(t)
}

func Err(t string) {
	fmt.Println(Red(t))
}

func Div() string {
	return Gry("────────────────────────────────────────")
}
