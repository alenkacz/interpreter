package object

import (
	"bytes"
	"fmt"
	"github.com/alenkacz/interpreter-book/pkg/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER = "INTEGER"
	STRING = "STRING"
	BOOLEAN = "BOOLEAN"
	NULL_TYPE = "NULL"
	ERROR = "ERROR"
	RETURN_TYPE = "RETURN"
	FUNCTION = "FUNCTION"
	)

var (
	TRUE = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NULL = &Null{}
)

type Object interface {
	Type() ObjectType
	Print() string
}

type Integer struct {
	Value int64
}

func (*Integer) Type() ObjectType { return INTEGER }
func (i *Integer) Print() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (*Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Print() string  { return fmt.Sprintf("%t", b.Value) }

type String struct {
	Value string
}

func (*String) Type() ObjectType { return STRING }
func (i *String) Print() string  { return i.Value }

type Error struct {
	Message string
}

func (*Error) Type() ObjectType { return ERROR }
func (b *Error) Print() string  { return fmt.Sprintf("%s", b.Message) }

type Null struct {
}

func (*Null) Type() ObjectType { return NULL_TYPE }
func (*Null) Print() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (*ReturnValue) Type() ObjectType { return RETURN_TYPE }
func (r *ReturnValue) Print() string  { return r.Value.Print() }

type Function struct {
	Environment *Environment
	Params []*ast.Identifier
	Block *ast.BlockStatement
}

func (*Function) Type() ObjectType { return FUNCTION }
func (f *Function) Print() string  {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Block.String())
	out.WriteString("\n}")

	return out.String() }