package mg

type MapContext struct {
	Map map[string]interface{}
}

var _ Context = (*MapContext)(nil)

func (m *MapContext) Get(name string) (interface{}, bool) {
	v, ok := m.Map[name]
	return v, ok
}

func (m *MapContext) Set(name string, value interface{}) {
	m.Map[name] = value
}

func (m *MapContext) IsEmpty() bool {
	return len(m.Map) == 0
}

func (m *MapContext) mixInto(other Context) {
	for k, v := range m.Map {
		other.Set(k, v)
	}
}
