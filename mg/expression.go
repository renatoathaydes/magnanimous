package mg

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"io"
	"log"
	"strings"
)

type ExpressionContent struct {
	Expression *govaluate.EvaluableExpression
	MarkDown   bool
	Text       string
	Location   Location
}

func NewExpression(arg string, location Location, isMarkDown bool, original string) Content {
	expr, err := govaluate.NewEvaluableExpression(arg)
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval: %s (%s)", location.String(), arg, err.Error())
		return unevaluatedExpression(original)
	}
	return &ExpressionContent{Expression: expr, MarkDown: isMarkDown, Location: location, Text: original}
}

func NewVariable(arg string, location Location, original string, ctx *WebFileContext) Content {
	parts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
	if len(parts) == 2 {
		variable, rawExpr := parts[0], parts[1]
		expr, err := govaluate.NewEvaluableExpression(rawExpr)
		if err != nil {
			log.Printf("WARNING: (%s) Unable to eval (defining %s): %s (%s)",
				location.String(), variable, rawExpr, err.Error())
			return unevaluatedExpression(original)
		}
		v, err := expr.Evaluate(*ctx)
		if err != nil {
			log.Printf("WARNING: (%s) eval failure: %s", location.String(), err.Error())
			return unevaluatedExpression(original)
		}
		(*ctx)[variable] = v
		return nil
	}
	log.Printf("WARNING: (%s) malformed define expression: %s", location.String(), arg)
	return unevaluatedExpression(original)
}

func unevaluatedExpression(original string) Content {
	return &StringContent{Text: fmt.Sprintf("{{%s}}", original)}
}

func (e *ExpressionContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) *MagnanimousError {
	r, err := e.Expression.Eval(magParams{
		webFiles:       files,
		origin:         e.Location,
		inclusionChain: inclusionChain,
	})
	if err == nil {
		writer.Write([]byte(fmt.Sprintf("%v", r)))
	} else {
		log.Printf("WARNING: (%s) eval failure: %s", e.Location.String(), err.Error())
		writer.Write([]byte(fmt.Sprintf("{{%s}}", e.Text)))
	}
	return nil
}

func (e *ExpressionContent) String() string {
	return fmt.Sprintf("ExpressionContent{%s}", e.Text)
}

func (e *ExpressionContent) IsMarkDown() bool {
	return e.MarkDown
}
