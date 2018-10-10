package main

import (
	"github.com/renatoathaydes/magnanimous/mg"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	procFiles := getFilesAt("source/processed")
	staticFiles := getFilesAt("source/static")
	otherFiles := getFilesAt("source", "source/processed/", "source/static/")
	webFiles := make(mg.WebFilesMap, len(procFiles)+len(staticFiles))
	mg.Process(&procFiles, "source/processed", webFiles)
	mg.CopyAll(&staticFiles, "source/static", webFiles)
	mg.AddNonWritables(&otherFiles, "source", webFiles)
	mg.WriteTo("target", webFiles)
}

func getFilesAt(root string, exclusions ...string) []string {
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

	mg.ExitIfError(&err, 2)

	return files
}
