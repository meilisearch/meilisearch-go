package meilisearch

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Sample struct for decoding
type TestDoc struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Dummy custom marshal/unmarshal for testing
func customMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func customUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func failingMarshal(v interface{}) ([]byte, error) {
	return nil, errors.New("marshal failed")
}

func failingUnmarshal(data []byte, v interface{}) error {
	return errors.New("unmarshal failed")
}

func TestHit_Decode_Success(t *testing.T) {
	hit := Hit{
		"id":    json.RawMessage(`1`),
		"title": json.RawMessage(`"Golang Rocks"`),
	}

	var doc TestDoc
	err := hit.Decode(&doc)

	require.NoError(t, err)
	assert.Equal(t, 1, doc.ID)
	assert.Equal(t, "Golang Rocks", doc.Title)
}

func TestHit_Decode_Error_Nil(t *testing.T) {
	hit := Hit{}
	err := hit.Decode(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vPtr must be a non-nil pointer")
}

func TestHit_DecodeWith_Success(t *testing.T) {
	hit := Hit{
		"id":    json.RawMessage(`2`),
		"title": json.RawMessage(`"Custom Decoder"`),
	}

	var doc TestDoc
	err := hit.DecodeWith(&doc, customMarshal, customUnmarshal)

	require.NoError(t, err)
	assert.Equal(t, 2, doc.ID)
	assert.Equal(t, "Custom Decoder", doc.Title)
}

func TestHit_DecodeWith_Error_NilPtr(t *testing.T) {
	hit := Hit{}
	err := hit.DecodeWith(nil, customMarshal, customUnmarshal)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vPtr must be a non-nil pointer")
}

func TestHit_DecodeWith_Error_MarshalFail(t *testing.T) {
	hit := Hit{}
	var doc TestDoc
	err := hit.DecodeWith(&doc, failingMarshal, customUnmarshal)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marshal failed")
}

func TestHit_DecodeWith_Error_UnmarshalFail(t *testing.T) {
	hit := Hit{
		"id": json.RawMessage(`"not-an-int"`), // invalid for int field
	}
	var doc TestDoc
	err := hit.DecodeWith(&doc, customMarshal, failingUnmarshal)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal failed")
}

func TestHits_Decode_Success(t *testing.T) {
	hits := Hits{
		{
			"id":    json.RawMessage(`1`),
			"title": json.RawMessage(`"First"`),
		},
		{
			"id":    json.RawMessage(`2`),
			"title": json.RawMessage(`"Second"`),
		},
	}

	var docs []TestDoc
	err := hits.Decode(&docs)
	require.NoError(t, err)
	assert.Len(t, docs, 2)
	assert.Equal(t, "First", docs[0].Title)
	assert.Equal(t, 2, docs[1].ID)
}

func TestHits_Decode_Error_NotPointer(t *testing.T) {
	hits := Hits{}
	var docs []TestDoc
	err := hits.Decode(docs) // pass by value
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "v must be a pointer to a slice")
}

func TestHits_DecodeWith_Success(t *testing.T) {
	hits := Hits{
		{
			"id":    json.RawMessage(`10`),
			"title": json.RawMessage(`"Hit 10"`),
		},
	}

	var docs []TestDoc
	err := hits.DecodeWith(&docs, customMarshal, customUnmarshal)
	require.NoError(t, err)
	assert.Equal(t, 10, docs[0].ID)
	assert.Equal(t, "Hit 10", docs[0].Title)
}

func TestHits_DecodeWith_Error_MarshalFail(t *testing.T) {
	hits := Hits{}
	var docs []TestDoc
	err := hits.DecodeWith(&docs, failingMarshal, customUnmarshal)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marshal failed")
}

func TestHits_DecodeWith_Error_UnmarshalFail(t *testing.T) {
	hits := Hits{
		{
			"id":    json.RawMessage(`"bad-int"`),
			"title": json.RawMessage(`"Bad"`),
		},
	}
	var docs []TestDoc
	err := hits.DecodeWith(&docs, customMarshal, failingUnmarshal)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal failed")
}

func TestHits_Len(t *testing.T) {
	hits := Hits{
		Hit{},
		Hit{},
	}
	assert.Equal(t, 2, hits.Len())
}
