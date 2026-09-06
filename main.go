package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"s/c"
	"s/u"
)

var b = `
  _________     .__       ____         .___      
 /   _____/__ __|__| ____/_   |      __| _/____  
 \_____  \|  |  \  |/ ___\|   |     / __ |/ __ \ 
 /        \  |  /  \  \___|   |    / /_/ \  ___/ 
/_______  /____/|__|\___  >___| /\ \____ |\___  >
        \/              \/      \/      \/    \/ 
`

func main() {
	fmt.Println(u.B(b))
	fmt.Println(u.S("s\n"))
	r := bufio.NewReader(os.Stdin)
	u.M()
	for {
		fmt.Print(u.P("\n> "))
		x, _ := r.ReadString('\n')
		x = strings.TrimSpace(x)
		switch x {
		case "1":
			c.N(r)
		case "2":
			c.P()
		case "3":
			c.K(r)
		case "4":
			c.D(r)
		case "5":
			c.U(r)
		case "6":
			c.CS(r)
		case "7":
			c.IS(r)
		case "8":
			c.IP(r)
		case "9":
			c.CIP()
		case "10":
			c.DI()
		case "11":
			c.ST()
		case "12":
			c.KP(r)
		case "13":
			c.NP(r)
		case "14":
			c.PV(r)
		case "0", "exit", "quit":
			fmt.Println(u.E("Exit."))
			os.Exit(0)
		case "help", "menu":
			u.M()
		default:
			u.M()
		}
	}
}
