package meilisearch

import "time"

// Task indicates information about a task resource
//
// Documentation: https://www.meilisearch.com/docs/learn/advanced/asynchronous_operations
type Task struct {
	Status         TaskStatus          `json:"status"`
	UID            int64               `json:"uid,omitempty"`
	TaskUID        int64               `json:"taskUid,omitempty"`
	IndexUID       string              `json:"indexUid"`
	Type           TaskType            `json:"type"`
	Error          meilisearchApiError `json:"error,omitempty"`
	TaskNetwork    TaskNetwork         `json:"network,omitempty"`
	Duration       string              `json:"duration,omitempty"`
	EnqueuedAt     time.Time           `json:"enqueuedAt"`
	StartedAt      time.Time           `json:"startedAt,omitempty"`
	FinishedAt     time.Time           `json:"finishedAt,omitempty"`
	Details        Details             `json:"details,omitempty"`
	CanceledBy     int64               `json:"canceledBy,omitempty"`
	CustomMetadata string              `json:"customMetadata,omitempty"`
}

// TaskNetwork indicates information about a task network
//
// Documentation: https://www.meilisearch.com/docs/reference/api/tasks#network
type TaskNetwork struct {
	Origin  *Origin                `json:"origin,omitempty"`
	Remotes map[string]*TaskRemote `json:"remotes,omitempty"`
}

type Origin struct {
	RemoteName string `json:"remoteName,omitempty"`
	TaskUID    string `json:"taskUid,omitempty"`
}

type TaskRemote struct {
	TaskUID *string `json:"task_uid,omitempty"`
	Error   *string `json:"error,omitempty"`
}

// TaskInfo indicates information regarding a task returned by an asynchronous method
//
// Documentation: https://www.meilisearch.com/docs/reference/api/tasks#tasks
type TaskInfo struct {
	Status     TaskStatus `json:"status"`
	TaskUID    int64      `json:"taskUid"`
	IndexUID   string     `json:"indexUid"`
	Type       TaskType   `json:"type"`
	EnqueuedAt time.Time  `json:"enqueuedAt"`
}

// TasksQuery is a list of filter available to send as query parameters
type TasksQuery struct {
	UIDS             []int64
	Limit            int64
	From             int64
	IndexUIDS        []string
	Statuses         []TaskStatus
	Types            []TaskType
	CanceledBy       []int64
	BeforeEnqueuedAt time.Time
	AfterEnqueuedAt  time.Time
	BeforeStartedAt  time.Time
	AfterStartedAt   time.Time
	BeforeFinishedAt time.Time
	AfterFinishedAt  time.Time
	Reverse          bool
}

// CancelTasksQuery is a list of filter available to send as query parameters
type CancelTasksQuery struct {
	UIDS             []int64
	IndexUIDS        []string
	Statuses         []TaskStatus
	Types            []TaskType
	BeforeEnqueuedAt time.Time
	AfterEnqueuedAt  time.Time
	BeforeStartedAt  time.Time
	AfterStartedAt   time.Time
}

// DeleteTasksQuery is a list of filter available to send as query parameters
type DeleteTasksQuery struct {
	UIDS             []int64
	IndexUIDS        []string
	Statuses         []TaskStatus
	Types            []TaskType
	CanceledBy       []int64
	BeforeEnqueuedAt time.Time
	AfterEnqueuedAt  time.Time
	BeforeStartedAt  time.Time
	AfterStartedAt   time.Time
	BeforeFinishedAt time.Time
	AfterFinishedAt  time.Time
}

type Details struct {
	ReceivedDocuments    int64               `json:"receivedDocuments,omitempty"`
	IndexedDocuments     int64               `json:"indexedDocuments,omitempty"`
	DeletedDocuments     int64               `json:"deletedDocuments,omitempty"`
	PrimaryKey           string              `json:"primaryKey,omitempty"`
	ProvidedIds          int64               `json:"providedIds,omitempty"`
	RankingRules         []string            `json:"rankingRules,omitempty"`
	DistinctAttribute    *string             `json:"distinctAttribute,omitempty"`
	SearchableAttributes []string            `json:"searchableAttributes,omitempty"`
	DisplayedAttributes  []string            `json:"displayedAttributes,omitempty"`
	StopWords            []string            `json:"stopWords,omitempty"`
	Synonyms             map[string][]string `json:"synonyms,omitempty"`
	FilterableAttributes []interface{}       `json:"filterableAttributes,omitempty"`
	SortableAttributes   []string            `json:"sortableAttributes,omitempty"`
	TypoTolerance        *TypoTolerance      `json:"typoTolerance,omitempty"`
	Pagination           *Pagination         `json:"pagination,omitempty"`
	Faceting             *Faceting           `json:"faceting,omitempty"`
	MatchedTasks         int64               `json:"matchedTasks,omitempty"`
	CanceledTasks        int64               `json:"canceledTasks,omitempty"`
	DeletedTasks         int64               `json:"deletedTasks,omitempty"`
	OriginalFilter       string              `json:"originalFilter,omitempty"`
	Swaps                []SwapIndexesParams `json:"swaps,omitempty"`
	DumpUid              string              `json:"dumpUid,omitempty"`
}

// TaskResult return of multiple tasks is wrap in a TaskResult
type TaskResult struct {
	Results []Task `json:"results"`
	Limit   int64  `json:"limit"`
	From    int64  `json:"from"`
	Next    int64  `json:"next"`
	Total   int64  `json:"total"`
}

// Batch gives information about the progress of batch of asynchronous operations.
type Batch struct {
	UID           int                    `json:"uid"`
	Progress      *BatchProgress         `json:"progress,omitempty"`
	Details       map[string]interface{} `json:"details,omitempty"`
	Stats         *BatchStats            `json:"stats,omitempty"`
	Duration      string                 `json:"duration,omitempty"`
	StartedAt     time.Time              `json:"startedAt,omitempty"`
	FinishedAt    time.Time              `json:"finishedAt,omitempty"`
	BatchStrategy string                 `json:"batchStrategy,omitempty"`
}

type BatchProgress struct {
	Steps      []*BatchProgressStep `json:"steps"`
	Percentage float64              `json:"percentage"`
}

type BatchProgressStep struct {
	CurrentStep string `json:"currentStep"`
	Finished    int    `json:"finished"`
	Total       int    `json:"total"`
}

type BatchStats struct {
	TotalNbTasks           int                               `json:"totalNbTasks"`
	Status                 map[string]int                    `json:"status"`
	Types                  map[string]int                    `json:"types"`
	IndexedUIDs            map[string]int                    `json:"indexUids"`
	ProgressTrace          map[string]string                 `json:"progressTrace"`
	WriteChannelCongestion *BatchStatsWriteChannelCongestion `json:"writeChannelCongestion"`
	InternalDatabaseSizes  *BatchStatsInternalDatabaseSize   `json:"internalDatabaseSizes"`
}

type BatchStatsWriteChannelCongestion struct {
	Attempts         int     `json:"attempts"`
	BlockingAttempts int     `json:"blocking_attempts"`
	BlockingRatio    float64 `json:"blocking_ratio"`
}

type BatchStatsInternalDatabaseSize struct {
	ExternalDocumentsIDs    string `json:"externalDocumentsIds"`
	WordDocIDs              string `json:"wordDocids"`
	WordPairProximityDocIDs string `json:"wordPairProximityDocids"`
	WordPositionDocIDs      string `json:"wordPositionDocids"`
	WordFidDocIDs           string `json:"wordFidDocids"`
	FieldIdWordCountDocIDs  string `json:"fieldIdWordCountDocids"`
	Documents               string `json:"documents"`
}

type BatchesResults struct {
	Results []*Batch `json:"results"`
	Total   int64    `json:"total"`
	Limit   int64    `json:"limit"`
	From    int64    `json:"from"`
	Next    int64    `json:"next"`
}

// BatchesQuery represents the query parameters for listing batches.
type BatchesQuery struct {
	UIDs             []int64
	BatchUIDs        []int64
	IndexUIDs        []string
	Statuses         []string
	Types            []string
	Limit            int64
	From             int64
	Reverse          bool
	BeforeEnqueuedAt time.Time
	BeforeStartedAt  time.Time
	BeforeFinishedAt time.Time
	AfterEnqueuedAt  time.Time
	AfterStartedAt   time.Time
	AfterFinishedAt  time.Time
}
