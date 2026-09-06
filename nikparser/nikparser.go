package nikparser

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"s/utils"
)

type Region struct {
	ID string `json:"id"`
	Nm string `json:"name"`
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

func Parse(n string) (map[string]interface{}, error) {
	if len(n) != 16 {
		return nil, errors.New("invalid NIK")
	}
	ps, err := fetchRegions("https://emsifa.github.io/api-wilayah-indonesia/api/provinces.json")
	if err != nil {
		return nil, err
	}
	pm := map[string]string{}
	for _, p := range ps {
		pm[p.ID] = p.Nm
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
		rm[r.ID] = r.Nm
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
		dm[d.ID[:len(d.ID)-1]] = d.Nm
	}
	if dm[n[0:6]] == "" {
		return nil, errors.New("invalid district")
	}
	day, _ := strconv.Atoi(n[6:8])
	mo, _ := strconv.Atoi(n[8:10])
	yc, _ := strconv.Atoi(n[10:12])
	g := "MALE"
	bd := day
	if day > 40 {
		g = "FEMALE"
		bd = day - 40
	}
	yy := 1900 + yc
	if yc < time.Now().Year()%100 {
		yy = 2000 + yc
	}
	birth := time.Date(yy, time.Month(mo), bd, 0, 0, 0, 0, time.UTC)
	age := time.Since(birth).Hours() / 24 / 365
	return map[string]interface{}{
		"nik":      n,
		"gender":   g,
		"birth":    birth.Format("02/01/2006"),
		"province": pm[n[0:2]],
		"city":     rm[n[0:4]],
		"district": dm[n[0:6]],
		"age":      int(age),
		"serial":   n[12:16],
	}, nil
}
