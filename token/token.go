package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// IDENTEFIERS +  LITERALS
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN  = "="
	BANG    = "!"
	PLUS    = "+"
	MINUS   = "-"
	ASTERIK = "*"
	SLASH   = "/"
	GT      = ">"
	LT      = "<"
	EQ      = "=="
	NEQ     = "!="
	// Delimters
	COMMA     = ","
	COLON     = ":"
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"

	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// KEYWORD
	FUNCTION  = "FUNCTION"
	AWAIT     = "AWAIT"
	YIELD     = "YIELD"
	IF        = "if"
	ELSE      = "else"
	RET       = "return"
	TRUE      = "true"
	FALSE     = "false"
	LET       = "LET"
	ASYNC     = "ASYNC"
	GENERATOR = "GENERATOR"
)

type (
	TokenType string
	Token     struct {
		Type    TokenType
		Literal string
	}
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RET,
	"false":  FALSE,
	"true":   TRUE,
	"async":  ASYNC,
	"await":  AWAIT,
	"gen":    GENERATOR,
	"yield":  YIELD,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
