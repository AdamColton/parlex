package parlex

// Lexeme is a unit of output from a lexer
type Lexeme interface {
	Kind() Symbol
	Value() string
	Pos() (line int, col int)
}

// Lexer is fulfilled by a type that can convert a string into a slice of
// Lexemes.
type Lexer interface {
	Lex(string) []Lexeme
}

// ParseNode is a node in the parse tree.
type ParseNode interface {
	Lexeme
	Parent() ParseNode
	Children() int
	Child(int) ParseNode
}

// Parser takes a slice of Lexemes and returns a ParseNode. If the parse fails,
// ParseNode will be nil.
type Parser interface {
	Parse([]Lexeme) ParseNode
}

type ParserConstructor func(Grammar) (Parser, error)

// Grammar represents a context free Grammar.
type Grammar interface {
	Productions(symbol Symbol) Productions
	NonTerminals() []Symbol // The first NonTerminal should be the start symbol
}

// Reducer is used to reduce a ParseTree to something more useful, generally
// clearing away symbols that are now represeneted by the tree structure.
type Reducer interface {
	Reduce(ParseNode) ParseNode
}
