package types

type ArtistInput []struct {
	ArtistName string `json:"artistName"`
	Albums     []struct {
		Name        string `json:"name"`
		ReleaseYear int    `json:"releaseYear"`
	} `json:"albums"`
}
