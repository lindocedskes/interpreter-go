package ast

import (
	"monkey/token"
)

// The base Node interfac
type Node interface { //返回关联词法单元的字面量
	TokenLiteral() string
}

type Statement interface { //ast中一些实现语句接口
	Node
	statementNode()
}

type Expression interface { //ast中一些实现表达式接口
	Node
	expressionNode() string
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

type LetStatement struct {
	Token token.Token //token.LET 词法单元
	Name  *Identifier //标识符
	Value Expression  //let语句产生值的表达式
}

// LetStatement句子节点需要实现的接口
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token //token.IDENT 词法单元
	Value string      //
}

// Identifier 标识符实现的也是表达式节点的接口
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
