package meilisearch

import (
	"encoding/json"
	"testing"
)

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
