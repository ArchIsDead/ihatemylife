package utils

import "fmt"

func ShowBanner(b string) { fmt.Println(GW(b)) }
func Title(t string) string { return Bld(Wht(t)) }
func Sub(t string) { fmt.Println(Gry(t)) }
func Prompt(t string) string { return Cyn(t) }
func Input(t string) string { return Wht(t) }
func Err(t string) { fmt.Println(Red(t)) }
func Acc(t string) string { return Ylw(t) }
func JSON(t string) string { return BW(t) }
func Div() string { return Gry("────────────────────────────────────────") }
