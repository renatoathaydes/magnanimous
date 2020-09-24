package mg

import (
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/renatoathaydes/magnanimous/mg/expression"
)

type webFileWithContext struct {
	file    WebFile
	context Context
}

type iterable struct {
	files []webFileWithContext
	items []interface{}
}

type iterationContent struct {
	UnscopedContent
	variable string
	item     interface{}
	contents []Content
	location *Location
}

var _ Content = (*iterationContent)(nil)

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

type parsedIterable struct {
	arg             string
	resolver        FileResolver
	location        *Location
	subInstructions []forLoopSubInstruction
}

type arrayIterable struct {
	array           []interface{}
	location        *Location
	subInstructions []forLoopSubInstruction
}

type directoryIterable struct {
	path            string
	resolver        FileResolver
	location        *Location
	subInstructions []forLoopSubInstruction
}

func (c *iterationContent) GetLocation() *Location {
	return c.location
}

func (c *iterationContent) Write(writer io.Writer, context Context) ([]Content, error) {
	context.Set(c.variable, c.item)
	return c.contents, nil
}

func parseIterable(arg string, location *Location, resolver FileResolver) (*parsedIterable, error) {
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
	return &parsedIterable{arg: forArg, location: location, resolver: resolver,
		subInstructions: subInstructions}, nil
}

func (e *arrayIterable) getItems(context Context) []interface{} {
	array := e.array
	for _, subInstruction := range e.subInstructions {
		if sortBy := subInstruction.sortBy; sortBy != nil {
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
	return array
}

func (e *directoryIterable) getItems(context Context) ([]webFileWithContext, error) {
	webFilesCtx, err := e.filesWithContext(context)
	if err != nil {
		return nil, err
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

	return webFilesCtx, nil
}

func (e *directoryIterable) filesWithContext(context Context) ([]webFileWithContext, error) {
	_, webFiles, err := e.resolver.FilesIn(e.path, e.location)
	if err != nil {
		return nil, err
	}

	webFilesCtx := make([]webFileWithContext, len(webFiles))
	for i, wf := range webFiles {
		ctx := wf.Processed.ResolveContext(context, false)
		// we must create a new ref here otherwise the file ref will point to the loop ref, which changes!
		refToFile := wf
		webFilesCtx[i] = webFileWithContext{file: refToFile, context: ctx}
	}
	return webFilesCtx, nil
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
