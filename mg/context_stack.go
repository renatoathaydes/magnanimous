package mg

import "fmt"

// ContextStack is a stack of InclusionChainItems.
//
// Used to keep state when writing nested Content.
type ContextStack struct {
	locations []*Location
	contexts  []Context
}

var _ Context = (*ContextStack)(nil)

// PushResult is the result of calling the Push method.
//
// Each Push call must have a matching Pop call, which receives the PushResult from the Push call.
type PushResult struct {
	popLocation bool
	popContext  bool
}

// NewContextStack creates a stack with a single context in it.
func NewContextStack(context Context) ContextStack {
	ctxs := make([]Context, 1, 10)
	ctxs[0] = context
	return ContextStack{contexts: ctxs}
}

func (c *ContextStack) ToStack() *ContextStack {
	return c
}

// Push a new item on the scope stack.
//
// Only provide a location if this scope is including another file, otherwise provide nil.
// If createScope is true, push a new [Context] onto the stack, otherwise just update the inclusion location stack.
func (c *ContextStack) Push(location *Location, createScope bool) PushResult {
	pushedLocation := false
	if location != nil {
		if len(c.locations) == 0 || location.Origin != c.locations[len(c.locations)-1].Origin {
			c.locations = append(c.locations, location)
			pushedLocation = true
		}
	}
	if createScope {
		c.push(NewContext())
	}
	return PushResult{popLocation: pushedLocation, popContext: createScope}
}

func (c *ContextStack) push(ctx Context) {
	c.contexts = append(c.contexts, ctx)
}

// Pop both location and context.
func (c *ContextStack) Pop(pushResult PushResult) {
	if pushResult.popLocation {
		c.locations = c.locations[:len(c.locations)-1]
	}
	if pushResult.popContext {
		c.contexts = c.contexts[:len(c.contexts)-1]
	}
}

// Top gives the top element on the stack.
func (c *ContextStack) Top() Context {
	return c.contexts[len(c.contexts)-1]
}

// NearestLocation finds the nearest Location to the top of the Stack. Returns nil if no Location has been pushed.
func (c *ContextStack) NearestLocation() *Location {
	l := len(c.locations)
	if l == 0 {
		return nil
	}
	return c.locations[l-1]
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

func (c *ContextStack) Get(name string) (interface{}, bool) {
	for i := 0; i < c.Size(); i++ {
		ctx := c.GetContextAt(i)
		if v, ok := ctx.Get(name); ok {
			return v, true
		}
	}
	return nil, false
	return c.Top().Get(name)
}

// Set the value for the given name.
func (c *ContextStack) Set(name string, value interface{}) {
	c.Top().Set(name, value)
}

// Remove the value with the given name.
func (c *ContextStack) Remove(name string) interface{} {
	return c.Top().Remove(name)
}

// IsEmpty returns whether this context contains no values.
func (c *ContextStack) IsEmpty() bool {
	return c.Top().IsEmpty()
}

func (c *ContextStack) String() string {
	if len(c.contexts) == 0 {
		return "Empty ContextStack"
	}
	// try to use the top-context string just because it may resolve to a file-path's context
	// see mapContext
	top := c.Top()
	if s, ok := top.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprintf("ContextStack{%v, %v}", c.contexts, c.locations)
}
