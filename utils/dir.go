package utils

import (
	"os"
	"path/filepath"
)

func GetCurrentDir() string {
	pathExecutable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(pathExecutable)
}
