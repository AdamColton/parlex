package tree

import (
	"errors"
	"fmt"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/lexeme"
	"github.com/adamcolton/parlex/symbol/stringsymbol"
	"regexp"
	"strings"
)

var reTreeLine = regexp.MustCompile(`\s*(\}?)\s*([^\}\{:\'\s]+|$)\s*(?:\:\s*\'((?:[^\'\\]|(?:\\.))*)\')?\s*(\{?)$`)

// ErrBadTreeString is thrown if a tree definition cannot be parsed.
var ErrBadTreeString = errors.New("Bad Tree String")

// New creates a *PN from a string. They are more commonly created from a Parser
// but this can be useful for testing.
func New(str string) (*PN, error) {
	cur := &PN{}
	for _, line := range strings.Split(str, "\n") {
		m := reTreeLine.FindStringSubmatch(line)
		if len(line) != len(m[0]) {
			return nil, ErrBadTreeString
		}
		if m[1] == "}" {
			cur = cur.P
		}
		if m[2] != "" {
			kind := stringsymbol.Symbol(m[2])
			val := m[3]
			ch := &PN{
				Lexeme: lexeme.New(kind).Set(val),
				P:      cur,
			}
			cur.C = append(cur.C, ch)
			if m[4] == "{" {
				cur = ch
			}
		}
	}
	cur = cur.C[0]
	cur.P = nil
	return cur, nil
}

// PN parse node forms a tree
type PN struct {
	parlex.Lexeme
	P *PN
	C []*PN
}

// Parent returns a reference to the nodes parent. If parent is nil, this is the
// root of the tree. This is part of the parlex.ParseNode interface.
func (p *PN) Parent() parlex.ParseNode {
	return p.P
}

// Children returns the number of children the node has. This is part of the
// parlex.ParseNode interface.
func (p *PN) Children() int {
	return len(p.C)
}

// Child returns a child by index. This is part of the parlex.ParseNode
// interface.
func (p *PN) Child(cIdx int) parlex.ParseNode {
	cIdx, _, ok := p.GetIdx(cIdx)
	if !ok {
		return nil
	}
	return p.C[cIdx]
}

// String converts the entire tree (starting a *PN) to a string. This string can
// be used to create a copy of the tree.
func (p *PN) String() string {
	// find the right slice length first to efficiently allot a []string
	segs := make([]string, 0, p.sliceReq())
	return strings.Join(p.string("", segs), "")
}

func (p *PN) string(pad string, s []string) []string {
	if p == nil {
		return append(s, pad, "NIL\n") // A.2
	}
	s = append(s, pad, p.Lexeme.Kind().String()) // B.2/3
	if v := p.Lexeme.Value(); v != "" {
		s = append(s, fmt.Sprintf(": %q", p.Lexeme.Value())) // C.3
	}
	if len(p.C) > 0 {
		s = append(s, " {\n") // B.1/3
		for _, child := range p.C {
			s = child.string(pad+"\t", s)
		}
		s = append(s, pad, "}\n") // D.2
	} else {
		s = append(s, "\n") // B.1/3
	}
	return s
}

// sliceReq computes the slice capacity requirement for String.
// the comments help line up the size requirements
// X.2 means this line requres 2 and lines up with comment X
// X.1/2 means this line 1 but it was combined with other appends at X.
func (p *PN) sliceReq() int {
	if p == nil {
		return 2 // A
	}
	r := 3 // B
	if p.Lexeme.Value() != "" {
		r += 1 // C
	}
	if len(p.C) > 0 {
		r += 2 // D
		for _, child := range p.C {
			r += child.sliceReq()
		}
	}
	return r
}

// Size counts the number of nodes in a tree
func (p *PN) Size() int {
	if p == nil {
		return 0
	}
	size := 1
	for _, child := range p.C {
		size += child.Size()
	}
	return size
}

// ChildAt takes an index and a list of symbols and returns true if there is a
// child at that index and it matches one of the symbols.
func (p *PN) ChildAt(cIdx int, symbs ...string) bool {
	cIdx, _, ok := p.GetIdx(cIdx)
	if !ok {
		return false
	}
	cs := p.C[cIdx].Kind().String()
	for _, s := range symbs {
		if s == cs {
			return true
		}
	}
	return false
}
