package object

// 建立对象系统,原始数据类型

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

type ObjectType string

const (
	//类型被封装，对应一个封装结构体
	INTEGER_OBJ      = "INTEGER" //整数类型
	BOOLEAN_OBJ      = "BOOLEAN" //布尔类型
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION" //函数封装
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

// 函数类 封装 形参，函数体，局部域
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
