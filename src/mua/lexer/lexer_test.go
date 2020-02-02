package lexer

import (
	"testing"

	"mua/token"
)

func TestNextToken(t *testing.T) {
	// input := `=+(){},;`
	input := `let five = 5
let ten = 10

let add = fn(x, y) {
	x + y
}

let result = add(five, ten)

if 5 < 10 {
	return true
} else {
	return false
}

5 == 10
5 != 20
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
`
	tests := []struct {
		expectedType	token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.ID, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.LET, "let"},
		{token.ID, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.LET, "let"},
		{token.ID, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.L_PAREN, "("},
		{token.ID, "x"},
		{token.COMMA, ","},
		{token.ID, "y"},
		{token.R_PAREN, ")"},
		{token.L_BRACE, "{"},
		{token.ID, "x"},
		{token.PLUS, "+"},
		{token.ID, "y"},
		{token.R_BRACE, "}"},
		{token.LET, "let"},
		{token.ID, "result"},
		{token.ASSIGN, "="},
		{token.ID, "add"},
		{token.L_PAREN, "("},
		{token.ID, "five"},
		{token.COMMA, ","},
		{token.ID, "ten"},
		{token.R_PAREN, ")"},
		{token.IF, "if"},
		{token.INT, "5"},
		{token.LESS, "<"},
		{token.INT, "10"},
		{token.L_BRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.R_BRACE, "}"},
		{token.ELSE, "else"},
		{token.L_BRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.R_BRACE, "}"},
		{token.INT, "5"},
		{token.EQUAL, "=="},
		{token.INT, "10"},
		{token.INT, "5"},
		{token.NOT_EQ, "!="},
		{token.INT, "20"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.L_BRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.R_BRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.L_BRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.R_BRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}