package meilisearch

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

type docTest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type docTestBooks struct {
	BookID int    `json:"book_id"`
	Title  string `json:"title"`
	Tag    string `json:"tag"`
	Year   int    `json:"year"`
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func deleteAllIndexes(client ClientInterface) (ok bool, err error) {
	list, err := client.GetIndexes(nil)
	if err != nil {
		return false, err
	}

	for _, index := range list.Results {
		task, _ := client.DeleteIndex(index.UID)
		_, err := client.WaitForTask(task.TaskUID)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func deleteAllKeys(client ClientInterface) (ok bool, err error) {
	list, err := client.GetKeys(nil)
	if err != nil {
		return false, err
	}

	for _, key := range list.Results {
		if strings.Contains(key.Description, "Test") || (key.Description == "") {
			_, err = client.DeleteKey(key.Key)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func cleanup(c ClientInterface) func() {
	return func() {
		_, _ = deleteAllIndexes(c)
		_, _ = deleteAllKeys(c)
	}
}

func testWaitForTask(t *testing.T, i *Index, u *TaskInfo) {
	_, err := i.WaitForTask(u.TaskUID)
	require.NoError(t, err)
}

func testWaitForBatchTask(t *testing.T, i *Index, u []TaskInfo) {
	for _, id := range u {
		_, err := i.WaitForTask(id.TaskUID)
		require.NoError(t, err)
	}
}

func GetPrivateKey() (key string) {
	list, err := defaultClient.GetKeys(nil)
	if err != nil {
		return ""
	}
	for _, key := range list.Results {
		if strings.Contains(key.Name, "Default Admin API Key") || (key.Description == "") {
			return key.Key
		}
	}
	return ""
}

func GetPrivateUIDKey() (key string) {
	list, err := defaultClient.GetKeys(nil)
	if err != nil {
		return ""
	}
	for _, key := range list.Results {
		if strings.Contains(key.Name, "Default Admin API Key") || (key.Description == "") {
			return key.UID
		}
	}
	return ""
}

func SetUpEmptyIndex(index *IndexConfig) (resp *Index, err error) {
	client := NewClient(ClientConfig{
		Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
		APIKey: masterKey,
	})
	task, err := client.CreateIndex(index)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	finalTask, _ := client.WaitForTask(task.TaskUID)
	if finalTask.Status != "succeeded" {
		os.Exit(1)
	}
	return client.GetIndex(index.Uid)
}

func SetUpIndexWithVector(indexUID string) (resp *Index, err error) {
	client := NewClient(ClientConfig{
		Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
		APIKey: masterKey,
	})

	req := internalRequest{
		endpoint:    "/experimental-features",
		method:      http.MethodPatch,
		contentType: "application/json",
		withRequest: map[string]interface{}{
			"vectorStore": true,
		},
	}

	if err := client.executeRequest(req); err != nil {
		return nil, err
	}

	index := client.Index(indexUID)

	documents := []map[string]interface{}{
		{"book_id": 123, "title": "Pride and Prejudice", "_vectors": []float64{0.1, 0.2, 0.3}},
		{"book_id": 456, "title": "Le Petit Prince", "_vectors": []float64{2.4, 8.5, 1.6}},
	}

	task, err := index.AddDocuments(documents)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	finalTask, _ := index.WaitForTask(task.TaskUID)
	if finalTask.Status != "succeeded" {
		os.Exit(1)
	}

	return client.GetIndex(indexUID)
}

func SetUpBasicIndex(indexUID string) {
	client := NewClient(ClientConfig{
		Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
		APIKey: masterKey,
	})
	index := client.Index(indexUID)

	documents := []map[string]interface{}{
		{"book_id": 123, "title": "Pride and Prejudice"},
		{"book_id": 456, "title": "Le Petit Prince"},
		{"book_id": 1, "title": "Alice In Wonderland"},
		{"book_id": 1344, "title": "The Hobbit"},
		{"book_id": 4, "title": "Harry Potter and the Half-Blood Prince"},
		{"book_id": 42, "title": "The Hitchhiker's Guide to the Galaxy"},
	}
	task, err := index.AddDocuments(documents)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	finalTask, _ := index.WaitForTask(task.TaskUID)
	if finalTask.Status != "succeeded" {
		os.Exit(1)
	}
}

func SetUpIndexWithNestedFields(indexUID string) {
	client := NewClient(ClientConfig{
		Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
		APIKey: masterKey,
	})
	index := client.Index(indexUID)

	documents := []map[string]interface{}{
		{"id": 1, "title": "Pride and Prejudice", "info": map[string]interface{}{"comment": "A great book", "reviewNb": 50}},
		{"id": 2, "title": "Le Petit Prince", "info": map[string]interface{}{"comment": "A french book", "reviewNb": 600}},
		{"id": 3, "title": "Le Rouge et le Noir", "info": map[string]interface{}{"comment": "Another french book", "reviewNb": 700}},
		{"id": 4, "title": "Alice In Wonderland", "comment": "A weird book", "info": map[string]interface{}{"comment": "A weird book", "reviewNb": 800}},
		{"id": 5, "title": "The Hobbit", "info": map[string]interface{}{"comment": "An awesome book", "reviewNb": 900}},
		{"id": 6, "title": "Harry Potter and the Half-Blood Prince", "info": map[string]interface{}{"comment": "The best book", "reviewNb": 1000}},
		{"id": 7, "title": "The Hitchhiker's Guide to the Galaxy"},
	}
	task, err := index.AddDocuments(documents)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	finalTask, _ := index.WaitForTask(task.TaskUID)
	if finalTask.Status != "succeeded" {
		os.Exit(1)
	}
}

func SetUpIndexForFaceting() {
	client := NewClient(ClientConfig{
		Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
		APIKey: masterKey,
	})
	index := client.Index("indexUID")

	booksTest := []docTestBooks{
		{BookID: 123, Title: "Pride and Prejudice", Tag: "Romance", Year: 1813},
		{BookID: 456, Title: "Le Petit Prince", Tag: "Tale", Year: 1943},
		{BookID: 1, Title: "Alice In Wonderland", Tag: "Tale", Year: 1865},
		{BookID: 1344, Title: "The Hobbit", Tag: "Epic fantasy", Year: 1937},
		{BookID: 4, Title: "Harry Potter and the Half-Blood Prince", Tag: "Epic fantasy", Year: 2005},
		{BookID: 42, Title: "The Hitchhiker's Guide to the Galaxy", Tag: "Epic fantasy", Year: 1978},
		{BookID: 742, Title: "The Great Gatsby", Tag: "Tragedy", Year: 1925},
		{BookID: 834, Title: "One Hundred Years of Solitude", Tag: "Tragedy", Year: 1967},
		{BookID: 17, Title: "In Search of Lost Time", Tag: "Modernist literature", Year: 1913},
		{BookID: 204, Title: "Ulysses", Tag: "Novel", Year: 1922},
		{BookID: 7, Title: "Don Quixote", Tag: "Satiric", Year: 1605},
		{BookID: 10, Title: "Moby Dick", Tag: "Novel", Year: 1851},
		{BookID: 730, Title: "War and Peace", Tag: "Historical fiction", Year: 1865},
		{BookID: 69, Title: "Hamlet", Tag: "Tragedy", Year: 1598},
		{BookID: 32, Title: "The Odyssey", Tag: "Epic", Year: 1571},
		{BookID: 71, Title: "Madame Bovary", Tag: "Novel", Year: 1857},
		{BookID: 56, Title: "The Divine Comedy", Tag: "Epic", Year: 1303},
		{BookID: 254, Title: "Lolita", Tag: "Novel", Year: 1955},
		{BookID: 921, Title: "The Brothers Karamazov", Tag: "Novel", Year: 1879},
		{BookID: 1032, Title: "Crime and Punishment", Tag: "Crime fiction", Year: 1866},
	}
	task, err := index.AddDocuments(booksTest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	finalTask, _ := index.WaitForTask(task.TaskUID)
	if finalTask.Status != "succeeded" {
		os.Exit(1)
	}
}

var (
	masterKey     = "masterKey"
	defaultClient = NewClient(ClientConfig{
		Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
		APIKey: masterKey,
	})
	defaultRankingRules = []string{
		"words", "typo", "proximity", "attribute", "sort", "exactness",
	}
	defaultTypoTolerance = TypoTolerance{
		Enabled: true,
		MinWordSizeForTypos: MinWordSizeForTypos{
			OneTypo:  5,
			TwoTypos: 9,
		},
		DisableOnWords:      []string{},
		DisableOnAttributes: []string{},
	}
	defaultPagination = Pagination{
		MaxTotalHits: 1000,
	}
	defaultFaceting = Faceting{
		MaxValuesPerFacet: 100,
	}
)

var customClient = NewFastHTTPCustomClient(ClientConfig{
	Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
	APIKey: masterKey,
},
	&fasthttp.Client{
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Name:      "custom-client",
	})

var timeoutClient = NewClient(ClientConfig{
	Host:    getenv("MEILISEARCH_URL", "http://localhost:7700"),
	APIKey:  masterKey,
	Timeout: 1,
})

var privateClient = NewClient(ClientConfig{
	Host:   getenv("MEILISEARCH_URL", "http://localhost:7700"),
	APIKey: GetPrivateKey(),
})

func TestMain(m *testing.M) {
	_, _ = deleteAllIndexes(defaultClient)
	code := m.Run()
	_, _ = deleteAllIndexes(defaultClient)
	os.Exit(code)
}

func Test_deleteAllIndexes(t *testing.T) {
	indexUIDS := []string{
		"Test_deleteAllIndexes",
		"Test_deleteAllIndexes2",
		"Test_deleteAllIndexes3",
	}
	_, _ = deleteAllIndexes(defaultClient)

	for _, uid := range indexUIDS {
		task, err := defaultClient.CreateIndex(&IndexConfig{
			Uid: uid,
		})
		if err != nil {
			t.Fatal(err)
		}
		_, err = defaultClient.WaitForTask(task.TaskUID)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, _ = deleteAllIndexes(defaultClient)

	for _, uid := range indexUIDS {
		resp, err := defaultClient.GetIndex(uid)
		if resp != nil {
			t.Fatal(resp)
		}
		if err == nil {
			t.Fatal("deleteAllIndexes: One or more indexes were not deleted")
		}
	}
}
