package utils

import (
	"os"
	"path/filepath"
)

func SetPath(dir string) string {
	baseDir, _ := os.Getwd()
	dirPath, _ := filepath.Abs(filepath.Join(baseDir, dir))
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	}
	return dirPath
}
