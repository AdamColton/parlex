package topdown

import (
	"errors"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/lexeme"
	"github.com/adamcolton/parlex/symbol/setsymbol"
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
	set := setsymbol.New()
	set.LoadGrammar(t.Grammar)
	op := &tdOp{
		Topdown: t,
		lxs:     set.LoadLexemes(lexemes),
		memo:    make(map[treeKey]*acceptResp),
		set:     set,
	}
	start := op.set.Symbol(nts[0]).Idx()
	return op.accept(treeKey{start, 0}, true).node()
}

type treeKey struct {
	idx int
	pos int
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
	lxs  []*lexeme.Lexeme
	memo map[treeKey]*acceptResp
	set  *setsymbol.Set
}

func (op *tdOp) accept(key treeKey, all bool) *acceptResp {
	if resp, ok := op.memo[key]; ok {
		return resp
	}
	resp := op.tryAccept(key, all)
	op.memo[key] = resp
	return resp
}

// Tries to accept the lexemes into the grammper from a starting symbol and
// position. If end == -1, then it will return the first accepting rule. If
// end > -1, it the rule must end on that position. This is used at the outer
// most level to assure accept consumes all the lexemes
func (op *tdOp) tryAccept(key treeKey, all bool) *acceptResp {
	symbol := op.set.ByIdx(key.idx)
	productions := op.Productions(symbol)

	if productions == nil {
		if key.pos < len(op.lxs) && key.idx == op.lxs[key.pos].K.(*setsymbol.Symbol).Idx() {
			return resp(op.lxs[key.pos], key.pos+1)
		}
		return nil
	}

	for i := productions.Iter(); i.Next(); {
		accepts := op.acceptProd(key, i.Production)
		if accepts != nil && (!all || accepts.end == len(op.lxs)) {
			return accepts
		}
	}

	return nil
}

func (op *tdOp) acceptProd(key treeKey, prod parlex.Production) *acceptResp {
	children := make([]*tree.PN, prod.Symbols())
	pos := key.pos

	for i := prod.Iter(); i.Next(); {
		symbol := op.set.Symbol(i.Symbol)
		resp := op.accept(treeKey{symbol.Idx(), pos}, false)
		if resp == nil {
			return nil
		}
		children[i.Idx], pos = resp.PN, resp.end
	}

	symbol := op.set.ByIdx(key.idx)
	return resp(lexeme.New(symbol), pos, children...)
}
