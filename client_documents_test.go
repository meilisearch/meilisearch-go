package meilisearch

import (
	"log"
	"testing"
)

var documents clientDocuments

type docTest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func TestClientDocuments_Get(t *testing.T) {
	updateIdRes, err := documents.AddOrUpdate([]interface{}{
		docTest{Id: "leslapins", Name: "nestle"},
	})

	if err != nil {
		t.Fatal(err)
	}

	AwaitAsyncUpdateId(documents, updateIdRes)

	var doc docTest
	err = documents.Get("leslapins", &doc)

	if err != nil {
		t.Fatal(err)
	}
}

func TestClientDocuments_Delete(t *testing.T) {
	updateIdRes, err := documents.AddOrUpdate([]interface{}{
		docTest{Id: "bloubiboulga2", Name: "nestle"},
	})

	if err != nil {
		t.Fatal(err)
	}

	AwaitAsyncUpdateId(documents, updateIdRes)

	updateIdRes, err = documents.Delete("bloubiboulga2")
	if err != nil {
		t.Fatal(err)
	}

	AwaitAsyncUpdateId(documents, updateIdRes)

	var doc docTest
	err = documents.Get("bloubiboulga2", &doc)

	if err.(*Error).ErrCode != ErrCodeResponseStatusCode {
		t.Fatal(err)
	}
}

func TestClientDocuments_Deletes(t *testing.T) {
	updateIdRes, err := documents.AddOrUpdate([]interface{}{
		docTest{Id: "bloubiboulga", Name: "nestle"},
		docTest{Id: "bloubiboulga1", Name: "nestle"},
	})

	if err != nil {
		t.Fatal(err)
	}

	AwaitAsyncUpdateId(documents, updateIdRes)

	updateIdRes, err = documents.Deletes([]string{"bloubiboulga", "bloubiboulga1"})
	if err != nil {
		t.Fatal(err)
	}

	AwaitAsyncUpdateId(documents, updateIdRes)

	var doc docTest
	err = documents.Get("bloubiboulga", &doc)

	if err.(*Error).ErrCode != ErrCodeResponseStatusCode {
		t.Fatal(err)
	}
}

func TestClientDocuments_List(t *testing.T) {
	updateIdRes, err := documents.AddOrUpdate([]interface{}{
		docTest{Id: "chocapic3", Name: "nestle"},
	})

	if err != nil {
		t.Fatal(err)
	}

	AwaitAsyncUpdateId(documents, updateIdRes)

	var list []docTest
	err = documents.List(ListDocumentsRequest{
		Offset: 0,
		Limit:  100,
	}, &list)

	if err != nil {
		t.Fatal(err)
	}

	// tests are running in parallel so there can be more than 1 docs
	if len(list) < 1 {
		t.Fatal("number of doc should be at least 1")
	}
}

func TestClientDocuments_AddOrUpdate(t *testing.T) {
	updateIdRes, err := documents.AddOrUpdate([]interface{}{
		docTest{Id: "chocapic", Name: "nestle"},
		docTest{Id: "chocapic2", Name: "nestle2"},
	})

	if err != nil {
		t.Fatal(err)
	}

	AwaitAsyncUpdateId(documents, updateIdRes)

	var list []docTest
	err = documents.List(ListDocumentsRequest{
		Offset: 0,
		Limit:  100,
	}, &list)

	if err != nil {
		t.Fatal(err)
	}

	// tests are running in parallel so there can be more than 2 docs
	if len(list) < 2 {
		t.Fatal("number of doc should at least 2")
	}
}

func TestClientDocuments_ClearAllDocuments(t *testing.T) {
	updateIdRes, err := documents.AddOrUpdate([]interface{}{
		docTest{Id: "chocapic", Name: "nestle"},
		docTest{Id: "chocapic2", Name: "nestle2"},
	})

	if err != nil {
		t.Fatal(err)
	}

	AwaitAsyncUpdateId(documents, updateIdRes)

	_, err = documents.DeleteAllDocuments()

	if err != nil {
		t.Fatal(err)
	}
}

func init() {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	resp, err := newClientIndexes(client).Create(CreateIndexRequest{
		Name: "documents_tests",
		UID:  "documents_tests",
		Schema: Schema{
			"id":   {"identifier", "indexed", "displayed"},
			"name": {"indexed", "displayed"},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	documents = newClientDocuments(client, resp.UID)
}
