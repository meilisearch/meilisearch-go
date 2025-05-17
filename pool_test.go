package meilisearch

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestPooledBuffer_Read(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	buf := pool.Get().(*bytes.Buffer)
	buf.WriteString("hello world")

	pb := &pooledBuffer{
		Buffer: buf,
		pool:   pool,
	}

	readBuf := make([]byte, 5)
	n, err := pb.Read(readBuf)

	require.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "hello", string(readBuf[:n]))
}

func TestPooledBuffer_Close(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	buf := pool.Get().(*bytes.Buffer)
	buf.WriteString("data to reset")

	bufAddr := buf

	pb := &pooledBuffer{
		Buffer: buf,
		pool:   pool,
	}

	err := pb.Close()
	require.NoError(t, err)

	got := pool.Get().(*bytes.Buffer)
	assert.Equal(t, bufAddr, got, "buffer address should be the same")
	assert.Zero(t, got.Len(), "buffer should be reset to length 0")
	assert.Equal(t, "", got.String(), "buffer string should be empty")
}
