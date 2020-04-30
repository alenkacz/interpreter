package object

type Environment struct {
	outer *Environment
	values map[string]Object
}

func NewEnvironment(outer *Environment) *Environment {
	return &Environment{
		values: make(map[string]Object),
		outer: outer,
	}
}

func (e *Environment) Get(key string) (Object, bool) {
	val, ok := e.values[key]
	if !ok && e.outer != nil {
		return e.outer.Get(key)
	}
	return val, ok
}

func (e *Environment) Set(key string, value Object) {
	e.values[key] = value
}