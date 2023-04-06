package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const PORMPT = ">> "

// REPL 实现读取-求值-打印 循环
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in) //为文本 I/O 提供了缓冲区，读入一行给扫描器
	env := object.NewEnviroment()   //创建标识符的环境-域

	for {
		fmt.Fprintf(out, PORMPT)
		scanned := scanner.Scan() //读取缓冲区内容

		if !scanned {
			return
		}
		line := scanner.Text() //读取??
		l := lexer.New(line)   //字符串转为 lexer结构，l.NextToken()才会转换词法单元
		p := parser.New(l)     //语法解析 传入lexer结构文本 并初始化

		io.WriteString(out, "语法解析过程可视化输出：\n")
		program := p.ParseProgram() //开始语法解析处理程序
		if len(p.Errors()) != 0 {   //错误输出
			printParserErrors(out, p.Errors())
			continue
		}

		//ast树遍历求值
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, "\n求值结果:\n")
			io.WriteString(out, evaluated.Inspect()) //查看求值结果
			io.WriteString(out, "\n")
		}
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-""lyt""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " 语法解析错误:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
