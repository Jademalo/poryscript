package token

// Type distinguishes between different types of tokens in the Poryscript lexer.
type Type string

// Token represents a single token in the Poryscript lexer.
type Token struct {
	Type       Type
	Literal    string
	LineNumber int
}

// Token types
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT     = "IDENT"
	INT       = "INT"
	STRING    = "STRING"
	RAWSTRING = "RAWSTRING"

	// Operators
	EQ  = "=="
	NEQ = "!="
	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	// Delimeters
	COMMA = ","

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	SCRIPT    = "SCRIPT"
	RAW       = "RAW"
	RAWGLOBAL = "RAWGLOBAL"
	VAR       = "VAR"
	FLAG      = "FLAG"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	IF        = "IF"
	ELSE      = "ELSE"
	ELSEIF    = "ELSEIF"
)

// If statement comparison types
const (
	CMPVAR  = "CMPVAR"
	CMPFLAG = "CMPFLAG"
)

var keywords = map[string]Type{
	"script":     SCRIPT,
	"raw":        RAW,
	"raw_global": RAWGLOBAL,
	"var":        VAR,
	"flag":       FLAG,
	"TRUE":       TRUE,
	"FALSE":      FALSE,
	"true":       TRUE,
	"false":      FALSE,
	"if":         IF,
	"else":       ELSE,
	"elif":       ELSEIF,
}

// GetIdentType looks up the token type for the given identifier
func GetIdentType(ident string) Type {
	if tokType, ok := keywords[ident]; ok {
		return tokType
	}
	return IDENT
}
