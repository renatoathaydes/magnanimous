package mg

import (
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"io"
	"log"
	"sort"
	"strings"
)

type DefineContent struct {
	Name     string
	Expr     *expression.Expression
	Location Location
	scope    Scope
}

type ExpressionContent struct {
	expr     *expression.Expression
	MarkDown bool
	Text     string
	Location Location
	scope    Scope
}

type iterableExpression struct {
	array    *expression.Expression
	path     string
	resolver FileResolver
	location Location
}

type fileConsumer func(file *WebFile) error
type itemConsumer func(interface{}) error

type iterable interface {
	forEach(files WebFilesMap, inclusionChain []InclusionChainItem,
		parameters magParams, fc fileConsumer, ic itemConsumer) error
}

func NewExpression(arg string, location Location, isMarkDown bool, original string, scope Scope) Content {
	expr, err := expression.ParseExpr(arg)
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval: %s (%s)", location.String(), arg, err.Error())
		return unevaluatedExpression(original)
	}
	return &ExpressionContent{expr: &expr, MarkDown: isMarkDown, Location: location, Text: original, scope: scope}
}

func NewVariable(arg string, location Location, original string, scope Scope) Content {
	parts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
	if len(parts) == 2 {
		variable, rawExpr := parts[0], parts[1]
		expr, err := expression.ParseExpr(rawExpr)
		if err != nil {
			log.Printf("WARNING: (%s) Unable to eval (defining %s): %s (%s)",
				location.String(), variable, rawExpr, err.Error())
			return unevaluatedExpression(original)
		}
		return &DefineContent{Name: variable, Expr: &expr, Location: location, scope: scope}
	}
	log.Printf("WARNING: (%s) malformed define expression: %s", location.String(), arg)
	return unevaluatedExpression(original)
}

func unevaluatedExpression(original string) Content {
	return &StringContent{Text: fmt.Sprintf("{{%s}}", original)}
}

func asIterable(arg string, location Location, resolver FileResolver) (iterable, error) {
	if strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]") {
		expr, err := expression.ParseExpr(fmt.Sprintf("[]interface{}{%s}", arg[1:len(arg)-1]))
		if err != nil {
			return nil, err
		}
		return &iterableExpression{array: &expr, location: location}, nil
	}
	return &iterableExpression{path: arg, location: location, resolver: resolver}, nil
}

func (e *iterableExpression) forEach(files WebFilesMap, inclusionChain []InclusionChainItem,
	parameters magParams, fc fileConsumer, ic itemConsumer) error {
	if e.array != nil {
		v, err := expression.EvalExpr(*e.array, parameters)
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
		_, webFiles, err := e.resolver.FilesIn(e.path, e.location)
		if err != nil {
			return err
		}

		sort.Slice(webFiles, func(i, j int) bool {
			return webFiles[i].Name < webFiles[j].Name
		})

		for _, item := range webFiles {
			item.evalDefinitions(files, inclusionChain)
			err := fc(&item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var _ Content = (*ExpressionContent)(nil)

func (e *ExpressionContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	r, err := expression.EvalExpr(*e.expr, magParams{
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

var _ Content = (*DefineContent)(nil)

func (d *DefineContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	d.Run(files, inclusionChain)
	return nil
}

func (d *DefineContent) Run(files WebFilesMap, inclusionChain []InclusionChainItem) {
	v, err := expression.EvalExpr(*d.Expr, magParams{
		webFiles:       files,
		scope:          d.scope,
		inclusionChain: inclusionChain,
	})
	if err != nil {
		log.Printf("WARNING: (%s) define failure: %s", d.Location.String(), err.Error())
	}
	d.scope.Context()[d.Name] = v
}

func (DefineContent) IsMarkDown() bool {
	return false
}
