package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
	"time"
)

func TestIncludeFile(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("ABCDEF"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 6, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("OUTER\n{{ include /processed/hi.txt }}\nEND"))
	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11, &resolver, time.Now())

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})

	checkParsing(t, otherProcessed, expectedCtx, []string{
		"OUTER\n",
		"ABCDEF",
		"\nEND"})
}

func TestIncludeEvalFile(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("ABCDEF"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 6, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	// let 'p' and the appended String contain `/` to cause common double-slash problem (magnanimous should accept it)
	r = bufio.NewReader(strings.NewReader("{{define p `/processed/`}}" +
		"OUTER\n{{ include eval p + `/hi.txt` }}\nEND"))

	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11, &resolver, time.Now())

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})
	expectedCtx["p"] = "/processed/"

	checkParsing(t, otherProcessed, expectedCtx, []string{
		"",
		"OUTER\n",
		"ABCDEF",
		"\nEND"})
}

func TestIncludeImplicitEvalFile(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("ABCDEF"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 6, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("{{define p `/processed`}}" +
		"OUTER\n{{ include \"/processed\" + `/hi.txt` }}\nEND"))

	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11, &resolver, time.Now())

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})
	expectedCtx["p"] = "/processed"

	checkParsing(t, otherProcessed, expectedCtx, []string{
		"",
		"OUTER\n",
		"ABCDEF",
		"\nEND"})
}

func TestIncludeFileNested(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("A1"))
	processed1, err := mg.ProcessReader(r, "source/a1.txt", 6, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("A2\n{{include /a1.txt}}"))
	processed2, err := mg.ProcessReader(r, "source/processed/a2.txt", 6, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	// include relative path
	r = bufio.NewReader(strings.NewReader("A3\n{{ include a2.txt }}\nEND"))
	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", 11, &resolver, time.Now())

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/a1.txt"] = mg.WebFile{Processed: processed1}
	files["source/processed/a2.txt"] = mg.WebFile{Processed: processed2}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})

	checkParsing(t, otherProcessed, expectedCtx, []string{
		"A3\n",
		"A2\nA1",
		"\nEND"})
}

func TestIncludeUpPath(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("A1{{include .../_msg}}A1"))
	processed1, err := mg.ProcessReader(r, "source/processed/en/a1.txt", 12, &resolver, time.Now())
	check(err)

	r = bufio.NewReader(strings.NewReader("A2{{include .../_msg}}A2"))
	processed2, err := mg.ProcessReader(r, "source/processed/pt/abc/a2.txt", 12, &resolver, time.Now())
	check(err)

	r = bufio.NewReader(strings.NewReader("English"))
	english, err := mg.ProcessReader(r, "source/_msg", 7, &resolver, time.Now())
	check(err)

	r = bufio.NewReader(strings.NewReader("Portuguese"))
	portuguese, err := mg.ProcessReader(r, "source/processed/pt/_msg", 7, &resolver, time.Now())
	check(err)

	files["source/processed/en/a1.txt"] = mg.WebFile{Processed: processed1}
	files["source/processed/pt/abc/a2.txt"] = mg.WebFile{Processed: processed2}
	files["source/_msg"] = mg.WebFile{Processed: english}
	files["source/processed/pt/_msg"] = mg.WebFile{Processed: portuguese}

	expectedCtx := make(map[string]interface{})

	checkParsing(t, processed1, expectedCtx, []string{
		"A1",
		"English",
		"A1"})
	checkParsing(t, processed2, expectedCtx, []string{
		"A2",
		"Portuguese",
		"A2"})
}
