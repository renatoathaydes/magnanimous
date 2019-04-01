package mg

import (
	"fmt"
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"io"
	"log"
	"strings"
)

type IncludeInstruction struct {
	Text     string
	Path     string
	Origin   *Location
	Resolver FileResolver
}

func NewIncludeInstruction(arg string, location *Location, original string, resolver FileResolver) *IncludeInstruction {
	return &IncludeInstruction{Text: original, Path: arg, Origin: location, Resolver: resolver}
}

func (inc *IncludeInstruction) String() string {
	return fmt.Sprintf("IncludeInstruction{%s, %v, %v}", inc.Path, inc.Origin, inc.Resolver)
}

func (inc *IncludeInstruction) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	maybePath := pathOrEval(inc.Path, magParams{stack: stack, webFiles: files})
	var actualPath string
	if s, ok := maybePath.(string); ok {
		actualPath = s
	} else {
		log.Printf("WARNING: path expression evaluated to invalid value: %v", maybePath)
		_, err := writer.Write([]byte(inc.Text))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
		return nil
	}
	path := inc.Resolver.Resolve(actualPath, inc.Origin, stack.NearestLocation())
	//fmt.Printf("Including %s from %v : %s\n", inc.Path, inc.Origin, path)
	webFile, ok := files.WebFiles[path]
	if !ok {
		log.Printf("WARNING: (%s) include non-existent resource: %s", inc.Origin.String(), actualPath)
		_, err := writer.Write([]byte(inc.Text))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	} else {
		stack = stack.Push(inc.Origin, false)
		err := detectCycle(stack, actualPath, path, inc.Origin)
		if err != nil {
			return err
		}
		err = webFile.Write(writer, files, stack)
		if err != nil {
			return err
		}
	}
	return nil
}

func pathOrEval(path string, params magParams) interface{} {
	startIndex := -1
	if strings.HasPrefix(path, "eval ") {
		startIndex = 5
	}
	if strings.HasPrefix(path, "\"") ||
		strings.HasPrefix(path, "`") ||
		strings.HasPrefix(path, "[") {
		startIndex = 0
	}
	if startIndex != -1 {
		// treat rest of argument as an expression that evaluates to a path
		res, err := expression.Eval(path[startIndex:], params)
		if err != nil {
			log.Printf("WARNING: eval expression error: %v", err)
		} else {
			return res
		}
	}
	return path
}

func detectCycle(stack ContextStack, includedPath, absPath string, location *Location) error {
	for _, loc := range stack.locations {
		if loc.Origin == absPath {
			chain := inclusionChainToString(stack.locations)
			return &MagnanimousError{
				Code: InclusionCycleError,
				message: fmt.Sprintf(
					"Cycle detected! Inclusion of %s at %s comes back into itself via %s",
					includedPath, location.String(), chain),
			}
		}
	}
	return nil
}
