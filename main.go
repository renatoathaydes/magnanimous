package main

import (
	"flag"
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	SourceDir = "source/"
	TargetDir = "target/"
)

func main() {
	start := time.Now()

	var rootDir string
	switch len(os.Args) {
	case 0:
		fallthrough
	case 1:
		rootDir = ""
	case 2:
		rootDir = os.Args[1]
	default:
		log.Printf("ERROR: too many arguments provided")
	}

	globalCtx := parseOptions()

	mag := mg.Magnanimous{SourcesDir: filepath.Join(rootDir, SourceDir), GlobalContex: *globalCtx}
	webFiles, err := mag.ReadAll()
	if err != nil {
		log.Printf("ERROR: %s", err)
		panic(err)
	}

	if len(webFiles.WebFiles) == 0 {
		fmt.Printf("No files found in the %s directory, nothing to do.\n", mag.SourcesDir)
		return
	}

	err = mag.WriteTo(filepath.Join(rootDir, TargetDir), webFiles)
	if err != nil {
		log.Printf("ERROR: %s", err)
		panic(err)
	}

	log.Printf("Magnanimous generated website in %s\n", time.Since(start))
}

func parseOptions() (globalContext *string) {
	globalContext = flag.String("globalctx", "",
		"Path to the global context file relative to the processed/ dir")

	flag.Parse()

	return
}
