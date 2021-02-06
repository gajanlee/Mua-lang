package lexer

import "muc/token"

type Lexer struct {
	input		 string
	position	 int		// current position in input (point to current char)
	readPosition int		// current reading position (after current char)
	char		 byte		// current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.char {
	case '=':
		// Check if it is `==`
		if l.peekChar() == '=' {
			ch := l.char
			l.readChar()
			tok = token.Token{Type: token.EQUAL, Literal: string(ch) + string(l.char)}
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '(':
		tok = newToken(token.L_PAREN, l.char)
	case ')':
		tok = newToken(token.R_PAREN, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case ':':
		tok = newToken(token.COLON, l.char)
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '!':
		// Check if it is `!=`
		if l.peekChar() == '=' {
			ch := l.char
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.char)}
		} else {
			tok = newToken(token.BANG, l.char)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.char)
	case '/':
		tok = newToken(token.SLASH, l.char)
	case '<':
		tok = newToken(token.LESS, l.char)
	case '>':
		tok = newToken(token.GREATER, l.char)
	case '[':
		tok = newToken(token.L_BRACKET, l.char)
	case ']':
		tok = newToken(token.R_BRACKET, l.char)
	case '{':
		tok = newToken(token.L_BRACE, l.char)
	case '}':
		tok = newToken(token.R_BRACE, l.char)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()

	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
} 

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (l *Lexer) readString() string {
	position := l.position + 1		// skip the "
	for {
		l.readChar()
		if l.char == '"' {
			break
		}
	}
	return l.input[position:l.position]
}