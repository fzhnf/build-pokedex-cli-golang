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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex> ")
		if scanner.Scan() {
			text := scanner.Text()
			textCleaned := cleanInput(text)
			fmt.Printf("Your command was: %s\n", textCleaned[0])
		}
	}
}
