package mg

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type componentContext struct {
	contents       []Content
	cachedContents *string
	Map            map[string]interface{}
	writeContext   *writeContext
}

type writeContext struct {
	files          WebFilesMap
	inclusionChain []ContextStackItem
}

type Component struct {
	Location Location
	context  componentContext
	Text     string
	Path     string
	Origin   Location
	Resolver FileResolver
}

var _ Context = (*componentContext)(nil)
var _ Content = (*Component)(nil)
var _ ContentContainer = (*Component)(nil)

func (c *Component) AppendContent(content Content) {
	c.context.contents = append(c.context.contents, content)
}

func (c *Component) GetContents() []Content {
	return c.context.contents
}

func (c *Component) Context() Context {
	return &c.context
}

func (c *Component) Parent() Scope {
	return c.parent
}

func (c *componentContext) Get(name string) (interface{}, bool) {
	if name == "__contents__" {
		return c.content()
	}
	v, ok := c.Map[name]
	return v, ok
}

func (c *componentContext) Set(name string, value interface{}) {
	c.Map[name] = value
}

func (c *componentContext) IsEmpty() bool {
	return len(c.Map) == 0 && len(c.contents) == 0
}

func (c *componentContext) mixInto(other Context) {
	// components do not mix into their surrounding scope
}

func (c *componentContext) content() (string, bool) {
	if c.cachedContents != nil {
		return *c.cachedContents, true
	}

	if wctx := c.writeContext; wctx != nil {
		// resolve contents and cache it
		var writer strings.Builder
		for _, content := range c.contents {
			err := content.Write(&writer, wctx.files, wctx.inclusionChain)
			if err != nil {
				return "", false
			}
		}
		contents := writer.String()
		c.cachedContents = &contents
		return contents, true
	} else {
		panic("Tried to write component before write context was set")
	}
}

func NewComponentInstruction(arg string, location Location, original string,
	scope Scope, resolver FileResolver) Content {
	compContext := componentContext{Map: make(map[string]interface{}, 2)}
	return &Component{
		Path:     arg,
		Location: location,
		Resolver: resolver,
		parent:   scope,
		context:  compContext,
		Text:     original,
		Origin:   location,
	}
}

func (c *Component) Write(writer io.Writer, files WebFilesMap, stack ContextStack) error {
	path := c.Resolver.Resolve(c.Path, c.Origin)
	//fmt.Printf("Including %s from %v : %s\n", c.Path, c.Origin, path)
	componentFile, ok := files.WebFiles[path]
	if !ok {
		log.Printf("WARNING: (%s) refers to a non-existent Component: %s", c.Origin.String(), c.Path)
		_, err := writer.Write([]byte(c.Text))
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	} else {
		inclusionChain = append(inclusionChain, ContextStackItem{Location: &c.Origin, scope: c})
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

		c.context.writeContext = &writeContext{files: files, inclusionChain: inclusionChain}
		runSideEffects(c, &files, inclusionChain)
		err := writeContents(componentFile.Processed, writer, files, inclusionChain, true)
		if err != nil {
			return err
		}
	}
	return nil
}
