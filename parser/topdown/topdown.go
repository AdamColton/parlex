package topdown

import (
	"errors"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexeme"
	"github.com/adamcolton/parlex/tree"
)

// TD is a Top Down parser
type TD struct {
	parlex.Grammar
}

// ErrLeftRecursion is thrown if the grammar is left recursive. Top down parsing
// cannot handle left recursion.
var ErrLeftRecursion = errors.New("Top Down parser cannot handle left recursion")

// New returns a topdown parser
func New(grmr parlex.Grammar) (*TD, error) {
	if grammar.IsLeftRecursive(grmr) {
		return nil, ErrLeftRecursion
	}
	return &TD{
		Grammar: grmr,
	}, nil
}

// Parse implements parlex.Parser
func (t *TD) Parse(lexemes []parlex.Lexeme) parlex.ParseNode {
	nts := t.NonTerminals()
	if len(nts) == 0 {
		return nil
	}
	op := &tdOp{
		TD:   t,
		lxs:  lexemes,
		memo: make(map[ParseOp]acceptResp),
	}
	return op.accept(ParseOp{nts[0], 0}, true).PN
}

// ParseOp represents a parser operation of accepting a symbol at a position.
type ParseOp struct {
	parlex.Symbol
	Pos int
}

type acceptResp struct {
	*tree.PN
	end int
}

// top-down parse operation
type tdOp struct {
	*TD
	lxs  []parlex.Lexeme
	memo map[ParseOp]acceptResp
}

func (op *tdOp) accept(pop ParseOp, all bool) acceptResp {
	if resp, ok := op.memo[pop]; ok {
		return resp
	}
	resp := op.tryAccept(pop, all)
	op.memo[pop] = resp
	return resp
}

// Tries to accept the lexemes into the grammper from a starting symbol and
// position. If end == -1, then it will return the first accepting rule. If
// end > -1, it the rule must end on that position. This is used at the outer
// most level to assure accept consumes all the lexemes
func (op *tdOp) tryAccept(pop ParseOp, all bool) acceptResp {
	if pop.Pos >= len(op.lxs) {
		return acceptResp{
			PN:  nil,
			end: pop.Pos,
		}
	}

	productions := op.Productions(pop.Symbol)
	if productions == nil {
		if pop.Pos < len(op.lxs) && pop.Symbol == op.lxs[pop.Pos].Kind() {
			return acceptResp{
				PN: &tree.PN{
					Lexeme: op.lxs[pop.Pos],
				},
				end: pop.Pos + 1,
			}
		}
		if pop.Symbol == "NIL" {
			return acceptResp{
				PN: &tree.PN{
					Lexeme: lexeme.New("NIL"),
				},
				end: pop.Pos,
			}
		}
	}

	for _, prod := range productions {
		children := make([]*tree.PN, len(prod))
		accepts := true
		pos := pop.Pos
		for i, symbol := range prod {
			if resp := op.accept(ParseOp{symbol, pos}, false); resp.PN != nil {
				children[i], pos = resp.PN, resp.end
			} else {
				accepts = false
				break
			}
		}
		if accepts && (!all || pos == len(op.lxs)) {
			return acceptResp{
				PN: &tree.PN{
					Lexeme: lexeme.New(pop.Symbol),
					C:      children,
				},
				end: pos,
			}
		}
	}

	return acceptResp{
		PN:  nil,
		end: pop.Pos,
	}
}
