package mg

func NewContextStack(context Context) ContextStack {
	items := make([]ContextStackItem, 1)
	items[0] = ContextStackItem{Context: context}
	return ContextStack{items}
}

func (c *ContextStack) Push(location *Location) ContextStack {
	item := ContextStackItem{Location: location, Context: &MapContext{}}
	items := append(c.chain, item)
	return ContextStack{items}
}

func (c *ContextStack) Top() *ContextStackItem {
	if len(c.chain) == 0 {
		return nil
	}
	return &c.chain[len(c.chain)-1]
}

func (c *ContextStack) GetContextAt(index int) Context {
	if len(c.chain) > 0 {
		return c.chain[0].Context
	}
	return nil
}

func (c *ContextStack) Size() int {
	return len(c.chain)
}
