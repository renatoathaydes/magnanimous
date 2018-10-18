package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestEvalString(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Hello {{ eval \"Joe\" }}"))
	processed, err := mg.ProcessReader(r, "", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{"Hello ", "Joe"})
}

func TestEvalArithmetic(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("2 + 2 == {{ eval 2 + 2 }}"))
	processed, err := mg.ProcessReader(r, "", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{"2 + 2 == ", "4"})
}

func TestEvalNonExistingParameter(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ eval 2 * a }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.html", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{"{{ eval 2 * a }}"})
}

func TestEvalWithExistingParameter(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define a 3 }}{{ eval 2 * a }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["a"] = float64(3)

	files := mg.WebFilesMap{}
	files["source/processed/hi.md"] = mg.WebFile{Processed: processed}

	checkParsing(t, processed.Context(), files, processed, expectedCtx, []string{"<p>6</p>\n"})
}

func TestEvalWithExistingParameterFromAnotherFile(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("A = {{ eval 2 * hello }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11)

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("OUTER\n{{ define hello 7 }}{{ include /processed/hi.txt }}\nEND"))
	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11)

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files := mg.WebFilesMap{}
	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})
	expectedCtx["hello"] = float64(7)

	checkParsing(t, otherProcessed.Context(), files, otherProcessed, expectedCtx, []string{
		"OUTER\n",
		"A = 14",
		"\nEND"})
}
