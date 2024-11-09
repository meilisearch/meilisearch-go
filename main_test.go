package meilisearch

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

var (
	masterKey           = "masterKey"
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
		SortFacetValuesBy: map[string]SortFacetType{
			"*": SortFacetTypeAlpha,
		},
	}
)

var testNdjsonDocuments = []byte(`{"id": 1, "name": "Alice In Wonderland"}
{"id": 2, "name": "Pride and Prejudice"}
{"id": 3, "name": "Le Petit Prince"}
{"id": 4, "name": "The Great Gatsby"}
{"id": 5, "name": "Don Quixote"}
`)

var testCsvDocuments = []byte(`id,name
1,Alice In Wonderland
2,Pride and Prejudice
3,Le Petit Prince
4,The Great Gatsby
5,Don Quixote
`)

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

func setup(t *testing.T, host string, options ...Option) ServiceManager {
	t.Helper()

	opts := make([]Option, 0)
	opts = append(opts, WithAPIKey(masterKey))
	opts = append(opts, options...)

	if host == "" {
		host = getenv("MEILISEARCH_URL", "http://localhost:7700")
	}

	sv := New(host, opts...)
	return sv
}

func cleanup(services ...ServiceManager) func() {
	return func() {
		for _, s := range services {
			_, _ = deleteAllIndexes(s)
			_, _ = deleteAllKeys(s)
		}
	}
}

func getPrivateKey(sv ServiceManager) (key string) {
	list, err := sv.GetKeys(nil)
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

func getPrivateUIDKey(sv ServiceManager) (key string) {
	list, err := sv.GetKeys(nil)
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

func deleteAllIndexes(sv ServiceManager) (ok bool, err error) {
	list, err := sv.ListIndexes(nil)
	if err != nil {
		return false, err
	}

	for _, index := range list.Results {
		task, _ := sv.DeleteIndex(index.UID)
		_, err := sv.WaitForTask(task.TaskUID, 0)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func deleteAllKeys(sv ServiceManager) (ok bool, err error) {
	list, err := sv.GetKeys(nil)
	if err != nil {
		return false, err
	}

	for _, key := range list.Results {
		if strings.Contains(key.Description, "Test") || (key.Description == "") {
			_, err = sv.DeleteKey(key.Key)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func testWaitForTask(t *testing.T, i IndexManager, u *TaskInfo) {
	t.Helper()
	r, err := i.WaitForTask(u.TaskUID, 0)
	require.NoError(t, err)
	require.Equal(t, TaskStatusSucceeded, r.Status, fmt.Sprintf("Task failed: %#+v", r))
}

func testWaitForBatchTask(t *testing.T, i IndexManager, u []TaskInfo) {
	for _, id := range u {
		_, err := i.WaitForTask(id.TaskUID, 0)
		require.NoError(t, err)
	}
}

func setUpEmptyIndex(sv ServiceManager, index *IndexConfig) (resp *IndexResult, err error) {
	task, err := sv.CreateIndex(index)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	finalTask, _ := sv.WaitForTask(task.TaskUID, 0)
	if finalTask.Status != "succeeded" {
		cleanup(sv)
		return setUpEmptyIndex(sv, index)
	}
	return sv.GetIndex(index.Uid)
}

func setUpBasicIndex(sv ServiceManager, indexUID string) {
	index := sv.Index(indexUID)

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
	finalTask, _ := index.WaitForTask(task.TaskUID, 0)
	if finalTask.Status != "succeeded" {
		os.Exit(1)
	}
}

func setupMovieIndex(t *testing.T, client ServiceManager) IndexManager {
	t.Helper()

	idx := client.Index("indexUID")

	testdata, err := os.Open("./testdata/movies.json")
	require.NoError(t, err)
	defer testdata.Close()

	tests := make([]map[string]interface{}, 0)

	require.NoError(t, json.NewDecoder(testdata).Decode(&tests))

	task, err := idx.AddDocuments(tests)
	require.NoError(t, err)
	testWaitForTask(t, idx, task)

	task, err = idx.UpdateFilterableAttributes(&[]string{"id"})
	require.NoError(t, err)
	testWaitForTask(t, idx, task)

	return idx
}

func setUpIndexForFaceting(client ServiceManager) {
	idx := client.Index("indexUID")

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
		{BookID: 1039, Title: "The Girl in the white shirt", Tag: "white shirt", Year: 1999},
		{BookID: 1050, Title: "星の王子さま", Tag: "物語", Year: 1943},
	}
	task, err := idx.AddDocuments(booksTest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	finalTask, _ := idx.WaitForTask(task.TaskUID, 0)
	if finalTask.Status != "succeeded" {
		os.Exit(1)
	}
}

func setUpIndexWithNestedFields(client ServiceManager, indexUID string) {
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
	finalTask, _ := index.WaitForTask(task.TaskUID, 0)
	if finalTask.Status != "succeeded" {
		os.Exit(1)
	}
}

func setUpIndexWithVector(client *meilisearch, indexUID string) (resp *IndexResult, err error) {
	req := &internalRequest{
		endpoint:    "/experimental-features",
		method:      http.MethodPatch,
		contentType: "application/json",
		withRequest: map[string]interface{}{
			"vectorStore": true,
		},
	}

	if err := client.client.executeRequest(context.Background(), req); err != nil {
		return nil, err
	}

	idx := client.Index(indexUID)
	taskInfo, err := idx.UpdateSettings(&Settings{
		Embedders: map[string]Embedder{
			"default": {
				Source:     "userProvided",
				Dimensions: 3,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	settingsTask, err := idx.WaitForTask(taskInfo.TaskUID, 0)
	if err != nil {
		return nil, err
	}
	if settingsTask.Status != TaskStatusSucceeded {
		return nil, fmt.Errorf("Update settings task failed: %#+v", settingsTask)
	}

	documents := []map[string]interface{}{
		{"book_id": 123, "title": "Pride and Prejudice", "_vectors": map[string]interface{}{"default": []float64{0.1, 0.2, 0.3}}},
		{"book_id": 456, "title": "Le Petit Prince", "_vectors": map[string]interface{}{"default": []float64{2.4, 8.5, 1.6}}},
	}

	taskInfo, err = idx.AddDocuments(documents)
	if err != nil {
		return nil, err
	}

	finalTask, _ := idx.WaitForTask(taskInfo.TaskUID, 0)
	if finalTask.Status != TaskStatusSucceeded {
		return nil, fmt.Errorf("Add documents task failed: %#+v", finalTask)
	}

	return client.GetIndex(indexUID)
}

func setUpDistinctIndex(client ServiceManager, indexUID string) {
	idx := client.Index(indexUID)

	atters := []string{"product_id", "title", "sku", "url"}
	task, err := idx.UpdateFilterableAttributes(&atters)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	finalTask, _ := idx.WaitForTask(task.TaskUID, 0)
	if finalTask.Status != "succeeded" {
		os.Exit(1)
	}

	documents := []map[string]interface{}{
		{"product_id": 123, "title": "white shirt", "sku": "sku1234", "url": "https://example.com/products/p123"},
		{"product_id": 456, "title": "red shirt", "sku": "sku213", "url": "https://example.com/products/p456"},
		{"product_id": 1, "title": "green shirt", "sku": "sku876", "url": "https://example.com/products/p1"},
		{"product_id": 1344, "title": "blue shirt", "sku": "sku963", "url": "https://example.com/products/p1344"},
		{"product_id": 4, "title": "yellow shirt", "sku": "sku9064", "url": "https://example.com/products/p4"},
		{"product_id": 42, "title": "gray shirt", "sku": "sku964", "url": "https://example.com/products/p42"},
	}
	task, err = idx.AddDocuments(documents)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	finalTask, _ = idx.WaitForTask(task.TaskUID, 0)
	if finalTask.Status != "succeeded" {
		os.Exit(1)
	}
}

func testParseCsvDocuments(t *testing.T, documents io.Reader) []map[string]interface{} {
	var (
		docs   []map[string]interface{}
		header []string
	)
	r := csv.NewReader(documents)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		if header == nil {
			header = record
			continue
		}
		doc := make(map[string]interface{})
		for i, key := range header {
			doc[key] = record[i]
		}
		docs = append(docs, doc)
	}
	return docs
}

func testParseNdjsonDocuments(t *testing.T, documents io.Reader) []map[string]interface{} {
	var docs []map[string]interface{}
	scanner := bufio.NewScanner(documents)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		doc := make(map[string]interface{})
		err := json.Unmarshal([]byte(line), &doc)
		require.NoError(t, err)
		docs = append(docs, doc)
	}
	require.NoError(t, scanner.Err())
	return docs
}
