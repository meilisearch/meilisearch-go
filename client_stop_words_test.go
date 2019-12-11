package meilisearch

import (
	"testing"
)

var stopWords clientStopWords

func TestClientStopWords_Add(t *testing.T) {
	resp, err := stopWords.Add([]string{"it", "is"})
	if err != nil {
		t.Fatal(err)
	}
	AwaitAsyncUpdateId(stopWords, resp)

	resps, err := stopWords.List()
	if err != nil {
		t.Fatal(err)
	}

	if len(resps) < 2 {
		t.Fatal("should be at leat 2")
	}
}

func TestClientStopWords_Deletes(t *testing.T) {
	resp, err := stopWords.Deletes([]string{"it", "is"})
	if err != nil {
		t.Fatal(err)
	}
	AwaitAsyncUpdateId(stopWords, resp)

	resps, err := stopWords.List()
	if err != nil {
		t.Fatal(err)
	}

	if len(resps) != 0 {
		t.Fatal("should be eq to 0")
	}
}

func TestClientStopWords_List(t *testing.T) {
	_, err := stopWords.List()
	if err != nil {
		t.Fatal(err)
	}
}

func init() {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	resp, err := newClientIndexes(client).Create(CreateIndexRequest{
		Name: "stop_words_tests",
		UID:  "stop_words_tests",
		Schema: Schema{
			"id":   {"identifier", "indexed", "displayed"},
			"name": {"indexed", "indexed", "displayed"},
		},
	})

	if err != nil {
		panic(err)
	}

	stopWords = newClientStopWords(client, resp.UID)
}
