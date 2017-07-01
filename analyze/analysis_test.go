package analyze

import (
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/symbol/stringsymbol"
	"github.com/stretchr/testify/assert"
	"testing"
)

var A = stringsymbol.Symbol("A")
var B = stringsymbol.Symbol("B")
var C = stringsymbol.Symbol("C")
var x = stringsymbol.Symbol("x")
var y = stringsymbol.Symbol("y")
var z = stringsymbol.Symbol("z")

func TestAnalyze(t *testing.T) {
	g, err := grammar.New(`
    A -> B C
         x
    B -> y
    C -> z 
  `)
	assert.NoError(t, err)
	assert.NotNil(t, g)

	a := Analyze(g)

	AIdx := a.set.Symbol(A).Idx()
	BIdx := a.set.Symbol(B).Idx()
	CIdx := a.set.Symbol(C).Idx()

	if !assert.Equal(t, 2, len(a.nonterm2firsts[AIdx])) {
		t.Error(a.nonterm2firsts[AIdx])
	}
	if !assert.Equal(t, 1, len(a.nonterm2firsts[BIdx])) {
		t.Error(a.nonterm2firsts[BIdx])
	}
	if !assert.Equal(t, 1, len(a.nonterm2firsts[CIdx])) {
		t.Error(a.nonterm2firsts[CIdx])
	}

	assert.True(t, a.HasFirst(A, x))
	assert.True(t, a.HasFirst(A, y))
	assert.False(t, a.HasFirst(A, z))
	assert.False(t, a.HasFirst(B, x))
	assert.True(t, a.HasFirst(B, y))
	assert.False(t, a.HasFirst(B, z))
	assert.False(t, a.HasFirst(C, x))
	assert.False(t, a.HasFirst(C, y))
	assert.True(t, a.HasFirst(C, z))

	g, err = grammar.New(`
	    A -> B C
	         x
	    B -> y
	      ->
	    C -> z
	  `)
	assert.NoError(t, err)
	assert.NotNil(t, g)

	a = Analyze(g)
	if !assert.Equal(t, 3, len(a.nonterm2firsts[AIdx])) {
		t.Error(a.nonterm2firsts[AIdx])
	}
	if !assert.Equal(t, 1, len(a.nonterm2firsts[BIdx])) {
		t.Error(a.nonterm2firsts[BIdx])
	}
	if !assert.Equal(t, 1, len(a.nonterm2firsts[CIdx])) {
		t.Error(a.nonterm2firsts[CIdx])
	}

	assert.True(t, a.HasFirst(A, x))
	assert.True(t, a.HasFirst(A, y))
	assert.True(t, a.HasFirst(A, z))
	assert.False(t, a.HasFirst(B, x))
	assert.True(t, a.HasFirst(B, y))
	assert.False(t, a.HasFirst(B, z))
	assert.False(t, a.HasFirst(C, x))
	assert.False(t, a.HasFirst(C, y))
	assert.True(t, a.HasFirst(C, z))

	assert.True(t, a.HasNilInFirst(A))
	assert.True(t, a.HasNilInFirst(B))
	assert.False(t, a.HasNilInFirst(C))
}
