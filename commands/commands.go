package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"s/ai"
	"s/dapo"
	"s/decoder"
	"s/encoder"
	"s/googlesearch"
	"s/hash"
	"s/ipinfo"
	"s/kodepos"
	"s/music"
	"s/nikparser"
	"s/nsfw"
	"s/preset"
	"s/sflbypass"
	"s/shortener"
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
	fmt.Println(utils.Bld(utils.Wht("[ CHECK PTK ]")))
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

func SF(r *bufio.Reader) {
	u := utils.Ask(r, "Safelink URL: ")
	res := sflbypass.Do(u)
	res.Show()
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
	}
	back(r)
}

func AI(r *bufio.Reader) {
	client := ai.New()
	utils.Err("initializing the ai...")
	if err := client.Init(); err != nil {
		utils.Err("Error: " + err.Error())
		back(r)
		return
	}

	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ RFOUR ]")))
	fmt.Println(utils.Gry("Type 'exit' to leave. Type '.' on a line by itself to send multi-line message."))
	fmt.Println(utils.Div())

	for {
		fmt.Print(utils.Gry("You: "))
		var lines []string

		for {
			line, _ := r.ReadString('\n')
			line = strings.TrimRight(line, "\n")
			line = strings.TrimRight(line, "\r")

			if line == "." {
				break
			}
			if line == "exit" || line == "quit" || line == "0" {
				back(r)
				return
			}
			lines = append(lines, line)
			if line == "" {
				break
			}
		}

		prompt := strings.Join(lines, "\n")
		prompt = strings.TrimSpace(prompt)
		if prompt == "" {
			continue
		}

		fmt.Println(utils.Gry("rfour: "))
		reply, err := client.Chat(prompt)
		if err != nil {
			utils.Err("Error: " + err.Error())
			continue
		}
		fmt.Println(utils.Wht(reply))
		fmt.Println()
	}
}
