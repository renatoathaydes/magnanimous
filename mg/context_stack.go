package mg

// NewContextStack creates a stack with a single context in it.
func NewContextStack(context Context) ContextStack {
	ctxs := make([]Context, 1, 10)
	ctxs[0] = context
	return ContextStack{contexts: ctxs}
}

// Push a new item on the scope stack.
//
// Only provide a location if this scope is including another file, otherwise provide nil.
// If createScope is true, push a new [Context] onto the stack, otherwise just update the inclusion location stack.
func (c *ContextStack) Push(location *Location, createScope bool) ContextStack {
	if location != nil {
		c.locations = append(c.locations, *location)
	}
	if createScope {
		scopes := append(c.contexts, NewContext())
		return ContextStack{locations: c.locations, contexts: scopes}
	}
	return *c
}

// Top gives the top element on the stack.
func (c *ContextStack) Top() Context {
	return c.contexts[len(c.contexts)-1]
}

// GetContextAt returns the [Context] at the given index on the stack (0 is the top of the stack).
func (c *ContextStack) GetContextAt(index int) Context {
	l := len(c.contexts)
	return c.contexts[l-1-index]
}

// Size of the stack.
func (c *ContextStack) Size() int {
	return len(c.contexts)
}
