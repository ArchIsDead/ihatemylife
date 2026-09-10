package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"s/akinator"
	"s/bypass"
	"s/dapo"
	"s/decoder"
	"s/encoder"
	"s/ffstalk"
	"s/googlesearch"
	"s/hash"
	"s/ipinfo"
	"s/kodepos"
	"s/music"
	"s/nikparser"
	"s/nsfw"
	"s/pinstalk"
	"s/preset"
	"s/shortener"
	"s/simpkb"
	"s/songfinder"
	"s/tracemoe"
	"s/ttstalk"
	"s/utils"
	"s/web2apk"
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
	fmt.Println(utils.Bld(utils.Wht("[ CHECK PTK & GTK ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("1. Search PTK / GTK"))
	fmt.Println(utils.Gry("2. List Provinces"))
	fmt.Println(utils.Gry("3. List Cities (by Province)"))
	fmt.Println(utils.Div())

	opt := utils.Ask(r, "Option [1]: ")
	if opt == "" {
		opt = "1"
	}

	switch opt {
	case "1":
		fmt.Println(utils.Div())
		fmt.Println(utils.Gry("Examples:"))
		fmt.Println(utils.Gry("  Keyword: Novy"))
		fmt.Println(utils.Gry("  Province Code: 15 (Jambi)"))
		fmt.Println(utils.Gry("  City Code: 1505 (Kab. Kerinci)"))
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
	case "2":
		res, err := simpkb.Provinsi()
		if err != nil {
			utils.Err("Error: " + err.Error())
			back(r)
			return
		}
		res.Show()
	case "3":
		prov := utils.Ask(r, "Province Code: ")
		res, err := simpkb.Kota(prov)
		if err != nil {
			utils.Err("Error: " + err.Error())
			back(r)
			return
		}
		res.Show()
	}
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
	p := utils.Ask(r, "Image Path or URL: ")
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

func SH(r *bufio.Reader) {
	u := utils.Ask(r, "URL: ")
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Service:"))
	fmt.Println(utils.Gry("  all   - All services"))
	fmt.Println(utils.Gry("  uto   - u.to"))
	fmt.Println(utils.Gry("  walee - wal.ee"))
	fmt.Println(utils.Div())
	service := utils.Ask(r, "Service: ")
	res, err := shortener.Do(u, service)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func BP(r *bufio.Reader) {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ BYPASS LINK ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("1. BypassTools (Ad-Link)"))
	fmt.Println(utils.Gry("2. BypassLink (SFL)"))
	fmt.Println(utils.Div())

	method := utils.Ask(r, "Method [1]: ")
	if method == "" {
		method = "1"
	}
	u := utils.Ask(r, "URL: ")

	res, err := bypass.Bypass(u, method)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}

	mname := "BypassTools"
	if method == "2" {
		mname = "BypassLink"
	}

	bypass.Show(u, res, mname)
	back(r)
}

func HA(r *bufio.Reader) {
	input := utils.Ask(r, "Text: ")
	res := hash.Do(input)
	res.Show()
	back(r)
}

func EN(r *bufio.Reader) {
	input := utils.Ask(r, "Text: ")
	res := encoder.Do(input)
	res.Show()
	back(r)
}

func DE(r *bufio.Reader) {
	input := utils.Ask(r, "Text: ")
	res := decoder.Do(input)
	res.Show()
	back(r)
}

func PR(r *bufio.Reader) {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ PRESETS ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("1. Add Banner"))
	fmt.Println(utils.Gry("2. Select Banner"))
	fmt.Println(utils.Gry("3. List Banners"))
	fmt.Println(utils.Gry("4. Remove Banner"))
	fmt.Println(utils.Gry("5. Add Music"))
	fmt.Println(utils.Gry("6. Select Music"))
	fmt.Println(utils.Gry("7. List Music"))
	fmt.Println(utils.Gry("8. Remove Music"))
	fmt.Println(utils.Gry("9. Toggle Music"))
	fmt.Println(utils.Gry("10. Set Volume"))
	fmt.Println(utils.Div())

	opt := utils.Ask(r, "Option: ")

	switch opt {
	case "1":
		name := utils.Ask(r, "Preset Name: ")
		text := utils.Ask(r, "Banner Text or URL: ")
		if strings.HasPrefix(text, "http://") || strings.HasPrefix(text, "https://") {
			bd, err := utils.F(text, nil)
			if err != nil {
				utils.Err("Error: " + err.Error())
				back(r)
				return
			}
			text = string(bd)
		}
		if err := preset.AddBanner(name, text); err != nil {
			utils.Err(err.Error())
		} else {
			utils.Err("Added")
		}
	case "2":
		name := utils.Ask(r, "Preset Name: ")
		if err := preset.SelectBanner(name); err != nil {
			utils.Err(err.Error())
		} else {
			utils.Err("Selected")
		}
	case "3":
		for i, b := range preset.ListBanners() {
			fmt.Println(utils.Gry(fmt.Sprintf("[%d] %s", i+1, b.Name)))
		}
	case "4":
		name := utils.Ask(r, "Preset Name: ")
		if err := preset.RemoveBanner(name); err != nil {
			utils.Err(err.Error())
		} else {
			utils.Err("Removed")
		}
	case "5":
		name := utils.Ask(r, "Preset Name: ")
		path := utils.Ask(r, "File Path or URL: ")
		if err := preset.AddMusic(name, path); err != nil {
			utils.Err(err.Error())
		} else {
			utils.Err("Added")
		}
	case "6":
		name := utils.Ask(r, "Preset Name: ")
		if err := preset.SelectMusic(name); err != nil {
			utils.Err(err.Error())
		} else {
			if path, ok := preset.GetSelectedMusic(); ok {
				music.Play(path)
			}
			utils.Err("Selected")
		}
	case "7":
		for i, m := range preset.ListMusics() {
			fmt.Println(utils.Gry(fmt.Sprintf("[%d] %s - %s", i+1, m.Name, m.Path)))
		}
	case "8":
		name := utils.Ask(r, "Preset Name: ")
		if err := preset.RemoveMusic(name); err != nil {
			utils.Err(err.Error())
		} else {
			music.Stop()
			utils.Err("Removed")
		}
	case "9":
		on := preset.ToggleMusic()
		if on {
			if path, ok := preset.GetSelectedMusic(); ok {
				music.Play(path)
			}
			utils.Err("Music ON")
		} else {
			music.Stop()
			utils.Err("Music OFF")
		}
	case "10":
		v := utils.Ask(r, "Volume (0-100): ")
		vi, _ := strconv.Atoi(v)
		preset.SetVolume(vi)
		music.SetVolume(vi)
		if music.IsPlaying() {
			if path, ok := preset.GetSelectedMusic(); ok {
				music.Play(path)
			}
		}
		utils.Err("Volume set to " + v)
	}
	back(r)
}

func SF(r *bufio.Reader) {
	p := utils.Ask(r, "Audio Path or URL: ")
	res, err := songfinder.Identify(p)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func TT(r *bufio.Reader) {
	u := utils.Ask(r, "Username: ")
	res, err := ttstalk.Get(u)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func PN(r *bufio.Reader) {
	u := utils.Ask(r, "Username: ")
	res, err := pinstalk.Get(u)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func FF(r *bufio.Reader) {
	u := utils.Ask(r, "UID: ")
	res, err := ffstalk.Get(u)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func WA(r *bufio.Reader) {
	client := &web2apk.Client{}

	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ WEB2APK ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Example:"))
	fmt.Println(utils.Gry("  App Name: MyApp"))
	fmt.Println(utils.Gry("  Package: com.myapp.dev"))
	fmt.Println(utils.Gry("  URL: https://example.com"))
	fmt.Println(utils.Div())

	var req web2apk.BuildRequest

	req.AppName = utils.Ask(r, "App Name: ")
	req.PackageName = utils.Ask(r, "Package Name: ")
	req.URL = utils.Ask(r, "URL: ")
	req.VersionName = utils.Ask(r, "Version Name [1.0]: ")
	if req.VersionName == "" {
		req.VersionName = "1.0"
	}
	req.VersionCode = utils.Ask(r, "Version Code [1]: ")
	if req.VersionCode == "" {
		req.VersionCode = "1"
	}
	req.Orientation = utils.Ask(r, "Orientation [auto]: ")
	if req.Orientation == "" {
		req.Orientation = "auto"
	}
	req.SplashType = utils.Ask(r, "Splash Type [image]: ")
	if req.SplashType == "" {
		req.SplashType = "image"
	}
	req.IconPath = utils.Ask(r, "Icon Path or URL: ")
	req.SplashPath = utils.Ask(r, "Splash Path or URL: ")

	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Permission Presets:"))
	for k, v := range web2apk.PermissionPresets {
		fmt.Println(utils.Gry("  " + k + " - " + v.Name))
	}
	fmt.Println(utils.Gry("  7 - Manual"))
	fmt.Println(utils.Gry("  empty - None"))
	fmt.Println(utils.Div())

	pc := utils.Ask(r, "Preset (1-7, empty for none): ")

	if preset, ok := web2apk.PermissionPresets[pc]; ok {
		req.Perms = preset.Perms
	} else if pc == "7" {
		for _, p := range web2apk.Permissions {
			safe := "Safe"
			if !p.Safe {
				safe = "Sensitive"
			}
			fmt.Println(utils.Gry(fmt.Sprintf("  [%s] %s - %s (%s)", p.ID, p.Name, p.Desc, safe)))
		}
		input := utils.Ask(r, "Permission IDs (comma separated): ")
		ids := strings.Split(input, ",")
		var chosen []string
		for _, id := range ids {
			id = strings.TrimSpace(id)
			for _, p := range web2apk.Permissions {
				if p.ID == id {
					chosen = append(chosen, p.Name)
				}
			}
		}
		req.Perms = chosen
	} else {
		req.Perms = nil
	}

	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Building... Please wait"))
	fmt.Println(utils.Div())

	res, err := client.Build(req)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}
	res.Show()
	back(r)
}

func AK(r *bufio.Reader) {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ AKINATOR ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("Theme:"))
	fmt.Println(utils.Gry("  1. Characters (default)"))
	fmt.Println(utils.Gry("  2. Animals"))
	fmt.Println(utils.Gry("  3. Objects"))
	fmt.Println(utils.Div())

	themeOpt := utils.Ask(r, "Theme [1]: ")
	theme := "characters"
	switch themeOpt {
	case "2":
		theme = "animals"
	case "3":
		theme = "objects"
	}

	child := utils.Ask(r, "Child Mode? (true/false) [false]: ")
	childMode := strings.ToLower(child) == "true"

	session, err := akinator.Start(theme, childMode)
	if err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}

	fmt.Println(utils.Div())
	fmt.Println(utils.Grn("Game started!"))
	fmt.Println(utils.Div())

	for {
		if session.Question == "" {
			break
		}

		fmt.Println(utils.Bld(utils.Wht(session.Question)))
		fmt.Println(utils.Gry(fmt.Sprintf("Step: %d | Progression: %.1f%%", session.Step, session.Progression)))
		fmt.Println(utils.Div())
		fmt.Println(utils.Gry("  1. Yes"))
		fmt.Println(utils.Gry("  2. No"))
		fmt.Println(utils.Gry("  3. I Don't Know"))
		fmt.Println(utils.Gry("  4. Probably"))
		fmt.Println(utils.Gry("  5. Probably Not"))
		fmt.Println(utils.Gry("  6. Back"))
		fmt.Println(utils.Gry("  7. Exclude"))
		fmt.Println(utils.Gry("  0. Quit"))
		fmt.Println(utils.Div())

		ans := utils.Ask(r, "> ")
		ans = strings.TrimSpace(ans)

		if ans == "0" || ans == "quit" || ans == "exit" {
			break
		}

		var ansStr string
		switch ans {
		case "1":
			ansStr = "yes"
		case "2":
			ansStr = "no"
		case "3":
			ansStr = "idk"
		case "4":
			ansStr = "probably"
		case "5":
			ansStr = "probably not"
		case "6":
			res, err := akinator.Back(session)
			if err != nil {
				utils.Err("Error: " + err.Error())
				continue
			}
			session.Question = res.Question
			continue
		case "7":
			res, err := akinator.Exclude(session)
			if err != nil {
				utils.Err("Error: " + err.Error())
				continue
			}
			session.Question = res.Question
			continue
		default:
			ansStr = ans
		}

		res, err := akinator.Answer(session, ansStr)
		if err != nil {
			utils.Err("Error: " + err.Error())
			continue
		}

		if res.Won {
			fmt.Println(utils.Div())
			fmt.Println(utils.Bld(utils.Wht("[ AKINATOR GUESSED ]")))
			fmt.Println(utils.Div())
			fmt.Println(utils.Gry("Name: ") + utils.Wht(res.Name))
			if res.Description != "" {
				fmt.Println(utils.Gry("Description: ") + utils.Wht(res.Description))
			}
			if res.Photo != "" {
				fmt.Println(utils.Gry("Photo: ") + utils.Wht(res.Photo))
			}
			if res.Pseudo != "" {
				fmt.Println(utils.Gry("Pseudo: ") + utils.Wht(res.Pseudo))
			}
			break
		}

		session.Question = res.Question
	}
	back(r)
}
