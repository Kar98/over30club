package main

import (
	"bufio"
	"fmt"
	"os"
)

type Config struct {
	scanner *bufio.Scanner
	token   string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := Config{scanner: scanner}
	commands := GenerateCommands()
	fmt.Print("Enter command > ")
	for scanner.Scan() {
		input, err := CleanInput(scanner.Text())
		if err != nil {
			fmt.Print("\nEnter command > ")
			continue
		}
		fmt.Println("Your command was: " + input[0])
		cmd, ok := commands[input[0]]
		if ok {
			err := cmd.callback(&config, input[1:])
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("Unknown command")
		}
	}

	//Have a CLI set up where you can enter in your data

	// Enter in the token

	// Scan through the artists albums

	// Get the playcounts for each of the artists albums and tally them up. Get the album year

	// Visualise the playcounts for each album by year

}
