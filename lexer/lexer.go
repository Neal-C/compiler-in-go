package lexer

import "github.com/Neal-C/compiler-in-go/token"

type Lexer struct {
	input        string
	position     int  // current position in the input (points to current character) // and where we last read
	readPosition int  // current reading position in the input (after current character)
	ch           byte // current char under examination
}

const BLANK_WHITESPACE = ' '

func New(input string) *Lexer {
	lexer := &Lexer{
		input: input,
	}

	lexer.readChar()

	return lexer
}

func (lexer *Lexer) readChar() {
	//  is to check whether we have reached the end of input.
	if lexer.readPosition >= len(lexer.input) {
		// 0 is the ASCII code for 'NUL' character
		lexer.ch = 0
		// signifies either
		// “we haven’t read anything yet” or “end of file” for us.
	} else {
		// But if we haven’t reached the end of input yet,
		// it sets l.ch to the next character by accessing l.input[l.readPosition].
		lexer.ch = lexer.input[lexer.readPosition]
	}

	lexer.position = lexer.readPosition
	lexer.readPosition++
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.skipWhitespace()

	switch lexer.ch {
	case '=':
		if lexer.peekChar() == '=' {
			ch := lexer.ch
			lexer.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(lexer.ch)}
		} else {
			tok = newToken(token.ASSIGN, lexer.ch)
		}
	case '+':
		tok = newToken(token.PLUS, lexer.ch)
	case '-':
		tok = newToken(token.MINUS, lexer.ch)
	case '!':
		if lexer.peekChar() == '=' {
			ch := lexer.ch
			lexer.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(lexer.ch)}
		} else {

			tok = newToken(token.BANG, lexer.ch)
		}
	case '/':
		tok = newToken(token.SLASH, lexer.ch)
	case '*':
		tok = newToken(token.ASTERISK, lexer.ch)
	case '<':
		tok = newToken(token.LT, lexer.ch)
	case '>':
		tok = newToken(token.GT, lexer.ch)
	case ';':
		tok = newToken(token.SEMICOLON, lexer.ch)
	case ':':
		tok = newToken(token.COLON, lexer.ch)
	case ',':
		tok = newToken(token.COMMA, lexer.ch)
	case '(':
		tok = newToken(token.LPAREN, lexer.ch)
	case ')':
		tok = newToken(token.RPAREN, lexer.ch)
	case '{':
		tok = newToken(token.LBRACE, lexer.ch)
	case '}':
		tok = newToken(token.RBRACE, lexer.ch)
	case '[':
		tok = newToken(token.LBRACKET, lexer.ch)
	case ']':
		tok = newToken(token.RBRACKET, lexer.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = lexer.readString()
	default:
		if isLetter(lexer.ch) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookUpIdent(tok.Literal)
			return tok
		} else if isDigit(lexer.ch) {
			tok.Type = token.INT
			tok.Literal = lexer.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lexer.ch)
		}
	}

	lexer.readChar()
	return tok
}

func (lexer *Lexer) readIdentifier() string {
	initialPosition := lexer.position
	for isLetter(lexer.ch) {
		lexer.readChar()
	}

	return lexer.input[initialPosition:lexer.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.ch == BLANK_WHITESPACE || lexer.ch == '\t' || lexer.ch == '\n' || lexer.ch == '\r' {
		lexer.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (lexer *Lexer) readNumber() string {
	initialPosition := lexer.position
	for isDigit(lexer.ch) {
		lexer.readChar()
	}
	return lexer.input[initialPosition:lexer.position]
}

func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
}

func (self *Lexer) readString() string {
	position := self.position + 1

	for {
		self.readChar()
		if self.ch == '"' || self.ch == 0 {
			break
		}
	}

	return self.input[position:self.position]
}
