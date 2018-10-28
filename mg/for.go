package mg

import (
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"io"
	"log"
	"sort"
	"strings"
)

type ForLoop struct {
	Variable string
	iter     iterable
	MarkDown bool
	Text     string
	Location Location
	contents []Content
	context  map[string]interface{}
	parent   Scope
}

type forLoopSubInstructions struct {
	sortBy *sortBySubInstruction
}

type sortBySubInstruction struct {
	field string
}

type fileConsumer func(file *WebFile) error

type itemConsumer func(interface{}) error

type iterable interface {
	forEach(files WebFilesMap, inclusionChain []InclusionChainItem,
		parameters magParams, fc fileConsumer, ic itemConsumer) error
}

type arrayIterable struct {
	array           *expression.Expression
	location        Location
	subInstructions forLoopSubInstructions
}

type directoryIterable struct {
	path            string
	resolver        FileResolver
	location        Location
	subInstructions forLoopSubInstructions
}

func NewForInstruction(arg string, location Location, isMarkDown bool,
	original string, resolver FileResolver) Content {
	parts := strings.SplitN(arg, " ", 2)
	switch len(parts) {
	case 0:
		fallthrough
	case 1:
		log.Printf("WARNING: (%s) Malformed for loop instruction", location.String())
		return unevaluatedExpression(original)
	}
	iter, err := asIterable(parts[1], location, resolver)
	if err != nil {
		log.Printf("WARNING: (%s) Unable to eval iterable in for expression: %s (%s)",
			location.String(), arg, err.Error())
		return unevaluatedExpression(original)
	}
	return &ForLoop{Variable: parts[0], iter: iter, MarkDown: isMarkDown,
		Text: original, Location: location, context: make(map[string]interface{}, 2)}
}

var _ Content = (*ForLoop)(nil)
var _ Scope = (*ForLoop)(nil)
var _ ContentContainer = (*ForLoop)(nil)

func (f *ForLoop) GetContents() []Content {
	return f.contents
}

func (f *ForLoop) AppendContent(content Content) {
	f.contents = append(f.contents, content)
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

func (f *ForLoop) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	err := f.iter.forEach(files, inclusionChain, magParams{
		webFiles:       files,
		inclusionChain: inclusionChain,
		scope:          f.parent,
	}, func(webFile *WebFile) error {
		// use the file's context as the value of the bound variable
		f.context[f.Variable] = webFile.Processed.Context()
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

func (f *ForLoop) String() string {
	return fmt.Sprintf("ForLoop{%s}", f.Text)
}

func (f *ForLoop) IsMarkDown() bool {
	return f.MarkDown
}

func asIterable(arg string, location Location, resolver FileResolver) (iterable, error) {
	var forArg string
	var subInstructions forLoopSubInstructions
	if strings.HasPrefix(arg, "(") {
		idx := strings.Index(arg, ")")
		if idx > 0 {
			subInstructions = parseForLoopSubInstructions(strings.TrimSpace(arg[1:idx]))
			forArg = strings.TrimSpace(arg[idx+1:])
		}
	} else {
		forArg = arg
	}
	return iterableFrom(forArg, subInstructions, location, resolver)
}

func iterableFrom(forArg string, subInstructions forLoopSubInstructions,
	location Location, resolver FileResolver) (iterable, error) {

	if strings.HasPrefix(forArg, "[") && strings.HasSuffix(forArg, "]") {
		expr, err := expression.ParseExpr(fmt.Sprintf("[]interface{}{%s}", forArg[1:len(forArg)-1]))
		if err != nil {
			return nil, err
		}
		return &arrayIterable{array: &expr, location: location, subInstructions: subInstructions}, nil
	}
	return &directoryIterable{path: forArg, location: location, resolver: resolver, subInstructions: subInstructions}, nil
}

func (e *arrayIterable) forEach(files WebFilesMap, inclusionChain []InclusionChainItem,
	parameters magParams, fc fileConsumer, ic itemConsumer) error {
	v, err := expression.EvalExpr(*e.array, parameters)
	if err != nil {
		return err
	}
	array, ok := v.([]interface{})
	if ok {
		sortBy := e.subInstructions.sortBy
		if sortBy != nil {
			sortArray(array, sortBy)
		}
		for _, item := range array {
			err := ic(item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *directoryIterable) forEach(files WebFilesMap, inclusionChain []InclusionChainItem,
	parameters magParams, fc fileConsumer, ic itemConsumer) error {
	_, webFiles, err := e.resolver.FilesIn(e.path, e.location)
	if err != nil {
		return err
	}

	if e.subInstructions.sortBy == nil {
		sort.Slice(webFiles, func(i, j int) bool {
			return webFiles[i].Name < webFiles[j].Name
		})
	} else {
		sortField := e.subInstructions.sortBy.field
		sort.Slice(webFiles, func(i, j int) bool {
			webFiles[i].evalDefinitions(files, inclusionChain)
			iv, ok := webFiles[i].Processed.Context()[sortField]
			if !ok {
				log.Printf("WARN: cannot sortBy %s - file %s does not define such property",
					sortField, webFiles[i].Name)
				return true
			}
			webFiles[j].evalDefinitions(files, inclusionChain)
			jv, ok := webFiles[j].Processed.Context()[sortField]
			if !ok {
				log.Printf("WARN: cannot sortBy %s - file %s does not define such property",
					sortField, webFiles[j].Name)
				return true
			}
			res, err := expression.Less(iv, jv)
			if err != nil {
				log.Printf("WARN: sortBy %s error - %s", sortField, err)
				return true
			}
			return res.(bool)
		})
	}

	for _, item := range webFiles {
		item.evalDefinitions(files, inclusionChain)
		err := fc(&item)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseForLoopSubInstructions(text string) forLoopSubInstructions {
	result := forLoopSubInstructions{}
	parts := strings.Fields(text)
	for i := 0; i < len(parts); i++ {
		switch p := parts[i]; p {
		case "sort":
			result.sortBy = &sortBySubInstruction{field: "_"}
		case "sortBy":
			if i < len(parts)-1 {
				result.sortBy = &sortBySubInstruction{field: parts[i+1]}
				i++
			} else {
				log.Printf("WARN: missing argument for 'sortBy' for-loop sub-instruction")
				return result
			}
		default:
			log.Printf("Unrecognized for-loop sub-instruction: " + p)
		}
	}
	return result
}

func sortArray(array []interface{}, instruction *sortBySubInstruction) {
	if instruction.field != "_" {
		log.Printf("WARN: It is not possible to sort simple array by field (use '_' instead): %s",
			instruction.field)
	}
	sort.Slice(array, func(i, j int) bool {
		res, err := expression.Less(array[i], array[j])
		if err != nil {
			log.Printf("WARN: %s", err.Error())
			return false
		}
		return res.(bool)
	})
}

func sortFiles() {

}
