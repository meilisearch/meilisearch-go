package meilisearch

type (
	TaskType               string // TaskType is the type of a task
	SortFacetType          string // SortFacetType is type of facet sorting, alpha or count
	TaskStatus             string // TaskStatus is the status of a task.
	ProximityPrecisionType string // ProximityPrecisionType accepts one of the ByWord or ByAttribute
	MatchingStrategy       string // MatchingStrategy one of the Last, All, Frequency
)

const (
	// Last returns documents containing all the query terms first. If there are not enough results containing all
	// query terms to meet the requested limit, Meilisearch will remove one query term at a time,
	// starting from the end of the query.
	Last MatchingStrategy = "last"
	// All only returns documents that contain all query terms. Meilisearch will not match any more documents even
	// if there aren't enough to meet the requested limit.
	All MatchingStrategy = "all"
	// Frequency returns documents containing all the query terms first. If there are not enough results containing
	//all query terms to meet the requested limit, Meilisearch will remove one query term at a time, starting
	//with the word that is the most frequent in the dataset. frequency effectively gives more weight to terms
	//that appear less frequently in a set of results.
	Frequency MatchingStrategy = "frequency"
)

const (
	// ByWord calculate the precise distance between query terms. Higher precision, but may lead to longer
	// indexing time. This is the default setting
	ByWord ProximityPrecisionType = "byWord"
	// ByAttribute determine if multiple query terms are present in the same attribute.
	// Lower precision, but shorter indexing time
	ByAttribute ProximityPrecisionType = "byAttribute"
)

const (
	// TaskStatusUnknown is the default TaskStatus, should not exist
	TaskStatusUnknown TaskStatus = "unknown"
	// TaskStatusEnqueued the task request has been received and will be processed soon
	TaskStatusEnqueued TaskStatus = "enqueued"
	// TaskStatusProcessing the task is being processed
	TaskStatusProcessing TaskStatus = "processing"
	// TaskStatusSucceeded the task has been successfully processed
	TaskStatusSucceeded TaskStatus = "succeeded"
	// TaskStatusFailed a failure occurred when processing the task, no changes were made to the database
	TaskStatusFailed TaskStatus = "failed"
	// TaskStatusCanceled the task was canceled
	TaskStatusCanceled TaskStatus = "canceled"
)

const (
	SortFacetTypeAlpha SortFacetType = "alpha"
	SortFacetTypeCount SortFacetType = "count"
)

const (
	// TaskTypeIndexCreation represents an index creation
	TaskTypeIndexCreation TaskType = "indexCreation"
	// TaskTypeIndexUpdate represents an index update
	TaskTypeIndexUpdate TaskType = "indexUpdate"
	// TaskTypeIndexDeletion represents an index deletion
	TaskTypeIndexDeletion TaskType = "indexDeletion"
	// TaskTypeIndexSwap represents an index swap
	TaskTypeIndexSwap TaskType = "indexSwap"
	// TaskTypeDocumentAdditionOrUpdate represents a document addition or update in an index
	TaskTypeDocumentAdditionOrUpdate TaskType = "documentAdditionOrUpdate"
	// TaskTypeDocumentDeletion represents a document deletion from an index
	TaskTypeDocumentDeletion TaskType = "documentDeletion"
	// TaskTypeSettingsUpdate represents a settings update
	TaskTypeSettingsUpdate TaskType = "settingsUpdate"
	// TaskTypeDumpCreation represents a dump creation
	TaskTypeDumpCreation TaskType = "dumpCreation"
	// TaskTypeTaskCancelation represents a task cancelation
	TaskTypeTaskCancelation TaskType = "taskCancelation"
	// TaskTypeTaskDeletion represents a task deletion
	TaskTypeTaskDeletion TaskType = "taskDeletion"
	// TaskTypeSnapshotCreation represents a snapshot creation
	TaskTypeSnapshotCreation TaskType = "snapshotCreation"
	// TaskTypeExport represents a task exportation
	TaskTypeExport TaskType = "export"
)

type (
	ContentEncoding          string
	EncodingCompressionLevel int
)

const (
	GzipEncoding    ContentEncoding = "gzip"
	DeflateEncoding ContentEncoding = "deflate"
	BrotliEncoding  ContentEncoding = "br"

	NoCompression          EncodingCompressionLevel = 0
	BestSpeed              EncodingCompressionLevel = 1
	BestCompression        EncodingCompressionLevel = 9
	DefaultCompression     EncodingCompressionLevel = -1
	HuffmanOnlyCompression EncodingCompressionLevel = -2
	ConstantCompression    EncodingCompressionLevel = -2
	StatelessCompression   EncodingCompressionLevel = -3
)

func (c ContentEncoding) String() string { return string(c) }

func (c ContentEncoding) IsZero() bool { return c == "" }

func (c EncodingCompressionLevel) Int() int { return int(c) }

// EmbedderSource The source corresponds to a service that generates embeddings from your documents.
type EmbedderSource string

const (
	OpenaiEmbedderSource      EmbedderSource = "openAi"
	HuggingFaceEmbedderSource EmbedderSource = "huggingFace"
	// RestEmbedderSource rest is a generic source compatible with any embeddings provider offering a REST API.
	RestEmbedderSource   EmbedderSource = "rest"
	OllamaEmbedderSource EmbedderSource = "ollama"
	// UserProvidedEmbedderSource Use userProvided when you want to generate embeddings manually. In this case,
	// you must include vector data in your documents’ _vectors field. You must also generate
	//vectors for search queries.
	UserProvidedEmbedderSource EmbedderSource = "userProvided"
	// CompositeEmbedderSource Choose composite to use one embedder during indexing time, and another embedder at search time.
	//Must be used together with indexingEmbedder and searchEmbedder.
	CompositeEmbedderSource EmbedderSource = "composite"
)

// ChatSource The source corresponds to a service that generates chat completions from your messages.
type ChatSource string

const (
	OpenaiChatSource      ChatSource = "openAi"
	AzureOpenAiChatSource ChatSource = "azureOpenAi"
	MistralChatSource     ChatSource = "mistral"
	GeminiChatSource      ChatSource = "gemini"
	VLlmChatSource        ChatSource = "vLlm"
)

// EmbedderPooling Configure how Meilisearch should merge individual tokens into a single embedding.
//
// pooling is optional for embedders with the huggingFace source.
type EmbedderPooling string

const (
	// UseModelEmbedderPooling Meilisearch will fetch the pooling method from the model configuration. Default value for new embedders
	UseModelEmbedderPooling EmbedderPooling = "useModel"
	// ForceMeanEmbedderPooling always use mean pooling. Default value for embedders created in Meilisearch <=v1.13
	ForceMeanEmbedderPooling EmbedderPooling = "forceMean"
	// ForceClsEmbedderPooling always use CLS pooling
	ForceClsEmbedderPooling EmbedderPooling = "forceCls"
)

// ChatRole The role of a message in a chat conversation.
type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleSystem    ChatRole = "system"
)
