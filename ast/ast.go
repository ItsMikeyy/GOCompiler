package ast

import (
	"mscript/token"
)

type Node interface {
	TokenLiteral() string //Used for debugging
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type LetStatement struct {
	Token token.Token //Let token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Token       token.Token //Return token
	ReturnValue Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		//Return string of literal
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
