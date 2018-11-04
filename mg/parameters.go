package mg

type magParams struct {
	inclusionChain []InclusionChainItem
	scope          Scope
	webFiles       *WebFilesMap
}

func (m magParams) Get(name string) (interface{}, bool) {
	v, ok := searchParamInScope(m.scope, name)
	if ok {
		return v, true
	}
	for i := len(m.inclusionChain) - 1; i >= 0; i-- {
		v, ok = searchParamInScope(m.inclusionChain[i].scope, name)
		if ok {
			return v, true
		}
	}
	if m.webFiles.GlobalContext != nil {
		v, ok := m.webFiles.GlobalContext[name]
		return v, ok
	}
	return nil, false
}

func searchParamInScope(scope Scope, name string) (interface{}, bool) {
	for scope != nil {
		v, ok := scope.Context()[name]
		if ok && v != nil {
			return v, true
		}
		scope = scope.Parent()
	}
	return nil, false
}
