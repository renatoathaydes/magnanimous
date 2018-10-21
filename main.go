package main

import (
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

	err = mg.WriteTo(TargetDir, webFiles)
	if err != nil {
		log.Printf("ERROR: %s", err)
		panic(err)
	}

	log.Printf("Magnanimous generated website in %s\n", time.Since(start))
}
