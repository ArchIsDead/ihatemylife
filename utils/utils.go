package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var st = time.Now()

func M() {
	fmt.Println(H())
	fmt.Println(T("[ MENU ]"))
	fmt.Println(H())
	fmt.Println(G("[01]") + W(" NUPTK Search"))
	fmt.Println(G("[02]") + W(" List Provinsi"))
	fmt.Println(G("[03]") + W(" List Kota"))
	fmt.Println(G("[04]") + W(" DAPO Progress"))
	fmt.Println(G("[05]") + W(" Username Scan"))
	fmt.Println(G("[06]") + W(" Cari Sekolah"))
	fmt.Println(G("[07]") + W(" Info Sekolah"))
	fmt.Println(G("[08]") + W(" IP Check"))
	fmt.Println(G("[09]") + W(" Current IP"))
	fmt.Println(G("[10]") + W(" Device Info"))
	fmt.Println(G("[11]") + W(" Stats"))
	fmt.Println(G("[12]") + W(" Kodepos Search"))
	fmt.Println(G("[13]") + W(" NIK Parser"))
	fmt.Println(G("[14]") + W(" Instagram Viewer"))
	fmt.Println(G("[00]") + R(" Exit"))
	fmt.Println(H())
}

func O(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(J(string(b)))
}

func P(r *bufio.Reader, l string) string {
	fmt.Print(P(l))
	v, _ := r.ReadString('\n')
	return strings.TrimSpace(v)
}

func Up() string {
	d := time.Since(st)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%dh %dm %ds", h, m, s)
}

func Di() map[string]interface{} {
	hn, _ := os.Hostname()
	return map[string]interface{}{
		"hostname": hn,
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"cpus":     runtime.NumCPU(),
		"go":       runtime.Version(),
		"uptime":   Up(),
	}
}
