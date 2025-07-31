package meilisearch

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"github.com/andybalholm/brotli"
	"io"
	"sync"
)

type encoder interface {
	Encode(io.Reader) (io.ReadCloser, error)
	Decode([]byte, interface{}) error
}

func newEncoding(ce ContentEncoding, level EncodingCompressionLevel) encoder {
	switch ce {
	case GzipEncoding:
		return &gzipEncoder{
			gzWriterPool: &sync.Pool{
				New: func() interface{} {
					w, err := gzip.NewWriterLevel(io.Discard, level.Int())
					return &gzipWriter{writer: w, err: err}
				},
			},
			bufferPool: &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }},
		}
	case DeflateEncoding:
		return &flateEncoder{
			flWriterPool: &sync.Pool{
				New: func() interface{} {
					w, err := zlib.NewWriterLevel(io.Discard, level.Int())
					return &flateWriter{writer: w, err: err}
				},
			},
			bufferPool: &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }},
		}
	case BrotliEncoding:
		return &brotliEncoder{
			brWriterPool: &sync.Pool{
				New: func() interface{} {
					return brotli.NewWriterLevel(io.Discard, level.Int())
				},
			},
			bufferPool: &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }},
		}
	default:
		return nil
	}
}

type gzipEncoder struct {
	gzWriterPool *sync.Pool
	bufferPool   *sync.Pool
}

type gzipWriter struct {
	writer *gzip.Writer
	err    error
}

func (g *gzipEncoder) Encode(rc io.Reader) (io.ReadCloser, error) {
	w := g.gzWriterPool.Get().(*gzipWriter)
	defer g.gzWriterPool.Put(w)
	if w.err != nil {
		return nil, w.err
	}
	defer func() {
		_ = w.writer.Close()
	}()

	buf := g.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	w.writer.Reset(buf)

	if _, err := copyZeroAlloc(w.writer, rc); err != nil {
		g.bufferPool.Put(buf)
		return nil, err
	}

	return &pooledBuffer{Buffer: buf, pool: g.bufferPool}, nil
}

func (g *gzipEncoder) Decode(data []byte, vPtr interface{}) error {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
	}()
	return json.NewDecoder(r).Decode(vPtr)
}

type flateEncoder struct {
	flWriterPool *sync.Pool
	bufferPool   *sync.Pool
}

type flateWriter struct {
	writer *zlib.Writer
	err    error
}

func (f *flateEncoder) Encode(rc io.Reader) (io.ReadCloser, error) {
	w := f.flWriterPool.Get().(*flateWriter)
	defer f.flWriterPool.Put(w)
	if w.err != nil {
		return nil, w.err
	}
	defer func() {
		_ = w.writer.Close()
	}()

	buf := f.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	w.writer.Reset(buf)

	if _, err := copyZeroAlloc(w.writer, rc); err != nil {
		f.bufferPool.Put(buf)
		return nil, err
	}

	return &pooledBuffer{Buffer: buf, pool: f.bufferPool}, nil
}

func (f *flateEncoder) Decode(data []byte, vPtr interface{}) error {
	r, err := zlib.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
	}()
	return json.NewDecoder(r).Decode(vPtr)
}

type brotliEncoder struct {
	brWriterPool *sync.Pool
	bufferPool   *sync.Pool
}

func (b *brotliEncoder) Encode(rc io.Reader) (io.ReadCloser, error) {
	w := b.brWriterPool.Get().(*brotli.Writer)
	defer func() {
		_ = w.Close()
	}()
	defer b.brWriterPool.Put(w)

	buf := b.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	w.Reset(buf)

	if _, err := copyZeroAlloc(w, rc); err != nil {
		b.bufferPool.Put(buf)
		return nil, err
	}

	return &pooledBuffer{Buffer: buf, pool: b.bufferPool}, nil
}

func (b *brotliEncoder) Decode(data []byte, vPtr interface{}) error {
	r := brotli.NewReader(bytes.NewBuffer(data))
	return json.NewDecoder(r).Decode(vPtr)
}

var copyBufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4096)
	},
}

func copyZeroAlloc(w io.Writer, r io.Reader) (int64, error) {
	if wt, ok := r.(io.WriterTo); ok {
		return wt.WriteTo(w)
	}
	if rt, ok := w.(io.ReaderFrom); ok {
		return rt.ReadFrom(r)
	}

	vbuf := copyBufPool.Get()
	buf := vbuf.([]byte)
	n, err := io.CopyBuffer(w, r, buf)
	copyBufPool.Put(vbuf)
	return n, err
}
