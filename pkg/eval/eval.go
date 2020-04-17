package eval

import (
	"fmt"
	"github.com/alenkacz/interpreter-book/pkg/ast"
	"github.com/alenkacz/interpreter-book/pkg/object"
)

func Eval(node ast.Node) object.Object {
	switch node.(type) {
	case *ast.IntegerLiteral:
		integer, _ := node.(*ast.IntegerLiteral)
		return &object.Integer{ Value: integer.Value }
	case *ast.Boolean:
		boolean, _ := node.(*ast.Boolean)
		if boolean.Value == true {
			return object.TRUE
		} else {
			return object.FALSE
		}
	case *ast.PrefixExpression:
		prefix, _ := node.(*ast.PrefixExpression)
		value := Eval(prefix.Right)
		switch prefix.Operator {
		case "!":
			return evalBang(value)
		case "-":
			return evalPrefixMinus(value)
		default:
			newError("unknown operator %s. %v", prefix.Operator, prefix)
		}
	case *ast.ExpressionStatement:
		exp, _ := node.(*ast.ExpressionStatement)
		return Eval(exp.Expression)
	case *ast.Program:
		var result object.Object
		program, _ := node.(*ast.Program)
		for _, stmt := range program.Statements {
			result = Eval(stmt)
		}
		return result
	}
	return nil
}

func evalBang(value object.Object) object.Object {
	switch value {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalPrefixMinus(value object.Object) object.Object {
	if value.Type() != object.INTEGER {
		return newError("unknown operator: -%s", value.Type())
	}

	integer := value.(*object.Integer).Value
	return &object.Integer{Value: -integer}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}