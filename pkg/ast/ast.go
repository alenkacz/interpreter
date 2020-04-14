package ast

import (
	"bytes"
	"fmt"
	"github.com/alenkacz/interpreter-book/pkg/token"
	"strconv"
	"strings"
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
	return e.Expression.String()
}

type BlockStatement struct {
	Statements []Statement
}

func (*BlockStatement) statementNode() {}
func (b *BlockStatement) String() string {
	var buf bytes.Buffer
	for _, s := range b.Statements {
		buf.WriteString(s.String())
		buf.WriteString("\n")
	}
	return buf.String()
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

type Boolean struct {
	Value bool
}
func (*Boolean) expressionNode() {}
func (i *Boolean) String() string {
	return strconv.FormatBool(i.Value)
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

type PrefixExpression struct {
	Right Expression
	Operator string
}
func (*PrefixExpression) expressionNode() {}
func (i *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", i.Operator, i.Right.String())
}

type IfExpression struct {
	Condition Expression
	Block *BlockStatement
	Alternative *BlockStatement
}
func (*IfExpression) expressionNode() {}
func (i *IfExpression) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("if(%s) {\n", i.Condition.String()))
	buf.WriteString(i.Block.String())
	if i.Alternative != nil {
		buf.WriteString("} else {\n")
		buf.WriteString(i.Alternative.String())
		buf.WriteString("}\n")
	}
	return buf.String()
}

type FunctionLiteral struct {
	Params []*Identifier
	Block *BlockStatement
}
func (*FunctionLiteral) expressionNode() {}
func (f *FunctionLiteral) String() string {
	var buf bytes.Buffer
	var paramNames []string
	for _, i := range f.Params {
		paramNames = append(paramNames, i.String())
	}
	buf.WriteString(fmt.Sprintf("fn(%s){ %s }", strings.Join(paramNames, ","), f.Block.String()))
	return buf.String()
}