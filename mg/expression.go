package mg

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"io"
	"log"
	"path/filepath"
	"strings"
)

type ExpressionContent struct {
	expr     *govaluate.EvaluableExpression
	MarkDown bool
	Text     string
	Location Location
	scope    Scope
}

type iterableExpression struct {
	array    *govaluate.EvaluableExpression
	path     string
	location Location
}

type fileConsumer func(string) error
type itemConsumer func(interface{}) error

type iterable interface {
	forEach(parameters magParams, fc fileConsumer, ic itemConsumer) error
}

func NewExpression(arg string, location Location, isMarkDown bool, original string) Content {
	expr, err := govaluate.NewEvaluableExpression(arg)
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval: %s (%s)", location.String(), arg, err.Error())
		return unevaluatedExpression(original)
	}
	return &ExpressionContent{expr: expr, MarkDown: isMarkDown, Location: location, Text: original}
}

func NewVariable(arg string, location Location, original string, scope Scope) Content {
	parts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
	if len(parts) == 2 {
		variable, rawExpr := parts[0], parts[1]
		expr, err := govaluate.NewEvaluableExpression(rawExpr)
		if err != nil {
			log.Printf("WARNING: (%s) Unable to eval (defining %s): %s (%s)",
				location.String(), variable, rawExpr, err.Error())
			return unevaluatedExpression(original)
		}
		v, err := expr.Evaluate(scope.Context())
		if err != nil {
			log.Printf("WARNING: (%s) eval failure: %s", location.String(), err.Error())
			return unevaluatedExpression(original)
		}
		scope.Context()[variable] = v
		return nil
	}
	log.Printf("WARNING: (%s) malformed define expression: %s", location.String(), arg)
	return unevaluatedExpression(original)
}

func unevaluatedExpression(original string) Content {
	return &StringContent{Text: fmt.Sprintf("{{%s}}", original)}
}

func asIterable(arg string, location Location) (iterable, error) {
	if strings.HasPrefix(arg, "(") && strings.HasSuffix(arg, ")") {
		expr, err := govaluate.NewEvaluableExpression(arg)
		if err != nil {
			return nil, err
		}
		return &iterableExpression{array: expr, location: location}, nil
	}
	return &iterableExpression{path: arg, location: location}, nil
}

func (e *iterableExpression) forEach(parameters magParams, fc fileConsumer, ic itemConsumer) error {
	if e.array != nil {
		v, err := e.array.Eval(parameters)
		if err != nil {
			return err
		}
		for _, item := range v.([]interface{}) {
			err := ic(item)
			if err != nil {
				return err
			}
		}
	} else {
		dir, f, err := DefaultFileResolver.FilesIn(e.path, e.location)
		if err != nil {
			return err
		}
		for _, item := range f {
			if !item.IsDir() {
				path := filepath.Join(dir, item.Name())
				err := fc(path)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

var _ ScopeSensitive = (*ExpressionContent)(nil)
var _ Content = (*ExpressionContent)(nil)

func (e *ExpressionContent) setScope(holder Scope) {
	e.scope = holder
}

func (e *ExpressionContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) error {
	r, err := e.expr.Eval(magParams{
		webFiles:       files,
		scope:          e.scope,
		inclusionChain: inclusionChain,
	})
	if err == nil {
		_, err = writer.Write([]byte(fmt.Sprintf("%v", r)))
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

func (e *ExpressionContent) IsMarkDown() bool {
	return e.MarkDown
}
