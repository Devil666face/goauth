package utils

import (
	"os"
	"path/filepath"
)

func SetPath(dir string) string {
	BASE_DIR, _ := os.Getwd()
	dirPath, _ := filepath.Abs(filepath.Join(BASE_DIR, dir))
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	}
	return dirPath
}
