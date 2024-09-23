package meilisearch

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"errors"
	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"sync"
	"testing"
)

type mockData struct {
	Name string
	Age  int
}

type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write error")
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func Test_Encode_ErrorOnNewWriter(t *testing.T) {
	g := &gzipEncoder{
		gzWriterPool: &sync.Pool{
			New: func() interface{} {
				return &gzipWriter{
					writer: nil,
					err:    errors.New("new writer error"),
				}
			},
		},
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}

	d := &flateEncoder{
		flWriterPool: &sync.Pool{
			New: func() interface{} {
				return &flateWriter{
					writer: nil,
					err:    errors.New("new writer error"),
				}
			},
		},
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}

	_, err := g.Encode(bytes.NewReader([]byte("test")))
	require.Error(t, err)

	_, err = d.Encode(bytes.NewReader([]byte("test")))
	require.Error(t, err)
}

func Test_Encode_ErrorInCopyZeroAlloc(t *testing.T) {
	g := &gzipEncoder{
		gzWriterPool: &sync.Pool{
			New: func() interface{} {
				w, _ := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
				return &gzipWriter{
					writer: w,
					err:    nil,
				}
			},
		},
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}

	d := &flateEncoder{
		flWriterPool: &sync.Pool{
			New: func() interface{} {
				w, err := zlib.NewWriterLevel(io.Discard, flate.DefaultCompression)
				return &flateWriter{
					writer: w,
					err:    err,
				}
			},
		},
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}

	b := &brotliEncoder{
		brWriterPool: &sync.Pool{
			New: func() interface{} {
				return brotli.NewWriterLevel(io.Discard, brotli.DefaultCompression)
			},
		},
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}

	_, err := g.Encode(&errorReader{})
	require.Error(t, err)

	_, err = d.Encode(&errorReader{})
	require.Error(t, err)

	_, err = b.Encode(&errorReader{})
	require.Error(t, err)
}

func Test_InvalidContentType(t *testing.T) {
	enc := newEncoding("invalid", DefaultCompression)
	require.Nil(t, enc)
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

	var invalidType int
	err = encoder.Decode(encodedData.Bytes(), &invalidType)
	assert.Error(t, err)
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

	var invalidType int
	err = encoder.Decode(encodedData.Bytes(), &invalidType)
	assert.Error(t, err)
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

	var invalidType int
	err = encoder.Decode(encodedData.Bytes(), &invalidType)
	assert.Error(t, err)
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

func TestCopyZeroAlloc(t *testing.T) {
	t.Run("RegularCopy", func(t *testing.T) {
		src := strings.NewReader("hello world")
		dst := &bytes.Buffer{}

		n, err := copyZeroAlloc(dst, src)
		assert.NoError(t, err, "copy should not produce an error")
		assert.Equal(t, int64(11), n, "copy length should be 11")
		assert.Equal(t, "hello world", dst.String(), "destination should contain the copied data")
	})

	t.Run("EmptySource", func(t *testing.T) {
		src := strings.NewReader("")
		dst := &bytes.Buffer{}

		n, err := copyZeroAlloc(dst, src)
		assert.NoError(t, err, "copy should not produce an error")
		assert.Equal(t, int64(0), n, "copy length should be 0")
		assert.Equal(t, "", dst.String(), "destination should be empty")
	})

	t.Run("LargeDataCopy", func(t *testing.T) {
		data := strings.Repeat("a", 10000)
		src := strings.NewReader(data)
		dst := &bytes.Buffer{}

		n, err := copyZeroAlloc(dst, src)
		assert.NoError(t, err, "copy should not produce an error")
		assert.Equal(t, int64(len(data)), n, "copy length should match the source data length")
		assert.Equal(t, data, dst.String(), "destination should contain the copied data")
	})

	t.Run("ErrorOnWrite", func(t *testing.T) {
		src := strings.NewReader("hello world")
		dst := &errorWriter{}

		n, err := copyZeroAlloc(dst, src)
		assert.Error(t, err, "copy should produce an error")
		assert.Equal(t, int64(0), n, "copy length should be 0 due to the error")
		assert.Equal(t, "write error", err.Error(), "error should match expected error")
	})

	t.Run("ErrorOnRead", func(t *testing.T) {
		src := &errorReader{}
		dst := &bytes.Buffer{}

		n, err := copyZeroAlloc(dst, src)
		assert.Error(t, err, "copy should produce an error")
		assert.Equal(t, int64(0), n, "copy length should be 0 due to the error")
		assert.Equal(t, "read error", err.Error(), "error should match expected error")
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		var wg sync.WaitGroup
		data := "concurrent data"
		var mu sync.Mutex
		dst := &bytes.Buffer{}

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				src := strings.NewReader(data) // each goroutine gets its own reader
				buf := &bytes.Buffer{}         // each goroutine uses a separate buffer
				_, _ = copyZeroAlloc(buf, src)
				mu.Lock()
				defer mu.Unlock()
				dst.Write(buf.Bytes()) // safely combine results
			}()
		}
		wg.Wait()

		mu.Lock()
		assert.Equal(t, strings.Repeat(data, 10), dst.String(), "destination should contain the copied data")
		mu.Unlock()
	})
}
