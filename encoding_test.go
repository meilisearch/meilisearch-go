package meilisearch

import (
	"bytes"
	"encoding/json"
	"errors"
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
	g := newEncoding(GzipEncoding, DefaultCompression)
	d := newEncoding(DeflateEncoding, DefaultCompression)
	b := newEncoding(BrotliEncoding, DefaultCompression)

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

func testEncoder(t *testing.T, enc encoder, original *mockData) {
	originalJSON, err := json.Marshal(original)
	assert.NoError(t, err)

	readCloser := io.NopCloser(bytes.NewReader(originalJSON))
	encodedReader, err := enc.Encode(readCloser)
	assert.NoError(t, err)
	defer encodedReader.Close()

	encodedData, err := io.ReadAll(encodedReader)
	assert.NoError(t, err)
	assert.NotEmpty(t, encodedData)

	var decoded mockData
	err = enc.Decode(encodedData, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, original, &decoded)

	var invalidTarget int
	err = enc.Decode(encodedData, &invalidTarget)
	assert.Error(t, err)
}

func TestGzipEncoder(t *testing.T) {
	testEncoder(t, newEncoding(GzipEncoding, DefaultCompression), &mockData{Name: "John Doe", Age: 30})
}

func TestDeflateEncoder(t *testing.T) {
	testEncoder(t, newEncoding(DeflateEncoding, DefaultCompression), &mockData{Name: "Jane Doe", Age: 25})
}

func TestBrotliEncoder(t *testing.T) {
	testEncoder(t, newEncoding(BrotliEncoding, DefaultCompression), &mockData{Name: "Jane Doe", Age: 25})
}

func TestEncoder_EmptyData(t *testing.T) {
	testEncoder(t, newEncoding(GzipEncoding, DefaultCompression), &mockData{})
	testEncoder(t, newEncoding(DeflateEncoding, DefaultCompression), &mockData{})
	testEncoder(t, newEncoding(BrotliEncoding, DefaultCompression), &mockData{})
}

func TestEncoder_InvalidDecode(t *testing.T) {
	encoders := []encoder{
		newEncoding(GzipEncoding, DefaultCompression),
		newEncoding(DeflateEncoding, DefaultCompression),
		newEncoding(BrotliEncoding, DefaultCompression),
	}
	for _, enc := range encoders {
		var decoded mockData
		err := enc.Decode([]byte("invalid data"), &decoded)
		assert.Error(t, err)
	}
}

func TestCopyZeroAlloc(t *testing.T) {
	t.Run("RegularCopy", func(t *testing.T) {
		src := strings.NewReader("hello world")
		dst := &bytes.Buffer{}
		n, err := copyZeroAlloc(dst, src)
		assert.NoError(t, err)
		assert.Equal(t, int64(11), n)
		assert.Equal(t, "hello world", dst.String())
	})
	t.Run("EmptySource", func(t *testing.T) {
		src := strings.NewReader("")
		dst := &bytes.Buffer{}
		n, err := copyZeroAlloc(dst, src)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), n)
		assert.Equal(t, "", dst.String())
	})
	t.Run("LargeDataCopy", func(t *testing.T) {
		data := strings.Repeat("a", 10000)
		src := strings.NewReader(data)
		dst := &bytes.Buffer{}
		n, err := copyZeroAlloc(dst, src)
		assert.NoError(t, err)
		assert.Equal(t, int64(len(data)), n)
		assert.Equal(t, data, dst.String())
	})
	t.Run("ErrorOnWrite", func(t *testing.T) {
		src := strings.NewReader("hello world")
		dst := &errorWriter{}
		n, err := copyZeroAlloc(dst, src)
		assert.Error(t, err)
		assert.Equal(t, int64(0), n)
		assert.Equal(t, "write error", err.Error())
	})
	t.Run("ErrorOnRead", func(t *testing.T) {
		src := &errorReader{}
		dst := &bytes.Buffer{}
		n, err := copyZeroAlloc(dst, src)
		assert.Error(t, err)
		assert.Equal(t, int64(0), n)
		assert.Equal(t, "read error", err.Error())
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
				src := strings.NewReader(data)
				buf := &bytes.Buffer{}
				_, _ = copyZeroAlloc(buf, src)
				mu.Lock()
				defer mu.Unlock()
				dst.Write(buf.Bytes())
			}()
		}
		wg.Wait()
		mu.Lock()
		assert.Equal(t, strings.Repeat(data, 10), dst.String())
		mu.Unlock()
	})
}
