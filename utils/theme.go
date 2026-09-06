package utils

import "fmt"

func B(b string) { fmt.Println(GW(b)) }
func T(t string) string { return B(W(t)) }
func S(t string) { fmt.Println(G(t)) }
func P(t string) string { return C(t) }
func I(t string) string { return W(t) }
func E(t string) { fmt.Println(R(t)) }
func A(t string) string { return Y(t) }
func J(t string) string { return BW(t) }
func H() string { return G("────────────────────────────────────────") }
