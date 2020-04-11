package ast

import (
	"bytes"
	"fmt"
	"github.com/alenkacz/interpreter-book/pkg/token"
)

type Node interface {
	String() string
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, st := range p.Statements {
out.WriteString(st.String())
}
return out.String()
}

type Statement interface {
	Node
	statementNode() // just to be able to distinguish betweet statement and expression
}

type LetStatement struct {
	identifier *token.Token
	value *Expression
}

func NewLetStatement(identifier *token.Token) *LetStatement {
	return &LetStatement{
		identifier: identifier,
	}
}

func (*LetStatement) statementNode() {}
func (l *LetStatement) Name() string {
	return l.identifier.Literal
}
func (l *LetStatement) String() string {
	return fmt.Sprintf("let %s = ;", l.identifier)
}

type ReturnStatement struct {
	ReturnValue Expression
}

func (*ReturnStatement) statementNode() {}
func (l *ReturnStatement) String() string {
	return "return;"
}

type ExpressionStatement struct {
	Expression Expression
}

func (*ExpressionStatement) statementNode() {}
func (e *ExpressionStatement) String() string {
	return e.String()
}

type Expression interface {
	Node
	expressionNode() // just to be able to distinguish betweet statement and expression
}

type IntegerLiteral struct {
	Value int64
}
func (*IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type Identifier struct {
	Name string
}
func (*Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return i.Name
}

type InfixExpression struct {
	Left Expression
	Right Expression
	Operator string
}
func (*InfixExpression) expressionNode() {}
func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Operator, i.Right.String())
}