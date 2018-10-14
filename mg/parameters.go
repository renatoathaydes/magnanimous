package mg

import (
	"errors"
	"fmt"
)

type magParams struct {
	originFile     string
	inclusionChain []string
	webFiles       WebFilesMap
}

func (m magParams) Get(name string) (interface{}, error) {
	files := make([]string, 1, len(m.inclusionChain)+1)
	files[0] = m.originFile
	files = append(files, m.inclusionChain...)
	for _, file := range files {
		f, ok := m.webFiles[file]
		if ok {
			if v, y := f.Context[name]; y {
				return v, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Parameter '%s' cannot be resolved", name))
}
