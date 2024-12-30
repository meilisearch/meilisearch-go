package meilisearch

import (
	"bytes"
	"encoding/json"
	"testing"
)

func BenchmarkGzipEncoder(b *testing.B) {
	encoder := newEncoding(GzipEncoding, DefaultCompression)
	data := bytes.NewReader(make([]byte, 1024*1024)) // 1 MB of data
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf, err := encoder.Encode(data)
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
		_ = buf
	}
}

func BenchmarkDeflateEncoder(b *testing.B) {
	encoder := newEncoding(DeflateEncoding, DefaultCompression)
	data := bytes.NewReader(make([]byte, 1024*1024)) // 1 MB of data
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf, err := encoder.Encode(data)
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
		_ = buf
	}
}

func BenchmarkBrotliEncoder(b *testing.B) {
	encoder := newEncoding(BrotliEncoding, DefaultCompression)
	data := bytes.NewReader(make([]byte, 1024*1024)) // 1 MB of data
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf, err := encoder.Encode(data)
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
		_ = buf
	}
}

func BenchmarkGzipDecoder(b *testing.B) {
	encoder := newEncoding(GzipEncoding, DefaultCompression)

	// Prepare a valid JSON input
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 12345,
		"key3": []string{"item1", "item2", "item3"},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		b.Fatalf("JSON marshal failed: %v", err)
	}

	// Encode the valid JSON data
	input := bytes.NewReader(jsonData)
	encoded, err := encoder.Encode(input)
	if err != nil {
		b.Fatalf("Encode failed: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := encoder.Decode(encoded.Bytes(), &result); err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkFlateDecoder(b *testing.B) {
	encoder := newEncoding(DeflateEncoding, DefaultCompression)

	// Prepare valid JSON input
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 12345,
		"key3": []string{"item1", "item2", "item3"},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		b.Fatalf("JSON marshal failed: %v", err)
	}

	// Encode the valid JSON data
	input := bytes.NewReader(jsonData)
	encoded, err := encoder.Encode(input)
	if err != nil {
		b.Fatalf("Encode failed: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := encoder.Decode(encoded.Bytes(), &result); err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkBrotliDecoder(b *testing.B) {
	encoder := newEncoding(BrotliEncoding, DefaultCompression)

	// Prepare valid JSON input
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 12345,
		"key3": []string{"item1", "item2", "item3"},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		b.Fatalf("JSON marshal failed: %v", err)
	}

	// Encode the valid JSON data
	input := bytes.NewReader(jsonData)
	encoded, err := encoder.Encode(input)
	if err != nil {
		b.Fatalf("Encode failed: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := encoder.Decode(encoded.Bytes(), &result); err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}
