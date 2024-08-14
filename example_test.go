package meilisearch

import (
	"fmt"
	"os"
)

func ExampleNew() {
	// WithAPIKey is optional
	meili := New("http://localhost:7700", WithAPIKey("foobar"))

	// An index is where the documents are stored.
	idx := meili.Index("movies")

	// If the index 'movies' does not exist, Meilisearch creates it when you first add the documents.
	documents := []map[string]interface{}{
		{"id": 1, "title": "Carol", "genres": []string{"Romance", "Drama"}},
		{"id": 2, "title": "Wonder Woman", "genres": []string{"Action", "Adventure"}},
		{"id": 3, "title": "Life of Pi", "genres": []string{"Adventure", "Drama"}},
		{"id": 4, "title": "Mad Max: Fury Road", "genres": []string{"Adventure", "Science Fiction"}},
		{"id": 5, "title": "Moana", "genres": []string{"Fantasy", "Action"}},
		{"id": 6, "title": "Philadelphia", "genres": []string{"Drama"}},
	}
	task, err := idx.AddDocuments(documents)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(task.TaskUID)
}

func ExampleConnect() {
	meili, err := Connect("http://localhost:7700", WithAPIKey("foobar"))
	if err != nil {
		fmt.Println(err)
		return
	}

	ver, err := meili.Version()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ver.PkgVersion)
}
