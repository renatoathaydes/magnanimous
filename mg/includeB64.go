package mg

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
)

type IncludeB64Instruction struct {
	Text     string
	Path     string
	Origin   *Location
	Resolver FileResolver
}

func NewIncludeB64Instruction(arg string, location *Location, original string, resolver FileResolver) *IncludeB64Instruction {
	return &IncludeB64Instruction{Text: original, Path: arg, Origin: location, Resolver: resolver}
}

func (inc *IncludeB64Instruction) String() string {
	return fmt.Sprintf("IncludeB64Instruction{%s, %v, %v}", inc.Path, inc.Origin, inc.Resolver)
}

func (inc *IncludeB64Instruction) Write(writer io.Writer, stack ContextStack) error {
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
		err = writeb64(webFile, writer, stack)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeb64(webFile *WebFile, writer io.Writer, stack ContextStack) error {
	bytes, err := asBytes(webFile.Processed.GetContents(), stack)
	if err != nil {
		return err
	}
	encoder := base64.NewEncoder(base64.StdEncoding, writer)
	defer encoder.Close()
	_, err = encoder.Write(bytes)
	return err
}
