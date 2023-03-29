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
	fmt.Printf("比如：\n")
	fmt.Printf("let add = fn(x,y) { x + y };\n")
	fmt.Printf("add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8));\n")
	fmt.Printf("return a + b * (c + d) / e - f;\n")
	fmt.Printf("5 < 4 != 3 > 4;\n")
	fmt.Printf("3 < 5 == false;\n")
	fmt.Printf("!(true == true);\n")

	repl.Start(os.Stdin, os.Stdout) //参数为系统的标准输入输出，
}
