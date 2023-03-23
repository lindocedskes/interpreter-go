package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const ( //设置运算符优先级
	_ int = iota //iota 是一个预先声明的标识符,当前 const 规范的无类型整数序号
	LOWEST
	EQUALS      // ==
	LESSGREATER //> or<
	SUM         //+
	PRODUCT     //*
	PREFIX      //-X or !X
	CALL        //myFunction(X)
)

var precedences = map[token.TokenType]int{ //{类型：优先级}映射
	token.EQ:       EQUALS,      //=
	token.NOT_EQ:   EQUALS,      //!=
	token.LT:       LESSGREATER, //<
	token.GT:       LESSGREATER, //>
	token.PLUS:     SUM,         //+
	token.MINUS:    SUM,         //-
	token.SLASH:    PRODUCT,     // /
	token.ASTERISK: PRODUCT,     //*
}

type Parser struct {
	l         *lexer.Lexer //输入文本，调用l.nextToken读取词法单元
	errors    []string     //存放错误
	curToken  token.Token  //当前词法单元
	peekToken token.Token  //下一个词法单元

	prefixParseFns map[token.TokenType]prefixParseFn //检查token类型映射是否有管理的解析函数
	infixParseFns  map[token.TokenType]infixParseFn  //实现token类型映射对应执行函数类型
}

// 初始化
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}} //[]string{}空错误切片

	//读取2个词法单元，以设置curToken和peekToken
	p.nextToken() //0，0->0,1 ;1表示指向第一个token
	p.nextToken() //0,1->1,2

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn) //初始化前缀映射
	p.registerPrefix(token.IDENT, p.parseIdentifier)           //标识符添加{token类型:解析函数}映射
	p.registerPrefix(token.INT, p.parseIntegerLiteral)         //整数字面量添加{token类型:解析函数}映射
	p.registerPrefix(token.BANG, p.parsePrefixExpression)      //前缀运算符（!）{token类型:解析函数}映射
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)     //前缀运算符（-）{token类型:解析函数}映射

	p.infixParseFns = make(map[token.TokenType]infixParseFn) //初始化中缀映射
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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
	case token.RETURN:
		return p.parseReturnStatement() //调用对RETURN语句的语法分析
	default:
		return p.parseExpressionStatement() //调用对Expression语句的语法分析
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

// Expression语句的语法分析
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken} //return语句根节点
	p.nextToken()

	//TODO  先跳过表达式的处理,直到遇到分号结束';'
	if !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// Expression语句的语法分析
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement")) //添加跟踪语句，执行结束后输出

	stmt := &ast.ExpressionStatement{Token: p.curToken} //return语句根节点
	stmt.Expression = p.parseExpression(LOWEST)         //传入前一个运算符优先级，初始为最低 例：1+2  +与LOWEST比较

	//TODO  分号可选';'
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

// 定义函数类型，前缀解析函数和中缀解析函数，映射：map[token.TokenType]prefixParseFn
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// 向{token类型:解析函数}映射中添加内容
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// 检查前缀位置是否有token类型关联的解析函数
func (p *Parser) parseExpression(precedence int) ast.Expression { //传入前一个运算符优先级，例：1+2+3 的第一个+
	defer untrace(trace("parseExpression")) //添加跟踪语句，执行结束后输出

	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type) //前缀解析函数-没有加入error消息
		return nil
	}
	leftExp := prefix() //存储前缀解析函数的返回值指针，数字返回*ast.IntegerLiteral，标识符返回*ast.Identifier，前缀操作符返回*ast.PrefixExpression
	//不为； 且优先级高，例cur在1+2+3中指向2，前一个运算符优先级>=下一个运算符优先级，则不执行循环;也有前一个为数字 < 下一个运算符，执行
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() { //precedence由于递归是当前层的前一个优先级***容易❌，cur多次递归后移动可能很后面
		infix := p.infixParseFns[p.peekToken.Type] //查找下一个token类型对应中缀解析函数
		if infix == nil {
			return leftExp
		}

		p.nextToken() //移动1个，cur->peek 例 1+2+3，cur现在指向第1个+运算符

		leftExp = infix(leftExp) //执行中缀解析函数，leftExp为左节点，⚠️例：1+2+3 递归第二层返回后 指向是*ast.InfixExpression:(1+2)
	}
	return leftExp
}

// 前缀解析函数-没有加入error消息
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.Errors(), msg)
}

// 表达式-标识符解析函数-返回Identifier节点包含token和value值
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// 表达式-整数字面量解析函数-返回IntegerLiteral节点包含token和value值,value是int类型
func (p *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral")) //添加跟踪语句，执行结束后输出

	lit := &ast.IntegerLiteral{Token: p.curToken}
	//str转int64
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

// 表达式-前缀运算符解析函数
func (p *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression")) //添加跟踪语句，执行结束后输出

	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken() //导致解析完表达式后，指向表达式最后一个token

	expression.Right = p.parseExpression(PREFIX) //递归解析前缀表达式，PREFIX这个问题留给下一节解决
	return expression
}

// 查看下一个token优先级，返回int，未定义的默认最低
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// 查看当前token优先级，返回int，未定义的默认最低
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// 表达式-中缀运算符解析函数, 需传入左表达式
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression { //例1+2，传入1为*ast.IntegerLiteral， 节点
	defer untrace(trace("parseInfixExpression")) //添加跟踪语句，执行结束后输出

	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left, //左操作数通过传入参数，放入节点
	}

	precedence := p.curPrecedence()                  //cur当前指向运算符，保存前一个运算符优先级
	p.nextToken()                                    //cur指向右边操作数
	expression.Right = p.parseExpression(precedence) //前面一个运算符作为参数传入，并递归继续解析右端表达式

	return expression
}
