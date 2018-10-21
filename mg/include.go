package mg

import (
	"fmt"
	"io"
	"log"
)

type IncludeInstruction struct {
	Path     string
	Origin   Location
	MarkDown bool
	resolver FileResolver
}

func NewIncludeInstruction(arg string, location Location, resolver FileResolver) *IncludeInstruction {
	return &IncludeInstruction{Path: arg, Origin: location, resolver: resolver}
}

func (c *IncludeInstruction) IsMarkDown() bool {
	return c.MarkDown
}

func (c *IncludeInstruction) String() string {
	return fmt.Sprintf("IncludeInstruction{%s}", c.Path)
}

func (c *IncludeInstruction) Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) error {
	path := c.resolver.Resolve(c.Path, c.Origin)
	//fmt.Printf("Including %s from %v : %s\n", c.Path, c.Origin, path)
	webFile, ok := files[path]
	if !ok {
		log.Printf("WARNING: (%s) include non-existent resource: %s", c.Origin.String(), c.Path)
		_, err := writer.Write([]byte(fmt.Sprintf("{{ include %s }}", c.Path)))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	} else {
		inclusionChain = append(inclusionChain, c.Origin)
		//ss:= inclusionChainToString(inclusionChain)
		//fmt.Printf("Chain: %s", ss)
		for _, f := range inclusionChain {
			if f.Origin == path {
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
	}
	return nil
}
