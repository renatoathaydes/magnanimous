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

func TestStringOR(t *testing.T) {
	examples := [][3]string{
		{"", "", ""},
		{"a", "", "a"},
		{"", "a", "a"},
		{"a", "a", "a"},
		{"a", "b", "a"},
		{"b", "a", "b"},
	}

	nilIfEmpty := func(s interface{}) *string {
		if s == "" {
			return nil
		}
		r := fmt.Sprintf("%v", s)
		return &r
	}

	nullIfEmpty := func(s string) string {
		if s == "" {
			return "null"
		}
		return fmt.Sprintf(`"%s"`, s)
	}

	for i, ex := range examples {
		x := nullIfEmpty(ex[0])
		y := nullIfEmpty(ex[1])
		expr := fmt.Sprintf("%v || %v", x, y)
		v, err := expression.Eval(expr, nil)

		if err != nil {
			t.Fatalf("[%d] Could not evaluate: %v", i, err)
		}
		expected := nilIfEmpty(ex[2])

		if ok := (expected == nil && v == nil) || v == *expected; !ok {
			t.Errorf("[%d] Expected '%s' to evaluate to '%v' but got '%v'", i, expr, expected, v)
		}
	}
}

func TestFloatOR(t *testing.T) {
	examples := [][3]float64{
		{float64(0), float64(0), float64(0)},
		{float64(0), float64(1), float64(1)},
		{float64(1), float64(0), float64(1)},
		{float64(1), float64(1), float64(1)},
		{float64(0), float64(42), float64(42)},
		{float64(66), float64(44), float64(66)},
	}

	for i, ex := range examples {
		expr := fmt.Sprintf("%v || %v", ex[0], ex[1])
		v, err := expression.Eval(expr, nil)

		if err != nil {
			t.Fatalf("[%d] Could not evaluate: %v", i, err)
		}
		if v != ex[2] {
			t.Errorf("[%d] Expected '%s' to evaluate to '%v' but got '%v'", i, expr, ex[2], v)
		}
	}
}
