package path_util

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ModulePath(path string) string {
	rootCode := strings.Split(path, "/")[0]
	rootPath := ""

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	currentPath := filepath.Clean(pwd)
	if strings.Contains(currentPath, rootCode) {
		returnPath, ok := FindRoot(currentPath, rootCode, "go.mod")
		if ok {
			rootPath = returnPath
		}
	}

	binPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(binPath, rootCode) {
		returnPath, ok := FindRoot(binPath, rootCode, "go.mod")
		if ok {
			rootPath = returnPath
		}
	}

	if rootPath == "" {
		_, fpath, _, _ := runtime.Caller(0)
		pkgFilePath := filepath.Clean(fpath)
		rootStringLoc := strings.LastIndex(pkgFilePath, rootCode)
		rootPath = pkgFilePath[:rootStringLoc]
		if !Exists(rootPath + filepath.Clean(path)) {
			rootPath = currentPath[:strings.LastIndex(currentPath, rootCode)]
		}
	}

	target := rootPath + filepath.Clean(path)

	location, err := filepath.Rel(currentPath, target)
	if err != nil {
		log.Fatal(err)
	}

	return location
}

func Exists(fpath string) bool {
	_, err := os.Stat(fpath)
	return !os.IsNotExist(err)
}

func FindRoot(path string, rootCode string, objName string) (string, bool) {
	rootPath := path
	loc := strings.LastIndex(rootPath, rootCode)
	for loc != -1 {
		rootPath = rootPath[:loc+len(rootCode)]
		if Exists(rootPath + filepath.Clean("/"+objName)) {
			return rootPath[:loc], true
		}
		rootPath = rootPath[:loc]
		loc = strings.LastIndex(rootPath, rootCode)
	}
	return "", false
}
