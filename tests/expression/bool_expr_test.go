package expression

import (
	"fmt"
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

func TestBoolNegation(t *testing.T) {
	v, err := expression.Eval(`!false`, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != true {
		t.Errorf("Expected 'true' but got '%v'", v)
	}
}

func TestBoolAnd(t *testing.T) {
	examples := [][3]bool{
		{true, true, true},
		{true, false, false},
		{false, true, false},
		{false, false, false},
	}

	for i, ex := range examples {
		expr := fmt.Sprintf("%v && %v", ex[0], ex[1])
		v, err := expression.Eval(expr, nil)

		if err != nil {
			t.Fatalf("[%d] Could not evaluate: %v", i, err)
		}

		if v != ex[2] {
			t.Errorf("[%d] Expected 'true' but got '%v'", i, v)
		}
	}
}

func TestBoolOR(t *testing.T) {
	examples := [][3]bool{
		{true, true, true},
		{true, false, true},
		{false, true, true},
		{false, false, false},
	}

	for i, ex := range examples {
		expr := fmt.Sprintf("%v || %v", ex[0], ex[1])
		v, err := expression.Eval(expr, nil)

		if err != nil {
			t.Fatalf("[%d] Could not evaluate: %v", i, err)
		}

		if v != ex[2] {
			t.Errorf("[%d] Expected 'true' but got '%v'", i, v)
		}
	}
}
