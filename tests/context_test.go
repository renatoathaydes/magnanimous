package tests

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/renatoathaydes/magnanimous/mg"
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
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 6, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}

	stack := mg.NewContextStack(mg.NewContext())

	context := processed.ResolveContext(&stack, false)

	if !stack.Top().IsEmpty() {
		t.Error("Expected ResolveContext to not affect the context, but it did", stack.Top())
	}

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

func TestResolveContextInPlace(t *testing.T) {
	files := make(map[string]mg.WebFile)
	resolver := mg.DefaultFileResolver{BasePath: "source", Files: &mg.WebFilesMap{WebFiles: files}}

	r := bufio.NewReader(strings.NewReader("hello\n{{define x `Foo`}}\n{{define y x + `Bar`}}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", "source", 6, &resolver, time.Now())

	if err != nil {
		t.Fatal(err)
	}

	files["source/processed/hi.txt"] = mg.WebFile{Processed: processed}

	stack := mg.NewContextStack(mg.NewContext())

	context := processed.ResolveContext(&stack, true)

	if &stack != context {
		t.Error("Expected ResolveContext (inPlace=true) to return same context", context, stack)
	}

	x := context.Remove("x")
	y := context.Remove("y")

	if x != "Foo" {
		t.Errorf("Expected 'Foo' but got %v", x)
	}
	if y != "FooBar" {
		t.Errorf("Expected 'FooBar' but got %v", y)
	}
	if !context.IsEmpty() {
		t.Error("Expected empty context after x and y were removed, but got", context)
	}
}
