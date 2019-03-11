package mg

import (
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
)

type ForLoop struct {
	Variable string
	iter     iterable
	Text     string
	Location *Location
	contents []Content
}

type forLoopSubInstruction struct {
	sortBy  *sortBySubInstruction
	reverse *reverseSubInstruction
	limit   *limitSubInstruction
}

type sortBySubInstruction struct {
	field string
}

type limitSubInstruction struct {
	max int
}

type reverseSubInstruction struct {
}

type webFileWithContext struct {
	file    *WebFile
	context Context
}

type fileConsumer func(file *webFileWithContext) error

type itemConsumer func(interface{}) error

type iterable interface {
	forEach(files WebFilesMap, stack ContextStack,
		parameters magParams, fc fileConsumer, ic itemConsumer) error
}

type arrayIterable struct {
	array           *expression.Expression
	location        *Location
	subInstructions []forLoopSubInstruction
}

type directoryIterable struct {
	path            string
	resolver        FileResolver
	location        *Location
	subInstructions []forLoopSubInstruction
}

func NewForInstruction(arg string, location *Location, original string, resolver FileResolver) Content {
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
	return &ForLoop{Variable: parts[0], iter: iter, Text: original, Location: location}
}

var _ Content = (*ForLoop)(nil)
var _ ContentContainer = (*ForLoop)(nil)

func (f *ForLoop) GetContents() []Content {
	return f.contents
}

func (f *ForLoop) AppendContent(content Content) {
	f.contents = append(f.contents, content)
}

func (f *ForLoop) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	stack = stack.Push(nil, true)
	err := f.iter.forEach(files, stack, magParams{
		webFiles: files,
		stack:    stack,
	}, func(file *webFileWithContext) error {
		// use the file's context as the value of the bound variable
		stack.Top().Set(f.Variable, file.context)
		return writeContents(f, writer, files, stack)
	}, func(item interface{}) error {
		// use whatever was evaluated from the array as the bound variable
		stack.Top().Set(f.Variable, item)
		return writeContents(f, writer, files, stack)
	})

	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}
	return nil
}

func (f *ForLoop) String() string {
	return fmt.Sprintf("ForLoop{%s}", f.Text)
}

func asIterable(arg string, location *Location, resolver FileResolver) (iterable, error) {
	var forArg string
	var subInstructions []forLoopSubInstruction
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

func iterableFrom(forArg string, subInstructions []forLoopSubInstruction,
	location *Location, resolver FileResolver) (iterable, error) {

	if strings.HasPrefix(forArg, "[") && strings.HasSuffix(forArg, "]") {
		expr, err := expression.ParseExpr(fmt.Sprintf("[]interface{}{%s}", forArg[1:len(forArg)-1]))
		if err != nil {
			return nil, err
		}
		return &arrayIterable{array: &expr, location: location, subInstructions: subInstructions}, nil
	}
	return &directoryIterable{path: forArg, location: location, resolver: resolver, subInstructions: subInstructions}, nil
}

func (e *arrayIterable) forEach(files WebFilesMap, stack ContextStack,
	parameters magParams, fc fileConsumer, ic itemConsumer) error {
	v, err := expression.EvalExpr(*e.array, parameters)
	if err != nil {
		return err
	}
	array, ok := v.([]interface{})
	if ok {
		for _, subInstruction := range e.subInstructions {
			sortBy := subInstruction.sortBy
			if sortBy != nil {
				sortArray(array, sortBy)
			}
			if subInstruction.reverse != nil {
				reverseArray(array)
			}
			if subInstruction.limit != nil {
				limit := len(array)
				if subInstruction.limit.max < limit {
					limit = subInstruction.limit.max
				}
				array = array[0:limit]
			}
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

func (e *directoryIterable) filesWithContext(files WebFilesMap,
	stack ContextStack, parameters magParams) ([]webFileWithContext, error) {
	path := maybeEvalPath(e.path, parameters)
	_, webFiles, err := e.resolver.FilesIn(path, e.location)
	if err != nil {
		return nil, err
	}

	webFilesCtx := make([]webFileWithContext, len(webFiles))
	for i, wf := range webFiles {
		ctx := wf.Processed.ResolveContext(files, stack)
		// we must create a new ref here otherwise the file ref will point to the loop ref, which changes!
		refToFile := wf
		webFilesCtx[i] = webFileWithContext{file: &refToFile, context: ctx}
	}
	return webFilesCtx, nil
}

func (e *directoryIterable) forEach(files WebFilesMap, stack ContextStack,
	parameters magParams, fc fileConsumer, ic itemConsumer) error {
	webFilesCtx, err := e.filesWithContext(files, stack, parameters)
	if err != nil {
		return err
	}

	// start by sorting by name if there's no sub-instructions or if the first sub-instruction is not sortBy
	sortByName := len(e.subInstructions) == 0 || e.subInstructions[0].sortBy == nil

	if sortByName {
		sort.Slice(webFilesCtx, func(i, j int) bool {
			return webFilesCtx[i].file.Name < webFilesCtx[j].file.Name
		})
	}

	for _, subInstruction := range e.subInstructions {
		if subInstruction.sortBy != nil {
			sortField := subInstruction.sortBy.field
			sortFiles(webFilesCtx, sortField)
		}

		if subInstruction.reverse != nil {
			reverseFiles(webFilesCtx)
		}

		if subInstruction.limit != nil {
			limit := len(webFilesCtx)
			if subInstruction.limit.max < limit {
				limit = subInstruction.limit.max
			}
			webFilesCtx = webFilesCtx[:limit]
		}
	}

	for _, item := range webFilesCtx {
		err := fc(&item)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseForLoopSubInstructions(text string) []forLoopSubInstruction {
	parts := strings.Fields(text)
	result := make([]forLoopSubInstruction, len(parts), len(parts))
	resultIdx := 0
	for i := 0; i < len(parts); i++ {
		switch p := parts[i]; p {
		case "sort":
			result[resultIdx].sortBy = &sortBySubInstruction{field: "_"}
			resultIdx++
		case "sortBy":
			if i < len(parts)-1 {
				result[resultIdx].sortBy = &sortBySubInstruction{field: parts[i+1]}
				i++
				resultIdx++
			} else {
				log.Printf("WARN: missing argument for 'sortBy' in for-loop sub-instruction")
				return result
			}
		case "limit":
			if i < len(parts)-1 {
				maxItems, err := strconv.ParseUint(parts[i+1], 10, 32)
				if err != nil {
					log.Printf("WARN: invalid argument for 'limit' in for-loop sub-instruction. "+
						"Expected positive integer, found %s", parts[i+1])
				} else {
					result[resultIdx].limit = &limitSubInstruction{max: int(maxItems)}
					resultIdx++
				}
				i++
			} else {
				log.Printf("WARN: missing argument for 'limit' in for-loop sub-instruction")
				return result
			}
		case "reverse":
			result[resultIdx].reverse = &reverseSubInstruction{}
			resultIdx++
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

func sortFiles(webFiles []webFileWithContext, sortField string) {
	sort.Slice(webFiles, func(i, j int) bool {
		iv, ok := webFiles[i].context.Get(sortField)
		if !ok {
			log.Printf("WARN: cannot sortBy %s - file %s does not define such property",
				sortField, webFiles[i].file.Name)
			return true
		}
		jv, ok := webFiles[j].context.Get(sortField)
		if !ok {
			log.Printf("WARN: cannot sortBy %s - file %s does not define such property",
				sortField, webFiles[j].file.Name)
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

func reverseFiles(webFiles []webFileWithContext) {
	for i := len(webFiles)/2 - 1; i >= 0; i-- {
		opp := len(webFiles) - 1 - i
		webFiles[i], webFiles[opp] = webFiles[opp], webFiles[i]
	}
}

func reverseArray(array []interface{}) {
	for i := len(array)/2 - 1; i >= 0; i-- {
		opp := len(array) - 1 - i
		array[i], array[opp] = array[opp], array[i]
	}
}
