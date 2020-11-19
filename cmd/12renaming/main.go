package main

import (
	"flag"
	// "fmt"
	"os"
	_path "path"
	"path/filepath"
	"regexp"
)

func matcher(re *regexp.Regexp, err error) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err2 error) error {
		if err != nil {
			return err
		} else if err2 != nil {
			return err2
		}

		loc := re.FindStringIndex(info.Name())
		if len(loc) == 0 {
			return nil
		} else {
			match := info.Name()[loc[0]:loc[1]]
			pfx := info.Name()[0:loc[0]]
			sfx := info.Name()[loc[1]:]
			newFile := match[1:] + " - " + pfx + sfx
			dir, _ := _path.Split(path)
			if err := os.Rename(path, _path.Join(dir, newFile)); err != nil {
				return err
			}
		}
		return nil
	}
}

func main() {
	var path string
	flag.StringVar(&path, "d", "", "Directory to traverse for renaming")
	flag.Parse()

	if path == "" {
		panic("Please provide a directory")
	} else if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}

	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(path,
		matcher(regexp.Compile("_[0-9]{3}")))
	if err != nil {
		panic(err)
	}

}
