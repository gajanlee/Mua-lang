package ast

import (
	"bytes"
	"mua/token"
	"strings"
)

type Node interface {
	TokenLiteral() string	// for debugging
	String()	   string	// debug and compare with other AST nodes
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

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token		// token.ID
	Value string
}

func (i *Identifier) expressionNode()	   {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token token.Token		// token.INT
	Value int64
}

func (il *IntegerLiteral) expressionNode()		{}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type Boolean struct {
	Token token.Token		// token.TRUE, token.FALSE
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token		// token.L_BRACKET
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, ele := range al.Elements {
		elements = append(elements, ele.String())
	}

	out.WriteString("[" + strings.Join(elements, ", ") + "]")
	return out.String()
}

type HashLiteral struct {
	Token token.Token		// token.L_BRACE
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String() + ":" + value.String())
	}

	out.WriteString("{" + strings.Join(pairs, ", ") + "}")
	return out.String()
}

// array[0]
type IndexExpression struct {
	Token token.Token		// token.L_BRACKET
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + ie.Left.String() + "[" + ie.Index.String() + "])")
	return out.String()
}

type PrefixExpression struct {
	Token    token.Token		// The prefix token, like '!', '-'
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()       {}
func (pe *PrefixExpression) TokenLiteral() string  { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + pe.Operator + pe.Right.String() + ")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token		// The infix operator, like '+', '-'...
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()		 {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")")
	return out.String()
}

type ExpressionStatement struct {
	Token	   token.Token		// the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()		 {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}
func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token		// token.LET
	Name *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()		  {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " " + ls.Name.String() + " = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token 		token.Token		// token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

// if <cond-expr> <consequence> else <alternative>
type IfExpression struct {
	Token       token.Token		// token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if" + ie.Condition.String() + " " + ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else " + ie.Consequence.String())
	}
	return out.String()
}

// [let x = ]  fn(x, y) {x + y}
type FunctionLiteral struct {
	Token      token.Token		// token.FUNCTIOn
	Parameters []*Identifier
	Body	   *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(" + strings.Join(params, ", ") + ") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// fn(x, y) { x + y; }(2 + 3)
// call(2, 3, fn(x, y) {x+y})
type CallExpression struct {
	Token     token.Token		// token.L_PAREN
	Function  Expression		// Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	out.WriteString(ce.Function.String() + "(" + strings.Join(args, ", ") + ")")
	
	return out.String()
}