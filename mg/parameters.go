package mg

type magParams struct {
	inclusionChain []InclusionChainItem
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
	for i := len(m.inclusionChain) - 1; i >= 0; i-- {
		scope = m.inclusionChain[i].scope
		for scope != nil {
			v, ok := scope.Context()[name]
			if ok {
				return v, true
			}
			scope = scope.Parent()
		}
	}
	return nil, false
}
