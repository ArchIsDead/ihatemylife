func (r *Result) Show() {
	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ IP INFO ]")))
	fmt.Println(utils.Div())

	if !r.Success {
		fmt.Println(utils.Red("Invalid IP"))
		return
	}

	fmt.Println(utils.Gry("IP: ") + utils.Wht(r.IP))
	if r.Type != "" {
		fmt.Println(utils.Gry("Type: ") + utils.Wht(r.Type))
	}
	if r.Continent != "" {
		fmt.Println(utils.Gry("Continent: ") + utils.Wht(r.Continent))
	}
	if r.Country != "" {
		fmt.Println(utils.Gry("Country: ") + utils.Wht(r.Country))
	}
	if r.Region != "" {
		fmt.Println(utils.Gry("Region: ") + utils.Wht(r.Region))
	}
	if r.City != "" {
		fmt.Println(utils.Gry("City: ") + utils.Wht(r.City))
	}
	if r.Lat != "<nil>" && r.Lat != "" {
		fmt.Println(utils.Gry("Latitude: ") + utils.Wht(r.Lat))
	}
	if r.Lon != "<nil>" && r.Lon != "" {
		fmt.Println(utils.Gry("Longitude: ") + utils.Wht(r.Lon))
	}
	if r.Org != "" {
		fmt.Println(utils.Gry("Org: ") + utils.Wht(r.Org))
	}
	if r.ISP != "" {
		fmt.Println(utils.Gry("ISP: ") + utils.Wht(r.ISP))
	}
	if r.Timezone != "" {
		fmt.Println(utils.Gry("Timezone: ") + utils.Wht(r.Timezone))
	}
	fmt.Println()
	fmt.Print(utils.Gry("Press Enter to return..."))
	fmt.Scanln()
}
