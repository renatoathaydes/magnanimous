package mg

import (
	"errors"
	"fmt"
)

type magParams struct {
	originFile string
	webFiles   WebFilesMap
}

func (m magParams) Get(name string) (interface{}, error) {
	f, ok := m.webFiles[m.originFile]
	if ok {
		if v, y := f.Context[name]; y {
			return v, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Parameter '%s' cannot be resolved", name))
}
