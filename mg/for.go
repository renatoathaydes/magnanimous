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
	context  map[string]interface{}
	parent   Scope
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

var _ Content = (*ForLoop)(nil)
var _ Scope = (*ForLoop)(nil)

func (f *ForLoop) AppendContent(content Content) {
	f.Contents = append(f.Contents, content)
}

func (f *ForLoop) Context() map[string]interface{} {
	return f.context
}

func (f *ForLoop) Parent() Scope {
	return f.parent
}

func (f *ForLoop) setParent(scope Scope) {
	f.parent = scope
}

func (f *ForLoop) Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) *MagnanimousError {
	err := f.iter.forEach(magParams{
		webFiles:       files,
		inclusionChain: inclusionChain,
		scope:          f.parent,
	}, func(file string) error {
		// use the file's context as the value of the bound variable
		f.context[f.Variable] = files[file].Processed.Context()
		return writeContents(f, writer, files, inclusionChain)
	}, func(item interface{}) error {
		// use whatever was evaluated from the array as the bound variable
		f.Context()[f.Variable] = item
		return writeContents(f, writer, files, inclusionChain)
	})
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}

func writeContents(f *ForLoop, writer io.Writer, files WebFilesMap, inclusionChain []Location) *MagnanimousError {
	for _, c := range f.Contents {
		err := c.Write(writer, files, inclusionChain)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *ForLoop) String() string {
	return fmt.Sprintf("ForLoop{%s}", f.Text)
}

func (f *ForLoop) IsMarkDown() bool {
	return f.MarkDown
}
