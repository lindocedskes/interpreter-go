package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l         *lexer.Lexer //输入文本，调用l.nextToken读取词法单元
	errors    []string     //存放错误
	curToken  token.Token  //当前词法单元
	peekToken token.Token  //下一个词法单元
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}} //[]string{}空错误切片

	//读取2个词法单元，以设置curToken和peekToken
	p.nextToken() //0，0->0,1 ;1表示指向第一个token
	p.nextToken() //0,1->1,2
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken() //l.NextToken() 输入文本转换为词法单元返回，并+1
}

func (p *Parser) ParseProgram() *ast.Program { //调用语法分析器入口

	program := &ast.Program{}              //AST根节点
	program.Statements = []ast.Statement{} //子结构体初始化
	for p.curToken.Type != token.EOF {     //遍历词法单元
		stmt := p.parseStatement() //语法分析一句，返回指向该句生成的AST的指针（子节点）
		if stmt != nil {
			program.Statements = append(program.Statements, stmt) //加入AST根节点的切片
		}
		p.nextToken()
	}
	return program

}

// 语法分析一句，返回指向该句生成的AST的指针（子节点）
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	//按词法单元类型进行不同的处理
	case token.LET:
		return p.parseLetStatement() //调用对LET语句的语法分析
	default:
		return nil
	}
}

// LET语句的语法分析
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken} //生成一个*ast.LetStatement节点
	if !p.expectPeek(token.IDENT) {              //expectPeek()期待的词法单元类型
		return nil
	}
	//let节点name放标识符节点
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) { //下一个词法单元类型是否是'='
		return nil
	}

	//TODO  先跳过表达式的处理,直到遇到分号结束';'
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// 当前词法单元类型是否是t类型
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 下一个的词法单元类型是否是t类型
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek(t)期待的下一个词法单元类型是否为t
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t) //对每一次期待检查，不是则添加错误到分析错误数组中
		return false
	}
}

// 返回语法分析的错误
func (p *Parser) Errors() []string {
	return p.errors
}

// 当下一个token与预期不符时，报错并加入错误队列
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
