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
	fmt.Printf("Hello %s! 这是自定义语言的语法分析器!\n", user.Username)
	fmt.Printf("你可以自由输入类似下面语言的代码：\n")
	fmt.Printf("比如：let add = fn(x,y) { x + y };\n")
	repl.Start(os.Stdin, os.Stdout) //参数为系统的标准输入输出，
}
