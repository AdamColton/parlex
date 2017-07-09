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

func (r Reducer) Can(node parlex.ParseNode) bool {
	_, has := r[node.Kind().String()]
	return has
}

func Merge(r1, r2 Reducer) Reducer {
	merged := Reducer{}
	for symbol, r2fn := range r2 {
		if r1fn, found := r1[symbol]; !found {
			merged[symbol] = r2fn
		} else {
			merged[symbol] = merge(r1fn, r2fn)
		}
	}
	for symbol, r1fn := range r1 {
		if _, found := merged[symbol]; !found {
			merged[symbol] = r1fn
		}
	}
	return merged
}

func merge(fn1, fn2 Reduction) Reduction {
	return func(node *PN) {
		fn1(node)
		fn2(node)
	}
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

// PromoteChild removes the node with the child at cIdx and replaces it's own
// lexeme with the value. The grandchildren are spliced into the replaced childs
// position. The cIdx value uses GetIdx.
func (p *PN) PromoteChild(cIdx int) bool {
	cIdx, l, ok := p.GetIdx(cIdx)
	if !ok {
		return false
	}

	p.Lexeme = p.C[cIdx].Lexeme
	tail := p.C[cIdx].C
	if cIdx+1 < l {
		tail = append(tail, p.C[cIdx+1:]...)
	}
	p.C = append(p.C[0:cIdx], tail...)
	for _, c := range p.C {
		c.P = p
	}
	return true
}

// ReplaceWithChild replaces the node with the child at cIdx.
func (p *PN) ReplaceWithChild(cIdx int) bool {
	cIdx, _, ok := p.GetIdx(cIdx)
	if !ok {
		return false
	}
	*p = *(p.C[cIdx])
	return true
}

// GetIdx takes a child index and performs a conversion and check on it. If the
// index is negative, it will convert it to a positive position relative to the
// end (so -1 is the last child). The second int returned is the number of
// children. And the last value is bool indicating if cIdx is between 0 and len.
func (p *PN) GetIdx(cIdx int) (int, int, bool) {
	l := len(p.C)
	if cIdx < 0 {
		cIdx = l + cIdx
	}
	return cIdx, l, cIdx >= 0 && cIdx < l
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
func (p *PN) RemoveChild(cIdx int) bool {
	cIdx, l, ok := p.GetIdx(cIdx)
	if !ok {
		return false
	}
	if cIdx == l-1 {
		p.C = p.C[:cIdx]
	} else {
		p.C = append(p.C[0:cIdx], p.C[cIdx+1:]...)
	}
	return true
}

// PromoteSingleChild ; if the node has a single child, that child will be
// promoted to replace the node.
func (p *PN) PromoteSingleChild() bool {
	if len(p.C) == 1 {
		p.PromoteChild(0)
		return true
	}
	return false
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
func (p *PN) PromoteChildrenOf(cIdx int) bool {
	cIdx, l, ok := p.GetIdx(cIdx)
	if !ok {
		return false
	}
	if cIdx == l-1 {
		p.C = append(p.C[:cIdx], p.C[cIdx].C...)
	} else {
		p.C = append(p.C[:cIdx], append(p.C[cIdx].C, p.C[cIdx+1:]...)...)
	}
	return true
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
		p.Lexeme = lexeme.New(p.Kind()).Set(p.C[cIdx].Value())
	}
	p.RemoveChild(cIdx)
}

func (p *PN) ChildIs(cIdx int, kind string) bool {
	cIdx, _, ok := p.GetIdx(cIdx)
	if !ok {
		return false
	}
	return p.Child(cIdx).Kind().String() == kind
}

func PromoteChildrenOf(cIdx int) Reduction {
	return func(node *PN) {
		node.PromoteChildrenOf(cIdx)
	}
}

// PromoteChildValue returns a Reduction that will replace the value of the node
// with the value from the child at cIdx and remove the child at cIdx. If cIdx
// is negative, it will find the child relative to the end. If cIdx is out of
// bounds, no action will be taken.
func PromoteChildValue(cIdx int) Reduction {
	return func(node *PN) { node.PromoteChildValue(cIdx) }
}

func PromoteGrandChildren() func(*PN) {
	return func(node *PN) {
		node.PromoteGrandChildren()
	}
}

// RemoveChild produces a Reduction for removing a child at cIdx.
func RemoveChild(cIdx int) Reduction {
	return func(node *PN) { node.RemoveChild(cIdx) }
}

// RemoveChildren produces a Reduction for removing a child at cIdx.
func RemoveChildren(cIdxs ...int) Reduction {
	return func(node *PN) { node.RemoveChildren(cIdxs...) }
}

// PromoteSingleChild fulfills Reduce. If the node has a single child, that
// child will be promoted to replace the node.
func PromoteSingleChild(node *PN) {
	node.PromoteSingleChild()
}

// ReplaceWithChild replaces the node with the child at position cIdx, using
// GetIdx
func ReplaceWithChild(cIdx int) Reduction {
	return func(node *PN) { node.ReplaceWithChild(cIdx) }
}

func PromoteChild(cIdx int) func(*PN) {
	return func(node *PN) {
		node.PromoteChild(cIdx)
	}
}

func chain(r1, r2 Reduction) Reduction {
	return func(node *PN) {
		r1(node)
		r2(node)
	}
}

func (r Reduction) PromoteChildrenOf(cIdx int) Reduction {
	return chain(r, PromoteChildrenOf(cIdx))
}

// PromoteChildValue returns a Reduction that will replace the value of the node
// with the value from the child at cIdx and remove the child at cIdx. If cIdx
// is negative, it will find the child relative to the end. If cIdx is out of
// bounds, no action will be taken.
func (r Reduction) PromoteChildValue(cIdx int) Reduction {
	return chain(r, PromoteChildValue(cIdx))
}

func (r Reduction) PromoteGrandChildren() func(*PN) {
	return chain(r, PromoteGrandChildren())
}

// RemoveChild produces a Reduction for removing a child at cIdx.
func (r Reduction) RemoveChild(cIdx int) Reduction {
	return chain(r, RemoveChild(cIdx))
}

// RemoveChildren produces a Reduction for removing a child at cIdx.
func (r Reduction) RemoveChildren(cIdxs ...int) Reduction {
	return chain(r, RemoveChildren(cIdxs...))
}

// PromoteSingleChild fulfills Reduce. If the node has a single child, that
// child will be promoted to replace the node.
func (r Reduction) PromoteSingleChild() Reduction {
	return chain(r, PromoteSingleChild)
}

// ReplaceWithChild replaces the node with the child at position cIdx, using
// GetIdx
func (r Reduction) ReplaceWithChild(cIdx int) Reduction {
	return chain(r, ReplaceWithChild(cIdx))
}

func (r Reduction) PromoteChild(cIdx int) func(*PN) {
	return chain(r, PromoteChild(cIdx))
}
