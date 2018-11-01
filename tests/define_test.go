package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestDefineNumber(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define a 2 }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["a"] = float64(2)

	checkParsing(t, processed.Context(), emptyFilesMap, processed, expectedCtx, []string{""})
}

func TestDefineString(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define title \"My Site\" }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["title"] = "My Site"

	checkParsing(t, processed.Context(), emptyFilesMap, processed, expectedCtx, []string{""})
}

func TestDefineStringConcat(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define title \"My\" + \" Site\" }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["title"] = "My Site"

	checkParsing(t, processed.Context(), emptyFilesMap, processed, expectedCtx, []string{""})
}

func TestDefineFromExpression(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define n 10 + 10 * (2 + 4) }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["n"] = float64(70)

	checkParsing(t, processed.Context(), emptyFilesMap, processed, expectedCtx, []string{""})
}

func TestDefineBasedOnPreviousDefine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(
		"{{ define a 10 }}" +
			"{{ define b 4 }}" +
			"{{ define c a * b }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.html", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["a"] = float64(10)
	expectedCtx["b"] = float64(4)
	expectedCtx["c"] = float64(40)

	checkParsing(t, processed.Context(), emptyFilesMap, processed, expectedCtx, []string{"", "", ""})
}

func TestMalformedDefine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.html", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{"{{ define }}"})

	r = bufio.NewReader(strings.NewReader("{{ define abc }}"))
	processed, err = mg.ProcessReader(r, "source/processed/hi.html", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{"{{ define abc }}"})
}
