package spotifytypes

import (
	"strings"
	"time"
)

type SearchResponse struct {
	Artists Artists `json:"artists"`
	Albums  Albums  `json:"albums"`
}
type ExternalUrls struct {
	Spotify string `json:"spotify"`
}
type Followers struct {
	Href  any `json:"href"`
	Total int `json:"total"`
}
type Images struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
type ArtistItem struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Followers    Followers    `json:"followers"`
	Genres       []string     `json:"genres"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Images       []Images     `json:"images"`
	Name         string       `json:"name"`
	Popularity   int          `json:"popularity"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}
type Artists struct {
	Href     string       `json:"href"`
	Limit    int          `json:"limit"`
	Next     string       `json:"next"`
	Offset   int          `json:"offset"`
	Previous any          `json:"previous"`
	Total    int          `json:"total"`
	Items    []ArtistItem `json:"items"`
}

type Albums struct {
	Href     string      `json:"href"`
	Limit    int         `json:"limit"`
	Next     string      `json:"next"`
	Offset   int         `json:"offset"`
	Previous any         `json:"previous"`
	Total    int         `json:"total"`
	Items    []AlbumItem `json:"items"`
}

type GetAlbumResponse struct {
	Href     string      `json:"href"`
	Limit    int         `json:"limit"`
	Next     any         `json:"next"`
	Offset   int         `json:"offset"`
	Previous any         `json:"previous"`
	Total    int         `json:"total"`
	Items    []AlbumItem `json:"items"`
}

type AlbumItem struct {
	AlbumType            string       `json:"album_type"`
	TotalTracks          int          `json:"total_tracks"`
	AvailableMarkets     []string     `json:"available_markets"`
	ExternalUrls         ExternalUrls `json:"external_urls"`
	Href                 string       `json:"href"`
	ID                   string       `json:"id"`
	Images               []Images     `json:"images"`
	Name                 string       `json:"name"`
	ReleaseDate          string       `json:"release_date"`
	ReleaseDatePrecision string       `json:"release_date_precision"`
	Type                 string       `json:"type"`
	URI                  string       `json:"uri"`
	Artists              []MiniArtist `json:"artists"`
	AlbumGroup           string       `json:"album_group"`
}

type MiniArtist struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

// We don't want live albums, compilations, etc. Preferably only studio albums
func (i AlbumItem) IsUnwantedAlbum() bool {
	lowerName := strings.ToLower(i.Name)
	if strings.Contains(lowerName, "(live") {
		return true
	}
	if strings.Contains(lowerName, "live at ") {
		return true
	}
	if strings.Contains(lowerName, "live in ") {
		return true
	}
	if strings.Contains(lowerName, "live on ") {
		return true
	}
	if strings.Contains(lowerName, "live from ") {
		return true
	}
	if strings.Contains(lowerName, "(tour") {
		return true
	}
	return false
}

type Albumv2 struct {
	Data struct {
		AlbumUnion struct {
			Typename     string `json:"__typename"`
			CourtesyLine string `json:"courtesyLine"`
			Date         struct {
				IsoString time.Time `json:"isoString"`
				Precision string    `json:"precision"`
			} `json:"date"`
			IsPreRelease bool   `json:"isPreRelease"`
			Label        string `json:"label"`
			Name         string `json:"name"`
			Playability  struct {
				Playable bool   `json:"playable"`
				Reason   string `json:"reason"`
			} `json:"playability"`
			PreReleaseEndDateTime any  `json:"preReleaseEndDateTime"`
			Saved                 bool `json:"saved"`
			SharingInfo           struct {
				ShareID  string `json:"shareId"`
				ShareURL string `json:"shareUrl"`
			} `json:"sharingInfo"`
			TracksV2 struct {
				Items []struct {
					Track struct {
						Artists struct {
							Items []struct {
								Profile struct {
									Name string `json:"name"`
								} `json:"profile"`
								URI string `json:"uri"`
							} `json:"items"`
						} `json:"artists"`
						AssociationsV3 struct {
							VideoAssociations struct {
								TotalCount int `json:"totalCount"`
							} `json:"videoAssociations"`
						} `json:"associationsV3"`
						ContentRating struct {
							Label string `json:"label"`
						} `json:"contentRating"`
						DiscNumber int `json:"discNumber"`
						Duration   struct {
							TotalMilliseconds int `json:"totalMilliseconds"`
						} `json:"duration"`
						Name        string `json:"name"`
						Playability struct {
							Playable bool `json:"playable"`
						} `json:"playability"`
						Playcount            string `json:"playcount"`
						RelinkingInformation any    `json:"relinkingInformation"`
						Saved                bool   `json:"saved"`
						TrackNumber          int    `json:"trackNumber"`
						URI                  string `json:"uri"`
					} `json:"track"`
					UID string `json:"uid"`
				} `json:"items"`
				TotalCount int `json:"totalCount"`
			} `json:"tracksV2"`
			Type  string `json:"type"`
			URI   string `json:"uri"`
			Items []struct {
				Discography struct {
					PopularReleasesAlbums struct {
						Items []struct {
							CoverArt struct {
								Sources []struct {
									Height int    `json:"height"`
									URL    string `json:"url"`
									Width  int    `json:"width"`
								} `json:"sources"`
							} `json:"coverArt"`
							Date struct {
								Year int `json:"year"`
							} `json:"date"`
							ID          string `json:"id"`
							Name        string `json:"name"`
							Playability struct {
								Playable bool   `json:"playable"`
								Reason   string `json:"reason"`
							} `json:"playability"`
							SharingInfo struct {
								ShareID  string `json:"shareId"`
								ShareURL string `json:"shareUrl"`
							} `json:"sharingInfo"`
							Type string `json:"type"`
							URI  string `json:"uri"`
						} `json:"items"`
					} `json:"popularReleasesAlbums"`
				} `json:"discography"`
			} `json:"items"`
		} `json:"albumUnion"`
	} `json:"data"`
}

// If all tracks contain "live" then it's a live album
func (a Albumv2) IsLiveAlbum() bool {
	for _, track := range a.Data.AlbumUnion.TracksV2.Items {
		trackName := strings.ToLower(track.Track.Name)
		if !strings.Contains(trackName, "live") {
			return false
		}
	}
	return true
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
