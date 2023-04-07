package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	NULL = &object.Null{} //空值类型的实例
	//bool AST求值优化，用引用避免每次求值都要新建实例
	TRUE  = &object.Boolean{Value: true}  //新建bool实例，TRUE为实例的引用
	FALSE = &object.Boolean{Value: false} //返回都通过引用共用该实例
)

// 对ast语法树进行遍历求值,*object.Environment 求值对应的环境 -全局域和局部域
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) { //传入ast语法树的类型
	case *ast.Program: //开始都是Program节点
		return evalProgram(node, env) //开始都是Program节点，传入Statements，逐句解析
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env) //表达式节点：进一步解析表达式，ast往下
	case *ast.PrefixExpression: //前缀节点
		right := Eval(node.Right, env)
		if isError(right) { //如果Eval解析错误，返回Error节点，及时抛出
			return right
		}
		return evalPrefixExpression(node.Operator, right) //表达式节点：进一步解析表达式，ast往下
	case *ast.InfixExpression: //中缀节点
		left := Eval(node.Left, env)
		if isError(left) { //如果Eval解析错误，返回Error节点，及时抛出
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) { //如果Eval解析错误，返回Error节点，及时抛出
			return right
		}
		return evalInfixExpression(node.Operator, left, right) //表达式节点：进一步解析表达式，ast往下
	case *ast.BlockStatement: //表达式-区块节点{}
		return evalBlockStatement(node, env)
	case *ast.IfExpression: //表达式节点 -if-esle节点
		return evalIfExpression(node, env)
	case *ast.ReturnStatement: //return节点
		val := Eval(node.ReturnValue, env)
		if isError(val) { //如果Eval解析错误，返回Error节点，及时抛出
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env) //解析letAST的value指向的表达式节点
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.FunctionLiteral: //定义函数——函数字面量'fn' AST
		params := node.Parameters
		body := node.Body
		//封装 形参，函数体，局部域
		return &object.Function{Parameters: params, Env: env, Body: body} //仅是声明，返回封装的函数
	case *ast.CallExpression: //调用函数 AST
		function := Eval(node.Function, env) //函数字面量(fn)和函数名的标识符，封装为FUNCTION类型，函数名的标识符的value（也是*ast.FunctionLiteral）会被解析返回FUNCTION
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env) //1. 对参数求值，node.Arguments函数的参数
		if len(args) == 1 && isError(args[0]) {      //遇到错误，停止求值
			return args[0]
		}
		return applyFunction(function, args) //调用函数，给入函数名（封装的FUNCTION类型）和参数集

	//终端节点
	case *ast.IntegerLiteral: //终端节点整数，返回值，以对象系统-原始数据类型 封装返回
		return &object.Integer{Value: node.Value}
	case *ast.Boolean: //终端节点布尔，返回值，以对象系统-原始数据类型 封装返回
		//return &object.Boolean{Value: node.Value}
		return nativeboolToBooleanObject(node.Value) //bool AST求值返回，共用本地实例
	//从标识符获取对应的值
	case *ast.Identifier:
		return evalIdentifier(node, env)
	}

	return nil
}

// 顶层程序语句集合
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env) //目前只给一句

		switch result := result.(type) {
		case *object.ReturnValue: //解析到return语句，停止求值，返回return 表达式的解析结果
			return result.Value //返回嵌套语句的return类，解包执行return
		case *object.Error:
			return result
		}
	}
	return result
}

// 嵌套语句集合{}
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env) //目前只给一句

		if result != nil { //嵌套语句解析到return || error 语句，停止求值，返回return 表达式的解析结果 ||error信息
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

// bool AST求值返回，共用本地实例
func nativeboolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// 前缀节点AST 求值 支持！-
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!": //!
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	//错误处理
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// 中缀节点AST 求值 支持 +-*/
func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ: //两边都是数字
		return evalIntergerInfixExpression(operator, left, right)
	case operator == "==": //两边不全是数字，现在情况是都是布尔值的 ==运算支持
		return nativeboolToBooleanObject(left == right) //布尔值相同，指针指向同一个 ==运算为真
	case operator == "!=":
		return nativeboolToBooleanObject(left == right)
	//错误处理
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// 前缀节点AST 求值 ! 取非操作 逻辑实现：返回逻辑相反的值。
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// 前缀节点AST 求值 - 取反操作 逻辑实现：返回数值相反的值。
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ { //右节点必须是int
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// 中缀节点AST 求值 +-*/操作 逻辑实现
func evalIntergerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case ">":
		return nativeboolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeboolToBooleanObject(leftVal < rightVal)
	case "==":
		return nativeboolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeboolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// if节点AST 求值
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) { //如果Eval解析错误，返回Error节点，及时抛出
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env) //执行解析真值对应语句集
	} else if ie.Alternative != nil { //flase 且else存在
		return Eval(ie.Alternative, env)
	} else { //else不存在
		return NULL
	}
}

// if节点AST 求值 -辅助函数：判真值,非空和不是FALSE都为true
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// 错误处理
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// 判断错误
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// 从环境中查找标识符对应的值 map{标识符,值}
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value) //node.Value存标识符string
	if !ok {
		return newError("identifier notfound: " + node.Value)
	}
	return val
}

// 调用函数，对参数求值 参数是表达式集合
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps { //对调用函数的各个参数求值
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated} //error类
		}
		result = append(result, evaluated) //求值结果集合
	}
	return result
}

// 调用函数*ast.CallExpression处理返回，给入函数名（封装的FUNCTION类型或？？）和参数集
// 求值函数体
func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function) //??标识符
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args) //参数绑定，形参和实参，并扩展域
	evaluated := Eval(function.Body, extendedEnv)    //函数体求值
	return unwrapReturnValue(evaluated)              //有无return语句的处理
}

// 参数绑定，形参和实参，并扩展域
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclodedEnvironment(fn.Env) //创建基于外部域的新内部域
	//*object.Function 内部的形参，args[]外部的实参(已经被求值过)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

// 函数体有无return语句处理
func unwrapReturnValue(obj object.Object) object.Object {
	// 函数体有return语句，解包返回上层
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}
