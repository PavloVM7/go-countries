package provider

import (
	"os"
)

type FileCountryReader struct {
	filePath string
}

func (r *FileCountryReader) Read() ([]byte, error) {
	return os.ReadFile(r.filePath)
}

func NewFileCountryReader(filePath string) *FileCountryReader {
	return &FileCountryReader{
		filePath: filePath,
	}
}
