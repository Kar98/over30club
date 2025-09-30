package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Kar98/over30club/client"
	"github.com/Kar98/over30club/spotifytypes"
	"github.com/Kar98/over30club/types"
)

var ErrNoAlbums = errors.New("0 albums returned")

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
		"getinputs": {
			Name:        "get artists via input",
			Description: "Gets the artists based on the _input.json file provided",
			Callback:    GetViaInput,
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
		searchResponse, err := sc.SearchArtist(artistName)
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
		albumList := make([]spotifytypes.Albumv2, 0, len(albumsReturned))
		for _, album := range albumsReturned {
			if album.IsUnwantedAlbum() {
				continue
			}
			albumDetails, err := sc.GetAlbumDetails(album.ID)
			if err != nil {
				errorOut(err)
				return
			}
			if albumDetails.IsLiveAlbum() {
				continue
			}
			albumList = append(albumList, albumDetails)
			time.Sleep(500 * time.Millisecond) // avoid rate limiting
		}
		outArtist, err := sc.GenerateArtist(artist, albumList)
		if err != nil {
			errorOut(err)
			return
		}

		err = saveArtistToDisk(outArtist)
		if err != nil {
			errorOut(err)
			return
		}
		fmt.Printf("done %s\n > ", artistName)
	}()

	return nil
}

func artistNameToFilepath(name string) string {
	filesafeName := strings.ToLower(name)
	filesafeName = strings.ReplaceAll(filesafeName, " ", "_")
	filesafeName = strings.ReplaceAll(filesafeName, "/", "_")
	return filesafeName
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

func GetViaInput(config *client.Config, args []string) error {
	var inputFilepath string
	if len(args) == 1 {
		inputFilepath = args[0] // to make this testable
	} else {
		inputFilepath = client.ArtistInputFile
	}
	defer fmt.Println("done getinputs")
	// Get list from artistinput
	sc, err := NewSpotifyClient(config)
	if err != nil {
		return err
	}
	var inputFile types.ArtistInput
	fileBytes, err := os.ReadFile(inputFilepath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileBytes, &inputFile)
	if err != nil {
		return err
	}

	// We have a list of albums. For each album, search for it in Spotify.
	foundAlbums := []types.AlbumWithQuery{}
	for i, artistInput := range inputFile {
		filesafeName := artistNameToFilepath(artistInput.ArtistName) + ".json"
		f, err := os.Open(path.Join(client.ArtistDir, filesafeName))
		// If file exists, then we don't need to get extra data
		if err == nil || artistInput.Processed {
			fmt.Println("data exists for", filesafeName)
			f.Close()
			inputFile[i].Processed = true
			continue
		}

		// get artist ID
		searchResponse, err := sc.SearchArtist(artistInput.ArtistName)
		if err != nil {
			return err
		}
		artist := searchResponse.Artists.Items[0]
		fmt.Println("getting", artist.Name)

		for _, album := range artistInput.Albums {
			albumItem, err := sc.getAlbumFromSearch(album.Name, album.ReleaseYear, artist.ID)
			if errors.Is(err, ErrNoAlbums) {
				fmt.Printf("album not found: %s releaseYear: %d\n", album.Name, album.ReleaseYear)
				continue
			} else if err != nil {
				return err
			}
			foundAlbums = append(foundAlbums, types.AlbumWithQuery{AlbumItem: albumItem, QueryName: album.Name})
		}

		// All albums are gathered and have the correct IDs
		// Get the detailed album details along with the
		albumList := make([]types.Albumv2WithQuery, 0)
		for _, album := range foundAlbums {
			albumDetails, err := sc.GetAlbumDetails(album.ID)
			if err != nil {
				return err
			}
			albumList = append(albumList, types.Albumv2WithQuery{Albumv2: albumDetails, QueryName: album.QueryName})
			time.Sleep(500 * time.Millisecond) // avoid rate limiting
		}
		artistData, err := sc.GenerateArtistFromInput(artist, albumList, getAverageYearOfBirth(artistInput.ArtistYearOfBirth))
		if err != nil {
			return err
		}
		err = saveArtistToDisk(artistData)
		if err != nil {
			return err
		}
		inputFile[i].Processed = true
	}

	// Write back to the inputs file to update any processed artists
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err = enc.Encode(inputFile)
	if err != nil {
		fmt.Println("could not encode input file", err.Error())
		return nil
	}
	err = os.WriteFile(inputFilepath, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("could not update input file", err.Error())
	}

	return nil
}

func saveArtistToDisk(artist types.Artist) error {
	// save data to disk
	outJson, err := json.MarshalIndent(artist, "", "  ")
	if err != nil {
		return err
	}
	filesafeName := artistNameToFilepath(artist.Name)
	filename := fmt.Sprintf("%s/%s.json", client.ArtistDir, filesafeName)
	err = os.WriteFile(filename, outJson, 0644)
	return err
}

func getAverageYearOfBirth(years []int) int {
	if len(years) == 0 {
		return 0
	}
	sum := 0
	for _, y := range years {
		sum += y
	}
	return int(sum / len(years))
}

func (sc SpotifyClient) getAlbumFromSearch(albumName string, releaseYear int, artistId string) (spotifytypes.AlbumItem, error) {
	albumSearch, err := sc.SearchAlbums(albumName)
	if err != nil {
		return spotifytypes.AlbumItem{}, err
	}
	// Album names are not unique at all, need to filter down to only the artist we want
	relevantAlbums := filterByArtistId(albumSearch, artistId)
	if len(relevantAlbums) == 0 {
		return spotifytypes.AlbumItem{}, ErrNoAlbums
	}

	for _, album := range relevantAlbums {
		// Exact match takes the highest preference. There shouldn't be 2 albums released in the same year by the 1 artist
		if strings.EqualFold(album.Name, albumName) {
			return album, nil
		}
	}

	for _, album := range relevantAlbums {
		// However exact match isn't likely due to punctuation/releases/etc.
		// Try to get a partial match + release year
		var albumYear time.Time
		var err error
		if album.ReleaseDatePrecision == "year" {
			albumYear, err = time.Parse("2006", album.ReleaseDate)
		} else {
			albumYear, err = time.Parse("2006-01-02", album.ReleaseDate)
		}
		if err != nil {
			fmt.Println("album.ReleaseDate: ", album.ReleaseDate, album.Name)
			return spotifytypes.AlbumItem{}, err
		}
		relevantAlbumLower := strings.ToLower(album.Name)
		albumNameLower := strings.ToLower(albumName)
		if strings.Contains(relevantAlbumLower, albumNameLower) && albumYear.Year() == releaseYear {
			return album, nil
		}
		// Flip the partial string check
		if strings.Contains(albumNameLower, relevantAlbumLower) && albumYear.Year() == releaseYear {
			return album, nil
		}
		// If still can't find then use release year only. If 2 albums released in same year then the first one will be grabbed.
		if albumYear.Year() == releaseYear {
			return album, nil
		}
	}
	// If still can't find, then there is a data issue so error out nicely
	return spotifytypes.AlbumItem{}, ErrNoAlbums
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

func filterByArtistId(searchResponse spotifytypes.SearchResponse, artistId string) []spotifytypes.AlbumItem {
	var data []spotifytypes.AlbumItem
	for _, album := range searchResponse.Albums.Items {
		// Even though we specify to only search albums in the search params, spotify will still return singles
		// example: Beyoncé -> albumId = 45BFNKQ0VGAXACLAEOy9Mv
		if album.AlbumType != "album" {
			continue
		}
		for _, artist := range album.Artists {
			if artist.ID == artistId {
				data = append(data, album)
				break
			}
		}
	}
	return data
}
