package meilisearch

import (
	"testing"
)

func TestClientIndexes_Create(t *testing.T) {
	var indexUID = "TestClientIndexes_Create"

	resp, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	if resp.UID != "TestClientIndexes_Create" {
		t.Fatal("response index does not have the same index")
	}
}

func TestClientIndexes_Get(t *testing.T) {
	var indexUID = "TestClientIndexes_Get"

	resp, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	i, err := client.Indexes().Get(resp.UID)
	if err != nil {
		t.Fatal(err)
	}

	if i.Name != resp.Name {
		t.Fatal("index name not eq", i.Name, resp.Name)
	}
}

func TestClientIndexes_Delete(t *testing.T) {
	var indexUID = "TestClientIndexes_Delete"

	resp, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	ok, err := client.Indexes().Delete(resp.UID)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("delete fail")
	}
}

func TestClientIndexes_List(t *testing.T) {
	var indexUID = "TestClientIndexes_List"

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	list, err := client.Indexes().List()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 {
		t.Fatal("len of indexes should be at list 1, found ", len(list))
	}
}

func TestClientIndexes_UpdateName(t *testing.T) {
	var indexUID = "TestClientIndexes_UpdateName"

	resp, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	update, err := client.Indexes().UpdateName(resp.UID, "TestClientIndexes_Update2")
	if err != nil {
		t.Fatal(err)
	}

	if update.Name != "TestClientIndexes_Update2" {
		t.Fatal("name of the index should be TestClientIndexes_Update2, found ", update.Name)
	}
}

func TestClientIndexes_UpdatePrimaryKey(t *testing.T) {
	var indexUID = "TestClientIndexes_UpdatePrimaryKey"

	resp, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	update, err := client.Indexes().UpdatePrimaryKey(resp.UID, "identifier")
	if err != nil {
		t.Fatal(err)
	}

	if update.PrimaryKey != "identifier" {
		t.Fatal("name of the index should be TestClientIndexes_Update2, found ", update.Name)
	}
}
