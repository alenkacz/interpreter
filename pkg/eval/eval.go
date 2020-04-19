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
	case *ast.InfixExpression:
		infix, _ := node.(*ast.InfixExpression)
		left := Eval(infix.Left)
		right := Eval(infix.Right)
		return evalInfixOperator(left, right, infix.Operator)
	case *ast.ExpressionStatement:
		exp, _ := node.(*ast.ExpressionStatement)
		return Eval(exp.Expression)
	case *ast.Program:
		var result object.Object
		program, _ := node.(*ast.Program)
		for _, stmt := range program.Statements {
			result = Eval(stmt)
			switch result.(type) {
			case *object.Error:
				return result
			}
		}
		return result
	}
	return nil
}

func evalInfixOperator(left object.Object, right object.Object, operator string) object.Object {
	switch operator {
	case "+":
		leftInt, leftok := left.(*object.Integer)
		rightInt, rightok := right.(*object.Integer)
		if !leftok || !rightok {
			return newError("infix operator + works only with integers. Got %s+%s", left.Type(), right.Type())
		}
		return &object.Integer{leftInt.Value + rightInt.Value}
	case "-":
		leftInt, leftok := left.(*object.Integer)
		rightInt, rightok := right.(*object.Integer)
		if !leftok || !rightok {
			return newError("infix operator - works only with integers. Got %s-%s", left.Type(), right.Type())
		}
		return &object.Integer{leftInt.Value - rightInt.Value}
	case "*":
		leftInt, leftok := left.(*object.Integer)
		rightInt, rightok := right.(*object.Integer)
		if !leftok || !rightok {
			return newError("infix operator * works only with integers. Got %s*%s", left.Type(), right.Type())
		}
		return &object.Integer{leftInt.Value * rightInt.Value}
	case "/":
		leftInt, leftok := left.(*object.Integer)
		rightInt, rightok := right.(*object.Integer)
		if !leftok || !rightok {
			return newError("infix operator * works only with integers. Got %s*%s", left.Type(), right.Type())
		}
		return &object.Integer{leftInt.Value / rightInt.Value}
	default:
		return evalEqualityExpression(left, right, operator)
	}
}

func evalEqualityExpression(left object.Object, right object.Object, operator string) object.Object {
	if left.Type() != right.Type() {
		return object.FALSE
	}
	if left.Type() == object.INTEGER {
		leftVal := left.(*object.Integer).Value
		rightVal := right.(*object.Integer).Value
		switch operator {
		case "==":
			return boolResultToObject(leftVal == rightVal)
		case "!=":
			return boolResultToObject(leftVal != rightVal)
		case ">":
			return boolResultToObject(leftVal > rightVal)
		case "<":
			return boolResultToObject(leftVal < rightVal)
		default:
			return newError("unsupported operator %s%s%s", left.Type(), operator, right.Type())
		}
	} else {
		switch operator {
		case "==":
			return boolResultToObject(left == right)
		case "!=":
			return boolResultToObject(left != right)
		default:
			return newError("unsupported operator %s%s%s", left.Type(), operator, right.Type())
		}
	}
	return nil
}

func boolResultToObject(b bool) object.Object {
	if b {
		return object.TRUE
	}
	return object.FALSE
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