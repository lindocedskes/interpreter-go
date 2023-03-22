package ast

import (
	"bytes"
	"monkey/token"
)

// The base Node interfac
type Node interface { //返回关联词法单元的字面量
	TokenLiteral() string
	String() string //调试时打印节点
}

type Statement interface { //ast中一些实现语句接口
	Node
	statementNode()
}

type Expression interface { //ast中一些实现表达式接口
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string { //AST根节点
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string { //AST根节点
	var out bytes.Buffer //创建一个缓冲区
	for _, s := range p.Statements {
		out.WriteString(s.String()) //s遍历每个Statements，调用每条语句的String方法的返回值写入缓冲区
	}
	return out.String() //缓冲区以字符串返回
}

/*LET 语句 AST结构：LET <标识符> = <表达式>*/
type LetStatement struct {
	Token token.Token //token.LET 词法单元
	Name  *Identifier //标识符
	Value Expression  //let语句产生值的表达式
}

// LetStatement句子节点需要实现的接口
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer                     //创建一个缓冲区
	out.WriteString(ls.TokenLiteral() + " ") //let
	out.WriteString(ls.Name.String())        //x
	out.WriteString(" = ")                   //=
	if ls.Value != nil {
		out.WriteString(ls.Value.String()) //5
	}
	return out.String()

}

// 标识符
type Identifier struct {
	Token token.Token //token.IDENT 词法单元
	Value string      //为了简单，有些标识符也有值
}

// Identifier 标识符实现的也是表达式节点的接口
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	return i.Value
}

/*return 语句 AST结构： return <表达式>*/
type ReturnStatement struct {
	Token       token.Token //return 词法单元
	ReturnValue Expression  //<表达式>
}

// 实现Statement全部方法，即实现该接口
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer                     //创建一个缓冲区
	out.WriteString(rs.TokenLiteral() + " ") //return
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String()) //5
	}
	return out.String()
}

/*解析表达式*/

// 表达式语句-结构
type ExpressionStatement struct {
	Token      token.Token //该表达式中第一个词法单元
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// 解析表达式-整型字面量
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// 解析表达式-前缀表达式 !-
type PrefixExpression struct {
	Token    token.Token //该表达式中第一个词法单元 !-
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// 解析表达式-中缀表达式
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
