package main

import (
	"github.com/renatoathaydes/magnanimous/mg"
	"os"
	"path/filepath"
)

func main() {
	procFiles := getFilesAt("source/processed")
	staticFiles := getFilesAt("source/static")
	webFiles := make(mg.WebFilesMap, len(procFiles)+len(staticFiles))
	mg.Process(&procFiles, "source/processed", &webFiles)

	// TODO just copy static files
	mg.Process(&staticFiles, "source/static", &webFiles)

	mg.WriteAt("target", &webFiles)
}

func getFilesAt(root string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			files = append(files, path)
		}
		return err
	})

	mg.ExitIfError(&err, 2)

	return files
}
