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
	assert.Equal(t, "ok", p.E)
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

type embeddedA struct {
	EmbStr string `json:"emb_str"`
}

type embeddedPtr struct {
	EmbNum int `json:"emb_num,string"`
}

type embeddedRenamed struct {
	R string `json:"r"`
}

type helperOuter struct {
	// Basic fields + tag variants
	ID        string  `json:"id"`
	Count     int     `json:"count"`
	CountStr  int     `json:"count_str,string"`
	Price     float64 `json:"price"`
	PriceStr  float64 `json:"price_str,string"`
	Active    bool    `json:"active"`
	ActiveStr bool    `json:"active_str,string"`
	Note      string  `json:"note"`
	// Containers
	Tags  []string          `json:"tags"`
	Attrs map[string]string `json:"attrs"`
	// Anonymous embedded (promoted)
	embeddedA
	// Anonymous pointer-embedded (promoted)
	*embeddedPtr
	// Named field (NOT embedded) => nested under "boxed"
	Boxed embeddedA `json:"boxed"`
	// Anonymous but renamed (NOT promoted) => nested under "renamed"
	embeddedRenamed `json:"renamed"`
}

func TestGetTypeInfoAndCollectFields(t *testing.T) {
	rt := reflect.TypeOf(helperOuter{})

	// First call builds and caches
	ti1 := getTypeInfo(rt)
	require.NotNil(t, ti1)
	require.NotEmpty(t, ti1.fields)
	require.NotEmpty(t, ti1.byNameIndex)

	// Expected top-level keys
	expectNames := []string{
		"id", "count", "count_str", "price", "price_str", "active", "active_str",
		"note", "tags", "attrs",
		"emb_str", // promoted from embeddedA
		"emb_num", // promoted from *embeddedPtr
		"boxed",   // named field (nested)
		"renamed", // anonymous but tagged => nested under "renamed"
	}
	for _, name := range expectNames {
		_, ok := ti1.byNameIndex[name]
		assert.Truef(t, ok, "expected byNameIndex to contain %q", name)
	}

	// Ensure hasString is set for string-tagged fields
	var fmCountStr, fmPriceStr, fmActiveStr fieldMeta
	{
		idx := ti1.byNameIndex["count_str"]
		fmCountStr = ti1.fields[idx]
		idx = ti1.byNameIndex["price_str"]
		fmPriceStr = ti1.fields[idx]
		idx = ti1.byNameIndex["active_str"]
		fmActiveStr = ti1.fields[idx]
	}
	assert.True(t, fmCountStr.hasString, "count_str should have hasString")
	assert.True(t, fmPriceStr.hasString, "price_str should have hasString")
	assert.True(t, fmActiveStr.hasString, "active_str should have hasString")

	// Promoted embedded keys must exist
	_, ok := ti1.byNameIndex["emb_str"]
	assert.True(t, ok, "embedded field emb_str should be promoted")
	_, ok = ti1.byNameIndex["emb_num"]
	assert.True(t, ok, "pointer-embedded field emb_num should be promoted")

	// Named nested keys must exist
	_, ok = ti1.byNameIndex["boxed"]
	assert.True(t, ok, "boxed should be present as named field")
	_, ok = ti1.byNameIndex["renamed"]
	assert.True(t, ok, "renamed should be present as named field (anonymous + tag)")

	// Second call hits cache (pointer equality)
	ti2 := getTypeInfo(rt)
	require.NotNil(t, ti2)
	assert.Equal(t, ti1, ti2, "expected getTypeInfo to return cached pointer")
}

func TestIndexByte(t *testing.T) {
	assert.Equal(t, 0, indexByte("abc", 'a'))
	assert.Equal(t, 1, indexByte("abc", 'b'))
	assert.Equal(t, 2, indexByte("abc", 'c'))
	assert.Equal(t, -1, indexByte("abc", 'z'))
	assert.Equal(t, -1, indexByte("", 'x'))
}

func TestHasJSONTagOption(t *testing.T) {
	assert.True(t, hasJSONTagOption("omitempty,string", "string"))
	assert.True(t, hasJSONTagOption("string,foo,bar", "string"))
	assert.True(t, hasJSONTagOption("foo,string", "string"))
	assert.False(t, hasJSONTagOption("omitempty", "string"))
	assert.False(t, hasJSONTagOption("", "string"))
	assert.False(t, hasJSONTagOption("strings", "string")) // substring should not match
}

func TestFieldByIndexPathAlloc_Simple(t *testing.T) {
	type Leaf struct {
		X int
	}
	type Mid struct {
		Leaf Leaf
	}
	type Root struct {
		Mid Mid
	}

	var r Root
	rv := reflect.ValueOf(&r).Elem()
	// index path: Root.Mid.Leaf.X => [0,0,0]
	fv, ok := fieldByIndexPathAlloc(rv, []int{0, 0, 0})
	require.True(t, ok)
	require.True(t, fv.CanAddr())
	// Set X via reflect
	require.True(t, fv.CanSet())
	fv.SetInt(42)
	assert.Equal(t, 42, r.Mid.Leaf.X)
}

func TestFieldByIndexPathAlloc_AllocatesIntermediatePointers(t *testing.T) {
	type Leaf struct {
		S string
	}
	type Mid struct {
		L *Leaf
	}
	type Root struct {
		M *Mid
	}

	var r Root // r.M == nil; r.M.L == nil
	rv := reflect.ValueOf(&r).Elem()
	// Path: Root.M (idx 0) -> *Mid.L (idx 0) -> *Leaf.S (idx 0)
	fv, ok := fieldByIndexPathAlloc(rv, []int{0, 0, 0})
	require.True(t, ok, "should allocate intermediate *struct pointers")
	require.True(t, fv.CanSet())
	fv.SetString("hi")

	// Ensure allocations occurred
	require.NotNil(t, r.M)
	require.NotNil(t, r.M.L)
	assert.Equal(t, "hi", r.M.L.S)
}

func TestFieldByIndexPathAlloc_FailsOnNonStructChain(t *testing.T) {
	type Bad struct {
		N int
	}
	type Root struct {
		B *Bad
	}
	var r Root
	rv := reflect.ValueOf(&r).Elem()
	// Path: Root.B (pointer to Bad) -> field index 1 (invalid: Bad has only field index 0)
	_, ok := fieldByIndexPathAlloc(rv, []int{0, 1})
	assert.False(t, ok, "accessing invalid index should fail")
}

func TestFieldByIndexPathAlloc_LeafPtrAlloc(t *testing.T) {
	type Leaf struct {
		X int
	}
	type Root struct {
		P *Leaf
	}
	var r Root
	rv := reflect.ValueOf(&r).Elem()
	// Path to P: [0]
	fv, ok := fieldByIndexPathAlloc(rv, []int{0})
	require.True(t, ok)
	require.NotNil(t, fv)
	// Because leaf is *Leaf and nil, helper will allocate it (since it’s a pointer to struct)
	require.NotNil(t, r.P)
	// You can now set subfields through the pointer
	r.P.X = 7
	assert.Equal(t, 7, r.P.X)
}

func TestUnmarshalSingleField_StringOptionParity(t *testing.T) {
	type S struct {
		I  int     `json:"i,string"`
		F  float64 `json:"f,string"`
		B  bool    `json:"b,string"`
		T  string  `json:"t"`  // NOTE: no ,string here (it's invalid on string)
		I2 int     `json:"i2"` // no ,string
	}

	var s S
	require.NoError(t, unmarshalSingleField(&s, "i", []byte(`"123"`)))
	require.NoError(t, unmarshalSingleField(&s, "f", []byte(`"5.5"`)))
	require.NoError(t, unmarshalSingleField(&s, "b", []byte(`"true"`)))
	require.NoError(t, unmarshalSingleField(&s, "t", []byte(`"hello"`)))
	require.NoError(t, unmarshalSingleField(&s, "i2", []byte(`456`)))

	assert.Equal(t, 123, s.I)
	assert.InDelta(t, 5.5, s.F, 1e-9)
	assert.True(t, s.B)
	assert.Equal(t, "hello", s.T)
	assert.Equal(t, 456, s.I2)
}

func TestUnmarshalSingleField_Errors(t *testing.T) {
	type S struct {
		I int `json:"i,string"` // expects quoted number
	}
	var s S
	// invalid quoted int
	err := unmarshalSingleField(&s, "i", []byte(`"abc"`))
	assert.Error(t, err)
	// malformed JSON for the mini-object
	err = unmarshalSingleField(&s, "i", []byte(`"123"`)) // ok
	assert.NoError(t, err)
}

func TestIsJSONNull(t *testing.T) {
	assert.True(t, isJSONNull([]byte("null")))
	assert.False(t, isJSONNull([]byte("nul")))
	assert.False(t, isJSONNull([]byte("NULL")))
	assert.False(t, isJSONNull([]byte(`"null"`)))
	assert.False(t, isJSONNull(nil))
}

func rm(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return json.RawMessage(b)
}

type covStringOpt struct {
	N  int `json:"n"`
	SN int `json:"sn,string"` // triggers f.hasString branch
}

type covBadType struct {
	Count int `json:"count"` // we will feed an object here to force field-level unmarshal error
}

func TestHitDecodeInto_Struct_SkipEmptyRaw(t *testing.T) {
	h := Hit{
		"n":  json.RawMessage{}, // len(raw) == 0 => skip
		"sn": rm("123"),         // ,string branch should set 123
	}
	var out covStringOpt
	require.NoError(t, h.DecodeInto(&out))
	assert.Equal(t, 0, out.N)    // untouched (skipped)
	assert.Equal(t, 123, out.SN) // set via unmarshalSingleField
}

func TestHitDecodeInto_Struct_UnmarshalErrorOnField(t *testing.T) {
	h := Hit{
		"count": rm(map[string]int{"oops": 1}), // object into int => error
	}
	var out covBadType
	err := h.DecodeInto(&out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `decode field "count"`)
	assert.Contains(t, err.Error(), "cannot unmarshal object")
}

func TestHitDecodeInto_Map_StringKeyRequired(t *testing.T) {
	h := Hit{
		"a": rm(1),
	}
	var bad map[int]int
	err := h.DecodeInto(&bad)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "map key must be string")
}

func TestHitDecodeInto_Map_NullZeroAndValueErrors(t *testing.T) {
	// null → zero value of elem type.
	// also force a value-unmarshal error on key "err".
	h := Hit{
		"ok":  rm(7),
		"nil": json.RawMessage("null"),
		"err": rm(map[string]int{"x": 1}), // object into int => error
	}

	// First, map[string]any: null -> nil interface{}
	var m1 map[string]any
	require.NoError(t, h.DecodeInto(&m1))
	assert.Equal(t, float64(7), m1["ok"]) // numbers into interface{} become float64
	assert.Nil(t, m1["nil"])
	// We didn't touch "err" yet because we returned early above; to test error path,
	// try a typed map where we decode all keys and catch the error.

	// Now, map[string]int: null -> 0; "err" should fail
	var m2 map[string]int
	err := h.DecodeInto(&m2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `decode map value for key "err"`)
	assert.Contains(t, err.Error(), "cannot unmarshal object")
}

func TestHitsDecodeInto_SliceOfMap_InterfaceHappy(t *testing.T) {
	h1 := Hit{"a": rm(1), "b": rm("x")}
	h2 := Hit{"a": rm(2), "b": rm("y")}
	hs := Hits{h1, h2}

	var out []map[string]any
	require.NoError(t, hs.DecodeInto(&out))
	require.Len(t, out, 2)
	assert.Equal(t, float64(1), out[0]["a"])
	assert.Equal(t, "x", out[0]["b"])
	assert.Equal(t, float64(2), out[1]["a"])
	assert.Equal(t, "y", out[1]["b"])
}

func TestHitsDecodeInto_SliceOfPtrMap_InterfaceHappy(t *testing.T) {
	h1 := Hit{"a": rm(1)}
	h2 := Hit{"a": rm(2)}
	hs := Hits{h1, h2}

	var out []*map[string]any
	require.NoError(t, hs.DecodeInto(&out))
	require.Len(t, out, 2)
	require.NotNil(t, out[0])
	require.NotNil(t, out[1])
	assert.Equal(t, float64(1), (*out[0])["a"])
	assert.Equal(t, float64(2), (*out[1])["a"])
}

type covElem struct {
	V int `json:"v"`
}

func TestHitsDecodeInto_ErrorIndex_PropagatesForStruct(t *testing.T) {
	// First ok, second bad (object into int)
	ok := Hit{"v": rm(1)}
	bad := Hit{"v": rm(map[string]int{"x": 1})}
	hs := Hits{ok, bad}

	var out []covElem
	err := hs.DecodeInto(&out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode hits[1]") // index included
	assert.Contains(t, err.Error(), "decode field \"v\"")
}

func TestHitsDecodeInto_ErrorIndex_PropagatesForMap(t *testing.T) {
	// First ok, second bad for typed map[string]int
	ok := Hit{"k": rm(5)}
	bad := Hit{"k": rm(map[string]int{"x": 1})}
	hs := Hits{ok, bad}

	var out []map[string]int
	err := hs.DecodeInto(&out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode hits[1]")
	assert.Contains(t, err.Error(), `decode map value for key "k"`)
}

func TestHitsDecodeInto_SliceElemPtrMap_NonStringKeyError(t *testing.T) {
	hs := Hits{Hit{"k": rm(1)}}
	var out []*map[int]int
	err := hs.DecodeInto(&out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "slice element must be map with string key")
}

func TestHitsDecodeInto_SliceElemMap_NonStringKeyError(t *testing.T) {
	hs := Hits{Hit{"k": rm(1)}}
	var out []map[int]int
	err := hs.DecodeInto(&out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "slice element must be map with string key")
}

func TestHitsDecodeInto_SliceElemPtrStruct_Happy(t *testing.T) {
	hs := Hits{
		{"v": rm(10)},
		{"v": rm(20)},
	}
	var out []*covElem
	require.NoError(t, hs.DecodeInto(&out))
	require.Len(t, out, 2)
	assert.Equal(t, 10, out[0].V)
	assert.Equal(t, 20, out[1].V)
}

func TestHitsDecodeInto_SliceElemStruct_Happy(t *testing.T) {
	hs := Hits{
		{"v": rm(3)},
		{"v": rm(4)},
	}
	var out []covElem
	require.NoError(t, hs.DecodeInto(&out))
	require.Len(t, out, 2)
	assert.Equal(t, 3, out[0].V)
	assert.Equal(t, 4, out[1].V)
}

func TestHitDecodeInto_Struct_StringTagBranch(t *testing.T) {
	h := Hit{"sn": rm("42")}
	var out covStringOpt
	require.NoError(t, h.DecodeInto(&out))
	assert.Equal(t, 42, out.SN)
}

func TestHitDecodeInto_Map_NullToZero_TypedInt(t *testing.T) {
	h := Hit{
		"x": rm(9),
		"z": json.RawMessage("null"),
	}
	var m map[string]int
	require.NoError(t, h.DecodeInto(&m))
	assert.Equal(t, 9, m["x"])
	assert.Equal(t, 0, m["z"]) // zero value for int
}

func TestHitDecodeInto_DispatchError(t *testing.T) {
	var s string
	err := Hit{"x": rm(1)}.DecodeInto(&s)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "struct or map")
}

func TestHitDecodeInto_EmptyRawUnknownKey_NoCrash(t *testing.T) {
	type S struct {
		A int `json:"a"`
	}
	h := Hit{
		"zzz": json.RawMessage{}, // unknown + empty
		"a":   rm(5),
	}
	var s S
	require.NoError(t, h.DecodeInto(&s))
	assert.Equal(t, 5, s.A)
}

func TestHitDecodeInto_Map_InterfaceTypes(t *testing.T) {
	h := Hit{
		"n": rm(1),
		"s": rm("str"),
		"o": rm(map[string]any{"k": 2}),
		"a": rm([]any{1, "x"}),
	}
	var m map[string]any
	require.NoError(t, h.DecodeInto(&m))
	assert.Equal(t, float64(1), m["n"])
	assert.Equal(t, "str", m["s"])
	assert.Equal(t, map[string]any{"k": float64(2)}, m["o"])
	assert.Equal(t, []any{float64(1), "x"}, m["a"])

	// Double-check types match what encoding/json would produce
	b, _ := json.Marshal(h) // {"n":1,"s":"str","o":{"k":2},"a":[1,"x"]}
	var std map[string]any
	require.NoError(t, json.Unmarshal(b, &std))
	assert.Equal(t, std, m)
}

func TestHitDecodeInto_Map_PreserveAndOverwrite(t *testing.T) {
	h := Hit{
		"a": rm(10),
		"b": rm(20),
	}
	m := map[string]int{"a": 1} // pre-populated
	require.NoError(t, h.DecodeInto(&m))
	// "a" should be overwritten, "b" added
	assert.Equal(t, 10, m["a"])
	assert.Equal(t, 20, m["b"])
}

func TestHitsDecodeInto_SliceElemInvalidKind(t *testing.T) {
	hs := Hits{{"x": rm(1)}}
	var out []int
	err := hs.DecodeInto(&out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "slice element must be struct, *struct, map[string]T, or *map[string]T")
}

func TestHitsDecodeInto_NonSlicePointer(t *testing.T) {
	hs := Hits{}
	var x int
	err := hs.DecodeInto(&x)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must point to a slice")
}

func TestHitsDecodeInto_BadPtr(t *testing.T) {
	var hs Hits
	err := hs.DecodeInto(nil)
	require.Error(t, err)

	var notPtr []map[string]any
	err = hs.DecodeInto(notPtr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "non-nil pointer")
}

func TestHitDecodeInto_BadKinds(t *testing.T) {
	h := Hit{}
	var ch chan int
	err := h.DecodeInto(&ch)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "struct or map")

	// Not pointer
	var s covStringOpt
	err = h.DecodeInto(s)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "non-nil pointer")
}

func TestHitDecodeInto_Map_CreatesMapIfNil(t *testing.T) {
	h := Hit{"x": rm(1)}
	var m map[string]any // nil
	require.NoError(t, h.DecodeInto(&m))
	require.NotNil(t, m)
	assert.Equal(t, float64(1), m["x"])
}

func TestHitsDecodeInto_PtrMap_AllocAndFill(t *testing.T) {
	hs := Hits{
		{"x": rm(1)},
		{"y": rm(2)},
	}
	var out []*map[string]any
	require.NoError(t, hs.DecodeInto(&out))
	require.Len(t, out, 2)
	assert.Equal(t, float64(1), (*out[0])["x"])
	assert.Equal(t, float64(2), (*out[1])["y"])
}

func TestHitDecodeInto_Struct_PointerEmbeddedAlloc(t *testing.T) {
	type Inner struct {
		Z int `json:"z"`
	}
	type Wrap struct{ *Inner }
	h := Hit{"z": rm(9)}
	var w Wrap
	require.NoError(t, h.DecodeInto(&w))
	require.NotNil(t, w.Inner)
	assert.Equal(t, 9, w.Z)
}

func TestHitDecodeInto_Struct_UnknownKeysIgnored(t *testing.T) {
	type S struct {
		A int `json:"a"`
	}
	h := Hit{"xxx": rm(1), "a": rm(2)}
	var s S
	require.NoError(t, h.DecodeInto(&s))
	assert.Equal(t, 2, s.A)
}
