package eparser

type Parser struct {
	lexer     *Lexer
	pos       int
	Variables map[string]Token
}
