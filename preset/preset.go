package preset

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Banner struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type Music struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Data struct {
	Banners  []Banner `json:"banners"`
	Musics   []Music  `json:"musics"`
	Selected string   `json:"selected_banner"`
	MusicOn  bool     `json:"music_on"`
	MusicSel string   `json:"selected_music"`
	Volume   int      `json:"volume"`
}

var dir string

func init() {
	home, _ := os.UserHomeDir()
	dir = filepath.Join(home, "suic1.de")
	os.MkdirAll(dir, 0755)
}

func dataPath() string {
	return filepath.Join(dir, "presets.json")
}

func Load() *Data {
	d := &Data{}

	if _, err := os.Stat(dataPath()); os.IsNotExist(err) {
		d.Banners = append(d.Banners, Banner{Name: "default", Text: defaultBanner})
		d.Musics = append(d.Musics, Music{Name: "1tap", Path: "https://files.catbox.moe/glyr3n.mp3"})
		d.Selected = "default"
		d.MusicOn = true
		d.MusicSel = "1tap"
		d.Volume = 100
		Save(d)
		return d
	}

	b, err := os.ReadFile(dataPath())
	if err != nil {
		return d
	}

	json.Unmarshal(b, d)

	if d.Volume == 0 {
		d.Volume = 100
		Save(d)
	}

	return d
}

func Save(d *Data) {
	b, _ := json.MarshalIndent(d, "", "  ")
	os.WriteFile(dataPath(), b, 0644)
}

func AddBanner(name, text string) error {
	d := Load()
	for _, b := range d.Banners {
		if b.Name == name {
			return fmt.Errorf("banner name exists")
		}
	}
	d.Banners = append(d.Banners, Banner{Name: name, Text: text})
	Save(d)
	return nil
}

func RemoveBanner(name string) error {
	d := Load()
	for i, b := range d.Banners {
		if b.Name == name {
			d.Banners = append(d.Banners[:i], d.Banners[i+1:]...)
			if d.Selected == name {
				d.Selected = "default"
			}
			Save(d)
			return nil
		}
	}
	return fmt.Errorf("banner not found")
}

func SelectBanner(name string) error {
	d := Load()
	for _, b := range d.Banners {
		if b.Name == name {
			d.Selected = name
			Save(d)
			return nil
		}
	}
	return fmt.Errorf("banner not found")
}

func GetSelectedBanner() string {
	d := Load()
	for _, b := range d.Banners {
		if b.Name == d.Selected {
			return b.Text
		}
	}
	return defaultBanner
}

func AddMusic(name, path string) error {
	d := Load()
	for _, m := range d.Musics {
		if m.Name == name {
			return fmt.Errorf("music name exists")
		}
	}
	d.Musics = append(d.Musics, Music{Name: name, Path: path})
	Save(d)
	return nil
}

func RemoveMusic(name string) error {
	d := Load()
	for i, m := range d.Musics {
		if m.Name == name {
			d.Musics = append(d.Musics[:i], d.Musics[i+1:]...)
			if d.MusicSel == name {
				d.MusicSel = ""
			}
			Save(d)
			return nil
		}
	}
	return fmt.Errorf("music not found")
}

func SelectMusic(name string) error {
	d := Load()
	for _, m := range d.Musics {
		if m.Name == name {
			d.MusicSel = name
			d.MusicOn = true
			Save(d)
			return nil
		}
	}
	return fmt.Errorf("music not found")
}

func ToggleMusic() bool {
	d := Load()
	d.MusicOn = !d.MusicOn
	Save(d)
	return d.MusicOn
}

func GetSelectedMusic() (string, bool) {
	d := Load()
	if !d.MusicOn {
		return "", false
	}
	for _, m := range d.Musics {
		if m.Name == d.MusicSel {
			return m.Path, true
		}
	}
	return "", false
}

func ListBanners() []Banner {
	d := Load()
	return d.Banners
}

func ListMusics() []Music {
	d := Load()
	return d.Musics
}

func CurrentBannerName() string {
	d := Load()
	return d.Selected
}

func CurrentMusicName() string {
	d := Load()
	return d.MusicSel
}

func MusicEnabled() bool {
	d := Load()
	return d.MusicOn
}

func SetVolume(v int) {
	if v < 0 {
		v = 0
	}
	if v > 100 {
		v = 100
	}
	d := Load()
	d.Volume = v
	Save(d)
}

func GetVolume() int {
	d := Load()
	if d.Volume == 0 {
		d.Volume = 100
		Save(d)
	}
	return d.Volume
}

const defaultBanner = `
  _________     .__       ____         .___      
 /   _____/__ __|__| ____/_   |      __| _/____  
 \_____  \|  |  \  |/ ___\|   |     / __ |/ __ \ 
 /        \  |  /  \  \___|   |    / /_/ \  ___/ 
/_______  /____/|__|\___  >___| /\ \____ |\___  >
        \/              \/      \/      \/    \/ 
`
