package token

type TokenType string

const (
	ILLEGAL = "illegal"
	EOF = "eof"

	IDENT = "ident"
	INT = "int"
	STRING = "string"

	SEMICOLON = ";"
	ASSIGN = "="
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACKET = "["
	RBRACKET = "]"
	COMMA = ","
	PLUS = "+"
	MINUS = "-"
	SLASH = "/"
	ASTERISK = "*"
	BANG = "!"
	LT = "<"
	GT = ">"
	EQ = "=="
	NOTEQ = "!="

	// keywords
	LET = "let"
	FUNC = "fn"
	IF = "if"
	ELSE = "else"
	RETURN = "return"
	TRUE = "true"
	FALSE = "false"
)

var keywords = map[string]Token {
	"let": Token{LET, "let"},
	"fn": Token{FUNC, "fn"},
	"if": Token{IF, "if"},
	"else": Token{ELSE, "else"},
	"return": Token{RETURN, "return"},
	"true": Token{TRUE, "true"},
	"false": Token{FALSE, "false"},
}

type Token struct {
	Type TokenType
	Literal string
}

func GetKeyword(token string) *Token {
	val, ok := keywords[token]
	if ok {
		return &val
	}
	return nil
}
