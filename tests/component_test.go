package tests

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/renatoathaydes/magnanimous/mg"
)

func TestIncludeComponent(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("Hi {{ eval slot1 }}! Your contents: {{ eval __contents__ }}."))
	processed, err := mg.ProcessReader(r, "source/processed/_comp.txt", "source", 6, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("OUTER\n{{ component /processed/_comp.txt }}" +
		"{{ slot slot1 }}Joe{{ end}}Four is {{ eval 2 + 2 }}" +
		"{{end}}\nEND"))
	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", "source", 11, &resolver, time.Now())

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/processed/_comp.txt"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})

	checkParsing(t, otherProcessed, expectedCtx, "OUTER\nHi Joe! Your contents: Four is 4.\nEND")
}
