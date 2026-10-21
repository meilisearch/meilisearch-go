package meilisearch

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/andybalholm/brotli"
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
	defer func() {
		_ = encoded.Close()
	}()
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
	defer func() {
		_ = encoded.Close()
	}()
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
	defer func() {
		_ = encoded.Close()
	}()
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

// decodeNonPooled mirrors what encoder.Decode used to do before readers were
// pooled: it allocates a brand-new gzip/zlib/brotli reader on every call
// instead of reusing one from a sync.Pool. It exists only to give the table
// benchmark below a "before" baseline to compare the pooled implementation
// against.
func decodeNonPooled(ce ContentEncoding, data []byte, vPtr interface{}) error {
	switch ce {
	case GzipEncoding:
		r, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return err
		}
		defer func() { _ = r.Close() }()
		return json.NewDecoder(r).Decode(vPtr)
	case DeflateEncoding:
		r, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			return err
		}
		defer func() { _ = r.Close() }()
		return json.NewDecoder(r).Decode(vPtr)
	case BrotliEncoding:
		r := brotli.NewReader(bytes.NewReader(data))
		return json.NewDecoder(r).Decode(vPtr)
	default:
		return fmt.Errorf("unsupported content encoding: %v", ce)
	}
}

func BenchmarkDecoderPooling(b *testing.B) {
	jsonData, err := json.Marshal(sampleMapData())
	if err != nil {
		b.Fatalf("Marshal failed: %v", err)
	}

	cases := []struct {
		name string
		ce   ContentEncoding
	}{
		{name: "Gzip", ce: GzipEncoding},
		{name: "Deflate", ce: DeflateEncoding},
		{name: "Brotli", ce: BrotliEncoding},
	}

	for _, tc := range cases {
		encoder := newEncoding(tc.ce, DefaultCompression)

		encoded, err := encoder.Encode(bytes.NewReader(jsonData))
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
		payload, err := io.ReadAll(encoded)
		_ = encoded.Close()
		if err != nil {
			b.Fatalf("ReadAll failed: %v", err)
		}

		b.Run(tc.name+"/Pooled", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var result map[string]interface{}
				if err := encoder.Decode(payload, &result); err != nil {
					b.Fatalf("Decode failed: %v", err)
				}
			}
		})

		b.Run(tc.name+"/NonPooled", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var result map[string]interface{}
				if err := decodeNonPooled(tc.ce, payload, &result); err != nil {
					b.Fatalf("Decode failed: %v", err)
				}
			}
		})
	}
}
