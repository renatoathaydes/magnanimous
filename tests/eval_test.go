package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"strings"
	"testing"
	"time"
)

func TestEvalString(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Hello {{ eval \"Joe\" }}"))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{"Hello ", "Joe"})
}

func TestEvalArray(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Numbers: {{ eval [ 1, 2, 3, 4 ] }}"))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{"Numbers: ", "[1 2 3 4]"})
}

func TestEvalDate(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Time: {{ eval date[\"2017-11-23T22:12:21\"] }}"))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{"Time: ", "23 Nov 2017, 10:12 PM"})
}

func TestEvalDateCustom(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Time: {{ eval date[\"2017-11-23T22:12:21\"][\"15:04:05 on 02 January 2006\"] }}"))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{"Time: ", "22:12:21 on 23 November 2017"})
}

func TestEvalDateOfFileUpdate(t *testing.T) {
	update := time.Date(1992, 12, 19, 8, 30, 0, 0, time.UTC)

	fileMap := make(map[string]mg.WebFile)
	fileMap["other.file"] = mg.WebFile{Processed: &mg.ProcessedFile{LastUpdated: update}}
	files := mg.WebFilesMap{WebFiles: fileMap}
	resolver := mg.DefaultFileResolver{Files: &files}

	r := bufio.NewReader(strings.NewReader("File updated on " +
		"{{define examplePath path[\"other.file\"]}}{{ eval date[examplePath][\"15:04:05 on 02 January 2006\"] }}"))
	processed, err := mg.ProcessReader(r, "", 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["examplePath"] = &expression.Path{Value: "other.file", LastUpdated: update}

	checkParsing(t, processed, expectedCtx, []string{"File updated on ", "", "08:30:00 on 19 December 1992"})
}

func TestEvalArithmetic(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("2 + 2 * 5 == {{ eval 2 + 2 * 5 }}"))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{"2 + 2 * 5 == ", "12"})
}

func TestEvalNonExistingParameter(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ eval 2 * a }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.html", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{"{{ eval 2 * a }}"})
}

func TestEvalWithExistingParameter(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define a 3 }}{{ eval 2 * a }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["a"] = float64(3)

	files := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	files.WebFiles["source/processed/hi.md"] = mg.WebFile{Processed: processed}

	checkParsing(t, processed, expectedCtx, []string{"<p>6</p>\n"})
}

func TestEvalWithOrExprParameterExists(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define a 3 }}{{ eval a || 100 }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["a"] = float64(3)

	files := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	files.WebFiles["source/processed/hi.md"] = mg.WebFile{Processed: processed}

	checkParsing(t, processed, expectedCtx, []string{"<p>3</p>\n"})
}

func TestEvalWithOrExprParameterDoesNotExist(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ eval a || 100 }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})

	files := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	files.WebFiles["source/processed/hi.md"] = mg.WebFile{Processed: processed}

	checkParsing(t, processed, expectedCtx, []string{"<p>100</p>\n"})
}

func TestEvalWithExistingParameterFromAnotherFile(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("A = {{ eval 2 * hello }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("OUTER\n{{ define hello 7 }}{{ include /processed/hi.txt }}\nEND"))
	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11, &resolver, time.Now())

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})
	expectedCtx["hello"] = float64(7)

	checkParsing(t, otherProcessed, expectedCtx, []string{
		"OUTER\n",
		"",
		"A = 14",
		"\nEND"})
}
