package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
	"time"
)

func TestMarkdownEngineAlwaysMakesTheSameThing(t *testing.T) {
	contents := []mg.Content{&mg.StringContent{Text: "# hello\n## world"}}

	writeContentToString := func(content mg.Content) string {
		stack := mg.NewContextStack(mg.NewContext())
		var w strings.Builder
		err := content.Write(&w, stack)
		check(err)
		return w.String()
	}

	expectedResult := "<h1>hello</h1>\n\n<h2>world</h2>\n"

	for i := 0; i < 10; i++ {
		markdown := mg.HtmlFromMarkdownContent{MarkDownContent: contents}
		r := writeContentToString(&markdown)
		if r != expectedResult {
			t.Errorf("[%d] '%s' != '%s'", i, string(r), expectedResult)
		}
	}
}

func TestProcessSimple(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello world"))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{"hello world"})
}

func TestProcessIncludeSimple(t *testing.T) {

	exampleFile := mg.ProcessedFile{}
	exampleFile.AppendContent(&mg.StringContent{Text: "from another file!"})

	m := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	m.WebFiles["source/processed/example.html"] = mg.WebFile{Processed: &exampleFile}

	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &m}

	r := bufio.NewReader(strings.NewReader("hello {{ include example.html }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hello.html", 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{"hello ", "from another file!"})
}

func TestMarkdownToHtml(t *testing.T) {
	header := mg.ProcessedFile{Path: "header"}
	header.AppendContent(&mg.StringContent{Text: "<html><body class=\"hello\">\n"})
	footer := mg.ProcessedFile{Path: "footer"}
	footer.AppendContent(&mg.StringContent{Text: "</body></html>"})

	m := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	m.WebFiles["header.html"] = mg.WebFile{Processed: &header, Name: "header.html"}
	m.WebFiles["footer.html"] = mg.WebFile{Processed: &footer, Name: "footer.html"}

	rsvr := mg.DefaultFileResolver{BasePath: "", Files: &m}

	loc := mg.Location{}
	file := mg.ProcessedFile{Path: "test_file"}
	file.AppendContent(mg.NewIncludeInstruction("header.html", &loc, "", &rsvr, true))
	file.AppendContent(&mg.StringContent{Text: "# Hello\n"})
	file.AppendContent(&mg.StringContent{Text: "## Mag"})
	file.AppendContent(mg.NewIncludeInstruction("footer.html", &loc, "", &rsvr, true))

	html := mg.MarkdownToHtml(file)

	stack := mg.NewContextStack(mg.NewContext())
	result, err := html.Bytes(stack)

	if err != nil {
		t.Fatal(err)
	}

	expectedHtml := "<html><body class=\"hello\">\n<h1>Hello</h1>\n\n<h2>Mag</h2>\n</body></html>"
	if string(result) != expectedHtml {
		t.Errorf("Expected '%s', but was '%s'", expectedHtml, result)
	}
}

func TestProcessIncludeMarkDown(t *testing.T) {
	// setup a md file to be included
	exampleFile := mg.ProcessedFile{}
	exampleFile.AppendContent(&mg.HtmlFromMarkdownContent{
		MarkDownContent: []mg.Content{&mg.StringContent{Text: "## header"}},
	})

	m := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	m.WebFiles["source/example.md"] = mg.WebFile{Processed: &exampleFile}

	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &m}

	// read a file that includes the previous file
	r := bufio.NewReader(strings.NewReader("# hello\n{{ include /example.md }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{"<h1>hello</h1>\n\n<h2>header</h2>\n"})
}

func TestProcessIgnoreEscapedBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello \\{{ include example.html }}"))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{"hello {{ include example.html }}"})
}

func TestProcessIgnoreEscapedClosingBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Hello {{\n  bad-instruction \"contains \\}} ignored\"\n}}. How are you?"))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{
		"Hello ",
		"{{\n  bad-instruction \"contains }} ignored\"\n}}",
		". How are you?",
	})
}

func TestProcessIgnoreEscapedNewLine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello {{ eval \"Joe\" }}\\\n, how are you\\\n, good?"))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext,
		[]string{"hello ", "Joe", ", how are you, good?"})
}

func TestProcessDoNotIgnoreEscapedEscapedNewLine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello {{ eval \"Joe\" }}\\\\\n, how are you\\\\\n, good?"))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext,
		[]string{"hello ", "Joe", "\\\n, how are you\\\n, good?"})
}

func TestProcessIgnoreEscapedWindowsNewLine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello {{ eval \"Joe\" }}\\\r\n, how are you\\\r\n, good?\r\nKeep line."))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext,
		[]string{"hello ", "Joe", ", how are you, good?\r\nKeep line."})
}

func TestProcessDoc(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello{{ doc this is ignored content }}{{ doc\n And this too\n}}"))
	processed, err := mg.ProcessReader(r, "", 11, nil, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed, emptyContext, []string{"hello"})
}
