package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Kar98/over30club/client"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*client.Config, []string) error
}

func GenerateCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit program",
			Callback:    Exit,
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    Help,
		},
		"settoken": {
			Name:        "set token",
			Description: "Sets the auth and client token for v2 requests",
			Callback:    SetTokens,
		},
		"getartist": {
			Name:        "get artist",
			Description: "Get an artists albums and playcounts for all tracks",
			Callback:    GetArtistInfo,
		},
		"test": {
			Name:        "test",
			Description: "A test command",
			Callback:    Test,
		},
	}
}

func Help(config *client.Config, _ []string) error {
	cmds := GenerateCommands()
	fmt.Print("Usage:\n\n")

	for k, v := range cmds {
		fmt.Printf("%s: %s\n", k, v.Description)
	}

	return nil
}

func Exit(config *client.Config, _ []string) error {
	os.Exit(0)
	return nil
}

func SetTokens(config *client.Config, _ []string) error {
	fmt.Printf("set v2 token > ")
	config.Scanner.Scan()
	token := config.Scanner.Text()
	if token != "" {
		config.V2.ClientToken = token
	} else {
		fmt.Println("no text entered")
		return nil
	}

	fmt.Printf("token saved = %s\n", config.V2.ClientToken)
	fmt.Print("Enter authorization > ")

	config.Scanner.Scan()
	auth := config.Scanner.Text()
	if auth != "" {
		config.V2.Authorization = auth
	}

	fileData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	os.WriteFile(client.UserDataFile, fileData, 0644)
	fmt.Println("tokens saved")

	return nil
}

func GetArtistInfo(config *client.Config, data []string) error {
	if len(data) != 1 {
		return errors.New("enter an artist name")
	}
	// Get the artist id
	sc, err := NewSpotifyClient(config)
	if err != nil {
		return err
	}
	searchResponse, err := sc.Search(data[0])
	if err != nil {
		return err
	}

	artistId := searchResponse.Artists.Items[0].ID

	// get a list of albums from the artist
	albumsResponse, err := sc.GetAlbumList(artistId)
	if err != nil {
		return err
	}

	// for each album, get the songs + their playcounts
	fmt.Printf("Arist: %s\n", searchResponse.Artists.Items[0].Name)
	fmt.Printf("Artist ID: %s\n", artistId)
	for _, album := range albumsResponse.Items {
		fmt.Printf("Album: %s\n", album.Name)
		fmt.Printf("Album ID: %s\n", album.ID)
	}
	albumDetails, err := sc.GetAlbumDetails(albumsResponse.Items[0].ID)
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Printf("album name: %s\n", albumDetails.Data.AlbumUnion.Name)
	for _, albumTrack := range albumDetails.Data.AlbumUnion.TracksV2.Items {
		fmt.Printf("Track: %s, Playcount: %s\n", albumTrack.Track.Name, albumTrack.Track.Playcount)
	}

	// save data to disk
	return nil
}

func Test(config *client.Config, data []string) error {
	//_, err := NewSpotifyClient(config)
	var cfg client.Config
	file, err := os.ReadFile(client.UserDataFile)
	if err != nil {
		return err
	}
	json.Unmarshal(file, &cfg)
	fmt.Println(cfg.V1.TokenExpiry)

	fmt.Println("pass")
	return err
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
