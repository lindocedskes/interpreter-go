package object

// 建立对象系统,原始数据类型

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER" //整数类型
	BOOLEAN_OBJ      = "BOOLEAN" //布尔类型
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
)

type Object interface { //
	Type() ObjectType //类型名方法
	Inspect() string  //值方法
}

// 整数类型
type Integer struct {
	Value int64
}

func (i Integer) Type() ObjectType { return INTEGER_OBJ }
func (i Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// 布尔类型
type Boolean struct {
	Value bool
}

func (b Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// 空值
type Null struct {
}

func (n Null) Type() ObjectType { return NULL_OBJ }
func (n Null) Inspect() string  { return "null" }

// 返回值
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// 异常处理ERROR
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR" + e.Message }

// 环境：存储 {标识符,值}
type Environment struct {
	store map[string]Object
}

// 环境——产生一个Environment-域 实例
func NewEnviroment() *Environment {
	s := make(map[string]Object)  //产生一个map
	return &Environment{store: s} //产生一个Environment实例
}

// 环境——域 中查找标识符对应的值
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// 环境——域 中存放 标识符对应的值
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
