package mg

import (
	"fmt"
	"path/filepath"
)

type mapContext struct {
	ctx map[string]interface{}
	str *string
}

var _ Context = (*mapContext)(nil)

// NewContext creates a simple [Context] based on a map.
func NewContext() Context {
	m := make(map[string]interface{}, 10)
	return &mapContext{ctx: m}
}

func (m *mapContext) ToStack() *ContextStack {
	stack := NewContextStack(m)
	return &stack
}

// newFileContext creates a simple Context that, when evaluated, resolves to a file path.
// This allows for-loops to evaluate all file paths in a directory.
func newFileContext(file *ProcessedFile, base Context) Context {
	var str string
	if file != nil {
		path := file.Path
		if file.NewExtension != "" {
			path = changeFileExt(path, file.NewExtension)
		}
		// try to figure out a valid link to the file
		s, err := filepath.Rel(file.BasePath, path)
		if err == nil {
			str = "/" + s
		} else {
			str = path
		}
	}
	stack := NewContextStack(base)
	stack.push(&mapContext{ctx: make(map[string]interface{}, 10), str: &str})
	return &stack
}

func (m *mapContext) Get(name string) (interface{}, bool) {
	v, ok := m.ctx[name]
	return v, ok
}

func (m *mapContext) Remove(name string) interface{} {
	v, _ := m.ctx[name]
	delete(m.ctx, name)
	return v
}

func (m *mapContext) Set(name string, value interface{}) {
	m.ctx[name] = value
}

func (m *mapContext) IsEmpty() bool {
	return len(m.ctx) == 0
}

func (m *mapContext) String() string {
	if m.str != nil {
		return *m.str
	}
	return fmt.Sprint(m.ctx)
}
