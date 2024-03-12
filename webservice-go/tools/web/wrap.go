package web

import (
	"io"
	"net/http"
	"sync/atomic"
)

type respWriterWrapper struct {
	http.ResponseWriter
	status  int
	written atomic.Int64
}

func newRespWriterWrapper(w http.ResponseWriter) *respWriterWrapper {
	return &respWriterWrapper{ResponseWriter: w, status: http.StatusOK}
}

func (rw *respWriterWrapper) Status() int {
	return rw.status
}

func (rw *respWriterWrapper) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *respWriterWrapper) Write(p []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(p)
	rw.written.Add(int64(n))
	return n, err
}

// bodyWrapper wraps a http.Request.Body (an io.ReadCloser) to track the number
// of bytes read and the last error.
type bodyWrapper struct {
	io.ReadCloser
	read atomic.Int64
}

func (w *bodyWrapper) Read(b []byte) (int, error) {
	n, err := w.ReadCloser.Read(b)
	w.read.Add(int64(n))
	return n, err
}

func (w *bodyWrapper) Close() error {
	return w.ReadCloser.Close()
}
