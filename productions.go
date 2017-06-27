package parlex

import (
	"strings"
)

// Symbol is base of a grammar.
type Symbol string

// Len returns the number of characters in the string (not necessarily the
// number of bytes)
func (s Symbol) Len() int { return len([]rune(string(s))) }

// Production is a slice of symbols. It is not actually the full grammatic
// production because it does not contain the left side of the production.
type Production []Symbol

// String joins the symbols of the prodction with spaces
func (p Production) String() string {
	strs := make([]string, len(p))
	for i, symbol := range p {
		strs[i] = string(symbol)
	}
	return strings.Join(strs, " ")
}

// Productions are used to represent the set of productions available from a
// non-terminal
type Productions []Production
