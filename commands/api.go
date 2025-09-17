package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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
	v1token, err := getV1Token()
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

func getV1Token() (string, error) {
	file, err := os.ReadFile(client.UserDataFile)
	if err != nil {
		return "", err
	}
	var config client.Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return "", err
	}
	if config.V1.Token == "" {
		return "", errors.New("no v1 token found")
	}
	// Check if token has expired
	if config.V1.TokenExpiry.After(time.Now()) {
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

func (sc *SpotifyClient) Search(query string) (types.SearchResponse, error) {
	client := &http.Client{}
	// Get artist id
	getArtistUrl := fmt.Sprintf("https://api.spotify.com/v1/search?query=%s&type=artist&market=AU&limit=1", query)
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

	var searchResponse types.SearchResponse
	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return types.SearchResponse{}, err
	}
	err = json.Unmarshal(bs, &searchResponse)
	if err != nil {
		return types.SearchResponse{}, err
	}
	return searchResponse, nil
}

func (sc *SpotifyClient) GetAlbumList(artistId string) (types.GetAlbumResponse, error) {
	getAlbumsUrl := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums?include_groups=album&market=AU", artistId)
	req, err := http.NewRequest("GET", getAlbumsUrl, nil)
	if err != nil {
		return types.GetAlbumResponse{}, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return types.GetAlbumResponse{}, err
	}
	defer res.Body.Close()

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
	return types.Albumv2{}, nil
}

func (sc *SpotifyClient) setV1Headers(req *http.Request) {
	req.Header.Set("Authorization", sc.v1Token)
	req.Header.Set("Content-Type", "application/json")
}

func (sc *SpotifyClient) setV2Headers(req *http.Request) {
	req.Header.Set("authorization", sc.v2Auth)
	req.Header.Set("client-token", sc.v2Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}
