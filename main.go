package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"s/commands"
	"s/music"
	"s/preset"
	"s/utils"
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

func autoUpdate() {
	cmd := exec.Command("git", "pull")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Run()
}

func startMusic() {
	vol := preset.GetVolume()
	music.SetVolume(vol)
	if path, ok := preset.GetSelectedMusic(); ok && path != "" {
		music.Play(path)
	}
}

func showMain(r *bufio.Reader) {
	clear()
	banner := preset.GetSelectedBanner()
	utils.ShowBanner(banner)

	for {
		utils.Menu()
		fmt.Print(utils.Prompt("\n> "))
		x, _ := r.ReadString('\n')
		x = strings.TrimSpace(x)
		x = strings.TrimLeft(x, "0")
		if x == "" {
			x = "0"
		}

		switch x {
		case "1":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.N(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "2":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.D(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "3":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.U(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "4":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.CS(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "5":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.IS(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "6":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.IP(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "7":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.CIP()
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "8":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.KP(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "9":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.NP(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "10":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.WZ(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "11":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.TM(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "12":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.NS(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "13":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.GS(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "14":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.SH(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "15":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.BP(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "16":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.HA(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "17":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.EN(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "18":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.DE(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "19":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.PR(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "20":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.SF(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "21":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.TT(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "22":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.PN(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "23":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.FF(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "24":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.WA(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "25":
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			commands.AK(r)
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
		case "0", "exit", "quit":
			clear()
			music.Stop()
			utils.Err("Exit.")
			os.Exit(0)
		default:
			clear()
			utils.ShowBanner(preset.GetSelectedBanner())
			utils.Err("Unknown command.")
		}
	}
}

func main() {
	autoUpdate()
	startMusic()
	r := bufio.NewReader(os.Stdin)
	showMain(r)
}
