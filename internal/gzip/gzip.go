package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w            http.ResponseWriter
	zw           *gzip.Writer
	contentType  map[string]bool
	compressable bool
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:            w,
		zw:           gzip.NewWriter(w),
		contentType:  map[string]bool{"application/json": true, "text/html": true},
		compressable: false,
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.writer().Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {

	contentType := c.w.Header().Get("Content-Type")

	if _, ok := c.contentType[contentType]; ok {
		c.w.Header().Set("Content-Encoding", "gzip")
		c.compressable = true
	}

	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	if c.compressable {
		return c.writer().(io.WriteCloser).Close()
	}
	return nil
}

// writer() определяет обработчик для записи с жатием или без.
func (c *compressWriter) writer() io.Writer {
	if c.compressable {
		return c.zw
	} else {
		return c.w
	}
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
