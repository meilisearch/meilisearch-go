package meilisearch

import (
	"bytes"
	"sync"
)

type pooledBuffer struct {
	*bytes.Buffer
	pool *sync.Pool
}

func (pb *pooledBuffer) Read(p []byte) (int, error) {
	return pb.Buffer.Read(p)
}

func (pb *pooledBuffer) Close() error {
	pb.Reset()
	pb.pool.Put(pb.Buffer)
	return nil
}
