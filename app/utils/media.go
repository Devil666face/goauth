package utils

import (
	"os"
	"path/filepath"
)

func SetMediaPath(Mediafs string) string {
	BASE_DIR, _ := os.Getwd()
	mediaPath, _ := filepath.Abs(filepath.Join(BASE_DIR, Mediafs))
	if _, err := os.Stat(mediaPath); os.IsNotExist(err) {
		os.MkdirAll(mediaPath, os.ModePerm)
	}
	return mediaPath
}
