package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestMarkdownEngineAlwaysMakesTheSameThing(t *testing.T) {
	contents := []mg.Content{&mg.StringContent{Text: "# hello\n## world"}}

	writeContentToString := func(content mg.Content) string {
		var w strings.Builder
		content.Write(&w, nil, nil)
		return w.String()
	}

	expectedResult := "<h1>hello</h1>\n\n<h2>world</h2>\n"

	for i := 0; i < 10; i++ {
		markdown := mg.HtmlFromMarkdownContent{MarkDownContent: contents}
		r := writeContentToString(&markdown)
		if r != expectedResult {
			t.Errorf("[%d] '%s' != %s", i, string(r), expectedResult)
		}
	}
}

func TestProcessSimple(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello world"))
	processed, err := mg.ProcessReader(r, "", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{"hello world"})
}

func TestProcessIncludeSimple(t *testing.T) {

	exampleFile := mg.ProcessedFile{}
	exampleFile.AppendContent(&mg.StringContent{Text: "from another file!"})

	m := mg.WebFilesMap{}
	m["source/processed/example.html"] = mg.WebFile{Processed: &exampleFile}

	resolver := mg.DefaultFileResolver{BasePath: "source", Files: m}

	r := bufio.NewReader(strings.NewReader("hello {{ include example.html }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hello.html", 11, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), m, processed, emptyContext, []string{"hello ", "from another file!"})
}

func TestMarkdownToHtml(t *testing.T) {
	header := mg.ProcessedFile{}
	header.AppendContent(&mg.StringContent{Text: "<html><body class=\"hello\">\n"})
	footer := mg.ProcessedFile{}
	footer.AppendContent(&mg.StringContent{Text: "</body></html>"})

	m := mg.WebFilesMap{}
	m["header.html"] = mg.WebFile{Processed: &header, Name: "header.html"}
	m["footer.html"] = mg.WebFile{Processed: &footer, Name: "footer.html"}

	rsvr := mg.DefaultFileResolver{BasePath: "", Files: m}

	file := mg.ProcessedFile{}
	file.AppendContent(&mg.IncludeInstruction{Path: "header.html", Resolver: &rsvr})
	file.AppendContent(&mg.StringContent{Text: "# Hello\n"})
	file.AppendContent(&mg.StringContent{Text: "## Mag"})
	file.AppendContent(&mg.IncludeInstruction{Path: "footer.html", Resolver: &rsvr})

	html := mg.MarkdownToHtml(file)

	result, err := html.Bytes(m, nil)

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

	m := mg.WebFilesMap{}
	m["source/example.md"] = mg.WebFile{Processed: &exampleFile}

	resolver := mg.DefaultFileResolver{BasePath: "source", Files: m}

	// read a file that includes the previous file
	r := bufio.NewReader(strings.NewReader("# hello\n{{ include /example.md }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), m, processed, emptyContext, []string{"<h1>hello</h1>\n\n<h2>header</h2>\n"})
}

func TestProcessIgnoreEscapedBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello \\{{ include example.html }}"))
	processed, err := mg.ProcessReader(r, "", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{"hello {{ include example.html }}"})
}

func TestProcessIgnoreEscapedClosingBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Hello {{\n  bad-instruction \"contains \\}} ignored\"\n}}. How are you?"))
	processed, err := mg.ProcessReader(r, "", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{
		"Hello ",
		"{{\n  bad-instruction \"contains }} ignored\"\n}}",
		". How are you?",
	})
}
