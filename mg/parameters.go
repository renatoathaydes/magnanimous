package mg

import "github.com/renatoathaydes/magnanimous/mg/expression"

type magParams struct {
	location     *Location
	stack        ContextStack
	fileResolver FileResolver
}

var _ expression.Context = (*magParams)(nil)
var _ expression.PathFinder = (*magParams)(nil)

func (m *magParams) Get(name string) (interface{}, bool) {
	for i := 0; i < m.stack.Size(); i++ {
		ctx := m.stack.GetContextAt(i)
		if v, ok := ctx.Get(name); ok {
			return v, true
		}
	}
	return nil, false
}

func (m *magParams) File(path string) (*WebFile, bool) {
	path = m.fileResolver.Resolve(path, m.location, m.stack.NearestLocation())
	return m.fileResolver.Get(path)
}

func (m *magParams) Path(path string) (*expression.Path, bool) {
	if f, ok := m.File(path); ok {
		return &expression.Path{Value: path, LastUpdated: f.Processed.LastUpdated}, true
	}
	return nil, false
}

func (w *WebFilesMap) findFile(path string) (*expression.Path, bool) {
	if f, ok := w.WebFiles[path]; ok {
		return &expression.Path{Value: f.Processed.Path, LastUpdated: f.Processed.LastUpdated}, true
	}
	return nil, false
}
