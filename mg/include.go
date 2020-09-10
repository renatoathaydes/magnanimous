package mg

import (
	"fmt"
	"io"
	"log"
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

func (inc *IncludeInstruction) Write(writer io.Writer, stack ContextStack) error {
	params := magParams{stack: stack, location: inc.Origin, fileResolver: inc.Resolver}
	maybePath := pathOrEval(inc.Path, &params)
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
	webFile, ok := params.File(actualPath)
	if !ok {
		log.Printf("WARNING: (%s) include non-existent resource: %s", inc.Origin.String(), actualPath)
		_, err := writer.Write([]byte(inc.Text))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	} else {
		stack = stack.Push(inc.Origin, false)
		err := detectCycle(stack, actualPath, webFile.Processed.Path, inc.Origin)
		if err != nil {
			return err
		}
		err = webFile.Write(writer, stack)
		if err != nil {
			return err
		}
	}
	return nil
}
