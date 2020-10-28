package meilisearch

import (
	"testing"
)

func TestClientDocuments_Get(t *testing.T) {
	var indexUID = "TestClientDocuments_Get"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]interface{}{
			docTest{ID: "123", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	var doc docTest
	if err = client.
		Documents(indexUID).
		Get("123", &doc); err != nil {
		t.Fatal(err)
	}

	expect := docTest{ID: "123", Name: "nestle"}
	if doc != expect {
		t.Errorf("%v != %v", doc, expect)
	}
}

func TestClientDocuments_Delete(t *testing.T) {
	var indexUID = "TestClientDocuments_Delete"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]interface{}{
			docTest{ID: "123", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	updateIDRes, err = client.Documents(indexUID).Delete("123")
	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	var doc docTest
	err = client.Documents(indexUID).Get("123", &doc)

	if err.(*Error).ErrCode != ErrCodeResponseStatusCode {
		t.Fatal(err)
	}
}

func TestClientDocuments_Deletes(t *testing.T) {
	var indexUID = "deletes"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]interface{}{
			docTest{ID: "123", Name: "nestle"},
			docTest{ID: "456", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	updateIDRes, err = client.Documents(indexUID).Deletes([]string{"123", "456"})

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	var doc docTest
	err = client.Documents(indexUID).Get("123", &doc)

	if err.(*Error).ErrCode != ErrCodeResponseStatusCode {
		t.Fatal(err)
	}
}

func TestClientDocuments_List(t *testing.T) {
	var indexUID = "TestClientDocuments_List"

	if _, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	}); err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]interface{}{
			docTest{ID: "123", Name: "nestle"},
			docTest{ID: "456", Name: "hershey"},
		})

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	var list []docTest
	err = client.Documents(indexUID).List(ListDocumentsRequest{
		Offset: 1,
		Limit:  1,
	}, &list)

	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 || list[0].ID != "456" {
		t.Fatal("expected to return the document[1]")
	}
}

func TestClientDocuments_AddOrReplace(t *testing.T) {
	var indexUID = "TestClientDocuments_AddOrReplace"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrReplace([]docTest{
			{ID: "123", Name: "nestle"},
			{ID: "456", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	var list []docTest
	err = client.Documents(indexUID).List(ListDocumentsRequest{
		Offset: 0,
		Limit:  100,
	}, &list)

	if err != nil {
		t.Fatal(err)
	}

	// tests are running in parallel so there can be more than 1 docs
	if len(list) < 2 {
		t.Fatal("number of doc should be at least 1")
	}
}

func TestClientDocuments_AddOrUpdate(t *testing.T) {
	var indexUID = "TestClientDocuments_AddOrUpdate"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]docTest{
			{ID: "123", Name: "nestle"},
			{ID: "456", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	var list []docTest
	err = client.Documents(indexUID).List(ListDocumentsRequest{
		Offset: 0,
		Limit:  100,
	}, &list)

	if err != nil {
		t.Fatal(err)
	}

	// tests are running in parallel so there can be more than 1 docs
	if len(list) < 2 {
		t.Fatal("number of doc should be at least 1")
	}
}

func TestClientDocuments_DeleteAllDocuments(t *testing.T) {
	var indexUID = "TestClientDocuments_DeleteAllDocuments"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]interface{}{
			docTest{ID: "123", Name: "nestle"},
			docTest{ID: "456", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	_, _ = client.DefaultWaitForPendingUpdate(indexUID, updateIDRes)

	_, err = client.Documents(indexUID).DeleteAllDocuments()

	if err != nil {
		t.Fatal(err)
	}
}
