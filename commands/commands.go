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
	k := utils.P(r, "Keyword: ")
	p := utils.P(r, "Province: ")
	kk := utils.P(r, "City: ")
	dp := utils.P(r, "Dapodik 0/1/empty: ")
	ps := utils.P(r, "Passport 0/1/empty: ")
	pg := utils.P(r, "Page: ")
	res, err := simpkb.Cari(k, p, kk, ps, dp, pg)
	if err != nil {
		fmt.Println(utils.E("Error: " + err.Error()))
		return
	}
	utils.O(res)
}

func P() {
	res, err := simpkb.Provinsi()
	if err != nil {
		fmt.Println(utils.E("Error: " + err.Error()))
		return
	}
	utils.O(res)
}

func K(r *bufio.Reader) {
	res, err := simpkb.Kota()
	if err != nil {
		fmt.Println(utils.E("Error: " + err.Error()))
		return
	}
	utils.O(res)
}

func D(r *bufio.Reader) {
	m := utils.P(r, "Mode (province/regency/district/school): ")
	switch m {
	case "province":
		j := utils.P(r, "Level: ")
		st := utils.P(r, "Status: ")
		res, _ := dapo.ProgressProvinsi(j, st)
		utils.O(res)
	case "regency":
		k := utils.P(r, "Province Code: ")
		j := utils.P(r, "Level: ")
		st := utils.P(r, "Status: ")
		res, _ := dapo.ProgressKabupaten(k, j, st)
		utils.O(res)
	case "district":
		k := utils.P(r, "Regency Code: ")
		j := utils.P(r, "Level: ")
		st := utils.P(r, "Status: ")
		res, _ := dapo.ProgressKecamatan(k, j, st)
		utils.O(res)
	case "school":
		k := utils.P(r, "District Code: ")
		j := utils.P(r, "Level: ")
		res, _ := dapo.ProgressSekolah(k, j)
		utils.O(res)
	}
}

func U(r *bufio.Reader) {
	n := utils.P(r, "Username: ")
	m := utils.P(r, "Mode all/exist/notexist: ")
	rc := utils.P(r, "Rescan true/false: ")
	res, err := whatsmyname.Scan(n, m, rc)
	if err != nil {
		fmt.Println(utils.E("Error: " + err.Error()))
		return
	}
	utils.O(res)
}

func CS(r *bufio.Reader) {
	q := utils.P(r, "School Name: ")
	res, err := dapo.CariSekolah(q)
	if err != nil {
		fmt.Println(utils.E("Error: " + err.Error()))
		return
	}
	utils.O(res)
}

func IS(r *bufio.Reader) {
	n := utils.P(r, "NPSN: ")
	res, err := dapo.InfoSekolah(n)
	if err != nil {
		fmt.Println(utils.E("Error: " + err.Error()))
		return
	}
	utils.O(res)
}

func IP(r *bufio.Reader) {
	ip := utils.P(r, "IP Address: ")
	res, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		fmt.Println(utils.E("Error: " + err.Error()))
		return
	}
	defer res.Body.Close()
	var o map[string]interface{}
	json.NewDecoder(res.Body).Decode(&o)
	utils.O(o)
}

func CIP() {
	res, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		fmt.Println(utils.E("Error: " + err.Error()))
		return
	}
	defer res.Body.Close()
	var o map[string]interface{}
	json.NewDecoder(res.Body).Decode(&o)
	utils.O(o)
}

func DI() {
	utils.O(utils.Di())
}

func ST() {
	host, _ := os.Hostname()
	ips, _ := net.LookupIP(host)
	fmt.Println(utils.T("[ STATS ]"))
	fmt.Println(utils.G("Hostname: ") + utils.W(host))
	for _, ip := range ips {
		fmt.Println(utils.G("IP: ") + utils.W(ip.String()))
	}
	fmt.Println(utils.G("Uptime: ") + utils.W(utils.UpStr()))
}

func KP(r *bufio.Reader) {
	q := utils.P(r, "Postal Code / Area Name: ")
	pg := utils.P(r, "Page: ")
	pi, _ := strconv.Atoi(pg)
	if pi < 1 {
		pi = 1
	}
	res, tp, err := kodepos.Cari(q, pi)
	if err != nil {
		fmt.Println(utils.E("Error: " + err.Error()))
		return
	}
	fmt.Println(utils.A(fmt.Sprintf("Total: %d | Page: %d/%d", len(res), pi, tp)))
	utils.O(res)
}

func NP(r *bufio.Reader) {
	n := utils.P(r, "NIK: ")
	res, err := nikparser.Parse(n)
	if err != nil {
		fmt.Println(utils.E("Error: " + err.Error()))
		return
	}
	utils.O(res)
}

func PV(r *bufio.Reader) {
	n := utils.P(r, "Instagram Username: ")
	c := protonviewer.New()
	res, err := c.UserInfo(n)
	if err != nil {
		fmt.Println(utils.E("Error: " + err.Error()))
		return
	}
	id := c.ID(res)
	fmt.Println(utils.G("User ID: ") + utils.W(id))
	utils.O(res)
	if id != "" {
		st, _ := c.Stories(n)
		utils.O(st)
		hi, _ := c.Highlights(id)
		utils.O(hi)
	}
	ps, _ := c.Posts(n, "")
	utils.O(ps)
}
