package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PORMPT = ">> "

// REPL 实现读取-求值-打印 循环
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in) //为文本 I/O 提供了缓冲区，读入一行给扫描器

	for {
		fmt.Fprintf(out, PORMPT)
		scanned := scanner.Scan() //读取缓冲区内容??
		if !scanned {
			return
		}
		line := scanner.Text() //读取??
		l := lexer.New(line)   //按行转换词法单元

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok) //无错误则输出tok词法单元的值
		}
	}
}
