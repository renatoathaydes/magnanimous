package mg

import (
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"io"
	"log"
	"strings"
	"time"
)

type DefineContent struct {
	Name     string
	Expr     *expression.Expression
	Location *Location
}

type ExpressionContent struct {
	Expr     *expression.Expression
	Text     string
	Location *Location
}

func NewEvalInstruction(arg string, location *Location, original string) Content {
	expr, err := expression.ParseExpr(arg)
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval: %s (%s)", location.String(), arg, err.Error())
		return unevaluatedExpression(original)
	}
	return &ExpressionContent{Expr: &expr, Location: location, Text: original}
}

func NewDefineInstruction(arg string, location *Location, original string) Content {
	parts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
	if len(parts) == 2 {
		variable, rawExpr := parts[0], parts[1]
		expr, err := expression.ParseExpr(rawExpr)
		if err != nil {
			log.Printf("WARNING: (%s) Unable to eval (defining %s): %s (%s)",
				location.String(), variable, rawExpr, err.Error())
			return unevaluatedExpression(original)
		}
		return &DefineContent{Name: variable, Expr: &expr, Location: location}
	}
	log.Printf("WARNING: (%s) malformed define expression: %s", location.String(), arg)
	return unevaluatedExpression(original)
}

func unevaluatedExpression(original string) Content {
	return &StringContent{Text: fmt.Sprintf("{{%s}}", original)}
}

var _ Content = (*ExpressionContent)(nil)

func (e *ExpressionContent) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	r, err := expression.EvalExpr(*e.Expr, magParams{
		stack:    stack,
		webFiles: files,
	})
	if err == nil {
		// an expression can evaluate to a content container, such as a slot
		if c, ok := r.(ContentContainer); ok {
			err = writeContents(c, writer, files, stack)
			if err != nil {
				return err
			}
		} else {
			if r == nil {
				r = ""
			}
			if date, ok := r.(time.Time); ok {
				r = date.Format("02 Jan 2006, 03:04 PM")
			}
			_, err = writer.Write([]byte(fmt.Sprintf("%v", r)))
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

func (e *ExpressionContent) String() string {
	return fmt.Sprintf("ExpressionContent{%s}", e.Text)
}

var _ Content = (*DefineContent)(nil)

func (d *DefineContent) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	// DefineContent does not write anything!
	if v, ok := d.Eval(files, stack); ok {
		if v == nil {
			stack.Top().Remove(d.Name)
		} else {
			stack.Top().Set(d.Name, v)
		}
	}
	return nil
}

func (d *DefineContent) Eval(files WebFilesMap, stack ContextStack) (interface{}, bool) {
	v, err := expression.EvalExpr(*d.Expr, magParams{
		webFiles: files,
		stack:    stack,
	})
	if err != nil {
		log.Printf("WARNING: (%s) define failure: %s", d.Location.String(), err.Error())
		return nil, false
	}
	return v, true
}
