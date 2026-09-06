package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"s/dapo"
	"s/kodepos"
	"s/nikparser"
	"s/protonviewer"
	"s/simpkb"
	"s/utils"
	"s/whatsmyname"
)

func N(r *bufio.Reader) {
	k := utils.Ask(r, "Keyword: ")
	p := utils.Ask(r, "Province: ")
	kk := utils.Ask(r, "City: ")
	dp := utils.Ask(r, "Dapodik 0/1/empty: ")
	ps := utils.Ask(r, "Passport 0/1/empty: ")
	pg := utils.Ask(r, "Page: ")
	res, err := simpkb.Cari(k, p, kk, ps, dp, pg)
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	utils.PrintJSON(res)
}

func P() {
	res, err := simpkb.Provinsi()
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	utils.PrintJSON(res)
}

func K(r *bufio.Reader) {
	res, err := simpkb.Kota()
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	utils.PrintJSON(res)
}

func D(r *bufio.Reader) {
	m := utils.Ask(r, "Mode (province/regency/district/school): ")
	switch m {
	case "province":
		j := utils.Ask(r, "Level: ")
		st := utils.Ask(r, "Status: ")
		res, _ := dapo.ProgressProvinsi(j, st)
		utils.PrintJSON(res)
	case "regency":
		k := utils.Ask(r, "Province Code: ")
		j := utils.Ask(r, "Level: ")
		st := utils.Ask(r, "Status: ")
		res, _ := dapo.ProgressKabupaten(k, j, st)
		utils.PrintJSON(res)
	case "district":
		k := utils.Ask(r, "Regency Code: ")
		j := utils.Ask(r, "Level: ")
		st := utils.Ask(r, "Status: ")
		res, _ := dapo.ProgressKecamatan(k, j, st)
		utils.PrintJSON(res)
	case "school":
		k := utils.Ask(r, "District Code: ")
		j := utils.Ask(r, "Level: ")
		res, _ := dapo.ProgressSekolah(k, j)
		utils.PrintJSON(res)
	}
}

func U(r *bufio.Reader) {
	n := utils.Ask(r, "Username: ")
	m := utils.Ask(r, "Mode all/exist/notexist: ")
	rc := utils.Ask(r, "Rescan true/false: ")
	res, err := whatsmyname.Scan(n, m, rc)
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	utils.PrintJSON(res)
}

func CS(r *bufio.Reader) {
	q := utils.Ask(r, "School Name: ")
	res, err := dapo.CariSekolah(q)
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	utils.PrintJSON(res)
}

func IS(r *bufio.Reader) {
	n := utils.Ask(r, "NPSN: ")
	res, err := dapo.InfoSekolah(n)
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	utils.PrintJSON(res)
}

func IP(r *bufio.Reader) {
	ip := utils.Ask(r, "IP Address: ")
	res, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	defer res.Body.Close()
	var o map[string]interface{}
	json.NewDecoder(res.Body).Decode(&o)
	utils.PrintJSON(o)
}

func CIP() {
	res, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	defer res.Body.Close()
	var o map[string]interface{}
	json.NewDecoder(res.Body).Decode(&o)
	utils.PrintJSON(o)
}

func DI() {
	utils.PrintJSON(utils.DeviceInfo())
}

func ST() {
	host, _ := os.Hostname()
	ips, _ := net.LookupIP(host)
	fmt.Println(utils.Title("[ STATS ]"))
	fmt.Println(utils.Gry("Hostname: ") + utils.Wht(host))
	for _, ip := range ips {
		fmt.Println(utils.Gry("IP: ") + utils.Wht(ip.String()))
	}
	fmt.Println(utils.Gry("Uptime: ") + utils.Wht(utils.UpStr()))
}

func KP(r *bufio.Reader) {
	q := utils.Ask(r, "Postal Code / Area Name: ")
	pg := utils.Ask(r, "Page: ")
	pi, _ := strconv.Atoi(pg)
	if pi < 1 {
		pi = 1
	}
	res, tp, err := kodepos.Cari(q, pi)
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	fmt.Println(utils.Acc(fmt.Sprintf("Total: %d | Page: %d/%d", len(res), pi, tp)))
	utils.PrintJSON(res)
}

func NP(r *bufio.Reader) {
	n := utils.Ask(r, "NIK: ")
	res, err := nikparser.Parse(n)
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	utils.PrintJSON(res)
}

func PV(r *bufio.Reader) {
	n := utils.Ask(r, "Instagram Username: ")
	c := protonviewer.New()
	res, err := c.UserInfo(n)
	if err != nil {
		utils.Err("Error: " + err.Error())
		return
	}
	id := c.ID(res)
	fmt.Println(utils.Gry("User ID: ") + utils.Wht(id))
	utils.PrintJSON(res)
	if id != "" {
		st, _ := c.Stories(n)
		utils.PrintJSON(st)
		hi, _ := c.Highlights(id)
		utils.PrintJSON(hi)
	}
	ps, _ := c.Posts(n, "")
	utils.PrintJSON(ps)
}
