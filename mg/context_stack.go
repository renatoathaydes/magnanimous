package mg

func NewContextStack(context Context) ContextStack {
	items := make([]ContextStackItem, 1)
	items[0] = ContextStackItem{Context: context}
	return ContextStack{items}
}

// Push a new item on the scope stack.
// Only provide a location if this scope is including another file.
func (c *ContextStack) Push(location *Location) ContextStack {
	item := ContextStackItem{Location: location, Context: CreateContext()}
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
	l := len(c.chain)
	return c.chain[l-1-index].Context
}

func (c *ContextStack) Size() int {
	return len(c.chain)
}
