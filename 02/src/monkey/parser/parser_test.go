package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

/*LET语句测试*/
func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l) //该类型新分配的零值的指针
	program := p.ParseProgram()
	checkParseErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgrm return nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement) //传入参数s是不是LetStatement
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name { //??value这里应该放token的值，letStmt.Name标识符节点的值value!=预期的标识符
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name { //value这里应该放token的值
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

// 检测Parse解析errors中是否有错误，有输出并终止测试
func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error:%q", msg)
	}
	t.FailNow() //终止测试
}

/*return语句测试*/
func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 838383;
`
	l := lexer.New(input)
	p := New(l) //该类型新分配的零值的指针
	program := p.ParseProgram()
	checkParseErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgrm return nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement) //断言判断是否是ast.ReturnStatement语法树
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue //下一句
		}
		if returnStmt.TokenLiteral() != "return" { //类型正确，值错误
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}
