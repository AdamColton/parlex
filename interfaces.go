package parlex

// Symbol is base of a grammar. A symbol should always return the same string.
// Two symbols that return the same thing are considered to be the same.
type Symbol interface {
	String() string
}

// Production is a slice of symbols. It is not actually the full grammatic
// production because it does not contain the left side of the production.
type Production interface {
	Symbols() int
	Symbol(int) Symbol
	Iter() *ProductionIterator
}

// Productions are used to represent the set of productions available from a
// non-terminal
type Productions interface {
	Productions() int
	Production(int) Production
	Iter() *ProductionsIterator
}

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

// ParserConstructor is a function that takes a Grammar and returns a Parser
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
	Can(ParseNode) bool
}
