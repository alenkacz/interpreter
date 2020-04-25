package object

type Environment struct {
	outer *Environment
	values map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]Object),
	}
}

func (e *Environment) Get(key string) (Object, bool) {
	val, ok := e.values[key]
	return val, ok
}

func (e *Environment) Set(key string, value Object) {
	e.values[key] = value
}