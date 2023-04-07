package main

import (
	"fmt"
	"monkey/repl"
	"os"
	user2 "os/user"
)

func main() {
	user, err := user2.Current() //当前返回当前用户
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! 这是自定义语言的解释器!\n", user.Username)
	fmt.Printf("该语言支持函数、高阶函数、闭包、整数，以及算术运算\n")
	fmt.Printf("比如：\n")
	fmt.Printf("let add = fn(x,y) { return x + y };\n")
	fmt.Printf("add(1 + 2 * (3 + 4) / 5 - 6, add(6, 7 * 8));\n")
	fmt.Printf("fn(x) {x==10}(10);\n")
	fmt.Printf("闭包例子：\n")
	fmt.Printf("let newAddr = fn(x) {fn(y) { x+y }};\n")
	fmt.Printf("let addTwo = newAddr(2);\n")
	fmt.Printf("addTwo(2);\n")
	fmt.Printf("嵌套函数例子：\n")
	fmt.Printf("let add = fn(a,b) { a + b };\n")
	fmt.Printf("let applyFunc = fn(a,b,func) { func(a ,b) };\n")
	fmt.Printf("applyFunc(2,2,add);\n")

	repl.Start(os.Stdin, os.Stdout) //参数为系统的标准输入输出，
}
