package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestSimpleContext(t *testing.T) {
	ctx := mg.NewContext()

	if !ctx.IsEmpty() {
		t.Error("Should be empty on creation")
	}
	if v, ok := ctx.Get("something"); ok {
		t.Error("Should not contain something, but it does:", v)
	}

	ctx.Set("something", true)

	if ctx.IsEmpty() {
		t.Error("Should NOT be empty after setting one item:", ctx)
	}

	if v, ok := ctx.Get("something"); ok {
		if v != true {
			t.Error("something should be true, but was:", v)
		}
	} else {
		t.Error("Should contain something, but it does not")
	}
}

func TestResolveContext(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("hello\n{{define x 1}}\n{{define y x + 1}}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 6, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}

	stack := mg.NewContextStack(mg.NewContext())

	context := processed.ResolveContext(mg.WebFilesMap{WebFiles: files}, stack)

	x := context.Remove("x")
	y := context.Remove("y")

	if x != float64(1) {
		t.Errorf("Expected 1 but got %v", x)
	}
	if y != float64(2) {
		t.Errorf("Expected 2 but got %v", y)
	}
	if !context.IsEmpty() {
		t.Error("Expected empty context after x and y were removed, but got", context)
	}
}
