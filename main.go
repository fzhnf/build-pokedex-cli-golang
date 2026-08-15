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

func commandExit(*config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(*config) error {
	multilineHelp := `
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
`
	fmt.Println(multilineHelp)
	return nil
}

type config struct {
	commands map[string]cliCommand
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := config{
		commands: map[string]cliCommand{
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
		},
	}
	for {
		fmt.Print("Pokedex> ")
		if scanner.Scan() {
			text := scanner.Text()
			textCleaned := cleanInput(text)
			if len(textCleaned) == 0 || textCleaned[0] == "" {
				cfg.commands["help"].callback(&cfg)
				continue
			}

			if command, ok := cfg.commands[textCleaned[0]]; !ok {
				cfg.commands["help"].callback(&cfg)
			} else {
				command.callback(&cfg)
			}
		}
	}
}
