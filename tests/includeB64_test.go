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

func TestIncludeB64MdFile(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("# Hello World"))
	processed, err := mg.ProcessReader(r, "source/processed/hello.md", "source", 6, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	r = bufio.NewReader(strings.NewReader("Base64({{ includeB64 /processed/hello.md }})"))
	otherProcessed, otherErr := mg.ProcessReader(r, "source/processed/other.txt", "source", 11, &resolver, time.Now())

	if otherErr != nil {
		t.Fatal(otherErr)
	}

	files["source/processed/hello.md"] = mg.WebFile{Processed: processed}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}

	expectedCtx := make(map[string]interface{})

	checkParsing(t, otherProcessed, expectedCtx,
		"Base64(PGgxPkhlbGxvIFdvcmxkPC9oMT4K)")
}

func TestIncludeB64MdFilePlain(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("# Hello World"))
	helloMd, err := mg.ProcessReader(r, "source/processed/hello.md", "source", 6, &resolver, time.Now())
	check(err)

	r = bufio.NewReader(strings.NewReader("## Foo"))
	fooMd, err := mg.ProcessReader(r, "source/processed/foo.md", "source", 6, &resolver, time.Now())
	check(err)

	r = bufio.NewReader(strings.NewReader("<div>{{ include /processed/foo.md }}</div>"))
	aHtml, err := mg.ProcessReader(r, "source/processed/a.html", "source", 6, &resolver, time.Now())
	check(err)

	r = bufio.NewReader(strings.NewReader("Base64({{ includeB64 ( plain ) /processed/hello.md }})"))
	otherProcessed, err := mg.ProcessReader(r, "source/processed/other.txt", "source", 11, &resolver, time.Now())
	check(err)

	r = bufio.NewReader(strings.NewReader("{{ includeB64 (plain) /processed/a.html }}"))
	aB64, err := mg.ProcessReader(r, "source/processed/a.b64", "source", 11, &resolver, time.Now())
	check(err)

	files["source/processed/hello.md"] = mg.WebFile{Processed: helloMd}
	files["source/processed/foo.md"] = mg.WebFile{Processed: fooMd}
	files["source/processed/a.html"] = mg.WebFile{Processed: aHtml}
	files["source/processed/other.txt"] = mg.WebFile{Processed: otherProcessed}
	files["source/processed/a.b64"] = mg.WebFile{Processed: aB64}

	expectedCtx := make(map[string]interface{})

	checkParsing(t, otherProcessed, expectedCtx,
		"Base64(IyBIZWxsbyBXb3JsZA==)")

	checkParsing(t, aHtml, expectedCtx,
		"<div><h2>Foo</h2>\n</div>")

	checkParsing(t, aB64, expectedCtx,
		"PGRpdj4jIyBGb288L2Rpdj4=")

}
