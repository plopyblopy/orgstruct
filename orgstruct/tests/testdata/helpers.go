package testdata

import (
	"path/filepath"
	"runtime"
)

func GetCurDirPath() string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	return currentDir
}
