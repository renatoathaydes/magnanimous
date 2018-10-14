package mg

import (
	"errors"
	"fmt"
)

type magParams struct {
	origin         Location
	inclusionChain []Location
	webFiles       WebFilesMap
}

func (m magParams) Get(name string) (interface{}, error) {
	files := make([]Location, 1, len(m.inclusionChain)+1)
	files[0] = m.origin
	files = append(files, m.inclusionChain...)
	for _, file := range files {
		f, ok := m.webFiles[file.Origin]
		if ok {
			if v, y := f.Context[name]; y {
				return v, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Parameter '%s' cannot be resolved", name))
}
