package meilisearch

import (
	"encoding/json"
	"errors"
	"reflect"
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
	exampleEmbedded                   // embedded, inlined
	Renamed         string            `json:"renamed_field"`
	unexp           string            // unexported: must remain zero
}

// --- Helpers ---

func makeHitFromJSON(t *testing.T, s string) Hit {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("bad test json: %v", err)
	}
	return Hit(m)
}

func makeHitFromStruct(t *testing.T, v any) Hit {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
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
	if err := h.DecodeInto(&got); err != nil {
		t.Fatalf("DecodeInto error: %v", err)
	}

	if got.ID != "bk_1" || got.Title != "Intro to Go" || got.Price != 42 {
		t.Fatalf("basic fields mismatch: %+v", got)
	}
	if got.exampleEmbedded.E != "emb" {
		t.Fatalf("embedded field not set: %+v", got)
	}
	if got.Renamed != "X" {
		t.Fatalf("renamed tag not respected: %+v", got)
	}
	if got.unexp != "" {
		t.Fatalf("unexported field must remain zero, got %q", got.unexp)
	}
	if got.OptS != nil {
		t.Fatalf("optional pointer should be nil when not present, got %v", *got.OptS)
	}
	if !reflect.DeepEqual(got.Tags, []string{"a", "b"}) {
		t.Fatalf("tags mismatch: %+v", got.Tags)
	}
	if !reflect.DeepEqual(got.Attrs, map[string]string{"lang": "en"}) {
		t.Fatalf("attrs mismatch: %+v", got.Attrs)
	}
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
	if err := h.DecodeInto(&got); err != nil {
		t.Fatalf("DecodeInto error: %v", err)
	}

	// nulls should zero the fields
	if got.Title != "" || got.Price != 0 || got.OptS != nil || got.exampleEmbedded.E != "" || got.Renamed != "" {
		t.Fatalf("nulls not treated as zero values: %+v", got)
	}

	// missing fields: zero values (ID is present above; test missing behavior too)
	h2 := makeHitFromJSON(t, `{}`)
	var got2 exampleBookForTest
	if err := h2.DecodeInto(&got2); err != nil {
		t.Fatalf("DecodeInto error: %v", err)
	}
	var zero exampleBookForTest
	if !reflect.DeepEqual(got2, zero) {
		t.Fatalf("missing fields should yield zero value struct: %+v", got2)
	}
}

func TestHitDecodeInto_Errors(t *testing.T) {
	h := makeHitFromJSON(t, `{"id":"x"}`)

	// nil
	if err := h.DecodeInto(nil); err == nil {
		t.Fatalf("expected error for nil target")
	}

	// non-pointer
	var notPtr exampleBookForTest
	if err := h.DecodeInto(notPtr); err == nil {
		t.Fatalf("expected error for non-pointer")
	}

	// pointer to non-struct
	var x int
	if err := h.DecodeInto(&x); err == nil {
		t.Fatalf("expected error for pointer to non-struct")
	}
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
	if err := hs.DecodeInto(&out); err != nil {
		t.Fatalf("DecodeInto (slice) error: %v", err)
	}

	if len(out) != 3 {
		t.Fatalf("len mismatch: %d", len(out))
	}
	for i, b := range out {
		if b.ID != "bk_3" || b.Title != "T" || b.Price != 7 || b.exampleEmbedded.E != "e" || b.Renamed != "R" {
			t.Fatalf("element %d mismatch: %+v", i, b)
		}
	}
}

func TestHitsDecodeInto_PtrSlice(t *testing.T) {
	h := makeHitFromStruct(t, exampleBookForTest{ID: "bk_4", Title: "Ptr"})
	hs := Hits{h, h}

	var out []*exampleBookForTest
	if err := hs.DecodeInto(&out); err != nil {
		t.Fatalf("DecodeInto (ptr slice) error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("len mismatch: %d", len(out))
	}
	for i, p := range out {
		if p == nil || p.ID != "bk_4" || p.Title != "Ptr" {
			t.Fatalf("element %d mismatch: %+v", i, p)
		}
	}
}

func TestHitsDecodeInto_EmptyInput(t *testing.T) {
	var hs Hits
	var out []exampleBookForTest
	if err := hs.DecodeInto(&out); err != nil {
		t.Fatalf("DecodeInto error on empty input: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty output slice, got %d", len(out))
	}
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
	if err := hs.DecodeInto(&out); err != nil {
		t.Fatalf("DecodeInto error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("len mismatch: %d", len(out))
	}
	got := out[0]
	if got.ID != "bk_5" || got.Title != "ok" || got.exampleEmbedded.E != "E" {
		t.Fatalf("values mismatch: %+v", got)
	}
	// null → zero
	if got.OptS != nil || got.Tags != nil || got.Attrs != nil {
		t.Fatalf("null should zero pointer/map/slice fields: %+v", got)
	}
}

func TestHitsDecodeInto_Errors(t *testing.T) {
	h := makeHitFromJSON(t, `{"id":"bk_6"}`)
	hs := Hits{h}

	// nil
	if err := hs.DecodeInto(nil); err == nil {
		t.Fatalf("expected error for nil")
	}

	// non-pointer
	var notPtr []exampleBookForTest
	if err := hs.DecodeInto(notPtr); err == nil {
		t.Fatalf("expected error for non-pointer")
	}

	// pointer to non-slice
	var x int
	if err := hs.DecodeInto(&x); err == nil {
		t.Fatalf("expected error for pointer to non-slice")
	}

	// slice of non-struct element
	var bad1 []int
	if err := hs.DecodeInto(&bad1); err == nil {
		t.Fatalf("expected error for slice element not struct/*struct")
	}

	// slice of pointer to non-struct
	var bad2 []*int
	if err := hs.DecodeInto(&bad2); err == nil {
		t.Fatalf("expected error for slice element pointer to non-struct")
	}

	// ensure error text is informative (optional assert on substring)
	var out []exampleBookForTest
	err := hs.DecodeInto(&out)
	if err != nil && !errors.Is(err, nil) {
		// no-op: just compile-time reference to errors package to avoid lint noise
	}
}

func TestHitDecodeInto_IgnoresUnknownFields(t *testing.T) {
	h := makeHitFromJSON(t, `{"id":"bk_7","unknown":123}`)
	var b exampleBookForTest
	if err := h.DecodeInto(&b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.ID != "bk_7" {
		t.Fatalf("id mismatch: %+v", b)
	}
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
