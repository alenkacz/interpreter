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
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	tokenizer *tokenizer.Tokenizer

	currentToken *token.Token
	nextToken *token.Token

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
		switch p.currentToken.Type {
		case token.EOF:
			break
		case token.LET:
			// identifier
			if p.peekError(token.IDENT) {
				break
			}
			p.readNextToken()
			identifier := p.currentToken

			// assign sign
			if p.peekError(token.ASSIGN) {
				break
			}
			p.readNextToken()

			// expression
			// TODO real expression parsing
			for {
				if p.currentToken.Type == token.SEMICOLON {
					p.readNextToken()
					break
				}
				p.readNextToken()
			}
			result.Statements = append(result.Statements, ast.NewLetStatement(identifier))
		case token.RETURN:
			for {
				if p.currentToken.Type == token.SEMICOLON {
					p.readNextToken()
					break
				}
				p.readNextToken()
			}
			result.Statements = append(result.Statements, &ast.ReturnStatement{})
		default:
			result.Statements = append(result.Statements, p.parseExpressionStatement())
			p.readNextToken()
		}
	}
	return result
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{}

	stmt.Expression = p.parseExpression()

	if p.nextToken.Type == token.SEMICOLON {
		p.readNextToken()
	}

	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	parseFn, ok := p.parseFns[p.currentToken.Type]
	if !ok {
		p.errors = append(p.errors, fmt.Sprintf("Unknown token type %s, no parseFn found", p.currentToken.Type))
		return nil
	}
	left := parseFn()

	for p.nextToken.Type != token.SEMICOLON {
		p.readNextToken()

		expression := &ast.InfixExpression{
			Operator: p.currentToken.Literal,
			Left:     left,
		}

		p.readNextToken()
		expression.Right = p.parseExpression()

		return expression
	}
	p.readNextToken()
	return left
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, _ := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	return &ast.IntegerLiteral{
		Value: value,
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Name: p.currentToken.Literal,
	}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Value: p.currentToken.Type == token.TRUE}
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