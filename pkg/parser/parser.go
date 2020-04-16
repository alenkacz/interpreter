package parser

import (
	"fmt"
	"github.com/alenkacz/interpreter-book/pkg/ast"
	"github.com/alenkacz/interpreter-book/pkg/token"
	"github.com/alenkacz/interpreter-book/pkg/tokenizer"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOTEQ:    EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type Parser struct {
	tokenizer *tokenizer.Tokenizer

	currentToken *token.Token
	nextToken    *token.Token

	errors []string

	parseFns map[token.TokenType]prefixParseFn
}

func New(tokenizer *tokenizer.Tokenizer) *Parser {
	p := &Parser{
		tokenizer: tokenizer,
	}

	p.parseFns = make(map[token.TokenType]prefixParseFn)
	p.parseFns[token.INT] = p.parseIntegerLiteral
	p.parseFns[token.IDENT] = p.parseIdentifier
	p.parseFns[token.TRUE] = p.parseBoolean
	p.parseFns[token.FALSE] = p.parseBoolean
	p.parseFns[token.BANG] = p.parsePrefixExpression
	p.parseFns[token.MINUS] = p.parsePrefixExpression
	p.parseFns[token.PLUS] = p.parsePrefixExpression
	p.parseFns[token.LPAREN] = p.parseGroupedExpression
	p.parseFns[token.IF] = p.parseIfExpression
	p.parseFns[token.FUNC] = p.parseFuncExpression

	p.readNextToken()
	p.readNextToken()

	return p
}

func (p *Parser) readNextToken() {
	t := p.tokenizer.NextToken()
	p.currentToken = p.nextToken
	p.nextToken = &t
}

func (p *Parser) ParseProgram() *ast.Program {
	result := &ast.Program{
		Statements: make([]ast.Statement, 0),
	}
	for p.currentToken.Type != token.EOF {
		result.Statements = append(result.Statements, p.parseNextStatement())
		p.readNextToken()
	}
	return result
}

func (p *Parser) parseNextStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.EOF:
		break
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		for {
			if p.currentToken.Type == token.SEMICOLON {
				p.readNextToken()
				break
			}
			p.readNextToken()
		}
		return &ast.ReturnStatement{}
	default:
		stmt := p.parseExpressionStatement()
		return stmt
	}
	return nil
}

func (p *Parser) parseLetStatement() ast.Statement {
	if !p.readNextIfNextTypeIs(token.IDENT) {
		return nil
	}
	identifier := p.currentToken

	// assign sign
	if !p.readNextIfNextTypeIs(token.ASSIGN) {
		return nil
	}
	p.readNextToken()

	expression := p.parseExpression(LOWEST)

	if !p.readNextIfNextTypeIs(token.SEMICOLON) {
		return nil
	}

	return &ast.LetStatement{
		Identifier: identifier,
		Value: expression,
	}
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.nextToken.Type == token.SEMICOLON {
		p.readNextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	parseFn, ok := p.parseFns[p.currentToken.Type]
	if !ok {
		p.errors = append(p.errors, fmt.Sprintf("Unknown token type %s, no parseFn found", p.currentToken.Type))
		return nil
	}
	left := parseFn()

	for p.nextToken.Type != token.SEMICOLON && precedence < p.nextPrecedence() {
		p.readNextToken()
		precedence := p.currentPrecedence()

		expression := &ast.InfixExpression{
			Operator: p.currentToken.Literal,
			Left:     left,
		}

		p.readNextToken()
		expression.Right = p.parseExpression(precedence)

		left = expression
	}
	return left
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) nextPrecedence() int {
	if p, ok := precedences[p.nextToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, _ := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	return &ast.IntegerLiteral{
		Value: value,
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	identifier := &ast.Identifier{
		Name: p.currentToken.Literal,
	}
	if p.nextToken.Type == token.LPAREN {
		p.readNextToken()
		return p.parseCallExpression(identifier)
	}
	return identifier
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Value: p.currentToken.Type == token.TRUE}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	operator := p.currentToken
	p.readNextToken()
	return &ast.PrefixExpression{Operator: operator.Literal, Right: p.parseExpression(PREFIX)}
}

func (p *Parser) peekError(t token.TokenType) bool {
	if p.nextToken.Type != t {
		msg := fmt.Sprintf("expected next token to be %s, got %s instead",
			t, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return true
	}
	return false
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.readNextToken()

	exp := p.parseExpression(LOWEST)

	if p.nextToken.Type != token.RPAREN {
		return nil
	} else {
		p.readNextToken()
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	// we need to move the curToken pointer to start of expression
	if !p.readNextIfNextTypeIs(token.LPAREN) {
		return nil
	}
	p.readNextToken()

	condition := p.parseExpression(LOWEST)
	if !p.readNextIfNextTypeIs(token.RPAREN) {
		return nil
	}
	if !p.readNextIfNextTypeIs(token.LBRACE) {
		return nil
	}
	block := p.parseBlockStatement()
	expression := &ast.IfExpression{
		Condition: condition,
		Block: block,
	}
	if p.nextToken.Type == token.ELSE {
		p.readNextToken()

		if !p.readNextIfNextTypeIs(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	result := make([]ast.Statement, 0)
	p.readNextToken()
	for p.currentToken.Type != token.EOF && p.currentToken.Type != token.RBRACE {
		stmt := p.parseNextStatement()
		if stmt != nil {
			result = append(result, stmt)
		}
		p.readNextToken()
	}
	return &ast.BlockStatement{
		Statements: result,
	}
}

func (p *Parser) parseFuncExpression() ast.Expression {
	lit := &ast.FunctionLiteral{}

	if !p.readNextIfNextTypeIs(token.LPAREN) {
		return nil
	}

	lit.Params = p.parseFunctionParameters()

	if !p.readNextIfNextTypeIs(token.LBRACE) {
		return nil
	}

	lit.Block = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.nextToken.Type == token.RPAREN {
		p.readNextToken()
		return identifiers
	}

	p.readNextToken()

	ident := &ast.Identifier{Name: p.currentToken.Literal}
	identifiers = append(identifiers, ident)

	for p.nextToken.Type == token.COMMA {
		p.readNextToken()
		p.readNextToken()
		ident := &ast.Identifier{Name: p.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.readNextIfNextTypeIs(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	ident, ok := function.(*ast.Identifier)
	if !ok {
		p.errors = append(p.errors, fmt.Sprintf("expecting identifier in the beginning of call expression, got: %s", ident))
		return nil
	}
	exp := &ast.CallExpression{Function: ident}
	exp.Params = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.nextToken.Type == token.RPAREN {
		// empty args
		p.readNextToken()
		return args
	}

	p.readNextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.nextToken.Type == token.COMMA {
		p.readNextToken()
		p.readNextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.readNextIfNextTypeIs(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) readNextIfNextTypeIs(t token.TokenType) bool {
	if p.nextToken.Type != t {
		msg := fmt.Sprintf("expected next token to be %s, got %s instead",
			t, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return false
	}
	p.readNextToken()
	return true
}

func (p *Parser) readNextIfCurrentTypeIs(t token.TokenType) bool {
	if p.currentToken.Type != t {
		msg := fmt.Sprintf("expected next token to be %s, got %s instead",
			t, p.nextToken.Type)
		p.errors = append(p.errors, msg)
		return false
	}
	p.readNextToken()
	return true
}