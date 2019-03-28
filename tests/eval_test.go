package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestEvalString(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Hello {{ eval \"Joe\" }}"))
	processed, err := mg.ProcessReader(r, "", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, emptyFilesMap, processed, emptyContext, []string{"Hello ", "Joe"})
}

func TestEvalDate(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Time: {{ eval date[\"2017-11-23T22:12:21\"] }}"))
	processed, err := mg.ProcessReader(r, "", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, emptyFilesMap, processed, emptyContext, []string{"Time: ", "23 Nov 2017, 10:12 PM"})
}

func TestEvalArithmetic(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("2 + 2 * 5 == {{ eval 2 + 2 * 5 }}"))
	processed, err := mg.ProcessReader(r, "", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, emptyFilesMap, processed, emptyContext, []string{"2 + 2 * 5 == ", "12"})
}

func TestEvalNonExistingParameter(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ eval 2 * a }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.html", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, emptyFilesMap, processed, emptyContext, []string{"{{ eval 2 * a }}"})
}

func TestEvalWithExistingParameter(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define a 3 }}{{ eval 2 * a }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["a"] = float64(3)

	files := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	files.WebFiles["source/processed/hi.md"] = mg.WebFile{Processed: processed}

	checkParsing(t, files, processed, expectedCtx, []string{"<p>6</p>\n"})
}

func TestEvalWithOrExprParameterExists(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define a 3 }}{{ eval a || 100 }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["a"] = float64(3)

	files := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	files.WebFiles["source/processed/hi.md"] = mg.WebFile{Processed: processed}

	checkParsing(t, files, processed, expectedCtx, []string{"<p>3</p>\n"})
}

func TestEvalWithOrExprParameterDoesNotExist(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ eval a || 100 }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})

	files := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	files.WebFiles["source/processed/hi.md"] = mg.WebFile{Processed: processed}

	checkParsing(t, files, processed, expectedCtx, []string{"<p>100</p>\n"})
}

func TestEvalWithExistingParameterFromAnotherFile(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("A = {{ eval 2 * hello }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("OUTER\n{{ define hello 7 }}{{ include /processed/hi.txt }}\nEND"))
	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11, &resolver)

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})
	expectedCtx["hello"] = float64(7)

	checkParsing(t, *resolver.Files, otherProcessed, expectedCtx, []string{
		"OUTER\n",
		"",
		"A = 14",
		"\nEND"})
}
