package meilisearch

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type mockData struct {
	Name string
	Age  int
}

func TestGzipEncoder(t *testing.T) {
	encoder := newEncoding(GzipEncoding, DefaultCompression)
	assert.NotNil(t, encoder, "gzip encoder should not be nil")

	original := &mockData{Name: "John Doe", Age: 30}

	originalJSON, err := json.Marshal(original)
	assert.NoError(t, err, "marshalling original data should not produce an error")

	readCloser := io.NopCloser(bytes.NewReader(originalJSON))

	encodedData, err := encoder.Encode(readCloser)
	assert.NoError(t, err, "encoding should not produce an error")
	assert.NotNil(t, encodedData, "encoded data should not be nil")

	var decoded mockData
	err = encoder.Decode(encodedData.Bytes(), &decoded)
	assert.NoError(t, err, "decoding should not produce an error")
	assert.Equal(t, original, &decoded, "decoded data should match the original")
}

func TestDeflateEncoder(t *testing.T) {
	encoder := newEncoding(DeflateEncoding, DefaultCompression)
	assert.NotNil(t, encoder, "deflate encoder should not be nil")

	original := &mockData{Name: "Jane Doe", Age: 25}

	originalJSON, err := json.Marshal(original)
	assert.NoError(t, err, "marshalling original data should not produce an error")

	readCloser := io.NopCloser(bytes.NewReader(originalJSON))

	encodedData, err := encoder.Encode(readCloser)
	assert.NoError(t, err, "encoding should not produce an error")
	assert.NotNil(t, encodedData, "encoded data should not be nil")

	var decoded mockData
	err = encoder.Decode(encodedData.Bytes(), &decoded)
	assert.NoError(t, err, "decoding should not produce an error")
	assert.Equal(t, original, &decoded, "decoded data should match the original")
}

func TestBrotliEncoder(t *testing.T) {
	encoder := newEncoding(BrotliEncoding, DefaultCompression)
	assert.NotNil(t, encoder, "brotli encoder should not be nil")

	original := &mockData{Name: "Jane Doe", Age: 25}

	originalJSON, err := json.Marshal(original)
	assert.NoError(t, err, "marshalling original data should not produce an error")

	readCloser := io.NopCloser(bytes.NewReader(originalJSON))

	encodedData, err := encoder.Encode(readCloser)
	assert.NoError(t, err, "encoding should not produce an error")
	assert.NotNil(t, encodedData, "encoded data should not be nil")

	var decoded mockData
	err = encoder.Decode(encodedData.Bytes(), &decoded)
	assert.NoError(t, err, "decoding should not produce an error")
	assert.Equal(t, original, &decoded, "decoded data should match the original")
}

func TestGzipEncoder_EmptyData(t *testing.T) {
	encoder := newEncoding(GzipEncoding, DefaultCompression)
	assert.NotNil(t, encoder, "gzip encoder should not be nil")

	original := &mockData{}

	originalJSON, err := json.Marshal(original)
	assert.NoError(t, err, "marshalling original data should not produce an error")

	readCloser := io.NopCloser(bytes.NewReader(originalJSON))

	encodedData, err := encoder.Encode(readCloser)
	assert.NoError(t, err, "encoding should not produce an error")
	assert.NotNil(t, encodedData, "encoded data should not be nil")

	var decoded mockData
	err = encoder.Decode(encodedData.Bytes(), &decoded)
	assert.NoError(t, err, "decoding should not produce an error")
	assert.Equal(t, original, &decoded, "decoded data should match the original")
}

func TestDeflateEncoder_EmptyData(t *testing.T) {
	encoder := newEncoding(DeflateEncoding, DefaultCompression)
	assert.NotNil(t, encoder, "deflate encoder should not be nil")

	original := &mockData{}

	originalJSON, err := json.Marshal(original)
	assert.NoError(t, err, "marshalling original data should not produce an error")

	readCloser := io.NopCloser(bytes.NewReader(originalJSON))

	encodedData, err := encoder.Encode(readCloser)
	assert.NoError(t, err, "encoding should not produce an error")
	assert.NotNil(t, encodedData, "encoded data should not be nil")

	var decoded mockData
	err = encoder.Decode(encodedData.Bytes(), &decoded)
	assert.NoError(t, err, "decoding should not produce an error")
	assert.Equal(t, original, &decoded, "decoded data should match the original")
}

func TestBrotliEncoder_EmptyData(t *testing.T) {
	encoder := newEncoding(BrotliEncoding, DefaultCompression)
	assert.NotNil(t, encoder, "brotli encoder should not be nil")

	original := &mockData{}

	originalJSON, err := json.Marshal(original)
	assert.NoError(t, err, "marshalling original data should not produce an error")

	readCloser := io.NopCloser(bytes.NewReader(originalJSON))

	encodedData, err := encoder.Encode(readCloser)
	assert.NoError(t, err, "encoding should not produce an error")
	assert.NotNil(t, encodedData, "encoded data should not be nil")

	var decoded mockData
	err = encoder.Decode(encodedData.Bytes(), &decoded)
	assert.NoError(t, err, "decoding should not produce an error")
	assert.Equal(t, original, &decoded, "decoded data should match the original")
}

func TestGzipEncoder_InvalidData(t *testing.T) {
	encoder := newEncoding(GzipEncoding, DefaultCompression)
	assert.NotNil(t, encoder, "gzip encoder should not be nil")

	var decoded mockData
	err := encoder.Decode([]byte("invalid data"), &decoded)
	assert.Error(t, err, "decoding invalid data should produce an error")
}

func TestDeflateEncoder_InvalidData(t *testing.T) {
	encoder := newEncoding(DeflateEncoding, DefaultCompression)
	assert.NotNil(t, encoder, "deflate encoder should not be nil")

	var decoded mockData
	err := encoder.Decode([]byte("invalid data"), &decoded)
	assert.Error(t, err, "decoding invalid data should produce an error")
}

func TestBrotliEncoder_InvalidData(t *testing.T) {
	encoder := newEncoding(BrotliEncoding, DefaultCompression)
	assert.NotNil(t, encoder, "brotli encoder should not be nil")

	var decoded mockData
	err := encoder.Decode([]byte("invalid data"), &decoded)
	assert.Error(t, err, "decoding invalid data should produce an error")
}
