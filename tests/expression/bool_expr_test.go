package expression

import (
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"testing"
)

func TestBoolEquality(t *testing.T) {
	v, err := expression.Eval(`true == true`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestBoolInequality(t *testing.T) {
	v, err := expression.Eval(`true != false`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}
