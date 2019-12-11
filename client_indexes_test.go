package meilisearch

import (
	"testing"
	"time"
)

var indexes clientIndexes

func TestClientIndexes_Create(t *testing.T) {
	resp := createIndex(t, "TestClientIndexes_Create")

	if resp.Name != "TestClientIndexes_Create" {
		t.Fatal("response index does not have the same index")
	}
}

func TestClientIndexes_Get(t *testing.T) {
	resp := createIndex(t, "TestClientIndexes_Get")

	i, err := indexes.Get(resp.UID)
	if err != nil {
		t.Fatal(err)
	}

	if i.Name != resp.Name {
		t.Fatal("index name not eq", i.Name, resp.Name)
	}
}

func TestClientIndexes_Delete(t *testing.T) {
	resp := createIndex(t, "TestClientIndexes_Delete")
	ok, err := indexes.Delete(resp.UID)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("delete fail")
	}
}

func TestClientIndexes_List(t *testing.T) {
	createIndex(t, "TestClientIndexes_List")

	list, err := indexes.List()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) == 0 {
		t.Fatal("len of indexes should be at list 1, found ", len(list))
	}
}

func TestClientIndexes_Update(t *testing.T) {
	resp := createIndex(t, "TestClientIndexes_Update")

	update, err := indexes.Update(resp.UID, "TestClientIndexes_Update2")
	if err != nil {
		t.Fatal(err)
	}

	if update.Name != "TestClientIndexes_Update2" {
		t.Fatal("name of the index should be TestClientIndexes_Update2, found ", update.Name)
	}
}

func TestClientIndexes_GetSchema(t *testing.T) {
	resp := createIndex(t, "TestClientIndexes_GetSchema")

	time.Sleep(10 * time.Millisecond)
	_, err := indexes.GetSchema(resp.UID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientIndexes_GetSchemaRaw(t *testing.T) {
	resp := createIndex(t, "TestClientIndexes_GetSchemaRaw")

	time.Sleep(10 * time.Millisecond)
	_, err := indexes.GetRawSchema(resp.UID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientIndexes_UpdateSchema(t *testing.T) {
	resp := createIndex(t, "TestClientIndexes_GetSchemaRaw")

	time.Sleep(10 * time.Millisecond)
	_, err := indexes.UpdateSchema(resp.UID, Schema{
		"id":     resp.Schema["id"],
		"movies": []SchemaAttributes{SchemaAttributesDisplayed, SchemaAttributesIndexed},
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestClientIndexes_UpdateWithRawSchema(t *testing.T) {
	resp := createIndex(t, "TestClientIndexes_GetSchemaRaw")

	time.Sleep(10 * time.Millisecond)
	_, err := indexes.UpdateWithRawSchema(resp.UID, RawSchema{
		Identifier: "id",
		Attributes: map[string]RawAttribute{
			"id":    {Identifier: true, Indexed: true, Displayed: true},
			"title": {Indexed: true, Displayed: true},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func createIndex(t *testing.T, name string) *CreateIndexResponse {
	resp, err := indexes.Create(CreateIndexRequest{
		Name: name,
		Schema: Schema{
			"id": {"identifier", "indexed", "displayed"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	return resp
}

func init() {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	indexes = clientIndexes{client}
}
