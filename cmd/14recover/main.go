package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
)

type panicWriter struct {
	w          http.ResponseWriter
	Body       []byte
	StatusCode int
}

func (pw *panicWriter) Header() http.Header {
	return pw.w.Header()
}

func (pw *panicWriter) Write(body []byte) (int, error) {
	pw.Body = append(pw.Body, body...)
	return len(body), nil
}

func (pw *panicWriter) WriteHeader(status int) {
	pw.StatusCode = status
}

func (pw *panicWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := pw.w.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("ResponseWriter does not implement http.Hijacker")
	}
	return hj.Hijack()
}

func (pw *panicWriter) Flush() {
	f, ok := pw.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", panicMiddleware(mux, true)))
}

func panicMiddleware(next http.Handler, dev bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pw := &panicWriter{
			w:    w,
			Body: make([]byte, 0),
		}
		defer func() {
			if r := recover(); r != nil {
				// fmt.Fprint(w, "<h1>Something went wrong...</h1>")
				log.Println("Recovered from panic: ", r)
				stack := debug.Stack()
				log.Println(string(stack))
				if !dev {
					http.Error(w, "Something went wrong!", http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "<h1>panic: %v</h1>\n<pre>%s</pre>", r, string(stack))
				}
			} else {
				fmt.Fprint(w, string(pw.Body))
			}
		}()
		next.ServeHTTP(pw, r)
	})
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
