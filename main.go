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

	rootDir, globalCtx, ok := parseOptions()

	if !ok {
		return
	}

	mag := mg.Magnanimous{SourcesDir: filepath.Join(*rootDir, SourceDir), GlobalContex: *globalCtx}
	webFiles, err := mag.ReadAll()
	if err != nil {
		log.Printf("ERROR: %s", err)
		panic(err)
	}

	if len(webFiles.WebFiles) == 0 {
		fmt.Printf("No files found in the %s directory, nothing to do.\n", mag.SourcesDir)
		return
	}

	err = mag.WriteTo(filepath.Join(*rootDir, TargetDir), webFiles)
	if err != nil {
		log.Printf("ERROR: %s", err)
		panic(err)
	}

	log.Printf("Magnanimous generated website in %s\n", time.Since(start))
}

func parseOptions() (rootDir, globalContext *string, ok bool) {
	globalContext = flag.String("globalctx", "",
		"Path to the global context file relative to the 'processed' directory.")
	style := flag.String("style", "lovelace",
		"Style name for code highlighting. See https://xyproto.github.io/splash/docs/all.html.")

	help := flag.Bool("help", false, "Print usage help.")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n\n  %s [options...] [root-directory]\n\nOptions:\n",
			os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	mg.SetCodeStyle(*style)

	otherArgs := flag.Args()

	if *help {
		flag.Usage()
		return nil, nil, false
	}

	var rootDirValue string
	switch len(otherArgs) {
	case 0:
		rootDirValue = ""
	case 1:
		rootDirValue = otherArgs[0]
	default:
		log.Printf("ERROR: too many arguments provided")
		flag.Usage()
		return nil, nil, false
	}
	rootDir = &rootDirValue

	ok = true

	return
}
