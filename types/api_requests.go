package types

type PostQuery struct {
	Variables     Variables  `json:"variables"`
	OperationName string     `json:"operationName"`
	Extensions    Extensions `json:"extensions"`
}
type Variables struct {
	URI    string `json:"uri"`
	Locale string `json:"locale"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}
type PersistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}
type Extensions struct {
	PersistedQuery PersistedQuery `json:"persistedQuery"`
}
