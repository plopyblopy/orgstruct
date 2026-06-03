package testdata

import (
	"path/filepath"
	"runtime"
	"time"
)

func GetCurDirPath() string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	return currentDir
}

func CopyInt(v *int) *int {
	if v == nil {
		return nil
	}
	copy := *v
	return &copy
}

func CopyDateTime(v *time.Time) *time.Time {
	if v == nil {
		return nil
	}
	copy := *v
	return &copy
}
