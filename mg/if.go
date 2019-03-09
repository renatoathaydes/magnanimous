package mg

import (
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"io"
	"log"
)

type IfContent struct {
	Text      string
	condition *expression.Expression
	Location  *Location
	contents  []Content
}

var _ Content = (*IfContent)(nil)
var _ ContentContainer = (*IfContent)(nil)

func NewIfInstruction(arg string, location *Location, original string) Content {
	cond, err := expression.ParseExpr(arg)

	if err != nil {
		log.Printf("WARNING: (%s) Malformed if instruction: (%v)", location.String(), err)
		return unevaluatedExpression(original)
	}

	return &IfContent{
		Text:      original,
		condition: &cond,
		Location:  location,
	}
}

func (ic *IfContent) GetContents() []Content {
	return ic.contents
}

func (ic *IfContent) AppendContent(content Content) {
	ic.contents = append(ic.contents, content)
}

func (ic *IfContent) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	res, err := expression.EvalExpr(*ic.condition, magParams{
		webFiles: files,
		stack:    stack,
	})
	if err != nil {
		return err
	}

	switch res {
	case true:
		stack = stack.Push(nil, true)
		err = writeContents(ic, writer, files, stack)
		if err != nil {
			return err
		}
	case false:
	case nil:
		// nothing to write
	default:
		log.Printf("WARN: If condition evaluated to non-boolean value, assuming false: %v", res)
	}
	return nil
}
