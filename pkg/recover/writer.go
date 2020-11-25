package recover

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type writer struct {
	w          http.ResponseWriter
	Body       []byte
	StatusCode int
}

func (pw *writer) Header() http.Header {
	return pw.w.Header()
}

func (pw *writer) Write(body []byte) (int, error) {
	pw.Body = append(pw.Body, body...)
	return len(body), nil
}

func (pw *writer) WriteHeader(status int) {
	pw.StatusCode = status
}

func (pw *writer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := pw.w.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("ResponseWriter does not implement http.Hijacker")
	}
	return hj.Hijack()
}

func (pw *writer) Flush() {
	f, ok := pw.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}
