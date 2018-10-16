package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"reflect"
	"strings"
	"testing"
)

var emptyContext = make(map[string]interface{})
var emptyFilesMap = mg.WebFilesMap{}

func TestProcessSimple(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello world"))
	processed, err := mg.ProcessReader(r, "", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{"hello world"})
}

func TestProcessIncludeSimple(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello {{ include example.html }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hello.html", 11)

	if err != nil {
		t.Fatal(err)
	}

	exampleFile := mg.ProcessedFile{}
	exampleFile.AppendContent(&mg.StringContent{Text: "from another file!"})

	m := mg.WebFilesMap{}
	m["source/processed/example.html"] = mg.WebFile{Processed: &exampleFile}

	checkParsing(t, processed.Context(), m, processed, emptyContext, []string{"hello ", "from another file!"})
}

func TestMarkdownToHtml(t *testing.T) {
	file := mg.ProcessedFile{}
	file.AppendContent(&mg.StringContent{Text: "<html><body class=\"hello\">"})
	file.AppendContent(&mg.StringContent{Text: "# Hello\n", MarkDown: true})
	file.AppendContent(&mg.StringContent{Text: "## Mag", MarkDown: true})
	file.AppendContent(&mg.StringContent{Text: "</body></html>"})

	html := mg.MarkdownToHtml(file)

	m := mg.WebFilesMap{}
	result := string(html.Bytes(m, nil))

	expectedHtml := "<html><body class=\"hello\"><h1>Hello</h1>\n<h2>Mag</h2>\n</body></html>"
	if result != expectedHtml {
		t.Errorf("Expected '%s', but was '%s'", expectedHtml, result)
	}
}

func TestProcessIncludeMarkDown(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("## hello {{ include /example.md }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	if err != nil {
		t.Fatal(err)
	}

	exampleFile := mg.ProcessedFile{}
	exampleFile.AppendContent(&mg.HtmlFromMarkdownContent{
		MarkDownContent: &mg.StringContent{Text: "# header", MarkDown: true},
	})

	m := mg.WebFilesMap{}
	m["source/example.md"] = mg.WebFile{Processed: &exampleFile}

	checkParsing(t, processed.Context(), m, processed, emptyContext, []string{"<h2>hello</h2>\n", "<h1>header</h1>\n"})
}

func TestProcessIgnoreEscapedBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello \\{{ include example.html }}"))
	processed, err := mg.ProcessReader(r, "", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{"hello {{ include example.html }}"})
}

func TestProcessIgnoreEscapedClosingBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Hello {{\n  bad-instruction \"contains \\}} ignored\"\n}}. How are you?"))
	processed, err := mg.ProcessReader(r, "", 11)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{
		"Hello ",
		"{{\n  bad-instruction \"contains }} ignored\"\n}}",
		". How are you?",
	})
}

func checkParsing(t *testing.T,
	ctx map[string]interface{}, m mg.WebFilesMap, pf mg.ProcessedFile,
	expectedCtx map[string]interface{}, expectedContents []string) {

	if !reflect.DeepEqual(ctx, expectedCtx) {
		t.Errorf(
			"Expected Context: %v\n"+
				"Actual Context: %v", expectedCtx, ctx)
	}

	if len(pf.Contents) != len(expectedContents) {
		t.Fatalf("Expected %d content parts but got %d: %v",
			len(expectedContents), len(pf.Contents), pf.Contents)
	}

	for i, c := range pf.Contents {
		var result strings.Builder
		c.Write(&result, m, nil)

		if result.String() != expectedContents[i] {
			t.Errorf("Unexpected Content[%d]\nExpected: '%s'\nActual  : '%s'",
				i, expectedContents[i], result.String())
		}
	}
}
