package meilisearch

import (
	"testing"
	"time"
)

var keys clientKeys

func TestClientKeys_Create(t *testing.T) {
	_, err := keys.Create(CreateApiKeyRequest{
		Description: "This is a key",
		Acl: []ACL{
			AclDocumentsRead,
			AclIndexesWrite,
		},
		Indexes:  []string{"*"},
		ExpireAt: time.Now().Add(time.Hour * 10).Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientKeys_List(t *testing.T) {
	if _, err := keys.List(); err != nil {
		t.Fatal(err)
	}
}

func TestClientKeys_Get(t *testing.T) {
	respKey, err := keys.Create(CreateApiKeyRequest{
		Description: "This is a key",
		Acl: []ACL{
			AclDocumentsRead,
			AclIndexesWrite,
		},
		Indexes:  []string{"*"},
		ExpireAt: 1574332928,
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := keys.Get(respKey.Key); err != nil {
		t.Fatal(err)
	}
}

func TestClientKeys_Delete(t *testing.T) {
	respKey, err := keys.Create(CreateApiKeyRequest{
		Description: "This is a key",
		Acl: []ACL{
			AclDocumentsRead,
			AclIndexesWrite,
		},
		Indexes:  []string{"*"},
		ExpireAt: 1574332928,
	})
	if err != nil {
		t.Fatal(err)
	}

	deleted, err := keys.Delete(respKey.Key)
	if err != nil || !deleted {
		t.Fatal(err, deleted)
	}

	if _, err := keys.Get(respKey.Key); err == nil {
		t.Fatal(err)
	}
}

func TestClientKeys_Update(t *testing.T) {
	respKey, err := keys.Create(CreateApiKeyRequest{
		Description: "This is a key",
		Acl: []ACL{
			AclDocumentsRead,
			AclIndexesWrite,
		},
		Indexes:  []string{"*"},
		ExpireAt: 1574332928,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = keys.Update(respKey.Key, UpdateApiKeyRequest{
		Description: "new description",
		Acl:         respKey.Acl,
		Indexes:     respKey.Indexes,
		Revoked:     false,
	})
	if err != nil {
		t.Fatal(err)
	}

	k, err := keys.Get(respKey.Key)
	if err != nil {
		t.Fatal(err)
	}
	if k.Description != "new description" {
		t.Fatal("description not updated")
	}
}

func init() {
	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	keys = clientKeys{client}
}
