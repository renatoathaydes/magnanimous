package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestProcessSimple(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello world"))
	ctx, processed := mg.ProcessReader(r, 11)

	if len(*ctx) != 0 {
		t.Errorf("Expected empty context, but len(ctx) == %d", len(*ctx))
	}

	c := processed.Contents

	if len(c) != 1 {
		t.Errorf("Expected 1 Content, but len(Contents) == %d", len(c))
	}

	var result strings.Builder
	c[0].Write(&result)

	if result.String() != "hello world" {
		t.Errorf("Expected 'hello world', but was '%s'", result.String())
	}
}

func TestProcessIncludeSimple(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello {{ include /example.txt }}"))
	ctx, processed := mg.ProcessReader(r, 11)

	if len(*ctx) != 0 {
		t.Errorf("Expected empty context, but len(ctx) == %d", len(*ctx))
	}

	c := processed.Contents

	if len(c) != 2 {
		t.Errorf("Expected 2 Contents, but got %v", c)
	}

	var result strings.Builder
	c[0].Write(&result)

	if result.String() != "hello " {
		t.Errorf("Expected 'hello ', but was '%s'", result.String())
	}

	result.Reset()
	c[1].Write(&result)

	if result.String() != "example" {
		t.Errorf("Expected 'example', but was '%s'", result.String())
	}

}
