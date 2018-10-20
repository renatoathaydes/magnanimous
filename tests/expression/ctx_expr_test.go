package expression

import (
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"testing"
)

func TestSimpleIdentifierAccess(t *testing.T) {
	ctx := map[string]interface{}{"foo": "Hello"}
	v, err := expression.Eval(`foo`, ctx)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != "Hello" {
		t.Errorf("Expected 'Hello' but got '%v'", v)
	}
}

func TestIdentifierExpression(t *testing.T) {
	ctx := map[string]interface{}{"p1": "Hello", "p2": " ", "p3": "World"}
	v, err := expression.Eval(`p1 + p2 + p3`, ctx)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != "Hello World" {
		t.Errorf("Expected 'Hello World' but got '%v'", v)
	}
}
