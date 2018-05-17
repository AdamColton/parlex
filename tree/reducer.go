package tree

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/lexeme"
)

// Reduction is a function that reduces a node.
type Reduction func(node *PN)

// Reducer is used to reduce a ParseTree to something more useful, generally
// clearing away symbols that are now represeneted by the tree structure. It
// implements parlex.Reducer.
type Reducer map[string]Reduction

// Add a reduction.
func (r Reducer) Add(symbol string, reduction Reduction) {
	r[symbol] = reduction
}

// Can returns true if a reducer has a rule for the given node
func (r Reducer) Can(node parlex.ParseNode) bool {
	_, has := r[node.Kind().String()]
	return has
}

// Merge takes two Reducers and returns a single Reducer that is the merged
// result. If a Kind is present in both r1 and r2, the merged will behave as
// though running r1 then r2 on the node.
func Merge(r1, r2 Reducer) Reducer {
	merged := Reducer{}
	for symbol, r2fn := range r2 {
		if r1fn, found := r1[symbol]; !found {
			merged[symbol] = r2fn
		} else {
			merged[symbol] = Chain(r1fn, r2fn)
		}
	}
	for symbol, r1fn := range r1 {
		if _, found := merged[symbol]; !found {
			merged[symbol] = r1fn
		}
	}
	return merged
}

// Reduce performs a reduction on the tree. It makes a copy during the process
// and the result comes back as parlex.ParseNode. For this reason, you don't
// want to traverse up the tree during a reduction. Instead, use Reduce to
// traverse the tree once, then handle upward traversals in a second path or
// with a stack. Though often it can be avoided by adding the reduction logic
// further up the tree.
func (r Reducer) Reduce(node parlex.ParseNode) parlex.ParseNode {
	if node == nil {
		// RawReduce returning nil is not the same thing
		return nil
	}
	return r.RawReduce(node)
}

// RawReduce performs a reduction on the tree. It makes a copy during the
// process and the result comes back as the concrete type *PN. For this reason,
// you don't want to traverse up the tree during a reduction. Instead, use
// Reduce to traverse the tree once, then handle upward traversals in a second
// path or with a stack. Though often it can be avoided by adding the reduction
// logic further up the tree.
func (r Reducer) RawReduce(node parlex.ParseNode) *PN {
	if node == nil {
		return nil
	}
	cp := &PN{
		Lexeme: lexeme.Copy(node),
		C:      make([]*PN, node.Children()),
	}
	for i := range cp.C {
		cp.C[i] = r.RawReduce(node.Child(i))
	}

	if reduction := r[cp.Kind().String()]; reduction != nil {
		reduction(cp)
	}

	return cp
}
