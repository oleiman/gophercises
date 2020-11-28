package util

import (
	"os"
)

func FileExists(fname string) bool {
	info, err := os.Stat(fname)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
