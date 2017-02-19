package parlex

// L is a concrete implementation of Lexeme
type L struct {
	K    Symbol
	V    string
	L, C int
}

// Kind returns the token indicating what kind of lexeme this is
func (l *L) Kind() Symbol { return l.K }

// Value returns the string segment
func (l *L) Value() string { return l.V }

// Pos returns the position as (line, column) of where the lexeme started in
// the original string.
func (l *L) Pos() (int, int) { return l.L, l.C }
