package meilisearch

import (
	"os"
	"testing"
)

type docTest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type docTestBooks struct {
	BookID int    `json:"book_id"`
	Title  string `json:"title"`
	Tag    string `json:"tag"`
}

func deleteAllIndexes(client *Client) (ok bool, err error) {
	list, err := client.Indexes().List()

	if err != nil {
		return false, err
	}

	for _, index := range list {
		client.Indexes().Delete(index.UID)
	}

	return true, nil
}

func TestMain(m *testing.M) {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})
	deleteAllIndexes(client)
	m.Run()
	deleteAllIndexes(client)
	os.Exit(0)
}

func Test_deleteAllIndexes(t *testing.T) {
	var indexUIDs = []string{
		"Test_deleteAllIndexes",
		"Test_deleteAllIndexes2",
		"Test_deleteAllIndexes3",
	}

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	for _, uid := range indexUIDs {
		_, err := client.Indexes().Create(CreateIndexRequest{
			UID: uid,
		})

		if err != nil {
			t.Fatal(err)
		}
	}

	deleteAllIndexes(client)

	for _, uid := range indexUIDs {
		resp, err := client.Indexes().Get(uid)
		if err == nil {
			t.Fatal(err)
		}
		if resp != nil {
			t.Fatal("deleteAllIndexes: One or more indexes were not deleted")
		}
	}
}
