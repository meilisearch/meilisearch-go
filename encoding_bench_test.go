package meilisearch

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

func generate1MBData() []byte {
	return bytes.Repeat([]byte("a"), 1024*1024)
}

func BenchmarkGzipEncoder(b *testing.B) {
	encoder := newEncoding(GzipEncoding, DefaultCompression)
	raw := generate1MBData()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(raw)
		encoded, err := encoder.Encode(reader)
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
		_, _ = io.Copy(io.Discard, encoded)
		_ = encoded.Close()
	}
}

func BenchmarkDeflateEncoder(b *testing.B) {
	encoder := newEncoding(DeflateEncoding, DefaultCompression)
	raw := generate1MBData()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(raw)
		encoded, err := encoder.Encode(reader)
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
		_, _ = io.Copy(io.Discard, encoded)
		_ = encoded.Close()
	}
}

func BenchmarkBrotliEncoder(b *testing.B) {
	encoder := newEncoding(BrotliEncoding, DefaultCompression)
	raw := generate1MBData()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(raw)
		encoded, err := encoder.Encode(reader)
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
		_, _ = io.Copy(io.Discard, encoded)
		_ = encoded.Close()
	}
}

func BenchmarkGzipDecoder(b *testing.B) {
	encoder := newEncoding(GzipEncoding, DefaultCompression)
	jsonData, _ := json.Marshal(sampleMapData())

	encoded, _ := encoder.Encode(bytes.NewReader(jsonData))
	defer encoded.Close()
	payload, _ := io.ReadAll(encoded)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := encoder.Decode(payload, &result); err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkFlateDecoder(b *testing.B) {
	encoder := newEncoding(DeflateEncoding, DefaultCompression)
	jsonData, _ := json.Marshal(sampleMapData())

	encoded, _ := encoder.Encode(bytes.NewReader(jsonData))
	defer encoded.Close()
	payload, _ := io.ReadAll(encoded)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := encoder.Decode(payload, &result); err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkBrotliDecoder(b *testing.B) {
	encoder := newEncoding(BrotliEncoding, DefaultCompression)
	jsonData, _ := json.Marshal(sampleMapData())

	encoded, _ := encoder.Encode(bytes.NewReader(jsonData))
	defer encoded.Close()
	payload, _ := io.ReadAll(encoded)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := encoder.Decode(payload, &result); err != nil {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func sampleMapData() map[string]interface{} {
	return map[string]interface{}{
		"key1": "value1",
		"key2": 12345,
		"key3": []string{"item1", "item2", "item3"},
	}
}
