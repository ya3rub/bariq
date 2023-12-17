package lexer

import (
	"bariq/token"
)

type Lexer struct {
	input        string
	position     int  // points to current char
	readPosition int  // after cuurent char
	ch           byte // char being examined
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// gives the next char and increment pos to
// the next pos
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 0 is the ASCII code for nul
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peakChar() byte {
	if l.readPosition >= len(l.input) {
		return 0 // 0 is the ASCII code for nul
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '"':
		l.readChar()
		tok.Literal = l.readString()
		tok.Type = token.STRING
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '=':
		if l.peakChar() == '=' {
			tok.Literal = "=="
			tok.Type = token.EQ
			l.readChar()
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '!':
		if l.peakChar() == '=' {
			tok.Literal = "!="
			tok.Type = token.NEQ
			l.readChar()
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '*':
		tok = newToken(token.ASTERIK, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// you already did an l.readChar()
			return tok
		}
		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}
	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	// fmt.Println("start string", string(l.ch))
	position := l.position
	for {
		// fmt.Println(l.input[position])
		// WARN: it was buggy here, you wrote '0' instead of 0
		if l.ch == '"' || l.ch == 0 {
			break
		}
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	postition := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[postition:l.position]
}

func (l *Lexer) readIdentifier() string {
	postition := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[postition:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tt token.TokenType, ch byte) token.Token {
	return token.Token{Type: tt, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
