package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestDefineNumber(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define a 2 }}"))
	ctx, processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := mg.WebFileContext{}
	expectedCtx["a"] = float64(2)

	checkParsing(t, ctx, emptyFilesMap, processed, expectedCtx, []string{})
}

func TestDefineString(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define title \"My Site\" }}"))
	ctx, processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := mg.WebFileContext{}
	expectedCtx["title"] = "My Site"

	checkParsing(t, ctx, emptyFilesMap, processed, expectedCtx, []string{})
}

func TestDefineStringConcat(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define title \"My\" + \" Site\" }}"))
	ctx, processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := mg.WebFileContext{}
	expectedCtx["title"] = "My Site"

	checkParsing(t, ctx, emptyFilesMap, processed, expectedCtx, []string{})
}

func TestDefineBasedOnPreviousDefine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(
		"{{ define a 10 }}" +
			"{{ define b 4 }}" +
			"{{ define c a * b }}"))
	ctx, processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := mg.WebFileContext{}
	expectedCtx["a"] = float64(10)
	expectedCtx["b"] = float64(4)
	expectedCtx["c"] = float64(40)

	checkParsing(t, ctx, emptyFilesMap, processed, expectedCtx, []string{})
}

func TestMalformedDefine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define }}"))
	ctx, processed, err := mg.ProcessReader(r, "source/processed/hi.html", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, ctx, emptyFilesMap, processed, emptyContext, []string{"{{ define }}"})

	r = bufio.NewReader(strings.NewReader("{{ define abc }}"))
	ctx, processed, err = mg.ProcessReader(r, "source/processed/hi.html", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, ctx, emptyFilesMap, processed, emptyContext, []string{"{{ define abc }}"})
}
