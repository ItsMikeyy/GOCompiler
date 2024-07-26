package ast

import (
	"bytes"
	"mscript/token"
)

type Node interface {
	TokenLiteral() string //Used for debugging
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Let statements have the token 'let', The identifer (name) and the expression (value)
// let <identifier> = <expression>
type LetStatement struct {
	Token token.Token //Let token
	Name  *Identifier
	Value Expression
}

// return <expression>
type ReturnStatement struct {
	Token       token.Token //Return token
	ReturnValue Expression
}

type ExpressionStatement struct {
	Token      token.Token //First token of expression
	Expression Expression
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// <prefix><expression>
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

//<expression> <infix> <expression>

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

type Program struct {
	Statements []Statement
}

// Needed to satisfy the interface
func (ls *LetStatement) statementNode() {}

// Return the liteal of the let statement node
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Generating let statement string
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Keeps the ident token and the value of the identifer (name)
type Identifier struct {
	Token token.Token //token.IDENT
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// Return identifer name
func (i *Identifier) String() string {
	return i.Value
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

// Generating return statement string
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// Generating expression statement string
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// Return token literal for int
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}
func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

func (oe *InfixExpression) expressionNode() {}
func (oe *InfixExpression) TokenLiteral() string {
	return oe.Token.Literal
}

func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

// Converts Program object into string
func (p *Program) String() string {
	//Used for building strings
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
