package mg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strings"
)

type IncludeB64Instruction struct {
	Text     string
	Path     string
	IsPlain  bool
	Origin   *Location
	Resolver FileResolver
}

var _ Inclusion = (*IncludeB64Instruction)(nil)
var _ Content = (*IncludeB64Instruction)(nil)

func NewIncludeB64Instruction(arg string, location *Location, original string, resolver FileResolver) *IncludeB64Instruction {
	path, isPlain := parseIncludeB64Arg(arg, location)
	return &IncludeB64Instruction{Text: original, Path: path, IsPlain: isPlain, Origin: location, Resolver: resolver}
}

func parseIncludeB64Arg(arg string, location *Location) (path string, isPlain bool) {
	if strings.HasPrefix(arg, "(") {
		idx := strings.Index(arg, ")")
		if idx > 0 {
			subInstruction := strings.TrimSpace(arg[1:idx])
			if subInstruction == "plain" {
				isPlain = true
			} else if subInstruction != "" {
				log.Printf("WARNING: (%s) unrecognizable includeB64 sub-instruction: %s",
					location.String(), subInstruction)
			}
			path = strings.TrimSpace(arg[idx+1:])
		} else {
			path = arg
		}
	} else {
		path = arg
	}
	return
}

func (inc *IncludeB64Instruction) String() string {
	return fmt.Sprintf("IncludeB64Instruction{%s, %v, %v}", inc.Path, inc.Origin, inc.Resolver)
}

func (inc *IncludeB64Instruction) GetPath() string {
	return inc.Path
}

func (inc *IncludeB64Instruction) GetLocation() *Location {
	return inc.Origin
}

func (inc *IncludeB64Instruction) IsScoped() bool {
	return true
}

func (inc *IncludeB64Instruction) Write(writer io.Writer, context Context) ([]Content, error) {
	webFile, err := getInclusionByPath(inc, inc.Resolver, context, true)
	if err != nil {
		return nil, err
	}
	return nil, writeb64(webFile, writer, context, inc.IsPlain)
}

func writeb64(webFile *WebFile, writer io.Writer, context Context, isPlain bool) error {
	var b bytes.Buffer
	b.Grow(512)
	err := webFile.Write(&b, context.ToStack(), false, isPlain)
	if err != nil {
		return err
	}
	encoder := base64.NewEncoder(base64.StdEncoding, writer)
	defer encoder.Close()
	_, err = encoder.Write(b.Bytes())
	return err
}
