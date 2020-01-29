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
    ASSIGN   = "="
    PLUS     = "+"
    MINUS    = "-"
    BANG     = "!"
    ASTERISK = "*"
    SLASH    = "/"

    LESS    = "<"
    GREATER = ">"

    EQUAL  = "==" 
    NOT_EQ = "!="

    // Delimiters
    COMMA     = ","
    SEMICOLON = ";"

    L_PAREN = "("
    R_PAREN = ")"
    L_BRACE = "{"
    R_BRACE = "}"

    // keywords
    LET      = "LET"
    RETURN   = "RETURN"
    FUNCTION = "FUNCTION"

    IF    = "IF"
    ELSE  = "ELSE"
    TRUE  = "TRUE"
    FALSE = "FALSE"
)

var keywords = map[string]TokenType {
    "let": LET,
    "return": RETURN,
    "fn": FUNCTION,
    "if": IF,
    "else": ELSE,
    "true": TRUE,
    "false": FALSE,
}

func LookupIdentifier(ident string) TokenType {
    if tok, ok := keywords[ident]; ok {
        return tok
    }
    return ID
}