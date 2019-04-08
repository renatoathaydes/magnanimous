package mg

import (
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"io"
	"log"
	"strings"
)

type DefineContent struct {
	Name     string
	Expr     *expression.Expression
	Location *Location
	resolver FileResolver
}

type ExpressionContent struct {
	Expr     *expression.Expression
	Text     string
	Location *Location
	resolver FileResolver
}

func NewEvalInstruction(arg string, location *Location, original string, resolver FileResolver) Content {
	expr, err := expression.ParseExpr(arg)
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval: %s (%s)", location.String(), arg, err.Error())
		return unevaluatedExpression(original)
	}
	return &ExpressionContent{Expr: &expr, Location: location, Text: original, resolver: resolver}
}

func NewDefineInstruction(arg string, location *Location, original string, resolver FileResolver) Content {
	parts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
	if len(parts) == 2 {
		variable, rawExpr := parts[0], parts[1]
		expr, err := expression.ParseExpr(rawExpr)
		if err != nil {
			log.Printf("WARNING: (%s) Unable to eval (defining %s): %s (%s)",
				location.String(), variable, rawExpr, err.Error())
			return unevaluatedExpression(original)
		}
		return &DefineContent{Name: variable, Expr: &expr, Location: location, resolver: resolver}
	}
	log.Printf("WARNING: (%s) malformed define expression: %s", location.String(), arg)
	return unevaluatedExpression(original)
}

func unevaluatedExpression(original string) Content {
	return &StringContent{Text: fmt.Sprintf("{{%s}}", original)}
}

var _ Content = (*ExpressionContent)(nil)

func (e *ExpressionContent) Write(writer io.Writer, stack ContextStack) error {
	params := magParams{stack: stack, fileResolver: e.resolver, location: e.Location}
	r, err := expression.EvalExpr(*e.Expr, &params)
	if err == nil {
		// an expression can evaluate to a content container, such as a slot
		if c, ok := r.(ContentContainer); ok {
			err = writeContents(c, writer, stack)
			if err != nil {
				return err
			}
		} else {
			// evaluate special types to a simple string to write
			s := evalSpecialType(r, &params, stack, e.Location)
			_, err = writer.Write([]byte(s))
		}
	} else {
		log.Printf("WARNING: (%s) eval failure: %s", e.Location.String(), err.Error())
		_, err = writer.Write([]byte(fmt.Sprintf("{{%s}}", e.Text)))
	}
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}

func evalSpecialType(r interface{}, params *magParams, stack ContextStack, location *Location) string {
	switch v := r.(type) {
	case nil:
		return ""
	case *expression.DateTime:
		return v.Time.Format(v.Format)
	case *expression.Path:
		if f, ok := params.File(v.Value); ok {
			return f.Processed.Path
		}
		return v.Value
	case *expression.PathProperty:
		if f, ok := params.File(v.Path.Value); ok {
			ctx := f.Processed.ResolveContext(stack)
			if prop, ok := ctx.Get(v.Name); ok {
				return evalSpecialType(prop, params, stack, location)
			} else {
				log.Printf("WARNING: (%s) eval failure: File at path %s has no such property: %s",
					location.String(), v.Path.Value, prop)
				return fmt.Sprintf("%s.%s", v.Path.Value, v.Name)
			}
		}
	}
	// no special type found, stringify it
	return fmt.Sprintf("%v", r)
}

func (e *ExpressionContent) String() string {
	return fmt.Sprintf("ExpressionContent{%s}", e.Text)
}

var _ Content = (*DefineContent)(nil)

func (d *DefineContent) Write(writer io.Writer, stack ContextStack) error {
	// DefineContent does not write anything!
	if v, ok := d.Eval(stack); ok {
		if v == nil {
			stack.Top().Remove(d.Name)
		} else {
			stack.Top().Set(d.Name, v)
		}
	}
	return nil
}

func (d *DefineContent) Eval(stack ContextStack) (interface{}, bool) {
	v, err := expression.EvalExpr(*d.Expr, &magParams{
		fileResolver: d.resolver,
		stack:        stack,
		location:     d.Location,
	})
	if err != nil {
		log.Printf("WARNING: (%s) define failure: %s", d.Location.String(), err.Error())
		return nil, false
	}
	return v, true
}
