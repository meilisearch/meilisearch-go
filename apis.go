package meilisearch

type ApiIndexes interface {
	Get(uid string) (*Index, error)
	List() ([]Index, error)
	Create(request CreateIndexRequest) (*CreateIndexResponse, error)
	Update(uid string, name string) (*Index, error)
	Delete(uid string) (bool, error)

	GetRawSchema(uid string) (*SchemaRaw, error)
	GetSchema(uid string) (*Schema, error)
	UpdateSchema(uid string, schema Schema) (*UpdateIdResponse, error)
	UpdateWithRawSchema(uid string, schema SchemaRaw) (*UpdateIdResponse, error)
}

type ApiDocuments interface {
	Get(identifier string, documentPtr interface{}) error
	Delete(identifier string) (*UpdateIdResponse, error)
	Deletes(identifier []string) (*UpdateIdResponse, error)
	List(request ListDocumentsRequest, documentsPtr interface{}) error
	AddOrUpdate(documentsPtr interface{}) (*UpdateIdResponse, error)
	ClearAllDocuments() (*UpdateIdResponse, error)
}

type ApiSearch interface {
	Search(params SearchRequest) (*SearchResponse, error)
}

type ApiSynonyms interface {
	List(word string) ([]string, error)
	ListAll() ([]ListSynonymsResponse, error)
	Create(word string, synonyms []string) (*UpdateIdResponse, error)
	Update(word string, synonyms []string) (*UpdateIdResponse, error)
	Delete(word string) (*UpdateIdResponse, error)
	BatchCreate(request BatchCreateSynonymsRequest) (*UpdateIdResponse, error)
	DeleteAll() (*UpdateIdResponse, error)
}

type ApiStopWords interface {
	List() ([]string, error)
	Add(words []string) ([]UpdateIdResponse, error)
	Deletes(words []string) ([]UpdateIdResponse, error)
}

type ApiUpdates interface {
	Get(id int64) (*Unknown, error)
	List() ([]Unknown, error)
}

type ApiKey interface {
	Get(key string) (*APIKey, error)
	List() ([]APIKey, error)
	Create(request CreateApiKeyRequest) (*APIKey, error)
	Update(request UpdateApiKeyRequest) (*APIKey, error)
	Delete(key string) error
}

type ApiSettings interface {
	Get() (*Settings, error)
	AddOrUpdate(request Settings) (*UpdateIdResponse, error)
}

type ApiStats interface {
	Get() (*Stats, error)
	List() ([]Stats, error)
}

type ApiHealth interface {
	Get() error
	Set(health bool) error
}

type ApiVersion interface {
	Get() (*Version, error)
}

type ApiSystemInformation interface {
	Get() (*SystemInformation, error)
	GetPretty() (*SystemInformationPretty, error)
}
