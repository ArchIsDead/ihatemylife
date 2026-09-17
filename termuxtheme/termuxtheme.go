package termuxtheme

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"s/utils"
)

func userDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".termux-theme")
}

func repoDir() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, ".termux-theme")
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		os.MkdirAll(filepath.Dir(target), 0755)
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		return err
	})
}

func syncRepo() error {
	src := repoDir()
	dst := userDir()

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("repo folder .termux-theme not found")
	}

	os.MkdirAll(dst, 0755)

	return copyDir(src, dst)
}

func SyncSilent() {
	src := repoDir()
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return
	}
	syncRepo()
}

func ensureInstalled() error {
	dst := userDir()
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		return syncRepo()
	}
	return nil
}

func runScript(name string, args ...string) error {
	script := filepath.Join(userDir(), name)
	if _, err := os.Stat(script); os.IsNotExist(err) {
		return fmt.Errorf("script not found: %s", name)
	}
	cmd := exec.Command("bash", append([]string{script}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func Install() error {
	if err := syncRepo(); err != nil {
		return err
	}
	return runScript("install.sh")
}

func Apply(theme string) error {
	if err := ensureInstalled(); err != nil {
		return err
	}
	return runScript("apply.sh", theme)
}

func Banner(action, arg string) error {
	if err := ensureInstalled(); err != nil {
		return err
	}
	if arg != "" {
		return runScript("banner.sh", action, arg)
	}
	return runScript("banner.sh", action)
}

func Disable() error {
	if err := ensureInstalled(); err != nil {
		return err
	}

	home, _ := os.UserHomeDir()
	termuxDir := filepath.Join(home, ".termux")
	os.MkdirAll(termuxDir, 0755)

	propFile := filepath.Join(termuxDir, "termux.properties")
	content := `use-black-ui = true
extra-keys = [['ESC','/','-','HOME','UP','END','PGUP'],['TAB','CTRL','ALT','LEFT','DOWN','RIGHT','PGDN']]
terminal-margin-horizontal = 3
terminal-margin-vertical = 3
bell-character = ignore
`

	os.WriteFile(propFile, []byte(content), 0644)

	colorsFile := filepath.Join(termuxDir, "colors.properties")
	os.Remove(colorsFile)

	zshrc := filepath.Join(home, ".zshrc")
	os.Remove(zshrc)

	cmd := exec.Command("termux-reload-settings")
	cmd.Run()

	return nil
}

func ListThemes() []string {
	entries, err := os.ReadDir(filepath.Join(userDir(), "themes"))
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out
}

func Menu(r *bufio.Reader) {
	if err := ensureInstalled(); err != nil {
		utils.Err("Error: " + err.Error())
		return
	}

	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ TERMUX THEME ]")))
	fmt.Println(utils.Div())
	fmt.Println(utils.Gry("1. Install (deps + oh-my-zsh)"))
	fmt.Println(utils.Gry("2. Apply Theme"))
	fmt.Println(utils.Gry("3. List Themes"))
	fmt.Println(utils.Gry("4. Banner Manager"))
	fmt.Println(utils.Gry("5. Sync from repo"))
	fmt.Println(utils.Gry("6. Disable Theme"))
	fmt.Println(utils.Div())

	opt := utils.Ask(r, "Option: ")

	switch opt {
	case "1":
		if err := Install(); err != nil {
			utils.Err("Error: " + err.Error())
			return
		}
		utils.Err("Installed")
	case "2":
		themes := ListThemes()
		if len(themes) == 0 {
			utils.Err("No themes found. Run install first.")
			return
		}
		for i, t := range themes {
			fmt.Println(utils.Gry(fmt.Sprintf("  %d. %s", i+1, t)))
		}
		name := utils.Ask(r, "Theme name: ")
		if err := Apply(name); err != nil {
			utils.Err("Error: " + err.Error())
			return
		}
		utils.Err("Applied: " + name)
	case "3":
		for _, t := range ListThemes() {
			fmt.Println(utils.Gry("  - " + t))
		}
	case "4":
		fmt.Println(utils.Gry("1. Show"))
		fmt.Println(utils.Gry("2. Set (from file)"))
		fmt.Println(utils.Gry("3. Edit"))
		fmt.Println(utils.Gry("4. Reset"))
		a := utils.Ask(r, "Option: ")
		switch a {
		case "1":
			Banner("show", "")
		case "2":
			f := utils.Ask(r, "File path: ")
			Banner("set", f)
		case "3":
			Banner("edit", "")
		case "4":
			Banner("reset", "")
		}
	case "5":
		if err := syncRepo(); err != nil {
			utils.Err("Error: " + err.Error())
			return
		}
		utils.Err("Synced from repo")
	case "6":
		if err := Disable(); err != nil {
			utils.Err("Error: " + err.Error())
			return
		}
		utils.Err("Theme disabled")
	}
}
