package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var st = time.Now()

func M() {
	fmt.Println(H())
	fmt.Println(T("[ MENU ]"))
	fmt.Println(H())
	fmt.Println(G("[01]") + W(" NUPTK Search"))
	fmt.Println(G("[02]") + W(" List Provinces"))
	fmt.Println(G("[03]") + W(" List Cities"))
	fmt.Println(G("[04]") + W(" DAPO Progress"))
	fmt.Println(G("[05]") + W(" Username Scan"))
	fmt.Println(G("[06]") + W(" Search School"))
	fmt.Println(G("[07]") + W(" School Info"))
	fmt.Println(G("[08]") + W(" IP Check"))
	fmt.Println(G("[09]") + W(" Current IP"))
	fmt.Println(G("[10]") + W(" Device Info"))
	fmt.Println(G("[11]") + W(" Stats"))
	fmt.Println(G("[12]") + W(" Postal Code Search"))
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

func sh(cmd string) string {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func UpStr() string {
	raw := sh("cat /proc/uptime")
	p := strings.Fields(raw)
	if len(p) < 1 {
		return "0h 0m 0s"
	}
	var sec float64
	fmt.Sscanf(p[0], "%f", &sec)
	h := int(sec) / 3600
	m := (int(sec) % 3600) / 60
	s := int(sec) % 60
	return fmt.Sprintf("%dh %dm %ds", h, m, s)
}

func Up() {
	fmt.Println(G("Uptime: ") + W(UpStr()))
	fmt.Println(H())
}

func Di() map[string]interface{} {
	hn, _ := os.Hostname()
	wd, _ := os.Getwd()
	return map[string]interface{}{
		"hostname": hn,
		"workdir":  wd,
		"pid":      os.Getpid(),
		"uptime":   UpStr(),
		"app_up":   time.Since(st).String(),
	}
}
