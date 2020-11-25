package recover

import (
	"bufio"
	"strconv"
	"strings"
)

type line struct {
	File   bool
	Path   string
	Number int
	Text   string
}

type StackTrace struct {
	Message string
	Lines   []line
}

func newStackTrace(raw []byte, msg string) StackTrace {
	st := StackTrace{
		Message: msg,
		Lines:   make([]line, 0),
	}
	scanner := bufio.NewScanner(strings.NewReader(string(raw)))
	for scanner.Scan() {
		t := scanner.Text()
		if t[0] == '\t' {
			path := t[1:strings.IndexByte(t, ':')]
			rest := t[strings.IndexByte(t, ':')+1:]
			n, err := strconv.Atoi(rest)
			if err != nil {
				n, _ = strconv.Atoi(rest[:strings.IndexByte(rest, ' ')])
			}
			st.Lines = append(st.Lines,
				line{
					File:   true,
					Path:   path,
					Number: n,
					Text:   t,
				})
		} else {
			st.Lines = append(st.Lines,
				line{
					File: false,
					Text: t,
				})
		}
	}
	return st
}
