package lexer

import (
	"github.com/parsaakbari1209/interpreter/token"
)

// Lexer generates tokens out of a source code.
type Lexer struct {
	// Source code.
	input string

	// Current position in input (points to the current character).
	position int

	// Current reading position in input (after current character).
	readPosition int

	//
	// Current character under examination.
	// In order to fully support Unicode and UTF-8, "rune" should be used instead of "byte".
	//
	char byte
}

// New initializes and returns a new *Lexer.
// The new lexer is initiated with the first character.
func New(input string) *Lexer {
	lxr := &Lexer{input: input}
	lxr.readChar()
	return lxr
}

// NextToken reads and returns the next token and updates the lexer.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.char {
	case '=':
		tok = newToken(token.ASSIGN, l.char)
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			tok.Literal = string(char) + string(l.char)
			tok.Type = token.EQ
		}
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '!':
		tok = newToken(token.BANG, l.char)
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			tok.Literal = string(char) + string(l.char)
			tok.Type = token.NOT_EQ
		}
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '*':
		tok = newToken(token.ASTERISK, l.char)
	case '/':
		tok = newToken(token.SLASH, l.char)
	case '<':
		tok = newToken(token.LT, l.char)
	case '>':
		tok = newToken(token.GT, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}
	l.readChar()
	return tok
}

// newToken constructs a new token given the token type and character.
func newToken(typ token.TokenType, char byte) token.Token {
	return token.Token{Type: typ, Literal: string(char)}
}

// readChar gives the next character and advances the positions in the input string.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		//
		// 0 is the ASCII code for "NUL" character and signifies one of the options below:
		// - Haven't read anything yet.
		// - End of file.
		//
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// peerChar is similar to readChar but it doesn't move the position.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
