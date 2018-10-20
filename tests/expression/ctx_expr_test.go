package expression

import (
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"testing"
)

func TestSimpleIdentifierAccess(t *testing.T) {
	ctx := map[string]interface{}{"foo": "Hello"}
	v, err := expression.Eval(`foo`, &expression.MapContext{Map: ctx})

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != "Hello" {
		t.Errorf("Expected 'Hello' but got '%v'", v)
	}
}

func TestIdentifierExpression(t *testing.T) {
	ctx := map[string]interface{}{"p1": "Hello", "p2": " ", "p3": "World"}
	v, err := expression.Eval(`p1 + p2 + p3`, &expression.MapContext{Map: ctx})

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != "Hello World" {
		t.Errorf("Expected 'Hello World' but got '%v'", v)
	}
}

func TestNestedIdentifierAccess(t *testing.T) {
	person := map[string]interface{}{
		"name": "Joe",
		"age":  28,
		"address": map[string]interface{}{
			"street": "High Street",
			"city":   "Metropolis",
		},
	}
	ctx := map[string]interface{}{"person": person}

	v, err := expression.Eval(`"a person called " + person.name +
		", age " + person.age + ", lives on " + person.address.street +
        " in the city of " + person.address.city`, &expression.MapContext{Map: ctx})

	if err != nil {
		t.Fatalf("Could not evaluate: %v", err)
	}

	if v != "a person called Joe, age 28, lives on High Street in the city of Metropolis" {
		t.Errorf("Got unexpected: '%v'", v)
	}
}
