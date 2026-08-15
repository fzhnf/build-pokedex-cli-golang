package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	multilineHelp := `
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
`
	fmt.Println(multilineHelp)
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cliCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
	for {
		fmt.Print("Pokedex> ")
		if scanner.Scan() {
			text := scanner.Text()
			textCleaned := cleanInput(text)
			if textCleaned[0] == "" {
				cliCommands["help"].callback()
				continue
			}
			// if _, ok := cliCommands[textCleaned[0]]; !ok {
			// 	cliCommands["help"].callback()
			// }
			command, ok := cliCommands[textCleaned[0]]
			if !ok {
				cliCommands["help"].callback()
				continue
			}
			command.callback()
		}
	}
}
