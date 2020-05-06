package path_relative

import (
	"log"
	"path/filepath"
	"runtime"
)

// Get Absolute Path from relative path
func GetAbsPath(relativePath string) string {
	_, b, _, _ := runtime.Caller(1)
	path, err := filepath.Abs(filepath.Dir(b) + "/" + relativePath)
	if err != nil {
		log.Panic(err.Error())
	}
	return path
}
