package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Kar98/over30club/client"
	"github.com/Kar98/over30club/commands"
)

func main() {
	// Check if user data file exists. If not, create it.
	if _, err := os.Stat(client.UserDataFile); os.IsNotExist(err) {
		err := generateDefaultDataFile()
		if err != nil {
			fmt.Println("Error generating default data file:", err)
			return
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	config := client.Config{Scanner: scanner}
	cmds := commands.GenerateCommands()
	fmt.Print(" > ")
	for scanner.Scan() {
		refreshConfig(&config)
		input, err := commands.CleanInput(scanner.Text())
		if err != nil {
			fmt.Print("\n > ")
			continue
		}
		fmt.Println("Your command was: " + input[0])
		cmd, ok := cmds[input[0]]
		if ok {
			err := cmd.Callback(&config, input[1:])
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("Unknown command")
		}
		fmt.Print(" > ")
	}

}

func generateDefaultDataFile() error {
	err := os.MkdirAll("userdata", os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(client.UserDataFile)
	if err != nil {
		return err
	}
	defer file.Close()
	var config client.Config
	json.Unmarshal([]byte("{}"), &config)
	dummyFile, _ := json.MarshalIndent(config, "", "  ")
	_, err = file.Write(dummyFile)
	return err
}

func refreshConfig(config *client.Config) {
	var configJsonFile client.Config
	data, err := os.ReadFile(client.UserDataFile)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	err = json.Unmarshal(data, &configJsonFile)
	if err != nil {
		fmt.Println("Error unmarshalling config file:", err)
		return
	}
	config.V1 = configJsonFile.V1
	config.V2 = configJsonFile.V2
}
