package meilisearch

import "time"

type Unknown map[string]interface{}

//
// Internal types to Meilisearch
//

type SchemaAttributes string

const (
	SchemaAttributesDisplayed  SchemaAttributes = "displayed"
	SchemaAttributesIndexed    SchemaAttributes = "indexed"
	SchemaAttributesRanked     SchemaAttributes = "ranked"
	SchemaAttributesIdentifier SchemaAttributes = "identifier"
)

type Attributes map[string]bool

type RawAttribute struct {
	Displayed  bool `json:"displayed"`
	Indexed    bool `json:"indexed"`
	Ranked     bool `json:"ranked"`
	Identifier bool `json:"identifier"`
}
type RawSchema struct {
	Identifier string                  `json:"identifier"`
	Attributes map[string]RawAttribute `json:"attributes"`
}

type Schema map[string][]SchemaAttributes

type Index struct {
	Name      string    `json:"name"`
	Uid       string    `json:"uid"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Settings struct {
	RankingOrder  []string          `json:"rankingOrder,omitempty"`
	DistinctField string            `json:"distinctField,omitempty"`
	RankingRules  map[string]string `json:"rankingRules,omitempty"`
}

type Synonym struct {
	Input    string   `json:"input,omitempty"`
	Synonyms []string `json:"synonyms"`
}

type ACL string

const (
	IndexesRead    ACL = "IndexesRead"
	IndexesWrite   ACL = "IndexesWrite"
	DocumentsRead  ACL = "DocumentsRead"
	DocumentsWrite ACL = "DocumentsWrite"
	SettingsRead   ACL = "SettingsRead"
	SettingsWrite  ACL = "SettingsWrite"
	Admin          ACL = "Admin"
	All            ACL = "All"
)

type APIKey struct {
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Acl         []ACL     `json:"acl"`
	Indexes     []string  `json:"indexes"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	Revoked     bool      `json:"revoked"`
}

type Version struct {
	CommitSha  string    `json:"commitSha"`
	BuildDate  time.Time `json:"buildDate"`
	PkgVersion string    `json:"pkgVersion"`
}

type Stats struct {
	NumberOfDocuments int64            `json:"numberOfDocuments"`
	IsIndexing        bool             `json:"isIndexing"`
	FieldsFrequency   map[string]int64 `json:"fieldsFrequency"`
}

type SystemInformation struct {
	MemoryUsage    float64   `json:"memoryUsage"`
	ProcessorUsage []float64 `json:"processorUsage"`
	Global         struct {
		TotalMemory int64 `json:"totalMemory"`
		UsedMemory  int64 `json:"usedMemory"`
		UsedSwap    int64 `json:"usedSwap"`
		InputData   int64 `json:"inputData"`
		OutputData  int64 `json:"outputData"`
	} `json:"global"`
	Process struct {
		Memory int64 `json:"memory"`
		CPU    int64 `json:"cpu"`
	} `json:"process"`
}

type SystemInformationPretty struct {
	MemoryUsage    string   `json:"memoryUsage"`
	ProcessorUsage []string `json:"processorUsage"`
	Global         struct {
		TotalMemory string `json:"totalMemory"`
		UsedMemory  string `json:"usedMemory"`
		UsedSwap    string `json:"usedSwap"`
		InputData   string `json:"inputData"`
		OutputData  string `json:"outputData"`
	} `json:"global"`
	Process struct {
		Memory string `json:"memory"`
		CPU    string `json:"cpu"`
	} `json:"process"`
}

//
// Request/Response
//

type UpdateIdResponse struct {
	UpdateId int64 `json:"updateId"`
}

type CreateIndexRequest struct {
	Name   string `json:"name"`
	Uid    string `json:"uid,omitempty"`
	Schema Schema `json:"schema,omitempty"`
}

type CreateIndexResponse struct {
	Name      string    `json:"name"`
	Uid       string    `json:"uid"`
	Schema    Schema    `json:"schema,omitempty"`
	UpdateId  int64     `json:"updateId,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SearchRequest struct {
	Query                 string   `json:"q"`
	Offset                int64    `json:"offset,omitempty"`
	Limit                 int64    `json:"limit,omitempty"`
	AttributesToRetrieve  []string `json:"attributesToRetrieve,omitempty"`
	AttributesToSearchIn  []string `json:"attributesToSearchIn,omitempty"`
	AttributesToCrop      []string `json:"attributesYoCrop,omitempty"`
	CropLength            int64    `json:"cropLength,omitempty"`
	AttributesToHighlight []string `json:"attributesToHighlight,omitempty"`
	Filters               string   `json:"filters,omitempty"`
	TimeoutMs             int64    `json:"timeoutMs,omitempty"`
	Matches               bool     `json:"matches,omitempty"`
}

type realSearchRequest struct {
	Query                 string `json:"q"`
	Offset                int64  `json:"offset,omitempty"`
	Limit                 int64  `json:"limit,omitempty"`
	AttributesToRetrieve  string `json:"attributesToRetrieve,omitempty"`
	AttributesToSearchIn  string `json:"attributesToSearchIn,omitempty"`
	AttributesToCrop      string `json:"attributesYoCrop,omitempty"`
	CropLength            int64  `json:"cropLength,omitempty"`
	AttributesToHighlight string `json:"attributesToHighlight,omitempty"`
	Filters               string `json:"filters,omitempty"`
	TimeoutMs             int64  `json:"timeoutMs,omitempty"`
	Matches               bool   `json:"matches,omitempty"`
}

type SearchResponse struct {
	Hits             []interface{} `json:"hits"`
	Offset           int64         `json:"offset"`
	Limit            int64         `json:"limit"`
	ProcessingTimeMs int64         `json:"processingTimeMs"`
	Query            string        `json:"query"`
}

type ListDocumentsRequest struct {
	Offset               int64    `json:"offset,omitempty"`
	Limit                int64    `json:"limit,omitempty"`
	AttributesToRetrieve []string `json:"attributesToRetrieve,omitempty"`
}

type ListSynonymsResponse map[string][]string

type BatchCreateSynonymsRequest []Synonym

type CreateApiKeyRequest struct {
	Description string    `json:"description"`
	Acl         []ACL     `json:"acl"`
	Indexes     []string  `json:"indexes"`
	ExpireAt    time.Time `json:"expire_at"`
}

type UpdateApiKeyRequest struct {
	Description string   `json:"description"`
	Acl         []ACL    `json:"acl"`
	Indexes     []string `json:"indexes"`
	Revoked     bool     `json:"revoked"`
}
