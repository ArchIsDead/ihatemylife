package nikparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"s/utils"
)

type Region struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Result struct {
	NIK      string
	Gender   string
	Birth    string
	Province string
	City     string
	District string
	Age      int
	Serial   string
}

func Parse(n string) (*Result, error) {
	if len(n) != 16 {
		return nil, errors.New("invalid NIK")
	}

	ps, err := fetchRegions("https://emsifa.github.io/api-wilayah-indonesia/api/provinces.json")
	if err != nil {
		return nil, err
	}
	pm := map[string]string{}
	for _, p := range ps {
		pm[p.ID] = p.Name
	}
	if pm[n[0:2]] == "" {
		return nil, errors.New("invalid province")
	}

	rs, err := fetchRegions("https://emsifa.github.io/api-wilayah-indonesia/api/regencies/" + n[0:2] + ".json")
	if err != nil {
		return nil, err
	}
	rm := map[string]string{}
	for _, r := range rs {
		rm[r.ID] = r.Name
	}
	if rm[n[0:4]] == "" {
		return nil, errors.New("invalid city")
	}

	ds, err := fetchRegions("https://emsifa.github.io/api-wilayah-indonesia/api/districts/" + n[0:4] + ".json")
	if err != nil {
		return nil, err
	}
	dm := map[string]string{}
	for _, d := range ds {
		dm[d.ID[:len(d.ID)-1]] = d.Name
	}
	if dm[n[0:6]] == "" {
		return nil, errors.New("invalid district")
	}

	day, _ := strconv.Atoi(n[6:8])
	mo, _ := strconv.Atoi(n[8:10])
	yc, _ := strconv.Atoi(n[10:12])

	gender := "MALE"
	birthDay := day
	if day > 40 {
		gender = "FEMALE"
		birthDay = day - 40
	}

	year := 1900 + yc
	if yc < time.Now().Year()%100 {
		year = 2000 + yc
	}

	birth := time.Date(year, time.Month(mo), birthDay, 0, 0, 0, 0, time.UTC)
	age := int(time.Since(birth).Hours() / 24 / 365)

	return &Result{
		NIK:      n,
		Gender:   gender,
		Birth:    birth.Format("02/01/2006"),
		Province: pm[n[0:2]],
		City:     rm[n[0:4]],
		District: dm[n[0:6]],
		Age:      age,
		Serial:   n[12:16],
	}, nil
}

func fetchRegions(url string) ([]Region, error) {
	bd, err := utils.F(url, map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36",
	})
	if err != nil {
		return nil, err
	}
	var o []Region
	if err := json.Unmarshal(bd, &o); err != nil {
		return nil, err
	}
	return o, nil
}

func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ NIK PARSER ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("NIK: ") + utils.Wht(r.NIK))
	fmt.Println(utils.Gry("Gender: ") + utils.Wht(r.Gender))
	fmt.Println(utils.Gry("Birth: ") + utils.Wht(r.Birth))
	fmt.Println(utils.Gry("Age: ") + utils.Wht(fmt.Sprintf("%d", r.Age)))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Province: ") + utils.Wht(r.Province))
	fmt.Println(utils.Gry("City: ") + utils.Wht(r.City))
	fmt.Println(utils.Gry("District: ") + utils.Wht(r.District))
	fmt.Println(utils.Gry("Serial: ") + utils.Wht(r.Serial))
}
