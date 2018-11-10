package mg

import (
	"io"
)

type componentScope struct {
	context  map[string]interface{}
	contents []Content
	parent   Scope
}

type Component struct {
	Location Location
	include  *IncludeInstruction
	scope    componentScope
}

var _ Scope = (*componentScope)(nil)
var _ Scope = (*Component)(nil)
var _ Content = (*Component)(nil)
var _ ContentContainer = (*Component)(nil)

func (c *Component) AppendContent(content Content) {
	c.scope.AppendContent(content)
}

func (c *Component) Context() map[string]interface{} {
	return c.scope.Context()
}

func (c *Component) Parent() Scope {
	return c.scope.Parent()
}

func (c *Component) setParent(scope Scope) {
	c.scope.setParent(scope)
}

func (c *componentScope) AppendContent(content Content) {
	c.contents = append(c.contents, content)
}

func (c *componentScope) Context() map[string]interface{} {
	return c.context
}

func (c *componentScope) Parent() Scope {
	return c.parent
}

func (c *componentScope) setParent(scope Scope) {
	c.parent = scope
}

func (c *Component) GetContents() []Content {
	return c.scope.contents
}

func NewComponentInstruction(arg string, location Location, original string,
	scope Scope, resolver FileResolver) Content {
	compScope := componentScope{context: make(map[string]interface{}, 2), parent: scope}
	return &Component{
		Location: location,
		include:  NewIncludeInstruction(arg, location, original, &compScope, resolver),
		scope:    compScope,
	}
}

func (c *Component) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	runSideEffects(c, &files, inclusionChain)
	return c.include.Write(writer, files, inclusionChain)
}
