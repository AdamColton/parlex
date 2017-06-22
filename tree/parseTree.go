package tree

import (
	"errors"
	"github.com/adamcolton/parlex"
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
			ch := &PN{
				Lexeme: &parlex.L{
					K: parlex.Symbol(m[2]),
					V: m[3],
				},
				P: cur,
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
func (p *PN) Child(i int) parlex.ParseNode {
	if i >= len(p.C) {
		return nil
	}
	return p.C[i]
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
	s = append(s, pad, string(p.Lexeme.Kind())) // B.2/3
	if v := p.Lexeme.Value(); v != "" {
		s = append(s, ": '", p.Lexeme.Value(), "'") // C.3
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
		r += 3 // C
	}
	if len(p.C) > 0 {
		r += 2 // D
		for _, child := range p.C {
			r += child.sliceReq()
		}
	}
	return r
}

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
