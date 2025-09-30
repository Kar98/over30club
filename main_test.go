package main

import (
	"encoding/json"
	"testing"

	"github.com/Kar98/over30club/types"
)

func TestInputFile(t *testing.T) {
	jsonStr := `[{"artistName": "Elvis Presley","artistYearOfBirth": [1935],"albums": [{"name": "Elvis Presley","releaseYear": 1956},{"name": "Elvis","releaseYear": 1956},{"name": "Elvis' Christmas Album","releaseYear": 1957},{"name": "Elvis Is Back","releaseYear": 1960},{"name": "His Hand in Mine","releaseYear": 1960},{"name": "Something for Everybody","releaseYear": 1961},{"name": "Pot Luck","releaseYear": 1962},{"name": "Elvis for Everyone!","releaseYear": 1965},{"name": "How Great Thou Art","releaseYear": 1967},{"name": "From Elvis in Memphis","releaseYear": 1969},{"name": "That's the Way It Is","releaseYear": 1970},{"name": "Elvis Country","releaseYear": 1971},{"name": "Love Letters from Elvis","releaseYear": 1971},{"name": "Elvis Sings The Wonderful World of Christmas","releaseYear": 1971},{"name": "Elvis Now","releaseYear": 1972},{"name": "He Touched Me","releaseYear": 1972},{"name": "Elvis","releaseYear": 1973},{"name": "Raised on Rock","releaseYear": 1973},{"name": "Good Times","releaseYear": 1974},{"name": "Promised Land","releaseYear": 1975},{"name": "Today","releaseYear": 1975},{"name": "From Elvis Presley Boulevard, Memphis, Tennessee","releaseYear": 1976},{"name": "Moody Blue","releaseYear": 1977}]}]`
	var jsonObj types.ArtistInput
	err := json.Unmarshal([]byte(jsonStr), &jsonObj)
	assertNoError(t, err)
}

func TestDebug(t *testing.T) {
	rawData := types.ArtistInput{{
		ArtistName: "test",
	},
		{
			ArtistName: "test2",
		}}

	for _, artist := range rawData {
		artist.ArtistName = "newval"
	}

	t.Log("test")
}

func assert(t *testing.T, val1 any, val2 any) {
	if val1 != val2 {
		t.Fail()
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err.Error())
	}
}
