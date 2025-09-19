package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Kar98/over30club/client"
	"github.com/Kar98/over30club/types"
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
		"settokens": {
			Name:        "set token",
			Description: "Sets the client token and auth for v2 requests",
			Callback:    SetTokens,
		},
		"setauth": {
			Name:        "set auth",
			Description: "Sets the auth for v2 requests",
			Callback:    SetAuth,
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
	fmt.Print("Usage:\n")

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
		return errors.New("no text entered")
	}

	fmt.Print("enter v2auth > ")

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

func SetAuth(config *client.Config, _ []string) error {
	fmt.Printf("enter v2auth > ")
	config.Scanner.Scan()
	auth := config.Scanner.Text()
	if auth != "" {
		config.V2.Authorization = auth
	} else {
		return errors.New("no text entered")
	}

	fileData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	os.WriteFile(client.UserDataFile, fileData, 0644)
	fmt.Println("auth saved")

	return nil
}

func GetArtistInfo(config *client.Config, data []string) error {
	if len(data) == 0 {
		return errors.New("enter an artist name")
	}
	artistName := strings.Join(data, " ")
	errorOut := func(err error) {
		fmt.Println(err.Error())
		fmt.Print(" > ")
	}
	go func() {
		// Get the artist id
		sc, err := NewSpotifyClient(config)
		if err != nil {
			errorOut(err)
			return
		}
		searchResponse, err := sc.Search(artistName)
		if err != nil {
			errorOut(err)
			return
		}

		artist := searchResponse.Artists.Items[0]
		fmt.Printf("getting %s\n > ", artist.Name)

		// get a list of albums from the artist
		albumsReturned, err := sc.GetAlbumList(artist.ID)
		if err != nil {
			errorOut(err)
			return
		}

		// for each album, get the songs + their playcounts
		albumList := make([]types.Albumv2, 0, len(albumsReturned))
		for _, album := range albumsReturned {
			albumDetails, err := sc.GetAlbumDetails(album.ID)
			if err != nil {
				errorOut(err)
				return
			}
			albumList = append(albumList, albumDetails)
			time.Sleep(1 * time.Second) // avoid rate limiting
		}
		outArtist, err := sc.GenerateArtist(artist, albumList)
		if err != nil {
			errorOut(err)
			return
		}
		// save data to disk
		outJson, err := json.MarshalIndent(outArtist, "", "  ")
		if err != nil {
			errorOut(err)
			return
		}

		filesafeName := strings.ReplaceAll(strings.ToLower(outArtist.Name), " ", "_")
		filename := fmt.Sprintf("%s/%s.json", client.ArtistDir, filesafeName)
		err = os.WriteFile(filename, outJson, 0644)
		if err != nil {
			errorOut(err)
			return
		}
		fmt.Printf("done %s\n > ", artistName)
	}()

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
