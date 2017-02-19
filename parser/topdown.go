package parser

import (
	"errors"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/tree"
)

// TD is a Top Down parser
type TD struct {
	parlex.Grammar
}

// ErrLeftRecursion is thrown if the grammar is left recursive. Top down parsing
// cannot handle left recursion.
var ErrLeftRecursion = errors.New("Top Down parser cannot handle left recursion")

// TopDown returns a top down parser
func TopDown(grmr parlex.Grammar) (*TD, error) {
	if grammar.IsLeftRecursive(grmr) {
		return nil, ErrLeftRecursion
	}
	return &TD{
		Grammar: grmr,
	}, nil
}

// Parse implements parlex.Parser
func (t *TD) Parse(lexemes []parlex.Lexeme) parlex.ParseNode {
	var pn *tree.PN
	if nts := t.NonTerminals(); len(nts) > 0 {
		pn, _ = t.accept(nts[0], 0, len(lexemes), lexemes)
	}
	return pn
}

// Tries to accept the lexemes into the grammper from a starting symbol and
// position. If end == -1, then it will return the first accepting rule. If
// end > -1, it the rule must end on that position. This is used at the outer
// most level to assure accept consumes all the lexemes
func (t *TD) accept(symbol parlex.Symbol, cur, end int, lxs []parlex.Lexeme) (*tree.PN, int) {
	if cur >= len(lxs) {
		return nil, cur
	}
	productions := t.Productions(symbol)
	if productions == nil {
		if cur < len(lxs) && symbol == lxs[cur].Kind() {
			return &tree.PN{
				Lexeme: lxs[cur],
			}, cur + 1
		}
	}

	for _, prod := range productions {
		children := make([]*tree.PN, len(prod))
		accepts := true
		tc := cur
		for i, symbol := range prod {
			if pn, c := t.accept(symbol, tc, -1, lxs); pn != nil {
				children[i], tc = pn, c
			} else {
				accepts = false
				break
			}
		}
		if accepts && (end == -1 || end == tc) {
			return &tree.PN{
				Lexeme: &parlex.L{K: symbol},
				C:      children,
			}, tc
		}
	}

	return nil, cur
}
