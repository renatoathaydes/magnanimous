package mg

import (
	"errors"
	"fmt"
)

type magParams struct {
	inclusionChain []Location
	scope          Scope
	webFiles       WebFilesMap
}

func (m magParams) Get(name string) (interface{}, error) {
	scope := m.scope
	for scope != nil {
		v, ok := scope.Context()[name]
		if ok {
			return v, nil
		}
		scope = scope.Parent()
	}
	for _, f := range m.inclusionChain {
		file, ok := m.webFiles[f.Origin]
		if ok {
			// FIXME check the scopes within the including-file
			v, ok := file.Processed.Context()[name]
			if ok {
				return v, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Parameter '%s' cannot be resolved", name))
}
