package tree

import (
	"github.com/adamcolton/parlex"
)

// Reduction is a function that reduces a node.
type Reduction func(node *PN)

// Reducer is used to reduce a ParseTree to something more useful, generally
// clearing away symbols that are now represeneted by the tree structure. It
// implements parlex.Reducer.
type Reducer map[parlex.Symbol]Reduction

// Add a reduction.
func (r Reducer) Add(symbol parlex.Symbol, reduction Reduction) {
	r[symbol] = reduction
}

// Reduce performs a reduction on the tree. It makes a copy during the process
// and the result comes back as parlex.ParseNode. For this reason, you don't
// want to traverse up the tree during a reduction. Instead, use Reduce to
// traverse the tree once, then handle upward traversals in a second path or
// with a stack. Though often it can be avoided by adding the reduction logic
// further up the tree.
func (r Reducer) Reduce(node parlex.ParseNode) parlex.ParseNode {
	if node == nil {
		return node
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
	cp := &PN{
		Lexeme: &parlex.L{
			V: node.Value(),
			K: node.Kind(),
		},
		C: make([]*PN, node.Children()),
	}
	for i := range cp.C {
		cp.C[i] = r.RawReduce(node.Child(i))
	}

	if reduction := r[cp.Kind()]; reduction != nil {
		reduction(cp)
	}

	return cp
}

// PromoteChild replaces the node with the child at cIdx. If cIdx is negative,
// it will find the child relative to the end. If cIdx is out of bounds, no
// action will be taken.
func (p *PN) PromoteChild(cIdx int) {
	l := len(p.C)
	if cIdx < 0 {
		cIdx = l + cIdx
	}
	if cIdx >= 0 && cIdx < l {
		p.Lexeme = p.C[cIdx].Lexeme
		tail := p.C[cIdx].C
		if cIdx+1 < l {
			tail = append(tail, p.C[cIdx+1:]...)
		}
		p.C = append(p.C[0:cIdx], tail...)
		for _, c := range p.C {
			c.P = p
		}
	}
}

// RemoveChild produces a Reduction for removing a child at cIdx.
func RemoveChild(cIdx int) Reduction {
	return func(node *PN) { node.RemoveChild(cIdx) }
}

// RemoveChildren calls RemoveChild repeatedly for each index given. Note that
// the relative positions may change so if you wanted to remove what were
// initially indexes 1, 3 and 5, you would need to either call
// RemoveChildren(5,3,1) or RemoveChildren(1,2,3)
func (p *PN) RemoveChildren(cIdxs ...int) {
	for _, cIdx := range cIdxs {
		p.RemoveChild(cIdx)
	}
}

// RemoveChild remove the child at cIdx. If cIdx is negative, it will find the
// child relative to the end. If cIdx is out of bounds, no action will be taken.
func (p *PN) RemoveChild(cIdx int) {
	l := len(p.C)
	if cIdx < 0 {
		cIdx = l + cIdx
	}
	if cIdx == l-1 {
		p.C = p.C[:cIdx]
	} else if cIdx < l && cIdx >= 0 {
		p.C = append(p.C[0:cIdx], p.C[cIdx+1:]...)
	}
}

// PromoteSingleChild fulfills Reduce. If the node has a single child, that
// child will be promoted to replace the node.
func PromoteSingleChild(node *PN) {
	node.PromoteSingleChild()
}

// PromoteSingleChild ; if the node has a single child, that child will be
// promoted to replace the node.
func (p *PN) PromoteSingleChild() {
	if len(p.C) == 1 {
		p.PromoteChild(0)
	}
}

// PromoteGrandChildren will remove all the immediate children and replace them
// with the grand children.
func (p *PN) PromoteGrandChildren() {
	ct := 0
	for _, child := range p.C {
		ct += len(child.C)
	}
	newChildren := make([]*PN, 0, ct)
	for _, child := range p.C {
		for _, grandChild := range child.C {
			grandChild.P = p
		}
		newChildren = append(newChildren, child.C...)
	}
	p.C = newChildren
}

// PromoteChildrenOf will remove the child at cIdx and splice in all it's
// children. If cIdx is negative, it will find the child relative to the end. If
// cIdx is out of bounds, no action will be taken.
func (p *PN) PromoteChildrenOf(cIdx int) {
	l := len(p.C)
	if cIdx < 0 {
		cIdx = l + cIdx
	}
	if cIdx == l-1 {
		p.C = append(p.C[:cIdx], p.C[cIdx].C...)
	} else if cIdx >= 0 && l > cIdx {
		p.C = append(p.C[:cIdx], append(p.C[cIdx].C, p.C[:cIdx+1]...)...)
	}
}

// PromoteChildValue returns a Reduction that will replace the value of the node
// with the value from the child at cIdx and remove the child at cIdx. If cIdx
// is negative, it will find the child relative to the end. If cIdx is out of
// bounds, no action will be taken.
func PromoteChildValue(cIdx int) Reduction {
	return func(node *PN) { node.PromoteChildValue(cIdx) }
}

// PromoteChildValue returns a Reduction that will replace the value of the node
// with the value from the child at cIdx and remove the child at cIdx. If cIdx
// is negative, it will find the child relative to the end. If cIdx is out of
// bounds, no action will be taken.
func (p *PN) PromoteChildValue(cIdx int) {
	l := len(p.C)
	if cIdx < 0 {
		cIdx = l + cIdx
	}
	if cIdx >= 0 && l > cIdx {
		p.Lexeme = &parlex.L{
			K: p.Kind(),
			V: p.C[cIdx].Value(),
		}
	}
	p.RemoveChild(cIdx)
}
