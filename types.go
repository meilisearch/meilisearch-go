package meilisearch

const (
	contentTypeJSON   string = "application/json"
	contentTypeNDJSON string = "application/x-ndjson"
	contentTypeCSV    string = "text/csv"
	nullBody                 = "null"
)

// Version is the type that represents the versions in meilisearch
type Version struct {
	CommitSha  string `json:"commitSha"`
	CommitDate string `json:"commitDate"`
	PkgVersion string `json:"pkgVersion"`
}

// Health is the request body for set meilisearch health
type Health struct {
	Status string `json:"status"`
}

// JSONMarshal returns the JSON encoding of v.
type JSONMarshal func(v interface{}) ([]byte, error)

// JSONUnmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by v. If v is nil or not a pointer,
// Unmarshal returns an InvalidUnmarshalError.
type JSONUnmarshal func(data []byte, v interface{}) error
