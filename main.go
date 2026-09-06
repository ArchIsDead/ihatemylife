package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"s/commands"
	"s/utils"
)

var banner = `
  _________     .__       ____         .___      
 /   _____/__ __|__| ____/_   |      __| _/____  
 \_____  \|  |  \  |/ ___\|   |     / __ |/ __ \ 
 /        \  |  /  \  \___|   |    / /_/ \  ___/ 
/_______  /____/|__|\___  >___| /\ \____ |\___  >
        \/              \/      \/      \/    \/ 
`

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

func showMain(r *bufio.Reader) {
	clear()
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
			utils.ShowBanner(banner)
			commands.N(r)
			clear()
			utils.ShowBanner(banner)
		case "2":
			clear()
			utils.ShowBanner(banner)
			commands.D(r)
			clear()
			utils.ShowBanner(banner)
		case "3":
			clear()
			utils.ShowBanner(banner)
			commands.U(r)
			clear()
			utils.ShowBanner(banner)
		case "4":
			clear()
			utils.ShowBanner(banner)
			commands.CS(r)
			clear()
			utils.ShowBanner(banner)
		case "5":
			clear()
			utils.ShowBanner(banner)
			commands.IS(r)
			clear()
			utils.ShowBanner(banner)
		case "6":
			clear()
			utils.ShowBanner(banner)
			commands.IP(r)
			clear()
			utils.ShowBanner(banner)
		case "7":
			clear()
			utils.ShowBanner(banner)
			commands.CIP()
			clear()
			utils.ShowBanner(banner)
		case "8":
			clear()
			utils.ShowBanner(banner)
			commands.KP(r)
			clear()
			utils.ShowBanner(banner)
		case "9":
			clear()
			utils.ShowBanner(banner)
			commands.NP(r)
			clear()
			utils.ShowBanner(banner)
		case "10":
			clear()
			utils.ShowBanner(banner)
			commands.WZ(r)
			clear()
			utils.ShowBanner(banner)
		case "11":
			clear()
			utils.ShowBanner(banner)
			commands.TM(r)
			clear()
			utils.ShowBanner(banner)
		case "12":
			clear()
			utils.ShowBanner(banner)
			commands.NS(r)
			clear()
			utils.ShowBanner(banner)
		case "13":
			clear()
			utils.ShowBanner(banner)
			commands.GS(r)
			clear()
			utils.ShowBanner(banner)
		case "14":
			clear()
			utils.ShowBanner(banner)
			commands.UN(r)
			clear()
			utils.ShowBanner(banner)
		case "15":
			clear()
			utils.ShowBanner(banner)
			commands.SH(r)
			clear()
			utils.ShowBanner(banner)
		case "0", "exit", "quit":
			clear()
			utils.Err("Exit.")
			os.Exit(0)
		default:
			clear()
			utils.ShowBanner(banner)
			utils.Err("Unknown command.")
		}
	}
}

func main() {
	r := bufio.NewReader(os.Stdin)
	showMain(r)
}
