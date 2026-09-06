package utils

import (
	"fmt"
	"os"
	"runtime"
)

var Cc = true

func init() {
	if runtime.GOOS == "windows" {
		Cc = false
	}
	if os.Getenv("TERM") == "dumb" {
		Cc = false
	}
	if os.Getenv("NO_COLOR") != "" {
		Cc = false
	}
}

func cl(x, t string) string {
	if !Cc {
		return t
	}
	return x + t + "\033[0m"
}

func K(t string) string { return cl("\033[30m", t) }
func W(t string) string { return cl("\033[37m", t) }
func G(t string) string { return cl("\033[90m", t) }
func B(t string) string { return cl("\033[1m", t) }
func D(t string) string { return cl("\033[2m", t) }
func R(t string) string { return cl("\033[31m", t) }
func N(t string) string { return cl("\033[32m", t) }
func Y(t string) string { return cl("\033[33m", t) }
func L(t string) string { return cl("\033[34m", t) }
func M(t string) string { return cl("\033[35m", t) }
func C(t string) string { return cl("\033[36m", t) }

func BW(t string) string {
	if !Cc {
		return t
	}
	q := []rune(t)
	if len(q) == 0 {
		return t
	}
	var o string
	for i, r := range q {
		z := float64(i) / float64(len(q)-1)
		e := 232 + int(z*23)
		o += fmt.Sprintf("\033[38;5;%dm%c", e, r)
	}
	return o + "\033[0m"
}

func GW(t string) string {
	if !Cc {
		return t
	}
	q := []rune(t)
	if len(q) == 0 {
		return t
	}
	var o string
	for i, r := range q {
		z := float64(i) / float64(len(q)-1)
		e := 240 + int(z*15)
		o += fmt.Sprintf("\033[38;5;%dm%c", e, r)
	}
	return o + "\033[0m"
}

func WB(t string) string {
	if !Cc {
		return t
	}
	q := []rune(t)
	if len(q) == 0 {
		return t
	}
	var o string
	for i, r := range q {
		z := 1.0 - float64(i)/float64(len(q)-1)
		e := 232 + int(z*23)
		o += fmt.Sprintf("\033[38;5;%dm%c", e, r)
	}
	return o + "\033[0m"
}
