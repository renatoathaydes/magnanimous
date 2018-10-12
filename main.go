package main

import (
	"github.com/renatoathaydes/magnanimous/mg"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	procFiles, err1 := getFilesAt("source/processed")
	if err1 != nil {
		panic(err1)
	}
	staticFiles, err2 := getFilesAt("source/static")
	if err2 != nil {
		panic(err2)
	}
	otherFiles, err3 := getFilesAt("source", "source/processed/", "source/static/")
	if err3 != nil {
		panic(err3)
	}
	webFiles := make(mg.WebFilesMap, len(procFiles)+len(staticFiles)+len(otherFiles))
	mg.Process(&procFiles, "source/processed", webFiles)
	mg.CopyAll(&staticFiles, "source/static", webFiles)
	mg.AddNonWritables(&otherFiles, "source", webFiles)
	err := mg.WriteTo("target", webFiles)
	if err != nil {
		panic(err)
	}
}

func getFilesAt(root string, exclusions ...string) ([]string, error) {
	var files []string
	notExcluded := func(path string) bool {
		for _, e := range exclusions {
			if strings.HasPrefix(path, e) {
				return false
			}
		}
		return true
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && notExcluded(path) {
			files = append(files, path)
		}
		return err
	})
	return files, err
}
