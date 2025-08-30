// Package meilisearch is the official Meilisearch SDK for the Go programming language.
//
// The meilisearch-go SDK for Go provides APIs and utilities that developers can use to
// build Go applications that use meilisearch service.
//
// See the meilisearch package documentation for more information.
// https://www.meilisearch.com/docs/reference
//
//	Example:
//
//	meili := New("http://localhost:7700", WithAPIKey("foobar"))
//
//	idx := meili.Index("movies")
//
//	documents := []map[string]interface{}{
//		{"id": 1, "title": "Carol", "genres": []string{"Romance", "Drama"}},
//		{"id": 2, "title": "Wonder Woman", "genres": []string{"Action", "Adventure"}},
//		{"id": 3, "title": "Life of Pi", "genres": []string{"Adventure", "Drama"}},
//		{"id": 4, "title": "Mad Max: Fury Road", "genres": []string{"Adventure", "Science Fiction"}},
//		{"id": 5, "title": "Moana", "genres": []string{"Fantasy", "Action"}},
//		{"id": 6, "title": "Philadelphia", "genres": []string{"Drama"}},
//	}
//	task, err := idx.AddDocuments(documents, nil)
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	fmt.Println(task.TaskUID)
package meilisearch
