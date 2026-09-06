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

func main() {
	clear()
	utils.ShowBanner(banner)
	r := bufio.NewReader(os.Stdin)
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
			commands.N(r)
		case "2":
			clear()
			commands.D(r)
		case "3":
			clear()
			commands.U(r)
		case "4":
			clear()
			commands.CS(r)
		case "5":
			clear()
			commands.IS(r)
		case "6":
			clear()
			commands.IP(r)
		case "7":
			clear()
			commands.CIP()
		case "8":
			clear()
			commands.KP(r)
		case "9":
			clear()
			commands.NP(r)
		case "10":
			clear()
			commands.WZ(r)
		case "11":
			clear()
			commands.TM(r)
		case "12":
			clear()
			commands.DR(r)
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
