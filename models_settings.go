package meilisearch

type AttributeRule struct {
	AttributePatterns []string          `json:"attributePatterns"`
	Features          AttributeFeatures `json:"features"`
}

type AttributeFeatures struct {
	FacetSearch bool           `json:"facetSearch"`
	Filter      FilterFeatures `json:"filter"`
}

type FilterFeatures struct {
	Equality   bool `json:"equality"`
	Comparison bool `json:"comparison"`
}

// Settings is the type that represents the settings in meilisearch
type Settings struct {
	RankingRules         []string               `json:"rankingRules,omitempty"`
	DistinctAttribute    *string                `json:"distinctAttribute,omitempty"`
	SearchableAttributes []string               `json:"searchableAttributes,omitempty"`
	Dictionary           []string               `json:"dictionary,omitempty"`
	SearchCutoffMs       int64                  `json:"searchCutoffMs,omitempty"`
	ProximityPrecision   ProximityPrecisionType `json:"proximityPrecision,omitempty"`
	SeparatorTokens      []string               `json:"separatorTokens,omitempty"`
	NonSeparatorTokens   []string               `json:"nonSeparatorTokens,omitempty"`
	DisplayedAttributes  []string               `json:"displayedAttributes,omitempty"`
	StopWords            []string               `json:"stopWords,omitempty"`
	Synonyms             map[string][]string    `json:"synonyms,omitempty"`
	FilterableAttributes []string               `json:"filterableAttributes,omitempty"`
	SortableAttributes   []string               `json:"sortableAttributes,omitempty"`
	LocalizedAttributes  []*LocalizedAttributes `json:"localizedAttributes,omitempty"`
	TypoTolerance        *TypoTolerance         `json:"typoTolerance,omitempty"`
	Pagination           *Pagination            `json:"pagination,omitempty"`
	Faceting             *Faceting              `json:"faceting,omitempty"`
	Embedders            map[string]Embedder    `json:"embedders,omitempty"`
	PrefixSearch         *string                `json:"prefixSearch,omitempty"`
	FacetSearch          bool                   `json:"facetSearch,omitempty"`
	Chat                 *Chat                  `json:"chat,omitempty"`
	ForeignKeys          []ForeignKey           `json:"foreignKeys,omitempty"`
}

type ForeignKey struct {
	FieldName       string `json:"fieldName"`
	ForeignIndexUid string `json:"foreignIndexUid"`
}

type LocalizedAttributes struct {
	Locales           []string `json:"locales,omitempty"`
	AttributePatterns []string `json:"attributePatterns,omitempty"`
}

// TypoTolerance is the type that represents the typo tolerance setting in meilisearch
type TypoTolerance struct {
	Enabled             bool                `json:"enabled"`
	MinWordSizeForTypos MinWordSizeForTypos `json:"minWordSizeForTypos"`
	DisableOnWords      []string            `json:"disableOnWords"`
	DisableOnAttributes []string            `json:"disableOnAttributes"`
	DisableOnNumbers    bool                `json:"disableOnNumbers"`
}

// MinWordSizeForTypos is the type that represents the minWordSizeForTypos setting in the typo tolerance setting in meilisearch
type MinWordSizeForTypos struct {
	OneTypo  int64 `json:"oneTypo"`
	TwoTypos int64 `json:"twoTypos"`
}

// Pagination is the type that represents the pagination setting in meilisearch
type Pagination struct {
	MaxTotalHits int64 `json:"maxTotalHits"`
}

// Faceting is the type that represents the faceting setting in meilisearch
type Faceting struct {
	MaxValuesPerFacet int64 `json:"maxValuesPerFacet"`
	// SortFacetValuesBy index_name: alpha|count
	SortFacetValuesBy map[string]SortFacetType `json:"sortFacetValuesBy"`
}

type Chat struct {
	Description              string            `json:"description,omitempty"`
	DocumentTemplate         string            `json:"documentTemplate,omitempty"`
	DocumentTemplateMaxBytes int               `json:"documentTemplateMaxBytes,omitempty"`
	SearchParameters         *SearchParameters `json:"searchParameters,omitempty"`
}

type SearchParameters struct {
	Limit                 int64                `json:"limit,omitempty"`
	AttributesToSearchOn  []string             `json:"attributesToSearchOn,omitempty"`
	MatchingStrategy      MatchingStrategy     `json:"matchingStrategy,omitempty"`
	Sort                  []string             `json:"sort,omitempty"`
	Distinct              string               `json:"distinct,omitempty"`
	Hybrid                *SearchRequestHybrid `json:"hybrid,omitempty"`
	RankingScoreThreshold float64              `json:"rankingScoreThreshold,omitempty"`
}

// Embedder represents a unified configuration for various embedder types.
//
// Specs: https://www.meilisearch.com/docs/reference/api/settings#body-21
type Embedder struct {
	Source EmbedderSource `json:"source"` // The type of embedder: "openAi", "huggingFace", "userProvided", "rest", "ollama"
	// URL Meilisearch queries url to generate vector embeddings for queries and documents.
	// url must point to a REST-compatible embedder. You may also use url to work with proxies, such as when targeting openAi from behind a proxy.
	URL    string `json:"url,omitempty"`    // Optional for "openAi", "rest", "ollama"
	APIKey string `json:"apiKey,omitempty"` // Optional for "openAi", "rest", "ollama"
	// Model The model your embedder uses when generating vectors.
	// These are the officially supported models Meilisearch supports:
	//
	// - openAi: text-embedding-3-small, text-embedding-3-large, openai-text-embedding-ada-002
	//
	// - huggingFace: BAAI/bge-base-en-v1.5
	//
	// Other models, such as HuggingFace’s BERT models or those provided by Ollama and REST embedders
	// may also be compatible with Meilisearch.
	//
	// HuggingFace’s BERT models: https://huggingface.co/models?other=bert
	Model string `json:"model,omitempty"` // Optional for "openAi", "huggingFace", "ollama"
	// DocumentTemplate is a string containing a Liquid template. Meillisearch interpolates the template for each
	// document and sends the resulting text to the embedder. The embedder then generates document vectors based on this text.
	DocumentTemplate string `json:"documentTemplate,omitempty"` // Optional for most embedders
	// DocumentTemplateMaxBytes The maximum size of a rendered document template.
	//Longer texts are truncated to fit the configured limit.
	//
	// documentTemplateMaxBytes must be an integer. It defaults to 400.
	DocumentTemplateMaxBytes int                    `json:"documentTemplateMaxBytes,omitempty"`
	Dimensions               int                    `json:"dimensions,omitempty"`   // Optional for "openAi", "rest", "userProvided", "ollama"
	Revision                 string                 `json:"revision,omitempty"`     // Optional for "huggingFace"
	Distribution             *Distribution          `json:"distribution,omitempty"` // Optional for all embedders
	Request                  map[string]interface{} `json:"request,omitempty"`      // Optional for "rest"
	Response                 map[string]interface{} `json:"response,omitempty"`     // Optional for "rest"
	Headers                  map[string]string      `json:"headers,omitempty"`      // Optional for "rest"
	BinaryQuantized          bool                   `json:"binaryQuantized,omitempty"`
	Pooling                  EmbedderPooling        `json:"pooling,omitempty"`
	IndexingEmbedder         *Embedder              `json:"indexingEmbedder,omitempty"` // For Composite
	SearchEmbedder           *Embedder              `json:"searchEmbedder,omitempty"`   // For Composite
	IndexingFragments        map[string]Fragment    `json:"indexingFragments,omitempty"`
	SearchFragments          map[string]Fragment    `json:"searchFragments,omitempty"`
}

type Fragment struct {
	Value map[string]any `json:"value,omitempty"`
}

// Distribution represents a statistical distribution with mean and standard deviation (sigma).
type Distribution struct {
	Mean  float64 `json:"mean"`  // Mean of the distribution
	Sigma float64 `json:"sigma"` // Sigma (standard deviation) of the distribution
}

// ExperimentalFeaturesBase represents the experimental features result from the API.
type ExperimentalFeaturesBase struct {
	LogsRoute               *bool `json:"logsRoute,omitempty"`
	Metrics                 *bool `json:"metrics,omitempty"`
	EditDocumentsByFunction *bool `json:"editDocumentsByFunction,omitempty"`
	ContainsFilter          *bool `json:"containsFilter,omitempty"`
	Network                 *bool `json:"network,omitempty"`
	CompositeEmbedders      *bool `json:"compositeEmbedders,omitempty"`
	ChatCompletions         *bool `json:"chatCompletions,omitempty"`
	MultiModal              *bool `json:"multimodal,omitempty"`
	DynamicSearchRules      *bool `json:"dynamicSearchRules,omitempty"`
	GetTaskDocumentsRoute   *bool `json:"getTaskDocumentsRoute,omitempty"`
	RenderRoute             *bool `json:"renderRoute,omitempty"`
	ForeignKeys             *bool `json:"foreignKeys,omitempty"`
}

type ExperimentalFeaturesResult struct {
	LogsRoute               bool                      `json:"logsRoute"`
	Metrics                 bool                      `json:"metrics"`
	EditDocumentsByFunction bool                      `json:"editDocumentsByFunction"`
	ContainsFilter          bool                      `json:"containsFilter"`
	Network                 bool                      `json:"network"`
	CompositeEmbedders      bool                      `json:"compositeEmbedders"`
	ChatCompletions         bool                      `json:"chatCompletions"`
	MultiModal              bool                      `json:"multimodal"`
	DynamicSearchRules      bool                      `json:"dynamicSearchRules"`
	GetTaskDocumentsRoute   bool                      `json:"getTaskDocumentsRoute"`
	RenderRoute             bool                      `json:"renderRoute"`
	ForeignKeys             bool                      `json:"foreignKeys"`
	Personalize             *SearchRequestPersonalize `json:"personalize,omitempty"`
}
