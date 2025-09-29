package types

import "github.com/Kar98/over30club/spotifytypes"

type Member struct {
	Name        string `json:"name"`
	YearOfBirth int    `json:"age"`
}

type Track struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	PlayCount int    `json:"playCount"`
}

type Album struct {
	Name           string  `json:"name"`
	QueryName      string  `json:"queryName,omitempty"`
	ReleaseYear    int     `json:"releaseYear"`
	ID             string  `json:"id"`
	Tracks         []Track `json:"tracks"`
	TotalPlaycount int     `json:"totalPlaycount"`
}

type Artist struct {
	Name           string  `json:"name"`
	ID             string  `json:"id"`
	Albums         []Album `json:"albums"`
	AvgYearOfBirth int     `json:"avgYearOfBirth"`
}

// With Query is used for tracing how the album names were matched. Because it's very hard to match album names
// this will help track down any edge cases that appear with the albums (eg Rammstein album -> Untitled)
type AlbumWithQuery struct {
	spotifytypes.AlbumItem
	QueryName string `json:"queryName,omitempty"`
}

type Albumv2WithQuery struct {
	spotifytypes.Albumv2
	QueryName string `json:"queryName,omitempty"`
}
