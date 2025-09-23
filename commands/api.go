package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Kar98/over30club/client"
	"github.com/Kar98/over30club/types"
)

func NewSpotifyClient(config *client.Config) (*SpotifyClient, error) {
	if config == nil {
		return nil, errors.New("config required")
	}
	// load v1 token from file
	var v1token string
	v1token, err := getV1Token(config)
	// if not found then get new token from endpoint
	if err != nil {
		fmt.Println("Getting new v1 token")
		v1token, err = setV1Token(config)
		if err != nil {
			return nil, err
		}
	}
	// load v2 token
	return &SpotifyClient{
		v1Token: v1token,
		v2Token: config.V2.ClientToken,
		v2Auth:  config.V2.Authorization,
	}, nil
}

func getV1Token(config *client.Config) (string, error) {
	if config.V1.Token == "" {
		return "", errors.New("no v1 token found")
	}
	// Check if token has expired
	if config.V1.TokenExpiry.Before(time.Now()) {
		return "", errors.New("token expired")
	}

	return config.V1.Token, nil
}

func setV1Token(config *client.Config) (string, error) {
	if config.V1.Client == "" || config.V1.Secret == "" {
		return "", errors.New("set client and secret")
	}
	body := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", config.V1.Client, config.V1.Secret)
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewReader([]byte(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var tokenRes types.TokenResponse
	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(bs, &tokenRes)
	if err != nil {
		return "", err
	}
	// File will always be created in main.go
	file, err := os.ReadFile(client.UserDataFile)
	if err != nil {
		return "", err
	}
	var existingConfig client.Config
	err = json.Unmarshal(file, &existingConfig)
	if err != nil {
		return "", err
	}
	existingConfig.V1.Token = tokenRes.AccessToken
	existingConfig.V1.TokenExpiry = time.Now().Add(time.Duration(tokenRes.ExpiresIn) * time.Second)
	jsonFile, _ := json.MarshalIndent(existingConfig, "", "  ")
	err = os.WriteFile(client.UserDataFile, jsonFile, 0644)
	if err != nil {
		return "", err
	}
	return tokenRes.AccessToken, nil
}

type SpotifyClient struct {
	v1Token string
	v2Token string
	v2Auth  string
}

func (sc *SpotifyClient) Search(artistName string) (types.SearchResponse, error) {
	client := &http.Client{}
	// Get artist id
	encodedName := strings.ReplaceAll(artistName, " ", "%20")
	getArtistUrl := fmt.Sprintf("https://api.spotify.com/v1/search?query=%s&type=artist&market=AU&limit=1", encodedName)
	getArtistRes, err := http.NewRequest("GET", getArtistUrl, nil)
	if err != nil {
		return types.SearchResponse{}, err
	}
	sc.setV1Headers(getArtistRes)
	res, err := client.Do(getArtistRes)
	if err != nil {
		return types.SearchResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("request: ", getArtistUrl)
		body, _ := io.ReadAll(res.Body)
		fmt.Println("response: ", string(body))
		return types.SearchResponse{}, fmt.Errorf("error when searching, status: %s", res.Status)
	}

	var searchResponse types.SearchResponse
	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return types.SearchResponse{}, err
	}
	err = json.Unmarshal(bs, &searchResponse)
	if err != nil {
		return types.SearchResponse{}, err
	}
	if len(searchResponse.Artists.Items) == 0 {
		return types.SearchResponse{}, fmt.Errorf("no artist found for %s", artistName)
	}
	return searchResponse, nil
}

func (sc *SpotifyClient) GetAlbumList(artistId string) ([]types.AlbumItem, error) {
	var albums []types.AlbumItem
	getAlbumsUrl := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums?include_groups=album&market=AU&limit=50", artistId)
	getAlbumRes, err := sc.getAlbums(getAlbumsUrl)
	if err != nil {
		return nil, err
	}
	albums = append(albums, getAlbumRes.Items...)
	// won't be more than 100 albums
	_, isStr := getAlbumRes.Next.(string)
	if isStr {
		getAlbumsUrl = getAlbumRes.Next.(string)
		getAlbumRes, err := sc.getAlbums(getAlbumsUrl)
		if err != nil {
			return nil, err
		}
		albums = append(albums, getAlbumRes.Items...)
	}

	return albums, nil
}

func (sc *SpotifyClient) getAlbums(url string) (types.GetAlbumResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.GetAlbumResponse{}, err
	}
	sc.setV1Headers(req)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return types.GetAlbumResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("request: ", url)
		body, _ := io.ReadAll(res.Body)
		fmt.Println("response: ", string(body))
		return types.GetAlbumResponse{}, fmt.Errorf("error when getting album list, status: %s", res.Status)
	}

	var albumResponse types.GetAlbumResponse
	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return types.GetAlbumResponse{}, err
	}
	err = json.Unmarshal(bs, &albumResponse)
	if err != nil {
		return types.GetAlbumResponse{}, err
	}

	return albumResponse, nil
}

func (sc *SpotifyClient) GetAlbumDetails(albumId string) (types.Albumv2, error) {
	getAlbumDetails := "https://api-partner.spotify.com/pathfinder/v2/query"
	body := types.PostQuery{
		Variables: types.Variables{
			URI:    fmt.Sprintf("spotify:album:%s", albumId),
			Locale: "",
			Offset: 0,
			Limit:  50,
		},
		OperationName: "getAlbum",
		Extensions: types.Extensions{
			PersistedQuery: types.PersistedQuery{
				Version:    1,
				Sha256Hash: "b9bfabef66ed756e5e13f68a942deb60bd4125ec1f1be8cc42769dc0259b4b10",
			},
		},
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", getAlbumDetails, bytes.NewReader((jsonBody)))
	if err != nil {
		return types.Albumv2{}, err
	}
	sc.setV2Headers(req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return types.Albumv2{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("request: ", getAlbumDetails)
		body, _ := io.ReadAll(res.Body)
		fmt.Println("response: ", string(body))
		return types.Albumv2{}, fmt.Errorf("error when getting album details, status: %s", res.Status)
	}

	var albumResponse types.Albumv2
	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return types.Albumv2{}, err
	}
	err = json.Unmarshal(bs, &albumResponse)
	if err != nil {
		return types.Albumv2{}, err
	}
	return albumResponse, err
}

func (sc *SpotifyClient) setV1Headers(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+sc.v1Token)
}

func (sc *SpotifyClient) setV2Headers(req *http.Request) {
	req.Header.Set("authorization", sc.v2Auth)
	req.Header.Set("client-token", sc.v2Token)
	req.Header.Set("content-Type", "application/json;charset=UTF-8")
	req.Header.Set("accept", "application/json")
}

func (sc *SpotifyClient) GenerateArtist(artist types.ArtistItem, albums []types.Albumv2) (types.Artist, error) {
	convertedAlbums := make([]types.Album, 0, len(albums))
	for _, album := range albums {
		convertedAlbums = append(convertedAlbums, sc.toAlbum(album))
	}
	return types.Artist{
		Name:   artist.Name,
		ID:     artist.ID,
		Albums: convertedAlbums,
		// Members and AverageAge to be filled in later
		Members:        []types.Member{},
		AvgYearOfBirth: 0,
	}, nil
}

func (sc *SpotifyClient) toAlbum(album types.Albumv2) types.Album {
	totalPlaycount := 0
	tracks := make([]types.Track, 0, len(album.Data.AlbumUnion.TracksV2.Items))
	for _, track := range album.Data.AlbumUnion.TracksV2.Items {
		playcount, _ := strconv.Atoi(track.Track.Playcount)
		trackSplits := strings.Split(track.Track.URI, ":")
		if len(trackSplits) != 3 {
			panic("unexpected track URI format - " + track.Track.URI)
		}
		tracks = append(tracks, types.Track{
			Name:      track.Track.Name,
			ID:        trackSplits[2],
			PlayCount: playcount,
		})
		totalPlaycount += playcount
	}

	albumSplits := strings.Split(album.Data.AlbumUnion.URI, ":")
	if len(albumSplits) != 3 {
		panic("unexpected album URI format - " + album.Data.AlbumUnion.URI)
	}
	albumID := albumSplits[2]
	return types.Album{
		Name:           album.Data.AlbumUnion.Name,
		ReleaseYear:    album.Data.AlbumUnion.Date.IsoString.Year(),
		ID:             albumID,
		Tracks:         tracks,
		TotalPlaycount: totalPlaycount,
	}
}
