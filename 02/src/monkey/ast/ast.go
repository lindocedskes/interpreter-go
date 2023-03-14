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

/*LET 语句 AST结构：LET <标识符> = <表达式>*/
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
	Value string      //为了简单，有些标识符也有值
}

// Identifier 标识符实现的也是表达式节点的接口
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

/*return 语句 AST结构： return <表达式>*/
type ReturnStatement struct {
	Token       token.Token //return 词法单元
	ReturnValue Expression  //<表达式>
}

// 实现Statement全部方法，即实现该接口
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
