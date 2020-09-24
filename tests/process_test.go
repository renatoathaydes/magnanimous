package tests

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/renatoathaydes/magnanimous/mg"
)

var mdLocation = mg.Location{Origin: "a.md", Row: 0, Col: 0}
var nonMdLocation = mg.Location{Origin: "a.txt", Row: 0, Col: 0}

func TestMarkdownEngineAlwaysMakesTheSameThing(t *testing.T) {
	content := mg.NewStringContent("# hello\n## world", &mdLocation)
	pf := mg.ProcessedFile{NewExtension: ".html", Path: mdLocation.Origin}
	pf.AppendContent(content)
	wf := mg.WebFile{Processed: &pf, Name: "hello.md"}

	writeContentToString := func() string {
		stack := mg.NewContextStack(mg.NewContext())
		var w strings.Builder
		err := wf.Write(&w, &stack, false)
		check(err)
		return w.String()
	}

	expectedResult := "<h1>hello</h1>\n\n<h2>world</h2>\n"

	for i := 0; i < 10; i++ {
		r := writeContentToString()
		if r != expectedResult {
			t.Errorf("[%d] '%s' != '%s'", i, string(r), expectedResult)
		}
	}
}

func TestProcessSimple(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello world"))
	processed, err := mg.ProcessReader(r, "", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, "hello world")
}

func TestProcessIncludeSimple(t *testing.T) {

	exampleFile := mg.ProcessedFile{}
	content := mg.NewStringContent("from another file!", &nonMdLocation)
	exampleFile.AppendContent(content)

	m := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	m.WebFiles["source/processed/example.html"] = mg.WebFile{Processed: &exampleFile}

	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &m}

	r := bufio.NewReader(strings.NewReader("hello {{ include example.html }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hello.html", "source", 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, "hello from another file!")
}

func TestMarkdownToHtml(t *testing.T) {
	header := mg.ProcessedFile{Path: "header.html"}
	header.AppendContent(mg.NewStringContent("<html><body class=\"hello\">\n", &nonMdLocation))
	footer := mg.ProcessedFile{Path: "footer.html"}
	footer.AppendContent(mg.NewStringContent("</body></html>", &nonMdLocation))

	m := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	m.WebFiles["header.html"] = mg.WebFile{Processed: &header, Name: "header.html"}
	m.WebFiles["footer.html"] = mg.WebFile{Processed: &footer, Name: "footer.html"}

	rsvr := mg.DefaultFileResolver{BasePath: "", Files: &m}

	file := mg.ProcessedFile{Path: "test_file.md"}
	file.AppendContent(mg.NewIncludeInstruction("header.html", &mdLocation, "", &rsvr))
	file.AppendContent(mg.NewStringContent("# Hello\n", &mdLocation))
	file.AppendContent(mg.NewStringContent("## Mag", &mdLocation))
	file.AppendContent(mg.NewIncludeInstruction("footer.html", &mdLocation, "", &rsvr))

	checkParsing(t, &file, emptyContext,
		"<html><body class=\"hello\">\n"+
			"<h1>Hello</h1>\n\n<h2>Mag</h2>\n"+
			"</body></html>")
}

func TestProcessIncludeMarkDown(t *testing.T) {
	// setup a md file to be included
	exampleFile := mg.ProcessedFile{Path: "example.md", NewExtension: ".html"}
	exampleFile.AppendContent(mg.NewStringContent("## header", &mdLocation))

	m := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	m.WebFiles["source/example.md"] = mg.WebFile{Processed: &exampleFile, Name: "example.md"}

	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &m}

	// read a file that includes the previous file
	r := bufio.NewReader(strings.NewReader("# hello\n{{ include /example.md }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", "source", 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, "<h1>hello</h1>\n\n<h2>header</h2>\n")
}

func TestProcessIgnoreEscapedBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello \\{{ include example.html }}"))
	processed, err := mg.ProcessReader(r, "", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, "hello {{ include example.html }}")
}

func TestProcessIgnoreEscapedClosingBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Hello {{\n  bad-instruction \"contains \\}} ignored\"\n}}. How are you?"))
	processed, err := mg.ProcessReader(r, "", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext,
		"Hello "+
			"{{\n  bad-instruction \"contains }} ignored\"\n}}"+
			". How are you?")
}

func TestProcessIgnoreEscapedNewLine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello {{ eval \"Joe\" }}\\\n, how are you\\\n, good?"))
	processed, err := mg.ProcessReader(r, "", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, "hello Joe, how are you, good?")
}

func TestProcessDoNotIgnoreEscapedEscapedNewLine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello {{ eval \"Joe\" }}\\\\\n, how are you\\\\\n, good?"))
	processed, err := mg.ProcessReader(r, "", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext,
		"hello Joe\\\n, how are you\\\n, good?")
}

func TestProcessIgnoreEscapedWindowsNewLine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello {{ eval \"Joe\" }}\\\r\n, how are you\\\r\n, good?\r\nKeep line."))
	processed, err := mg.ProcessReader(r, "", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext,
		"hello Joe, how are you, good?\r\nKeep line.")
}

func TestProcessDoc(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello{{ doc this is ignored content }}{{ doc\n And this too\n}}"))
	processed, err := mg.ProcessReader(r, "", "source", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, "hello")
}
