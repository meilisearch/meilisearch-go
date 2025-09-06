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

type exampleEmbedded struct {
	E string `json:"e"`
}

type exampleBookForTest struct {
	ID              string            `json:"id"`
	Title           string            `json:"title"`
	Price           int               `json:"price"`
	OptS            *string           `json:"opt_s,omitempty"`
	Tags            []string          `json:"tags,omitempty"`
	Attrs           map[string]string `json:"attrs,omitempty"`
	exampleEmbedded                   // embedded, inlined (E is promoted)
	Renamed         string            `json:"renamed_field"`
	unexp           string            // unexported: must remain zero
}

func makeHitFromJSON(t *testing.T, s string) Hit {
	t.Helper()
	var m map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(s), &m), "bad test json")
	return Hit(m)
}

func makeHitFromStruct(t *testing.T, v any) Hit {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err, "marshal")
	var m map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(b, &m), "unmarshal to map")
	return Hit(m)
}

func TestHitDecodeInto_Basic(t *testing.T) {
	jsonStr := `{
		"id":"bk_1",
		"title":"Intro to Go",
		"price":42,
		"tags":["a","b"],
		"attrs":{"lang":"en"},
		"e":"emb",
		"renamed_field":"X",
		"unknown_ignored": true
	}`
	h := makeHitFromJSON(t, jsonStr)

	var got exampleBookForTest
	require.NoError(t, h.DecodeInto(&got))

	assert.Equal(t, "bk_1", got.ID)
	assert.Equal(t, "Intro to Go", got.Title)
	assert.Equal(t, 42, got.Price)
	assert.Equal(t, "emb", got.E, "embedded field should be promoted")
	assert.Equal(t, "X", got.Renamed)
	assert.Empty(t, got.unexp, "unexported field must remain zero")
	assert.Nil(t, got.OptS, "optional pointer should be nil when not present")
	assert.Equal(t, []string{"a", "b"}, got.Tags)
	assert.Equal(t, map[string]string{"lang": "en"}, got.Attrs)
}

func TestHitDecodeInto_NullAndMissing(t *testing.T) {
	jsonStr := `{
		"id":"bk_2",
		"title":null,
		"price":null,
		"opt_s":null,
		"e":null,
		"renamed_field":null
	}`
	h := makeHitFromJSON(t, jsonStr)

	var got exampleBookForTest
	require.NoError(t, h.DecodeInto(&got))

	// nulls should zero the fields
	assert.Equal(t, "", got.Title)
	assert.Equal(t, 0, got.Price)
	assert.Nil(t, got.OptS)
	assert.Equal(t, "", got.E)
	assert.Equal(t, "", got.Renamed)

	// Missing fields → zero value struct
	h2 := makeHitFromJSON(t, `{}`)
	var got2 exampleBookForTest
	require.NoError(t, h2.DecodeInto(&got2))
	assert.Equal(t, exampleBookForTest{}, got2)
}

func TestHitDecodeInto_Errors(t *testing.T) {
	h := makeHitFromJSON(t, `{"id":"x"}`)

	assert.Error(t, h.DecodeInto(nil))

	var notPtr exampleBookForTest
	assert.Error(t, h.DecodeInto(notPtr))

	var x int
	assert.Error(t, h.DecodeInto(&x))
}

func TestHitsDecodeInto_StructSlice(t *testing.T) {
	h := makeHitFromStruct(t, exampleBookForTest{
		ID:              "bk_3",
		Title:           "T",
		Price:           7,
		exampleEmbedded: exampleEmbedded{E: "e"},
		Renamed:         "R",
	})
	hs := Hits{h, h, h}

	var out []exampleBookForTest
	require.NoError(t, hs.DecodeInto(&out))
	require.Len(t, out, 3)

	for i, b := range out {
		assert.Equalf(t, "bk_3", b.ID, "idx=%d", i)
		assert.Equalf(t, "T", b.Title, "idx=%d", i)
		assert.Equalf(t, 7, b.Price, "idx=%d", i)
		assert.Equalf(t, "e", b.E, "idx=%d", i) // promoted embedded field
		assert.Equalf(t, "R", b.Renamed, "idx=%d", i)
	}
}

func TestHitsDecodeInto_PtrSlice(t *testing.T) {
	h := makeHitFromStruct(t, exampleBookForTest{ID: "bk_4", Title: "Ptr"})
	hs := Hits{h, h}

	var out []*exampleBookForTest
	require.NoError(t, hs.DecodeInto(&out))
	require.Len(t, out, 2)

	for i, p := range out {
		if assert.NotNilf(t, p, "idx=%d", i) {
			assert.Equalf(t, "bk_4", p.ID, "idx=%d", i)
			assert.Equalf(t, "Ptr", p.Title, "idx=%d", i)
		}
	}
}

func TestHitsDecodeInto_EmptyInput(t *testing.T) {
	var hs Hits
	var out []exampleBookForTest
	require.NoError(t, hs.DecodeInto(&out))
	assert.Len(t, out, 0)
}

func TestHitsDecodeInto_NullFields(t *testing.T) {
	jsonStr := `{
		"id":"bk_5",
		"title":"ok",
		"opt_s":null,
		"tags":null,
		"attrs":null,
		"e":"E"
	}`
	h := makeHitFromJSON(t, jsonStr)
	hs := Hits{h}

	var out []exampleBookForTest
	require.NoError(t, hs.DecodeInto(&out))
	require.Len(t, out, 1)

	got := out[0]
	assert.Equal(t, "bk_5", got.ID)
	assert.Equal(t, "ok", got.Title)
	assert.Equal(t, "E", got.E)
	assert.Nil(t, got.OptS)
	assert.Nil(t, got.Tags)
	assert.Nil(t, got.Attrs)
}

func TestHitsDecodeInto_Errors(t *testing.T) {
	h := makeHitFromJSON(t, `{"id":"bk_6"}`)
	hs := Hits{h}

	assert.Error(t, hs.DecodeInto(nil))

	var notPtr []exampleBookForTest
	assert.Error(t, hs.DecodeInto(notPtr))

	var x int
	assert.Error(t, hs.DecodeInto(&x))

	var bad1 []int
	assert.Error(t, hs.DecodeInto(&bad1))

	var bad2 []*int
	assert.Error(t, hs.DecodeInto(&bad2))
}

func TestHitDecodeInto_IgnoresUnknownFields(t *testing.T) {
	h := makeHitFromJSON(t, `{"id":"bk_7","unknown":123}`)
	var b exampleBookForTest
	require.NoError(t, h.DecodeInto(&b))
	assert.Equal(t, "bk_7", b.ID)
}

func TestHitDecodeInto_NestedStruct(t *testing.T) {
	type Child struct {
		X int `json:"x"`
	}
	type Parent struct {
		Child Child `json:"child"`
	}
	h := makeHitFromJSON(t, `{"child":{"x":7}}`)
	var p Parent
	require.NoError(t, h.DecodeInto(&p))
	assert.Equal(t, 7, p.Child.X)
}

func TestHitDecodeInto_NestedSlice(t *testing.T) {
	type Child struct {
		X int `json:"x"`
	}
	type Parent struct {
		Children []Child `json:"children"`
	}
	h := makeHitFromJSON(t, `{"children":[{"x":1},{"x":2}]}`)
	var p Parent
	require.NoError(t, h.DecodeInto(&p))
	require.Len(t, p.Children, 2)
	assert.Equal(t, 1, p.Children[0].X)
	assert.Equal(t, 2, p.Children[1].X)
}

func TestHitDecodeInto_MapOfStruct(t *testing.T) {
	type Child struct {
		X int `json:"x"`
	}
	type Parent struct {
		M map[string]Child `json:"m"`
	}
	h := makeHitFromJSON(t, `{"m":{"a":{"x":10},"b":{"x":20}}}`)
	var p Parent
	require.NoError(t, h.DecodeInto(&p))
	assert.Equal(t, 10, p.M["a"].X)
	assert.Equal(t, 20, p.M["b"].X)
}

func TestHitDecodeInto_PointerNested(t *testing.T) {
	type Child struct {
		X int `json:"x"`
	}
	type Parent struct {
		Child *Child `json:"child"`
	}
	h := makeHitFromJSON(t, `{"child":{"x":11}}`)
	var p Parent
	require.NoError(t, h.DecodeInto(&p))
	require.NotNil(t, p.Child)
	assert.Equal(t, 11, p.Child.X)
}

func TestHitDecodeInto_EmbeddedWithTag_AsNestedObject(t *testing.T) {
	type Embedded struct {
		E string `json:"e"`
	}
	type Parent struct {
		Embedded `json:"embedded"` // anonymous but tagged => NOT promoted; nested object
	}
	h := makeHitFromJSON(t, `{"embedded":{"e":"ok"}}`)
	var p Parent
	require.NoError(t, h.DecodeInto(&p))
	assert.Equal(t, "ok", p.Embedded.E)
}

func TestHitsDecodeInto_NestedSliceBatch(t *testing.T) {
	type Child struct {
		X int `json:"x"`
	}
	type Parent struct {
		Children []Child `json:"children"`
	}
	h := makeHitFromJSON(t, `{"children":[{"x":1},{"x":2},{"x":3}]}`)
	hs := Hits{h, h}
	var out []Parent
	require.NoError(t, hs.DecodeInto(&out))
	require.Len(t, out, 2)
	assert.Equal(t, 3, len(out[0].Children))
	assert.Equal(t, 2, out[1].Children[1].X)
}

type BookSmall struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Price int    `json:"price"`
}

type BookLarge struct {
	ID      string            `json:"id"`
	Title   string            `json:"title"`
	Author  string            `json:"author"`
	Price   float64           `json:"price"`
	InStock bool              `json:"in_stock"`
	Tags    []string          `json:"tags"`
	Ratings []float64         `json:"ratings"`
	Attrs   map[string]string `json:"attrs"`
	Nested  struct {
		ISBN      string   `json:"isbn"`
		PageCount int      `json:"page_count"`
		Editions  []string `json:"editions"`
	} `json:"nested"`
}

func makeHit(v any) Hit {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	return Hit(m)
}

func makeHitsSmall(n int) Hits {
	h := makeHit(BookSmall{ID: "bk_001", Title: "Intro to Go", Price: 42})
	hs := make(Hits, n)
	for i := range hs {
		hs[i] = h
	}
	return hs
}

func sampleLarge() BookLarge {
	return BookLarge{
		ID:      "bk_999",
		Title:   "The Complete Guide to High-Perf Go",
		Author:  "Gopher",
		Price:   129.99,
		InStock: true,
		Tags:    []string{"go", "performance", "concurrency", "json", "sdk"},
		Ratings: []float64{4.8, 4.9, 5.0, 4.7, 4.95, 4.85, 4.9},
		Attrs: map[string]string{
			"lang":   "en",
			"cover":  "hard",
			"series": "pro",
			"sku":    "GO-PERF-129",
		},
		Nested: struct {
			ISBN      string   `json:"isbn"`
			PageCount int      `json:"page_count"`
			Editions  []string `json:"editions"`
		}{
			ISBN:      "978-1-23456-789-7",
			PageCount: 864,
			Editions:  []string{"first", "second", "revised"},
		},
	}
}

func makeHitsLarge(n int) Hits {
	h := makeHit(sampleLarge())
	hs := make(Hits, n)
	for i := range hs {
		hs[i] = h
	}
	return hs
}

func BenchmarkHitsDecode_Small_1(b *testing.B)    { benchHitsDecodeSmall(b, 1, false) }
func BenchmarkHitsDecode_Small_100(b *testing.B)  { benchHitsDecodeSmall(b, 100, false) }
func BenchmarkHitsDecode_Small_1000(b *testing.B) { benchHitsDecodeSmall(b, 1000, false) }

func BenchmarkHitsDecodeInto_Small_1(b *testing.B)    { benchHitsDecodeSmall(b, 1, true) }
func BenchmarkHitsDecodeInto_Small_100(b *testing.B)  { benchHitsDecodeSmall(b, 100, true) }
func BenchmarkHitsDecodeInto_Small_1000(b *testing.B) { benchHitsDecodeSmall(b, 1000, true) }

func benchHitsDecodeSmall(b *testing.B, n int, fast bool) {
	hits := makeHitsSmall(n)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out []BookSmall
		if fast {
			if err := hits.DecodeInto(&out); err != nil {
				b.Fatal(err)
			}
		} else {
			if err := hits.Decode(&out); err != nil {
				b.Fatal(err)
			}
		}
		_ = out
	}
	// per-element metric (ns/op per hit)
	if n > 0 {
		b.SetBytes(int64(n)) // treat 1 "byte" == 1 element; useful to get ns/element in -benchmem output viewers
	}
}

func BenchmarkHitsDecode_Large_1(b *testing.B)    { benchHitsDecodeLarge(b, 1, false) }
func BenchmarkHitsDecode_Large_100(b *testing.B)  { benchHitsDecodeLarge(b, 100, false) }
func BenchmarkHitsDecode_Large_1000(b *testing.B) { benchHitsDecodeLarge(b, 1000, false) }

func BenchmarkHitsDecodeInto_Large_1(b *testing.B)    { benchHitsDecodeLarge(b, 1, true) }
func BenchmarkHitsDecodeInto_Large_100(b *testing.B)  { benchHitsDecodeLarge(b, 100, true) }
func BenchmarkHitsDecodeInto_Large_1000(b *testing.B) { benchHitsDecodeLarge(b, 1000, true) }

func benchHitsDecodeLarge(b *testing.B, n int, fast bool) {
	hits := makeHitsLarge(n)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out []BookLarge
		if fast {
			if err := hits.DecodeInto(&out); err != nil {
				b.Fatal(err)
			}
		} else {
			if err := hits.Decode(&out); err != nil {
				b.Fatal(err)
			}
		}
		_ = out
	}
	if n > 0 {
		b.SetBytes(int64(n)) // see note above
	}
}

func BenchmarkHitsDecodePtr_Small_1000(b *testing.B)     { benchHitsDecodePtrSmall(b, 1000, false) }
func BenchmarkHitsDecodeIntoPtr_Small_1000(b *testing.B) { benchHitsDecodePtrSmall(b, 1000, true) }

func benchHitsDecodePtrSmall(b *testing.B, n int, fast bool) {
	hits := makeHitsSmall(n)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out []*BookSmall
		if fast {
			if err := hits.DecodeInto(&out); err != nil {
				b.Fatal(err)
			}
		} else {
			if err := hits.Decode(&out); err != nil {
				b.Fatal(err)
			}
		}
		_ = out
	}
	if n > 0 {
		b.SetBytes(int64(n))
	}
}
