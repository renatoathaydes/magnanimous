package mg

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type ForLoop struct {
	Variable string
	iter     *parsedIterable
	Text     string
	Location *Location
	contents []Content
	resolver FileResolver
}

var _ Content = (*ForLoop)(nil)
var _ ContentContainer = (*ForLoop)(nil)

func NewForInstruction(arg string, location *Location, original string, resolver FileResolver) Content {
	parts := strings.SplitN(arg, " ", 2)
	switch len(parts) {
	case 0:
		fallthrough
	case 1:
		log.Printf("WARNING: (%s) Malformed for loop instruction", location.String())
		return unevaluatedExpression(original, location)
	}
	iter, err := parseIterable(parts[1], location, resolver)
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval iterable in for expression: %s (%s)",
			location.String(), arg, err.Error())
		return unevaluatedExpression(original, location)
	}
	return &ForLoop{Variable: parts[0], iter: iter, Text: original, Location: location, resolver: resolver}
}

func (f *ForLoop) AppendContent(content Content) {
	f.contents = append(f.contents, content)
}

func (f *ForLoop) GetLocation() *Location {
	return f.Location
}

func (f *ForLoop) IsScoped() bool {
	return true
}

func (f *ForLoop) Write(writer io.Writer, context Context) ([]Content, error) {
	var items []interface{}
	iterable, ok := f.resolveIterable(context)
	if !ok {
		return unevaluatedExpressions(f.Text, f.Location), nil
	}
	if iterable.files != nil {
		for _, file := range iterable.files {
			// use the file's context as the value of the bound variable
			items = append(items, file.context)
		}
	} else if iterable.items != nil {
		items = iterable.items
	} else if iterable.groups != nil {
		for _, item := range iterable.groups {
			items = append(items, item)
		}
	} else {
		return nil, fmt.Errorf("ForLoop iterable without files, groups or items: %s",
			strings.TrimSpace(f.Text))
	}

	var result []Content

	for _, item := range items {
		// use the file's context as the value of the bound variable
		result = append(result, &iterationContent{
			variable: f.Variable,
			contents: f.contents,
			item:     item,
			location: f.Location,
		})
	}
	return result, nil
}

func (f *ForLoop) resolveIterable(context Context) (iterable iterable, ok bool) {
	gIter := f.iter
	arg := pathOrEval(gIter.arg, context)
	switch a := arg.(type) {
	case string:
		dirIter := directoryIterable{path: a, location: gIter.location, resolver: gIter.resolver,
			subInstructions: gIter.subInstructions}
		files, groupedBy, err := dirIter.getItems(context)
		if err != nil {
			log.Printf("WARNING: (%s) for-loop expression error getting files to iterate over: %s", f.Location.String(), err.Error())
			return iterable, false
		}
		if groupedBy != nil {
			iterable.groups = groupedBy
		} else if files != nil {
			iterable.files = files
		}
	case []interface{}:
		itemsIter := arrayIterable{array: a, location: gIter.location, subInstructions: gIter.subInstructions}
		iterable.items = itemsIter.getItems(context)
	case []webFileWithContext:
		array := make([]interface{}, len(a))
		for i, item := range a {
			array[i] = item.context
		}
		iterable.items = array
	default:
		log.Printf("WARNING: (%s) invalid for-loop expression, cannot iterate over: %v", f.Location.String(), arg)
		return iterable, false
	}
	return iterable, true
}

func (f *ForLoop) String() string {
	return fmt.Sprintf("ForLoop{%s}", f.Text)
}
