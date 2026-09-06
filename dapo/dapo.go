package dapo

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	"s/utils"
)

const baseURL = "https://dapo.kemendikdasmen.go.id"

var tk struct {
	T string
	U string
}

func token() (string, string, error) {
	if tk.T != "" {
		return tk.T, tk.U, nil
	}
	_, _ = utils.F(baseURL+"/", map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36",
	})
	bd, err := utils.F(baseURL+"/env.js", map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36",
	})
	if err != nil {
		return "", "", err
	}
	raw := string(bd)
	tk.T = regexp.MustCompile(`VITE_API_TOKEN["']?\s*:\s*["']([^"']+)["']`).FindStringSubmatch(raw)[1]
	tk.U = regexp.MustCompile(`VITE_STRAPI_URL["']?\s*:\s*["']([^"']+)["']`).FindStringSubmatch(raw)[1]
	return tk.T, tk.U, nil
}

type ProgressResult struct {
	Title string
	Data  []ProgressItem
}

type ProgressItem struct {
	Code   string
	Name   string
	Total  string
	Status string
}

func fetchAPI(p string, q map[string]string) (map[string]interface{}, error) {
	t, u, err := token()
	if err != nil {
		return nil, err
	}
	qq := ""
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if qq != "" {
			qq += "&"
		}
		qq += k + "=" + q[k]
	}
	if qq != "" {
		p += "?" + qq
	}
	bd, err := utils.F(u+p, map[string]string{
		"Authorization": "Bearer " + t,
	})
	if err != nil {
		return nil, err
	}
	var o map[string]interface{}
	if err := json.Unmarshal(bd, &o); err != nil {
		return nil, err
	}
	return o, nil
}

func ProgressProvinsi(j, s string) (*ProgressResult, error) {
	q := map[string]string{}
	if j != "" {
		q["jenjang"] = j
	}
	if s != "" {
		q["status_sekolah"] = s
	}
	raw, err := fetchAPI("/api/progress-pengiriman/provinsi", q)
	if err != nil {
		return nil, err
	}
	r := &ProgressResult{Title: "PROVINCE PROGRESS"}
	r.Data = parseProgress(raw)
	return r, nil
}

func ProgressKabupaten(k, j, s string) (*ProgressResult, error) {
	raw, err := fetchAPI("/api/progress-pengiriman/kabupaten", map[string]string{
		"kode_provinsi":  k,
		"jenjang":        j,
		"status_sekolah": s,
	})
	if err != nil {
		return nil, err
	}
	r := &ProgressResult{Title: "REGENCY PROGRESS"}
	r.Data = parseProgress(raw)
	return r, nil
}

func ProgressKecamatan(k, j, s string) (*ProgressResult, error) {
	raw, err := fetchAPI("/api/progress-pengiriman/kecamatan", map[string]string{
		"kode_kabupaten": k,
		"jenjang":         j,
		"status_sekolah":  s,
	})
	if err != nil {
		return nil, err
	}
	r := &ProgressResult{Title: "DISTRICT PROGRESS"}
	r.Data = parseProgress(raw)
	return r, nil
}

func ProgressSekolah(k, j string) (*ProgressResult, error) {
	raw, err := fetchAPI("/api/progress-pengiriman/kecamatan/school", map[string]string{
		"kode_kecamatan": k,
		"jenjang":         j,
	})
	if err != nil {
		return nil, err
	}
	r := &ProgressResult{Title: "SCHOOL PROGRESS"}
	r.Data = parseProgress(raw)
	return r, nil
}

func parseProgress(raw map[string]interface{}) []ProgressItem {
	var items []ProgressItem
	data, _ := raw["data"].([]interface{})
	for _, d := range data {
		dd, _ := d.(map[string]interface{})
		items = append(items, ProgressItem{
			Code:   str(dd["kode"]),
			Name:   str(dd["nama"]),
			Total:  str(dd["total"]),
			Status: str(dd["status"]),
		})
	}
	return items
}

func (r *ProgressResult) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ " + r.Title + " ]")))
	fmt.Println(utils.Div())

	if len(r.Data) == 0 {
		fmt.Println(utils.Gry("No data"))
		return
	}

	for i, d := range r.Data {
		fmt.Println(utils.Gry(fmt.Sprintf("[%d] %s", i+1, d.Name)))
		if d.Code != "-" {
			fmt.Println(utils.Gry("Code: ") + utils.Wht(d.Code))
		}
		if d.Total != "-" {
			fmt.Println(utils.Gry("Total: ") + utils.Wht(d.Total))
		}
		if d.Status != "-" {
			fmt.Println(utils.Gry("Status: ") + utils.Wht(d.Status))
		}
		fmt.Println()
	}
}

type SchoolSearch struct {
	Query string
	Total int
	Data  []SchoolItem
}

type SchoolItem struct {
	NPSN     string
	Name     string
	Level    string
	Province string
	City     string
	District string
	Status   string
}

func CariSekolah(q string) (*SchoolSearch, error) {
	t, u, err := token()
	if err != nil {
		return nil, err
	}
	bd, err := utils.F(u+"/api/detail-sekolah/search?q="+q, map[string]string{
		"Authorization": "Bearer " + t,
	})
	if err != nil {
		return nil, err
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(bd, &raw); err != nil {
		return nil, err
	}

	r := &SchoolSearch{Query: q}
	data, _ := raw["data"].([]interface{})
	for _, d := range data {
		dd, _ := d.(map[string]interface{})
		r.Data = append(r.Data, SchoolItem{
			NPSN:     str(dd["npsn"]),
			Name:     str(dd["nama"]),
			Level:    str(dd["jenjang"]),
			Province: str(dd["provinsi"]),
			City:     str(dd["kabupaten"]),
			District: str(dd["kecamatan"]),
			Status:   str(dd["status_sekolah"]),
		})
	}
	r.Total = len(r.Data)
	return r, nil
}

func (r *SchoolSearch) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ SCHOOL SEARCH ]")))
	fmt.Println(utils.Gry(fmt.Sprintf("Query: %s", r.Query)))
	fmt.Println(utils.Gry(fmt.Sprintf("Total: %d", r.Total)))
	fmt.Println(utils.Div())

	if r.Total == 0 {
		fmt.Println(utils.Gry("No results"))
		return
	}

	for i, s := range r.Data {
		fmt.Println(utils.Gry(fmt.Sprintf("[%d] %s", i+1, s.Name)))
		if s.NPSN != "-" {
			fmt.Println(utils.Gry("NPSN: ") + utils.Wht(s.NPSN))
		}
		if s.Level != "-" {
			fmt.Println(utils.Gry("Level: ") + utils.Wht(s.Level))
		}
		if s.Status != "-" {
			fmt.Println(utils.Gry("Status: ") + utils.Wht(s.Status))
		}
		if s.District != "-" {
			fmt.Println(utils.Gry("District: ") + utils.Wht(s.District))
		}
		if s.City != "-" {
			fmt.Println(utils.Gry("City: ") + utils.Wht(s.City))
		}
		if s.Province != "-" {
			fmt.Println(utils.Gry("Province: ") + utils.Wht(s.Province))
		}
		fmt.Println()
	}
}

type SchoolInfo struct {
	NPSN     string
	Name     string
	Level    string
	Status   string
	Address  string
	Village  string
	District string
	City     string
	Province string
	Postal   string
	Principal string
	Foundation string
	Students int
	Male     int
	Female   int
	Teachers int
	Staff    int
	Classrooms int
	Accreditation string
	Electricity string
	Internet   string
	Provider   string
	Bandwidth  string
	Latitude   string
	Longitude  string
	Updated    string
	Semester   string
}

func InfoSekolah(n string) (*SchoolInfo, error) {
	t, u, err := token()
	if err != nil {
		return nil, err
	}
	bd, err := utils.F(u+"/api/detail-sekolah?npsn="+n, map[string]string{
		"Authorization": "Bearer " + t,
	})
	if err != nil {
		return nil, err
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(bd, &raw); err != nil {
		return nil, err
	}

	data, _ := raw["data"].([]interface{})
	if len(data) == 0 {
		return nil, fmt.Errorf("school not found")
	}

	dd, _ := data[0].(map[string]interface{})

	return &SchoolInfo{
		NPSN:         str(dd["npsn"]),
		Name:         str(dd["nama"]),
		Level:        str(dd["jenjang"]),
		Status:       str(dd["status_sekolah"]),
		Address:      str(dd["alamat_jalan"]),
		Village:      str(dd["desa_kelurahan"]),
		District:     str(dd["kecamatan"]),
		City:         str(dd["kabupaten"]),
		Province:     str(dd["provinsi"]),
		Postal:       str(dd["kode_pos"]),
		Principal:    str(dd["nama_kepsek"]),
		Foundation:   str(dd["nama_yayasan"]),
		Students:     intNum(dd["pd"]),
		Male:         intNum(dd["pd_l"]),
		Female:       intNum(dd["pd_p"]),
		Teachers:     intNum(dd["jum_guru"]),
		Staff:        intNum(dd["jum_tendik"]),
		Classrooms:   intNum(dd["ruang_kelas"]),
		Accreditation: str(dd["akreditasi"]),
		Electricity:  str(dd["sumber_listrik"]),
		Internet:     str(dd["akses_internet"]),
		Provider:     str(dd["internet_provider"]),
		Bandwidth:    str(dd["internet_bandwidth"]),
		Latitude:     str(dd["lintang"]),
		Longitude:    str(dd["bujur"]),
		Updated:      str(dd["tanggal_update"]),
		Semester:     str(dd["semester"]),
	}, nil
}

func (r *SchoolInfo) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ SCHOOL INFO ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Name: ") + utils.Wht(r.Name))
	fmt.Println(utils.Gry("NPSN: ") + utils.Wht(r.NPSN))
	fmt.Println(utils.Gry("Level: ") + utils.Wht(r.Level))
	fmt.Println(utils.Gry("Status: ") + utils.Wht(r.Status))
	fmt.Println(utils.Gry("Accreditation: ") + utils.Wht(r.Accreditation))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Address: ") + utils.Wht(r.Address))
	fmt.Println(utils.Gry("Village: ") + utils.Wht(r.Village))
	fmt.Println(utils.Gry("District: ") + utils.Wht(r.District))
	fmt.Println(utils.Gry("City: ") + utils.Wht(r.City))
	fmt.Println(utils.Gry("Province: ") + utils.Wht(r.Province))
	fmt.Println(utils.Gry("Postal Code: ") + utils.Wht(r.Postal))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Principal: ") + utils.Wht(r.Principal))
	fmt.Println(utils.Gry("Foundation: ") + utils.Wht(r.Foundation))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Students: ") + utils.Wht(fmt.Sprintf("%d", r.Students)))
	fmt.Println(utils.Gry("Male: ") + utils.Wht(fmt.Sprintf("%d", r.Male)))
	fmt.Println(utils.Gry("Female: ") + utils.Wht(fmt.Sprintf("%d", r.Female)))
	fmt.Println(utils.Gry("Teachers: ") + utils.Wht(fmt.Sprintf("%d", r.Teachers)))
	fmt.Println(utils.Gry("Staff: ") + utils.Wht(fmt.Sprintf("%d", r.Staff)))
	fmt.Println(utils.Gry("Classrooms: ") + utils.Wht(fmt.Sprintf("%d", r.Classrooms)))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Electricity: ") + utils.Wht(r.Electricity))
	fmt.Println(utils.Gry("Internet: ") + utils.Wht(r.Internet))
	fmt.Println(utils.Gry("Provider: ") + utils.Wht(r.Provider))
	fmt.Println(utils.Gry("Bandwidth: ") + utils.Wht(r.Bandwidth))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Latitude: ") + utils.Wht(r.Latitude))
	fmt.Println(utils.Gry("Longitude: ") + utils.Wht(r.Longitude))
	fmt.Println(utils.Gry("Updated: ") + utils.Wht(r.Updated))
	fmt.Println(utils.Gry("Semester: ") + utils.Wht(r.Semester))
}

func str(x interface{}) string {
	if x == nil {
		return "-"
	}
	s, _ := x.(string)
	if s == "" {
		return "-"
	}
	return s
}

func intNum(x interface{}) int {
	if x == nil {
		return 0
	}
	switch v := x.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		var n int
		fmt.Sscanf(v, "%d", &n)
		return n
	}
	return 0
}
