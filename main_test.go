package meilisearch

import (
	"crypto/tls"
	"fmt"
	"os"
	"testing"

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

func deleteAllIndexes(client ClientInterface) (ok bool, err error) {
	list, err := client.GetAllIndexes()
	if err != nil {
		return false, err
	}

	for _, index := range list {
		_, _ = client.DeleteIndex(index.UID)
	}

	return true, nil
}

func SetUpBasicIndex() {
	client := NewClient(ClientConfig{
		Host:   "http://localhost:7700",
		APIKey: masterKey,
	})
	index := client.Index("indexUID")

	documents := []map[string]interface{}{
		{"book_id": 123, "title": "Pride and Prejudice"},
		{"book_id": 456, "title": "Le Petit Prince"},
		{"book_id": 1, "title": "Alice In Wonderland"},
		{"book_id": 1344, "title": "The Hobbit"},
		{"book_id": 4, "title": "Harry Potter and the Half-Blood Prince"},
		{"book_id": 42, "title": "The Hitchhiker's Guide to the Galaxy"},
	}
	update, err := index.AddDocuments(documents)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	finalUpdateStatus, _ := index.DefaultWaitForPendingUpdate(update)
	if finalUpdateStatus != "processed" {
		os.Exit(1)
	}
}

func SetUpIndexForFaceting() {
	client := NewClient(ClientConfig{
		Host:   "http://localhost:7700",
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
	update, err := index.AddDocuments(booksTest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	finalUpdateStatus, _ := index.DefaultWaitForPendingUpdate(update)
	if finalUpdateStatus != "processed" {
		os.Exit(1)
	}
}

var masterKey = "masterKey"
var primaryKey = "primaryKey"
var defaultClient = NewClient(ClientConfig{
	Host:   "http://localhost:7700",
	APIKey: masterKey,
})

var customClient = NewFastHTTPCustomClient(ClientConfig{
	Host:   "http://localhost:7700",
	APIKey: masterKey,
},
	&fasthttp.Client{
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Name:      "custom-client",
	})

var timeoutClient = NewClient(ClientConfig{
	Host:    "http://localhost:7700",
	APIKey:  masterKey,
	Timeout: 1,
})

func TestMain(m *testing.M) {
	_, _ = deleteAllIndexes(defaultClient)
	code := m.Run()
	_, _ = deleteAllIndexes(defaultClient)
	os.Exit(code)
}

func Test_deleteAllIndexes(t *testing.T) {
	var indexUIDs = []string{
		"Test_deleteAllIndexes",
		"Test_deleteAllIndexes2",
		"Test_deleteAllIndexes3",
	}
	_, _ = deleteAllIndexes(defaultClient)

	for _, uid := range indexUIDs {
		_, err := defaultClient.CreateIndex(&IndexConfig{
			Uid: uid,
		})

		if err != nil {
			t.Fatal(err)
		}
	}

	_, _ = deleteAllIndexes(defaultClient)

	for _, uid := range indexUIDs {
		resp, err := defaultClient.GetIndex(uid)
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Fatal("deleteAllIndexes: One or more indexes were not deleted")
		}
	}
}
