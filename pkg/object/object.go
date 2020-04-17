package object

import "fmt"

type ObjectType string

const (
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	ERROR = "ERROR"
)

var (
	TRUE = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
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
