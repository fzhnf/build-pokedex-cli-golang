package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type locationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []locationArea `json:"results"`
}

func commandMap(cfg *config) error {
	var requestURL string
	if cfg.locationAreaNextPage != nil {
		requestURL = *cfg.locationAreaNextPage
	} else {
		requestURL = "https://pokeapi.co/api/v2/location-area"
	}

	res, err := http.Get(requestURL)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	locationAreas := locationAreaResponse{}
	if err := json.Unmarshal(body, &locationAreas); err != nil {
		log.Fatal(err)
	}

	cfg.locationAreaNextPage = locationAreas.Next
	cfg.locationAreaPrevPage = locationAreas.Previous

	for _, locationArea := range locationAreas.Results {
		fmt.Printf("%s\n", locationArea.Name)
	}

	return nil
}

func commandMapBack(cfg *config) error {
	if cfg.locationAreaPrevPage == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	requestURL := *cfg.locationAreaPrevPage
	cfg.locationAreaNextPage = &requestURL
	commandMap(cfg)

	return nil
}

type config struct {
	commands             map[string]cliCommand
	locationAreaNextPage *string
	locationAreaPrevPage *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func newConfig() *config {
	return &config{
		locationAreaNextPage: nil,
		locationAreaPrevPage: nil,
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
			"map": {
				name:        "map",
				description: "Displays the names of location areas in the Pokemon world",
				callback:    commandMap,
			},
			"mapb": {
				name:        "mapb",
				description: "Displays the previous names of location areas in the Pokemon world",
				callback:    commandMapBack,
			},
		},
	}
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex> ")
		if scanner.Scan() {
			text := scanner.Text()
			textCleaned := cleanInput(text)
			if len(textCleaned) == 0 || textCleaned[0] == "" {
				cfg.commands["help"].callback(cfg)
				continue
			}

			if command, ok := cfg.commands[textCleaned[0]]; !ok {
				cfg.commands["help"].callback(cfg)
			} else {
				command.callback(cfg)
			}
		}
	}
}

func main() {
	startRepl(newConfig())
}
