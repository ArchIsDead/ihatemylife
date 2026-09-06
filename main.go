package main

import (
	"bufio"
	"fmt"
	"os"
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

func main() {
	utils.B(banner)
	utils.S("s\n")
	r := bufio.NewReader(os.Stdin)
	for {
		utils.M()
		utils.Up()
		fmt.Print(utils.P("\n> "))
		x, _ := r.ReadString('\n')
		x = strings.TrimSpace(x)
		switch x {
		case "1":
			commands.N(r)
		case "2":
			commands.P()
		case "3":
			commands.K(r)
		case "4":
			commands.D(r)
		case "5":
			commands.U(r)
		case "6":
			commands.CS(r)
		case "7":
			commands.IS(r)
		case "8":
			commands.IP(r)
		case "9":
			commands.CIP()
		case "10":
			commands.DI()
		case "11":
			commands.ST()
		case "12":
			commands.KP(r)
		case "13":
			commands.NP(r)
		case "14":
			commands.PV(r)
		case "0", "exit", "quit":
			utils.E("Exit.")
			os.Exit(0)
		}
	}
}
