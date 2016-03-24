package eparser

type tokenType int

const (
	ILLEGAL tokenType = iota
	EOL

	literalsBegin
	IDENT // x
	INT   // 3
	FLOAT // 3.14
	literalsEnd

	operatorsBegin
	ADD // +
	SUB // -
	DIV // /
	MUL // *
	POW // **
	REM // %

	AND // &
	OR  // |
	XOR // ^
	LSH // <<
	RSH // >>
	NOT // ~

	EQ     // =
	ADD_EQ // +=
	SUB_EQ // -=
	DIV_EQ // /=
	MUL_EQ // *=
	POW_EQ // **=
	REM_EQ // %=

	AND_EQ // &=
	OR_EQ  // |=
	XOR_EQ // ^=
	LSH_EQ // <<=
	RSH_EQ // >>=
	operatorsEnd

	LPAREN // (
	RPAREN // )

	TOKEN_COUNT
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOL:     "EOL",

	IDENT: "IDENT",
	INT:   "INT",
	FLOAT: "FLOAT",

	ADD: "+",
	SUB: "-",
	DIV: "/",
	MUL: "*",
	POW: "**",
	REM: "%",

	AND: "&",
	OR:  "|",
	XOR: "^",
	LSH: "<<",
	RSH: ">>",
	NOT: "~",

	EQ:     "=",
	ADD_EQ: "+=",
	SUB_EQ: "-=",
	DIV_EQ: "/=",
	MUL_EQ: "*=",
	POW_EQ: "**=",
	REM_EQ: "%=",

	LPAREN: "(",
	RPAREN: ")",
}

type Token struct {
	Type  tokenType
	Value string
	Pos   int
}

func newToken(toktype tokenType, val string, pos int) *Token {
	return &Token{
		Type:  toktype,
		Value: val,
		Pos:   pos,
	}
}

func (tok *Token) String() string {
	if tok.Type < TOKEN_COUNT {
		return tokens[tok.Type]
	}

	return ""
}

func (tok *Token) IsOperator() bool {
	return tok.Type > operatorsBegin && tok.Type < operatorsEnd
}

func (tok *Token) IsLiteral() bool {
	return tok.Type > literalsBegin && tok.Type < literalsEnd
}
