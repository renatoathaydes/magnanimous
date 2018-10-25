package mg

import (
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"io"
	"log"
)

type IfContent struct {
	MarkDown  bool
	Text      string
	condition *expression.Expression
	Location  Location
	contents  []Content
	context   map[string]interface{}
	parent    Scope
}

var _ Content = (*IfContent)(nil)
var _ Scope = (*IfContent)(nil)
var _ ContentContainer = (*IfContent)(nil)

func NewIfInstruction(arg string, location Location, isMarkDown bool,
	original string) Content {

	cond, err := expression.ParseExpr(arg)

	if err != nil {
		log.Printf("WARNING: (%s) Malformed if instruction: (%v)", location.String(), err)
		return unevaluatedExpression(original)
	}

	return &IfContent{
		MarkDown:  isMarkDown,
		Text:      original,
		condition: &cond,
		Location:  location,
		context:   make(map[string]interface{}, 2),
	}
}

func (ic *IfContent) GetContents() []Content {
	return ic.contents
}

func (ic *IfContent) AppendContent(content Content) {
	ic.contents = append(ic.contents, content)
}

func (ic *IfContent) Context() map[string]interface{} {
	return ic.context
}

func (ic *IfContent) Parent() Scope {
	return ic.parent
}

func (ic *IfContent) setParent(scope Scope) {
	ic.parent = scope
}

func (ic *IfContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	res, err := expression.EvalExpr(*ic.condition, magParams{
		webFiles:       files,
		scope:          ic,
		inclusionChain: inclusionChain,
	})
	if err != nil {
		return err
	}

	switch res {
	case true:
		writeContents(ic, writer, files, inclusionChain)
	case false:
		// nothing to write
	default:
		log.Printf("WARN: If condition evaluated to non-boolean value, assuming false: %v", res)
	}
	return nil
}

func (ic *IfContent) IsMarkDown() bool {
	return ic.MarkDown
}
