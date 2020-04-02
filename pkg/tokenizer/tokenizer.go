package tokenizer

import (
	"github.com/alenkacz/interpreter-book/pkg/token"
)

type Tokenizer struct {
	input string

	nextPos     int
	currentChar byte
}

func New(input string) *Tokenizer {
	t := &Tokenizer{input: input}
	t.readChar()
	return t
}

func (t *Tokenizer) readChar() {
	if t.nextPos >= len(t.input) {
		t.currentChar = 0 // ascii code for NUL
	} else {
		t.currentChar = t.input[t.nextPos]
	}
	t.nextPos += 1
}

func (t *Tokenizer) peekChar() byte {
	if t.nextPos >= len(t.input) {
		return 0 // ascii code for NUL
	} else {
		return t.input[t.nextPos]
	}
}

func (t *Tokenizer) skipWhitespace() {
	for t.currentChar == ' ' || t.currentChar == '\t' || t.currentChar == '\n' || t.currentChar == '\r' {
		t.readChar()
	}
}

func (t *Tokenizer) NextToken() token.Token {
	result := token.Token{}
	t.skipWhitespace()

	switch t.currentChar {
	case '=':
		if t.peekChar() == '=' {
			t.readChar()
			result = token.Token{token.EQ, "=="}
		} else {
			result = token.Token{token.ASSIGN, "="}
		}
	case '!':
		if t.peekChar() == '=' {
			t.readChar()
			result = token.Token{token.NOTEQ, "!="}
		} else {
			result = token.Token{token.BANG, "!"}
		}
	case ';':
		result = token.Token{token.SEMICOLON, ";"}
	case '(':
		result = token.Token{token.LPAREN, "("}
	case ')':
		result = token.Token{token.RPAREN, ")"}
	case '{':
		result = token.Token{token.LBRACE, "{"}
	case '}':
		result = token.Token{token.RBRACE, "}"}
	case ',':
		result = token.Token{token.COMMA, ","}
	case '+':
		result = token.Token{token.PLUS, "+"}
	case '-':
		result = token.Token{token.MINUS, "-"}
	case '/':
		result = token.Token{token.SLASH, "/"}
	case '*':
		result = token.Token{token.ASTERISK, "*"}
	case '<':
		result = token.Token{token.LT, "<"}
	case '>':
		result = token.Token{token.GT, ">"}
	case 0:
		result = token.Token{Type: token.EOF}
	default:
		if isIdentifier(t.currentChar) {
			// possibility that the identifier consist of more than one char
			literal := t.readWhole(isIdentifier)
			keyword := token.GetKeyword(literal)
			if keyword != nil {
				result = *keyword
			} else {
				result = token.Token{token.IDENT, literal}
			}
		} else if isNumber(t.currentChar) {
			number := t.readWhole(isNumber)
			result = token.Token{token.INT, number}
		} else {
			result = token.Token{Type: token.ILLEGAL}
		}
	}

	t.readChar()

	return result
}

func isNumber(char byte) bool {
	return char >= '0' && char <= '9'
}

func (t *Tokenizer) readWhole(isType func(byte)bool) string {
	var result string

	for isType(t.currentChar) {
		result += string(t.currentChar)
		if isType(t.peekChar()) {
			t.readChar()
		} else {
			break
		}
	}
	return result
}

func isIdentifier(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_'
}


