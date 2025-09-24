// Struct for the display to filter out only the studio albums.
// Placed into artistdata/filter.json
type ArtistFitler = {
  artistName: string;
  albums: { name: string; releaseYear: number }[];
}[];
