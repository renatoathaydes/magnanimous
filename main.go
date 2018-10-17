package main

import (
	"github.com/renatoathaydes/magnanimous/mg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type filesCollector func() ([]string, error)

const (
	ProcessedDir = "source/processed/"
	StaticDir    = "source/static/"
	SourceDir    = "source/"
	TargetDir    = "target/"
)

func main() {
	start := time.Now()

	c1, c2, c3 := make(chan []string), make(chan []string), make(chan []string)
	go async(func() ([]string, error) { return getFilesAt(ProcessedDir) }, c1)
	go async(func() ([]string, error) { return getFilesAt(StaticDir) }, c2)
	go async(func() ([]string, error) {
		return getFilesAt(SourceDir, ProcessedDir, StaticDir)
	}, c3)

	procFiles, staticFiles, otherFiles := <-c1, <-c2, <-c3

	webFiles := make(mg.WebFilesMap, len(procFiles)+len(staticFiles)+len(otherFiles))
	mg.Process(&procFiles, ProcessedDir, webFiles)
	mg.CopyAll(&staticFiles, StaticDir, webFiles)
	mg.AddNonWritables(&otherFiles, SourceDir, webFiles)
	err := mg.WriteTo(TargetDir, webFiles)
	if err != nil {
		log.Printf("ERROR: (%s) %s", err.Code, err.Error())
		panic(*err)
	}

	log.Printf("Magnanimous generated website in %s\n", time.Since(start))
}

func async(fc filesCollector, c chan []string) {
	s, err := fc()
	if err != nil {
		panic(err)
	}
	c <- s
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
