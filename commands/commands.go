package commands

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"s/dapo"
	"s/googlesearch"
	"s/ipinfo"
	"s/kodepos"
	"s/nikparser"
	"s/nsfw"
	"s/simpkb"
	"s/tracemoe"
	"s/utils"
	"s/web2zip"
	"s/whatsmyname"
)

func clear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func back(r *bufio.Reader) {
	utils.Ask(r, "Press Enter to return...")
	clear()
}

func N(r *bufio.Reader) {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ NUPTK SEARCH ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Examples:"))
	fmt.Println(utils.Gry("  Keyword: Novy"))
	fmt.Println(utils.Gry("  Province Code: 15 (Jambi)"))
	fmt.Println(utils.Gry("  City Code: 1502 (Kab. Kerinci)"))
	fmt.Println(utils.Gry("  Dapodik: 1 (connected) / 0 (not) / empty (all)"))
	fmt.Println(utils.Gry("  Passport: 1 (registered) / 0 (not) / empty (all)"))
	fmt.Println(utils.Gry("  Page: 1"))
	fmt.Println(utils.Div())

	k := utils.Ask(r, "Keyword: ")
	p := utils.Ask(r, "Province Code: ")
	kk := utils.Ask(r, "City Code: ")
	dp := utils.Ask(r, "Dapodik 0/1/empty: ")
	ps := utils.Ask(r, "Passport 0/1/empty: ")
	pg := utils.Ask(r, "Page: ")
	res, err := simpkb.Cari(k, p, kk, ps, dp, pg)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func D(r *bufio.Reader) {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ DAPO PROGRESS ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Modes:"))
	fmt.Println(utils.Gry("  province  - Province level progress"))
	fmt.Println(utils.Gry("  regency   - Regency/City level progress"))
	fmt.Println(utils.Gry("  district  - District level progress"))
	fmt.Println(utils.Gry("  school    - School level progress"))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Examples:"))
	fmt.Println(utils.Gry("  Mode: province"))
	fmt.Println(utils.Gry("  Level: SMP"))
	fmt.Println(utils.Gry("  Status: Swasta"))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("  Mode: regency"))
	fmt.Println(utils.Gry("  Province Code: 050000"))
	fmt.Println(utils.Gry("  Level: SMP"))
	fmt.Println(utils.Gry("  Status: Swasta"))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("  Mode: district"))
	fmt.Println(utils.Gry("  Regency Code: 052000"))
	fmt.Println(utils.Gry("  Level: SMP"))
	fmt.Println(utils.Gry("  Status: Swasta"))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("  Mode: school"))
	fmt.Println(utils.Gry("  District Code: 052001"))
	fmt.Println(utils.Gry("  Level: SMP"))
	fmt.Println(utils.Div())

	m := utils.Ask(r, "Mode (province/regency/district/school): ")
	switch m {
	case "province":
		j := utils.Ask(r, "Level: ")
		st := utils.Ask(r, "Status: ")
		res, err := dapo.ProgressProvinsi(j, st)
		if err != nil {
			utils.Err("Error: " + err.Error())
			back(r)
			return
		}
		res.Show()
	case "regency":
		k := utils.Ask(r, "Province Code: ")
		j := utils.Ask(r, "Level: ")
		st := utils.Ask(r, "Status: ")
		res, err := dapo.ProgressKabupaten(k, j, st)
		if err != nil {
			utils.Err("Error: " + err.Error())
			back(r)
			return
		}
		res.Show()
	case "district":
		k := utils.Ask(r, "Regency Code: ")
		j := utils.Ask(r, "Level: ")
		st := utils.Ask(r, "Status: ")
		res, err := dapo.ProgressKecamatan(k, j, st)
		if err != nil {
			utils.Err("Error: " + err.Error())
			back(r)
			return
		}
		res.Show()
	case "school":
		k := utils.Ask(r, "District Code: ")
		j := utils.Ask(r, "Level: ")
		res, err := dapo.ProgressSekolah(k, j)
		if err != nil {
			utils.Err("Error: " + err.Error())
			back(r)
			return
		}
		res.Show()
	default:
		utils.Err("Invalid mode.")
	}
	back(r)
}

func U(r *bufio.Reader) {
	n := utils.Ask(r, "Username: ")
	m := utils.Ask(r, "Mode all/exist/notexist: ")
	rc := utils.Ask(r, "Rescan true/false: ")
	pg := utils.Ask(r, "Page: ")
	pi, _ := strconv.Atoi(pg)
	if pi < 1 {
		pi = 1
	}
	res, err := whatsmyname.Scan(n, m, rc, pi)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func CS(r *bufio.Reader) {
	q := utils.Ask(r, "School Name: ")
	res, err := dapo.CariSekolah(q)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func IS(r *bufio.Reader) {
	n := utils.Ask(r, "NPSN: ")
	res, err := dapo.InfoSekolah(n)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func IP(r *bufio.Reader) {
	ip := utils.Ask(r, "IP Address: ")
	res, err := ipinfo.Lookup(ip)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func CIP() {
	res, err := ipinfo.Lookup("")
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	res.Show()
}

func KP(r *bufio.Reader) {
	q := utils.Ask(r, "Postal Code / Area Name: ")
	pg := utils.Ask(r, "Page: ")
	pi, _ := strconv.Atoi(pg)
	if pi < 1 {
		pi = 1
	}
	res, err := kodepos.Cari(q, pi)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func NP(r *bufio.Reader) {
	n := utils.Ask(r, "NIK: ")
	res, err := nikparser.Parse(n)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func WZ(r *bufio.Reader) {
	u := utils.Ask(r, "URL: ")
	res, err := web2zip.Save(u)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func TM(r *bufio.Reader) {
	p := utils.Ask(r, "Image Path or URL: ")
	res, err := tracemoe.Search(p)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func NS(r *bufio.Reader) {
	p := utils.Ask(r, "Image Path: ")
	res, err := nsfw.Check(p)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func GS(r *bufio.Reader) {
	q := utils.Ask(r, "Query: ")
	res, err := googlesearch.Search(q)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}
