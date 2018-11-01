package mg

import (
	"fmt"
	"io"
	"log"
)

type IncludeInstruction struct {
	Path     string
	Origin   Location
	scope    Scope
	Resolver FileResolver
}

func NewIncludeInstruction(arg string, location Location, scope Scope, resolver FileResolver) *IncludeInstruction {
	return &IncludeInstruction{Path: arg, Origin: location, scope: scope, Resolver: resolver}
}

func (c *IncludeInstruction) String() string {
	return fmt.Sprintf("IncludeInstruction{%s, %v, %v}", c.Path, c.Origin, c.Resolver)
}

func (c *IncludeInstruction) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	path := c.Resolver.Resolve(c.Path, c.Origin)
	//fmt.Printf("Including %s from %v : %s\n", c.Path, c.Origin, path)
	webFile, ok := files[path]
	if !ok {
		log.Printf("WARNING: (%s) include non-existent resource: %s", c.Origin.String(), c.Path)
		_, err := writer.Write([]byte(fmt.Sprintf("{{ include %s }}", c.Path)))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	} else {
		inclusionChain = append(inclusionChain, InclusionChainItem{Location: &c.Origin, scope: c.scope})
		//ss:= inclusionChainToString(inclusionChain)
		//fmt.Printf("Chain: %s", ss)
		for _, f := range inclusionChain {
			if f.Location.Origin == path {
				chain := inclusionChainToString(inclusionChain)
				return &MagnanimousError{
					Code: InclusionCycleError,
					message: fmt.Sprintf(
						"Cycle detected! Inclusion of %s at %s comes back into itself via %s",
						c.Path, c.Origin.String(), chain),
				}
			}
		}
		err := webFile.Write(writer, files, inclusionChain)
		if err != nil {
			return err
		}

		// mix in the context of the include file into the surrounding context
		for k, v := range webFile.Processed.Context() {
			c.scope.Context()[k] = v
		}
	}
	return nil
}
