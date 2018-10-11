package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestProcessSimple(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello world"))
	ctx, processed := mg.ProcessReader(r, "", 11)

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
	ctx, processed := mg.ProcessReader(r, "", 11)

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
