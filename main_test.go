package main

import (
	"testing"

	"github.com/Kar98/over30club/types"
)

func TestMain(t *testing.T) {
	alItem := types.AlbumItem{
		AlbumType:            "",
		TotalTracks:          0,
		AvailableMarkets:     []string{},
		ExternalUrls:         types.ExternalUrls{},
		Href:                 "",
		ID:                   "",
		Images:               []types.Images{},
		Name:                 "Elton 60 - Live At Madison Square Garden",
		ReleaseDate:          "",
		ReleaseDatePrecision: "",
		Type:                 "",
		URI:                  "",
		Artists:              []types.Artists{},
		AlbumGroup:           "",
	}

	assert(t, alItem.IsUnwantedAlbum(), true)
}

func assert(t *testing.T, val1 any, val2 any) {
	if val1 != val2 {
		t.Fail()
	}
}
