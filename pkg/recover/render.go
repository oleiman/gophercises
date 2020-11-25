package recover

import (
	"fmt"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func fileExists(fname string) bool {
	info, err := os.Stat(fname)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func renderSource(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	if len(path) == 0 {
		http.Error(w, "Missing source file path for debug output", http.StatusBadRequest)
	} else if !fileExists(path) {
		http.Error(w, "File does not exist", http.StatusBadRequest)
	}
	highlight := make([][2]int, 0, 1)
	line, err := strconv.Atoi(r.FormValue("line"))
	if err == nil {
		highlight = append(highlight, [2]int{line, line})
	} else {
		fmt.Println(err)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "<h1>file: %s</h1>\n", path)

	lexer := lexers.Get("go")
	style := styles.Get("monokai")
	formatter := html.New(html.WithLineNumbers(true), html.HighlightLines(highlight))
	iterator, err := lexer.Tokenise(nil, string(content))
	err = formatter.Format(w, style, iterator)
}
