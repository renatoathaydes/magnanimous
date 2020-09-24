package tests

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/renatoathaydes/magnanimous/mg"
)

func TestIncludeB64File(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("Hello World"))
	processed, err := mg.ProcessReader(r, "source/processed/hello.txt", "source", 6, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("Base64({{ includeB64 /processed/hello.txt }})"))
	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", "source", 11, &resolver, time.Now())

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/processed/hello.txt"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})

	checkParsing(t, otherProcessed, expectedCtx,
		"Base64(SGVsbG8gV29ybGQ=)")
}
