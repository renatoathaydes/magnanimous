package expression

import (
	"reflect"
	"testing"

	"github.com/renatoathaydes/magnanimous/mg/expression"
)

func TestStringArray(t *testing.T) {
	e, err := expression.ParseExpr(`["abc", "def", "ghi"]`)

	if err != nil {
		t.Fatalf("Could not parse: %v", err)
	}

	v, err := expression.EvalExpr(&e, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if !reflect.DeepEqual(v, []interface{}{"abc", "def", "ghi"}) {
		t.Errorf("Expected '[abc, def, ghi]' but got '%v'", v)
	}
}
func TestIntArray(t *testing.T) {
	e, err := expression.ParseExpr(`[25, 42, 55, 62, 98]`)

	if err != nil {
		t.Fatalf("Could not parse: %v", err)
	}

	v, err := expression.EvalExpr(&e, nil)

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if !reflect.DeepEqual(v, []interface{}{25.0, 42.0, 55.0, 62.0, 98.0}) {
		t.Errorf("Expected '[25, 42, 55, 62, 98]' but got '%v'", v)
	}
}
