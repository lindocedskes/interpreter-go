package object

// 环境：存储 {标识符,值}
type Environment struct {
	store map[string]Object
	outer *Environment
}

// 环境——产生一个Environment-域 实例
func NewEnviroment() *Environment {
	s := make(map[string]Object)  //产生一个map
	return &Environment{store: s} //产生一个Environment实例
}

// 实现函数的局部域，传入函数为外部的域，内部新建一个
func NewEnclodedEnvironment(outer *Environment) *Environment {
	env := NewEnviroment()
	env.outer = outer
	return env
}

// 环境——域 中查找标识符对应的值
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil { //内部域没找到，且存在外部域
		obj, ok = e.outer.Get(name) //外部域找
	}
	return obj, ok
}

// 环境——域 中存放 标识符对应的值
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
