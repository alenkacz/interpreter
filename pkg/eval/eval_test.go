package eval

import (
	"github.com/alenkacz/interpreter-book/pkg/object"
	"github.com/alenkacz/interpreter-book/pkg/parser"
	"github.com/alenkacz/interpreter-book/pkg/tokenizer"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected, tt.input)
	}
}

func testEval(input string) object.Object {
	l := tokenizer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program, object.NewEnvironment())
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64, input string) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("%s: object is not Integer. got=%T (%+v)", input, obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("%s: object has wrong value. got=%d, want=%d",
			input, result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!false", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},

	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected, tt.input)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool, input string) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("%s: object is not Boolean. got=%T (%+v)", input, obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		expectedMessage string
	}{
		{
			"5 + true;",
			"infix operator + works only with integers. Got INTEGER+BOOLEAN",
		},
		{
			"5 + true; 5;",
			"infix operator + works only with integers. Got INTEGER+BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"infix operator + works only with integers. Got BOOLEAN+BOOLEAN",
		},
		{
			"5; true + false; 5",
			"infix operator + works only with integers. Got BOOLEAN+BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"infix operator + works only with integers. Got BOOLEAN+BOOLEAN",
		},
		{
			`
132
if (10 > 1) {
if (10 > 1) {
return true + false;
}
return 1;
}
`,
			"infix operator + works only with integers. Got BOOLEAN+BOOLEAN",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer), tt.input)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != object.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { return 10; }", 10},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
`,
			10,
		},
		/*{
			`
let f = fn(x) {
  return x;
  x + 10;
};
f(10);`,
			10,
		},
		{
			`
let f = fn(x) {
   let result = x + 10;
   return result;
   return 10;
};
f(10);`,
			20,
		},*/
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Params) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Params)
	}

	if fn.Params[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Params[0])
	}

	expectedBody := "(x + 2)"

	if fn.Block.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Block.String())
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, tt.input)
	}
}