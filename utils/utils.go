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

func Menu() {
	fmt.Println(Div())
	fmt.Println(Title("[ MENU ]"))
	fmt.Println(Div())
	fmt.Println(Gry("[01]") + Wht(" NUPTK Search"))
	fmt.Println(Gry("[02]") + Wht(" List Provinces"))
	fmt.Println(Gry("[03]") + Wht(" List Cities"))
	fmt.Println(Gry("[04]") + Wht(" DAPO Progress"))
	fmt.Println(Gry("[05]") + Wht(" Username Scan"))
	fmt.Println(Gry("[06]") + Wht(" Search School"))
	fmt.Println(Gry("[07]") + Wht(" School Info"))
	fmt.Println(Gry("[08]") + Wht(" IP Check"))
	fmt.Println(Gry("[09]") + Wht(" Current IP"))
	fmt.Println(Gry("[10]") + Wht(" Device Info"))
	fmt.Println(Gry("[11]") + Wht(" Stats"))
	fmt.Println(Gry("[12]") + Wht(" Postal Code Search"))
	fmt.Println(Gry("[13]") + Wht(" NIK Parser"))
	fmt.Println(Gry("[14]") + Wht(" Instagram Viewer"))
	fmt.Println(Gry("[00]") + Red(" Exit"))
	fmt.Println(Div())
}

func PrintJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(JSON(string(b)))
}

func Ask(r *bufio.Reader, l string) string {
	fmt.Print(Prompt(l))
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

func ShowUp() {
	fmt.Println(Gry("Uptime: ") + Wht(UpStr()))
	fmt.Println(Div())
}

func DeviceInfo() map[string]interface{} {
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
