package mg

type magParams struct {
	stack    ContextStack
	webFiles WebFilesMap
}

func (m magParams) Get(name string) (interface{}, bool) {
	for i := 0; i < m.stack.Size(); i++ {
		ctx := m.stack.GetContextAt(i)
		if v, ok := ctx.Get(name); ok {
			return v, true
		}
	}
	return nil, false
}
