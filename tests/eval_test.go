package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestEvalString(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Hello {{ eval \"Joe\" }}"))
	ctx, processed, err := mg.ProcessReader(r, "", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, ctx, emptyFilesMap, processed, emptyContext, []string{"Hello ", "Joe"})
}

func TestEvalArithmetic(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("2 + 2 == {{ eval 2 + 2 }}"))
	ctx, processed, err := mg.ProcessReader(r, "", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, ctx, emptyFilesMap, processed, emptyContext, []string{"2 + 2 == ", "4"})
}

func TestEvalNonExistingParameter(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ eval 2 * a }}"))
	ctx, processed, err := mg.ProcessReader(r, "source/processed/hi.html", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, ctx, emptyFilesMap, processed, emptyContext, []string{"{{ eval 2 * a }}"})
}

func TestEvalWithExistingParameter(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define a 3 }}{{ eval 2 * a }}"))
	ctx, processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := mg.WebFileContext{}
	expectedCtx["a"] = float64(3)

	files := mg.WebFilesMap{}
	files["source/processed/hi.md"] = mg.WebFile{Context: ctx, Processed: processed}

	checkParsing(t, ctx, files, processed, expectedCtx, []string{"<p>6</p>\n"})
}

func TestEvalWithExistingParameterFromAnotherFile(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("A = {{ eval 2 * hello }}"))
	_, processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11)

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("OUTER\n{{ define hello 7 }}{{ include /processed/hi.txt }}\nEND"))
	otherCtx, otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11)

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files := mg.WebFilesMap{}
	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed, Context: emptyContext}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed, Context: otherCtx}

	expectedCtx := mg.WebFileContext{}
	expectedCtx["hello"] = float64(7)

	checkParsing(t, otherCtx, files, otherProcessed, expectedCtx, []string{
		"OUTER\n",
		"A = 14",
		"\nEND"})
}
