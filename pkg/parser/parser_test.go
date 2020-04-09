package parser

import (
	"github.com/alenkacz/interpreter-book/pkg/ast"
	"github.com/alenkacz/interpreter-book/pkg/tokenizer"
	"strings"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := tokenizer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(p.errors) > 0 {
		t.Fatalf("Error(s) in ParseProgram(): %v", strings.Join(p.errors, ","))
	}
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := tokenizer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(p.errors) > 0 {
		t.Fatalf("Error(s) in ParseProgram(): %v", strings.Join(p.errors, ","))
	}
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}
	for _, stmt := range program.Statements {
		_, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := tokenizer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(p.errors) > 0 {
		t.Fatalf("Error(s) in ParseProgram(): %v", strings.Join(p.errors, ","))
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}
	if letStmt.Name() != name {
		t.Errorf("letStmt.Name() not '%s'. got=%s", name, letStmt.Name())
		return false
	}
	return true
}