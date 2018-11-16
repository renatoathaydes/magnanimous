package mg

import (
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"io"
	"log"
	"reflect"
	"strings"
)

type DefineContent struct {
	Name             string
	Expr             *expression.Expression
	Location         Location
	scope            Scope
	latestInclusions []InclusionChainItem
}

type ExpressionContent struct {
	Expr     *expression.Expression
	Text     string
	Location Location
	scope    Scope
}

func NewExpression(arg string, location Location, original string, scope Scope) Content {
	expr, err := expression.ParseExpr(arg)
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval: %s (%s)", location.String(), arg, err.Error())
		return unevaluatedExpression(original)
	}
	return &ExpressionContent{Expr: &expr, Location: location, Text: original, scope: scope}
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

var _ Content = (*ExpressionContent)(nil)

func (e *ExpressionContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	r, err := expression.EvalExpr(*e.Expr, magParams{
		webFiles:       &files,
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

var _ Content = (*DefineContent)(nil)
var _ SideEffectContent = (*DefineContent)(nil)

func (d *DefineContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	d.Run(&files, inclusionChain)
	return nil
}

func (d *DefineContent) Run(files *WebFilesMap, inclusionChain []InclusionChainItem) {
	if d.latestInclusions != nil &&
		reflect.ValueOf(d.latestInclusions).Pointer() == reflect.ValueOf(inclusionChain).Pointer() {
		// already evaluated for this inclusion chain
		return
	}

	v, err := expression.EvalExpr(*d.Expr, magParams{
		webFiles:       files,
		scope:          d.scope,
		inclusionChain: inclusionChain,
	})
	if err != nil {
		log.Printf("WARNING: (%s) define failure: %s", d.Location.String(), err.Error())
	}
	d.scope.Context().Set(d.Name, v)
	d.latestInclusions = inclusionChain
}
