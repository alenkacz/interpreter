package object

import "fmt"

type ObjectType string

const (
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	NULL_TYPE = "NULL"
	ERROR = "ERROR"
	RETURN_TYPE = "RETURN"
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