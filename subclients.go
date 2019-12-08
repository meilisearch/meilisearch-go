package meilisearch

type SubClientIndexes interface {
	Get(uid string) (Index, error)
	List() ([]Index, error)
	Create(request CreateIndexRequest) (CreateIndexResponse, error)
	Update(request UpdateIndexRequest) (Index, error)
	Delete(uid string) (UpdateIdResponse, error)

	GetRawSchema() (SchemaRaw, error)
	GetSchema() (Schema, error)
	UpdateSchema(schema Schema) (UpdateIdResponse, error)
	UpdateWithRawSchema(schema SchemaRaw) (UpdateIdResponse, error)
}

type SubClientDocuments interface {
	Get(identifier string) (interface{}, error)
	Delete(identifier string) (UpdateIdResponse, error)
	Deletes(identifier []string) (UpdateIdResponse, error)
	List(request ListDocumentsRequest) ([]interface{}, error)
	AddOrUpdate([]interface{}) (UpdateIdResponse, error)
	ClearAllDocuments() (UpdateIdResponse, error)
}

type SubClientSearch interface {
	Search(params SearchRequest) (SearchResponse, error)
}

type SubClientSynonyms interface {
	List(word string) ([]string, error)
	ListAll() ([]ListSynonymsResponse, error)
	Create(word string, synonyms []string) (UpdateIdResponse, error)
	Update(word string, synonyms []string) (UpdateIdResponse, error)
	Delete(word string) (UpdateIdResponse, error)
	BatchCreate(request BatchCreateSynonymsRequest) (UpdateIdResponse, error)
	DeleteAll() (UpdateIdResponse, error)
}

type SubClientStopWords interface {
	List() ([]string, error)
	Add(words []string) ([]UpdateIdResponse, error)
	Deletes(words []string) ([]UpdateIdResponse, error)
}

type SubClientUpdates interface {
	Get(id int64) (Unknown, error)
	List() ([]Unknown, error)
}

type SubClientKey interface {
	Get(key string) (APIKey, error)
	List() ([]APIKey, error)
	Create(request CreateApiKeyRequest) (APIKey, error)
	Update(request UpdateApiKeyRequest) (APIKey, error)
	Delete(key string) error
}

type SubClientSettings interface {
	Get() (Settings, error)
	AddOrUpdate(request Settings) (UpdateIdResponse, error)
}

type SubClientStats interface {
	Get() (Stats, error)
	List() ([]Stats, error)
}

type SubClientHealth interface {
	Get() error
	Set(health bool) error
}

type SubClientVersion interface {
	Get() (Version, error)
}

type SubClientSystemInformation interface {
	Get() (SystemInformation, error)
	GetPretty() (SystemInformationPretty, error)
}
