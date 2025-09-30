package types

// type ArtistInput []struct {
// 	ArtistName string `json:"artistName"`
// 	Albums     []struct {
// 		Name        string `json:"name"`
// 		ReleaseYear int    `json:"releaseYear"`
// 	} `json:"albums"`
// 	ArtistYearOfBirth []int `json:"artistYearOfBirth"`
// }

type ArtistInput []struct {
	ArtistName        string `json:"artistName"`
	ArtistYearOfBirth []int  `json:"artistYearOfBirth"`
	Albums            []struct {
		Name        string `json:"name"`
		ReleaseYear int    `json:"releaseYear"`
	} `json:"albums"`
	Processed bool `json:"processed,omitempty"`
}
