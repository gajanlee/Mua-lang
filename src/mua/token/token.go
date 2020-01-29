package token

type TokenType string

type Token struct {
    Type    TokenType
    Literal string
}

const (
    ILLEGAL = "ILLEGAL"
    EOF     = "EOF"
    
    // Identifier and Literals
    ID  = "ID"
    INT = "INT"

    // Operators
    ASSIGN = "="
    PLUS   = "+"
    
    // Delimiters
    COMMA     = ","
    SEMICOLON = ";"

    L_PAREN = "("
    R_PAREN = ")"
    L_BRACE = "{"
    R_BRACE = "}"

    // keywords
    FUNCTION = "FUNCTION"
    LET      = "LET"
)

var keywords = map[string]TokenType {
    "fn": FUNCTION,
    "let": LET,
}

func LookupIdentifier(ident string) TokenType {
    if tok, ok := keywords[ident]; ok {
        return tok
    }
    return ID
}