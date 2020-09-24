package mg

import (
	"io"
	"log"
	"strings"

	"github.com/renatoathaydes/magnanimous/mg/expression"
)

// DefineContent is a Content that defines a new variable in its context.
//
// Although DefineContent is Writable, it does not actually write anything where it's defined.
type DefineContent struct {
	UnscopedContent
	Name     string
	Text     string
	Expr     *expression.Expression
	Location *Location
	resolver FileResolver
}

var _ Content = (*DefineContent)(nil)
var _ Definition = (*DefineContent)(nil)

func NewDefineInstruction(arg string, location *Location, original string, resolver FileResolver) Content {
	parts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
	if len(parts) == 2 {
		variable, rawExpr := parts[0], parts[1]
		expr, err := expression.ParseExpr(rawExpr)
		if err != nil {
			log.Printf("WARNING: (%s) Unable to eval (defining %s): %s (%s)",
				location.String(), variable, rawExpr, err.Error())
			return unevaluatedExpression(original, location)
		}
		return &DefineContent{Name: variable, Text: original, Expr: &expr, Location: location, resolver: resolver}
	}
	log.Printf("WARNING: (%s) malformed define expression: %s", location.String(), arg)
	return unevaluatedExpression(original, location)
}

func (d *DefineContent) GetLocation() *Location {
	return d.Location
}

func (d *DefineContent) GetName() string {
	return d.Name
}

func (d *DefineContent) Write(writer io.Writer, context Context) ([]Content, error) {
	if v, ok := d.Eval(context); ok {
		if v == nil {
			context.Remove(d.Name)
		} else {
			context.Set(d.Name, v)
		}
		return nil, nil
	}
	return unevaluatedExpressions(d.Text, d.Location), nil
}

func (d *DefineContent) Eval(context Context) (interface{}, bool) {
	v, err := expression.EvalExpr(d.Expr, context)
	if err != nil {
		log.Printf("WARNING: (%s) define failure: %s", d.Location.String(), err.Error())
		return nil, false
	}
	return v, true
}
