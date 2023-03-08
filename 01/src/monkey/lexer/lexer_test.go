package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	//测试字符串，``表示字符串字面量，不转义
	input := `let five=5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
`
	tests := []struct { //结构体抽象结构,测试返回结果是否匹配
		expectedType    token.TokenType
		expectedLiteral string
	}{ //实例化，期望得到的值的集合，测试
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()             //转换当前*Lexer的正在查看的字符ch，返回为对应Token结构包含类型和值，postion+1
		if tok.Type != tt.expectedType { //类型不一致报错
			t.Fatalf("tests[%d]-tokentype wrong. expected=%q,got =%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral { //值不一致报错
			t.Fatalf("tests[%d]-tokentype wrong. expected=%q,got =%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
