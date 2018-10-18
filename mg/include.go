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
}

func NewIncludeInstruction(arg string, location Location) *IncludeInstruction {
	path := ResolveFile(arg, "source", location.Origin)
	return &IncludeInstruction{Path: path, Origin: location}
}

func (c *IncludeInstruction) IsMarkDown() bool {
	return c.MarkDown
}

func (c *IncludeInstruction) String() string {
	return fmt.Sprintf("IncludeInstruction{%s}", c.Path)
}

func (c *IncludeInstruction) Write(writer io.Writer, files WebFilesMap, inclusionChain []Location) error {
	webFile, ok := files[c.Path]
	if !ok {
		log.Printf("WARNING: (%s) include non-existent resource: %s", c.Origin.String(), c.Path)
		_, err := writer.Write([]byte(fmt.Sprintf("{{ include %s }}", c.Path)))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	} else {
		inclusionChain = append(inclusionChain, c.Origin)
		for _, f := range inclusionChain {
			if f.Origin == c.Path {
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
