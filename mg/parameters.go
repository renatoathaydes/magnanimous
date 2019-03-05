package mg

type magParams struct {
	stack    ContextStack
	webFiles WebFilesMap
}

func (m magParams) Get(name string) (interface{}, bool) {
	for i := 0; i < m.stack.Size(); i++ {
		ctx := m.stack.GetContextAt(i)
		v, ok := ctx.Get(name)
		if ok {
			return v, true
		}
	}
	if m.webFiles.GlobalContext != nil {
		return m.webFiles.GlobalContext.Get(name)
	}
	return nil, false
}
