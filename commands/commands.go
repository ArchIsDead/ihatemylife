func AI(r *bufio.Reader) {
	client := ai.New()
	utils.Err("Initializing rfour...")
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
			if len(lines) == 1 && line == "" {
				continue
			}
			if len(lines) > 0 && line == "" {
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
