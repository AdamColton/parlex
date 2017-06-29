package parlex

// Symbol is base of a grammar.
type Symbol interface {
	String() string
}

// Production is a slice of symbols. It is not actually the full grammatic
// production because it does not contain the left side of the production.
type Production interface {
	Symbols() int
	Symbol(int) Symbol
}

// Productions are used to represent the set of productions available from a
// non-terminal
type Productions interface {
	Productions() int
	Production(int) Production
}

func SymLen(s Symbol) int {
	return len([]rune(s.String()))
}
