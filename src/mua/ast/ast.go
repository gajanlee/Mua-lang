package ast

import "mua/token"

type Node interface {
	TokenLiteral() string	// for debugging
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// the root node Program holds all statements
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token		// token.LET
	Name *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()		  {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token		// token.ID
	Value string
}

func (i *Identifier) expressionNode()	   {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

