package topdown

import (
	"errors"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/lexeme"
	"github.com/adamcolton/parlex/symbol/stringsymbol"
	"github.com/adamcolton/parlex/tree"
)

// Topdown is a Top Down parser
type Topdown struct {
	parlex.Grammar
}

// ErrLeftRecursion is thrown if the grammar is left recursive. Top down parsing
// cannot handle left recursion.
var ErrLeftRecursion = errors.New("Top Down parser cannot handle left recursion")

// New returns a topdown parser
func New(grmr parlex.Grammar) (*Topdown, error) {
	if parlex.IsLeftRecursive(grmr) {
		return nil, ErrLeftRecursion
	}
	return &Topdown{
		Grammar: grmr,
	}, nil
}

// Parse implements parlex.Parser
func (t *Topdown) Parse(lexemes []parlex.Lexeme) parlex.ParseNode {
	nts := t.NonTerminals()
	if len(nts) == 0 {
		return nil
	}
	op := &tdOp{
		Topdown: t,
		lxs:     lexemes,
		memo:    make(map[treeKey]*acceptResp),
	}
	return op.accept(treeKey{stringsymbol.Symbol(nts[0].String()), 0, true}).node()
}

type treeKey struct {
	stringsymbol.Symbol
	pos int
	all bool
}

type acceptResp struct {
	*tree.PN
	end int
}

func resp(lx parlex.Lexeme, end int, children ...*tree.PN) *acceptResp {
	return &acceptResp{
		PN: &tree.PN{
			Lexeme: lx,
			C:      children,
		},
		end: end,
	}
}

func (a *acceptResp) node() *tree.PN {
	if a == nil {
		return nil
	}
	return a.PN
}

// top-down parse operation
type tdOp struct {
	*Topdown
	lxs  []parlex.Lexeme
	memo map[treeKey]*acceptResp
}

func (op *tdOp) accept(key treeKey) *acceptResp {
	if resp, ok := op.memo[key]; ok {
		return resp
	}
	resp := op.tryAccept(key)
	op.memo[key] = resp
	return resp
}

// Tries to accept the lexemes into the grammper from a starting symbol and
// position. If end == -1, then it will return the first accepting rule. If
// end > -1, it the rule must end on that position. This is used at the outer
// most level to assure accept consumes all the lexemes
func (op *tdOp) tryAccept(key treeKey) *acceptResp {
	productions := op.Productions(key.Symbol)

	if productions == nil {
		if key.pos < len(op.lxs) && key.Symbol.String() == op.lxs[key.pos].Kind().String() {
			return resp(op.lxs[key.pos], key.pos+1)
		}
		return nil
	}

	for i := productions.Iter(); i.Next(); {
		accepts := op.acceptProd(key, i.Production)
		if accepts != nil && (!key.all || accepts.end == len(op.lxs)) {
			return accepts
		}
	}

	return nil
}

func (op *tdOp) acceptProd(key treeKey, prod parlex.Production) *acceptResp {
	children := make([]*tree.PN, prod.Symbols())
	pos := key.pos

	for i := prod.Iter(); i.Next(); {
		symbol := stringsymbol.Symbol(i.Symbol.String())
		resp := op.accept(treeKey{symbol, pos, false})
		if resp == nil {
			return nil
		}
		children[i.Idx], pos = resp.PN, resp.end
	}

	return resp(lexeme.New(key.Symbol), pos, children...)
}
