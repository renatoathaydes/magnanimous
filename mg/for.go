package mg

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type ForLoop struct {
	Variable string
	iter     iterable
	MarkDown bool
	Text     string
	Location Location
	Contents []Content
	ctx      WebFileContext
}

func NewForInstruction(arg string, location Location, isMarkDown bool, original string) Content {
	parts := strings.SplitN(arg, " ", 2)
	switch len(parts) {
	case 0:
		fallthrough
	case 1:
		log.Printf("WARNING: (%s) Malformed for loop instruction", location.String())
		return unevaluatedExpression(original)
	}
	iter, err := asIterable(parts[1])
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval iterable in for expression: %s (%s)",
			location.String(), arg, err.Error())
		return unevaluatedExpression(original)
	}
	return &ForLoop{Variable: parts[0], iter: iter, MarkDown: isMarkDown, Text: original, Location: location}
}

// assert implementation of HasContent
var _ HasContent = (*ForLoop)(nil)

func (f *ForLoop) AppendContent(content Content) {
	f.Contents = append(f.Contents, content)
}

func (f *ForLoop) Context() WebFileContext {
	return f.ctx
}

// assert implementation of Content
var _ Content = (*ForLoop)(nil)

func (f *ForLoop) Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) *MagnanimousError {
	// TODO
	return nil
}

func (f *ForLoop) String() string {
	return fmt.Sprintf("ForLoop{%s}", f.Text)
}

func (f *ForLoop) IsMarkDown() bool {
	return f.MarkDown
}
