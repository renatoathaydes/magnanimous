package main

import (
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg"
	"log"
	"time"
)

const (
	SourceDir = "source/"
	TargetDir = "target/"
)

func main() {
	start := time.Now()
	mag := mg.Magnanimous{SourcesDir: SourceDir}
	webFiles, err := mag.ReadAll()
	if err != nil {
		log.Printf("ERROR: %s", err)
		panic(err)
	}

	if len(webFiles) == 0 {
		fmt.Printf("No files found in the %s directory, nothing to do.\n", SourceDir)
		return
	}

	err = mg.WriteTo(TargetDir, webFiles)
	if err != nil {
		log.Printf("ERROR: %s", err)
		panic(err)
	}

	log.Printf("Magnanimous generated website in %s\n", time.Since(start))
}
