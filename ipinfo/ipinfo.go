package ipinfo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"s/utils"
)

type Result struct {
	IP        string
	Success   bool
	Type      string
	Continent string
	Country   string
	Region    string
	City      string
	Lat       string
	Lon       string
	Org       string
	ISP       string
	Timezone  string
}

func Lookup(ip string) (*Result, error) {
	u := "https://ipwho.is/"
	if ip != "" {
		u += ip
	}
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	r := &Result{
		IP:      str(raw["ip"]),
		Success: raw["success"] == true,
		Type:    str(raw["type"]),
	}

	if c, ok := raw["continent"].(string); ok {
		r.Continent = c
	}
	if c, ok := raw["country"].(string); ok {
		r.Country = c
	}
	if c, ok := raw["region"].(string); ok {
		r.Region = c
	}
	if c, ok := raw["city"].(string); ok {
		r.City = c
	}
	r.Lat = fmt.Sprintf("%v", raw["latitude"])
	r.Lon = fmt.Sprintf("%v", raw["longitude"])

	if conn, ok := raw["connection"].(map[string]interface{}); ok {
		r.ISP = str(conn["isp"])
		r.Org = str(conn["org"])
	}
	if tz, ok := raw["timezone"].(map[string]interface{}); ok {
		r.Timezone = str(tz["id"])
	}

	return r, nil
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ IP INFO ]")))
	fmt.Println(utils.Div())

	if !r.Success {
		fmt.Println(utils.Red("Invalid IP"))
		return
	}

	fmt.Println(utils.Gry("IP: ") + utils.Wht(r.IP))
	if r.Type != "" {
		fmt.Println(utils.Gry("Type: ") + utils.Wht(r.Type))
	}
	if r.Continent != "" {
		fmt.Println(utils.Gry("Continent: ") + utils.Wht(r.Continent))
	}
	if r.Country != "" {
		fmt.Println(utils.Gry("Country: ") + utils.Wht(r.Country))
	}
	if r.Region != "" {
		fmt.Println(utils.Gry("Region: ") + utils.Wht(r.Region))
	}
	if r.City != "" {
		fmt.Println(utils.Gry("City: ") + utils.Wht(r.City))
	}
	if r.Lat != "<nil>" && r.Lat != "" {
		fmt.Println(utils.Gry("Latitude: ") + utils.Wht(r.Lat))
	}
	if r.Lon != "<nil>" && r.Lon != "" {
		fmt.Println(utils.Gry("Longitude: ") + utils.Wht(r.Lon))
	}
	if r.Org != "" {
		fmt.Println(utils.Gry("Org: ") + utils.Wht(r.Org))
	}
	if r.ISP != "" {
		fmt.Println(utils.Gry("ISP: ") + utils.Wht(r.ISP))
	}
	if r.Timezone != "" {
		fmt.Println(utils.Gry("Timezone: ") + utils.Wht(r.Timezone))
	}
}

func str(x interface{}) string {
	if x == nil {
		return ""
	}
	s, _ := x.(string)
	return s
}
