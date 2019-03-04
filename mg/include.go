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

func (inc *IncludeInstruction) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	path := inc.Resolver.Resolve(inc.Path, *inc.Origin)
	//fmt.Printf("Including %s from %v : %s\n", inc.Path, inc.Origin, path)
	webFile, ok := files.WebFiles[path]
	if !ok {
		log.Printf("WARNING: (%s) include non-existent resource: %s", inc.Origin.String(), inc.Path)
		_, err := writer.Write([]byte(inc.Text))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	} else {
		stack = stack.Push(inc.Origin)
		//ss:= inclusionChainToString(inclusionChain)
		//fmt.Printf("Chain: %s", ss)
		err := inc.detectCycle(stack, path)
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

func (inc *IncludeInstruction) detectCycle(stack ContextStack, path string) error {
	for _, f := range stack.chain {
		if f.Location.Origin == path {
			chain := inclusionChainToString(stack.chain)
			return &MagnanimousError{
				Code: InclusionCycleError,
				message: fmt.Sprintf(
					"Cycle detected! Inclusion of %s at %s comes back into itself via %s",
					inc.Path, inc.Origin.String(), chain),
			}
		}
	}
	return nil
}
