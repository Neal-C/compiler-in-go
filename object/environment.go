package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	nameToValueMap := make(map[string]Object)
	return &Environment{
		store: nameToValueMap,
		outer: nil,
	}
}

func (self *Environment) Get(name string) (Object, bool) {
	value, ok := self.store[name]
	if !ok && self.outer != nil {
		value, ok = self.outer.Get(name)
	}
	return value, ok
}

func (self *Environment) Set(name string, value Object) Object {
	self.store[name] = value
	return value
}

func NewEnclosedEnvironment(outerEnv *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outerEnv
	return env
}
