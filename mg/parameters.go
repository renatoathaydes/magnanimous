package mg

type magParams struct {
	inclusionChain []Location
	scope          Scope
	webFiles       WebFilesMap
}

func (m magParams) Get(name string) (interface{}, bool) {
	scope := m.scope
	for scope != nil {
		v, ok := scope.Context()[name]
		if ok {
			return v, true
		}
		scope = scope.Parent()
	}
	for _, f := range m.inclusionChain {
		file, ok := m.webFiles[f.Origin]
		if ok {
			// FIXME check the scopes within the including-file
			v, ok := file.Processed.Context()[name]
			if ok {
				return v, true
			}
		}
	}
	return nil, false
}
