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
//	sv, err := New("http://localhost:7700", WithAPIKey("foobar"))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println(sv.IsHealthy(context.Background()))
package meilisearch
