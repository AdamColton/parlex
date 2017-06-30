package setsymbol

import (
	"github.com/adamcolton/parlex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFulfillsInterface(t *testing.T) {
	s := New()
	A := s.Str("A")
	_ = parlex.Symbol(A)
	B := s.Str("B")
	C := s.Str("C")
	p1 := s.Production(A, B, C)
	_ = parlex.Production(p1)
	p2 := s.Production(C, A, B)
	ps := s.Productions(p1, p2)
	_ = parlex.Productions(ps)

	assert.Equal(t, A.Idx(), s.Str("A").Idx())
}
