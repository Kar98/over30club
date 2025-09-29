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
	Name           string   `json:"name"`
	ID             string   `json:"id"`
	Albums         []Album  `json:"albums"`
	Members        []Member `json:"members"`
	AvgYearOfBirth float64  `json:"avgYearOfBirth"`
}

type AlbumWithQuery struct {
	spotifytypes.AlbumItem
	QueryName string `json:"queryName,omitempty"`
}

type Albumv2WithQuery struct {
	spotifytypes.Albumv2
	QueryName string `json:"queryName,omitempty"`
}
