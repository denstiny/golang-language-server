package file

import (
	"io"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func Open(path string) (*os.File, error) {
	return os.Open(path)
}

func ReadAll(f *os.File) ([]byte, error) {
	return io.ReadAll(f)
}
