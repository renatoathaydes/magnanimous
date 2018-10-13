package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestProcessSimple(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello world"))
	ctx, processed, err := mg.ProcessReader(r, "", 11)

	if err != nil {
		t.Fatal(err)
	}

	if len(ctx) != 0 {
		t.Errorf("Expected empty context, but len(ctx) == %d", len(ctx))
	}

	c := processed.Contents

	if len(c) != 1 {
		t.Errorf("Expected 1 Content, but len(Contents) == %d", len(c))
	}

	var result strings.Builder
	m := mg.WebFilesMap{}
	c[0].Write(&result, m)

	if result.String() != "hello world" {
		t.Errorf("Expected 'hello world', but was '%s'", result.String())
	}
}

func TestProcessIncludeSimple(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello {{ include example.html }}"))
	ctx, processed, err := mg.ProcessReader(r, "", 11)

	if err != nil {
		t.Fatal(err)
	}

	if len(ctx) != 0 {
		t.Errorf("Expected empty context, but len(ctx) == %d", len(ctx))
	}

	c := processed.Contents

	if len(c) != 2 {
		t.Errorf("Expected 2 Contents, but got %v", c)
	}

	exampleFile := mg.ProcessedFile{}
	exampleFile.AppendContent(&mg.StringContent{Text: "from another file!"})

	m := mg.WebFilesMap{}
	m["source/example.html"] = mg.WebFile{Processed: exampleFile}

	var result strings.Builder
	c[0].Write(&result, m)

	if result.String() != "hello " {
		t.Errorf("Expected 'hello ', but was '%s'", result.String())
	}

	result.Reset()
	c[1].Write(&result, m)

	if result.String() != "from another file!" {
		t.Errorf("Expected 'from another file!', but was '%s'", result.String())
	}

}

func TestMarkdownToHtml(t *testing.T) {
	file := mg.ProcessedFile{}
	file.AppendContent(&mg.StringContent{Text: "<html><body class=\"hello\">"})
	file.AppendContent(&mg.StringContent{Text: "# Hello\n", MarkDown: true})
	file.AppendContent(&mg.StringContent{Text: "## Mag", MarkDown: true})
	file.AppendContent(&mg.StringContent{Text: "</body></html>"})

	html := mg.MarkdownToHtml(file)

	m := mg.WebFilesMap{}
	result := string(html.Bytes(m))

	expectedHtml := "<html><body class=\"hello\"><h1>Hello</h1>\n<h2>Mag</h2>\n</body></html>"
	if result != expectedHtml {
		t.Errorf("Expected '%s', but was '%s'", expectedHtml, result)
	}
}

func TestProcessIncludeMarkDown(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("## hello {{ include /example.md }}"))
	ctx, processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11)

	if err != nil {
		t.Fatal(err)
	}

	if len(ctx) != 0 {
		t.Errorf("Expected empty context, but len(ctx) == %d", len(ctx))
	}

	c := processed.Contents

	if len(c) != 2 {
		t.Errorf("Expected 2 Contents, but got %v", c)
	}

	exampleFile := mg.ProcessedFile{}
	exampleFile.AppendContent(&mg.HtmlFromMarkdownContent{
		MarkDownContent: &mg.StringContent{Text: "# header", MarkDown: true},
	})

	m := mg.WebFilesMap{}
	m["source/example.md"] = mg.WebFile{Processed: exampleFile}

	var result strings.Builder
	c[0].Write(&result, m)

	if result.String() != "<h2>hello</h2>\n" {
		t.Errorf("Expected '<h2>hello</h2>', but was '%s'", result.String())
	}

	result.Reset()
	c[1].Write(&result, m)

	if result.String() != "<h1>header</h1>\n" {
		t.Errorf("Expected '<h1>header</h1>', but was '%s'", result.String())
	}

}

func TestProcessIgnoreEscapedBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello \\{{ include example.html }}"))
	ctx, processed, err := mg.ProcessReader(r, "", 11)

	if err != nil {
		t.Fatal(err)
	}

	if len(ctx) != 0 {
		t.Errorf("Expected empty context, but len(ctx) == %d", len(ctx))
	}

	c := processed.Contents

	if len(c) != 1 {
		t.Errorf("Expected 1 Content, but got %v", c)
	}

	m := mg.WebFilesMap{}
	var result strings.Builder
	c[0].Write(&result, m)

	if result.String() != "hello {{ include example.html }}" {
		t.Errorf("Expected escaped '{' to be treated as text, but got unexpected result: '%s'", result.String())
	}

}

func TestProcessIgnoreEscapedClosingBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Hello {{\n  eval \"contains \\}} ignored\"\n}}. How are you?"))
	ctx, processed, err := mg.ProcessReader(r, "", 11)

	if err != nil {
		t.Fatal(err)
	}

	if len(ctx) != 0 {
		t.Errorf("Expected empty context, but len(ctx) == %d", len(ctx))
	}

	c := processed.Contents

	if len(c) != 3 {
		t.Fatalf("Expected 3 Contents, but got %v", c)
	}

	m := mg.WebFilesMap{}
	var result strings.Builder
	c[0].Write(&result, m)

	if result.String() != "Hello " {
		t.Errorf("Expected 'Hello ' but got '%s'", result.String())
	}

	result.Reset()
	c[1].Write(&result, m)

	if result.String() != "{{\n  eval \"contains }} ignored\"\n}}" {
		t.Errorf("Expected escaped '}' to be treated as text, but got unexpected result: '%s'", result.String())
	}

	result.Reset()
	c[2].Write(&result, m)

	if result.String() != ". How are you?" {
		t.Errorf("Expected '. How are you?' but got '%s'", result.String())
	}

}
