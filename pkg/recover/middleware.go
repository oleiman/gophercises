package recover

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
)

func RecoverMW(next *http.ServeMux, dev bool) http.Handler {
	next.HandleFunc("/debug/", renderSource)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pw := &writer{
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
					tmpl, err := template.ParseFiles("data/15/layout.html")
					if err == nil {
						msg := fmt.Sprintf("panic: %v", r)
						st := newStackTrace(stack, msg)
						tmpl.Execute(w, st)
					} else {
						fmt.Fprintf(w, "<h1>panic: %v</h1>\n<pre>%s</pre>", r, string(stack))
						log.Printf("%v", err)
					}
				}
			} else {
				fmt.Fprint(w, string(pw.Body))
			}
		}()
		next.ServeHTTP(pw, r)
	})
}
