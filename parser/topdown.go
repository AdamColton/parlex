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
		op := &tdOp{
			TD:  t,
			lxs: lexemes,
		}
		pn, _ = op.accept(nts[0], 0, true)
	}
	return pn
}

// top-down parse operation
type tdOp struct {
	*TD
	lxs []parlex.Lexeme
}

// Tries to accept the lexemes into the grammper from a starting symbol and
// position. If end == -1, then it will return the first accepting rule. If
// end > -1, it the rule must end on that position. This is used at the outer
// most level to assure accept consumes all the lexemes
func (op *tdOp) accept(symbol parlex.Symbol, pos int, all bool) (*tree.PN, int) {
	if pos >= len(op.lxs) {
		return nil, pos
	}
	productions := op.Productions(symbol)
	if productions == nil {
		if pos < len(op.lxs) && symbol == op.lxs[pos].Kind() {
			return &tree.PN{
				Lexeme: op.lxs[pos],
			}, pos + 1
		}
	}

	for _, prod := range productions {
		children := make([]*tree.PN, len(prod))
		accepts := true
		tc := pos
		for i, symbol := range prod {
			if pn, c := op.accept(symbol, tc, false); pn != nil {
				children[i], tc = pn, c
			} else {
				accepts = false
				break
			}
		}
		if accepts && (!all || tc == len(op.lxs)) {
			return &tree.PN{
				Lexeme: &parlex.L{K: symbol},
				C:      children,
			}, tc
		}
	}

	return nil, pos
}
