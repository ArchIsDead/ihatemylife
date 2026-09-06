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

func Blk(t string) string { return cl("\033[30m", t) }
func Wht(t string) string { return cl("\033[37m", t) }
func Gry(t string) string { return cl("\033[90m", t) }
func Bld(t string) string { return cl("\033[1m", t) }
func Dim(t string) string { return cl("\033[2m", t) }
func Red(t string) string { return cl("\033[31m", t) }
func Grn(t string) string { return cl("\033[32m", t) }

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
