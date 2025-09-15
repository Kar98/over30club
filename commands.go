package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func GenerateCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    Exit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    Help,
		},
		"settoken": {
			name:        "set token",
			description: "Sets the auth and client token for requests",
			callback:    SetTokens,
		},
	}
}

func Help(config *Config, _ []string) error {
	cmds := GenerateCommands()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for k, v := range cmds {
		fmt.Printf("%s: %s\n", k, v.description)
	}

	return nil
}

func Exit(config *Config, _ []string) error {
	os.Exit(0)
	return nil
}

func SetTokens(config *Config, _ []string) error {
	fmt.Printf("set auth > ")
	config.scanner.Scan()
	token := config.scanner.Text()
	if token != "" {
		config.token = token
	} else {
		fmt.Println("no text entered")
		return nil
	}

	fmt.Printf("token saved = %s\n", config.token)

	return nil
}

func GetAlbum(config *Config, data []string) error {
	// get album
	// save data to disk
	return nil
}

func CleanInput(text string) ([]string, error) {
	var words []string
	if text == "" {
		return []string{}, errors.New("no text")
	}
	splits := strings.Split(text, " ")
	for _, word := range splits {
		trimmedWord := strings.Trim(word, " ")
		trimmedWord = strings.ToLower(trimmedWord)
		if word != "" {
			words = append(words, trimmedWord)
		}
	}
	return words, nil
}
