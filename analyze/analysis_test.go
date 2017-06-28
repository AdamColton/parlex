package analyze

import (
	"github.com/adamcolton/parlex/grammar"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	assert.Equal(t, 3, len(a.nonterm2firsts))
	if !assert.Equal(t, 2, len(a.nonterm2firsts["A"])) {
		t.Error(a.nonterm2firsts["A"])
	}
	if !assert.Equal(t, 1, len(a.nonterm2firsts["B"])) {
		t.Error(a.nonterm2firsts["B"])
	}
	if !assert.Equal(t, 1, len(a.nonterm2firsts["C"])) {
		t.Error(a.nonterm2firsts["C"])
	}

	assert.True(t, a.HasFirst("A", "x"))
	assert.True(t, a.HasFirst("A", "y"))
	assert.False(t, a.HasFirst("A", "z"))
	assert.False(t, a.HasFirst("B", "x"))
	assert.True(t, a.HasFirst("B", "y"))
	assert.False(t, a.HasFirst("B", "z"))
	assert.False(t, a.HasFirst("C", "x"))
	assert.False(t, a.HasFirst("C", "y"))
	assert.True(t, a.HasFirst("C", "z"))

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
	assert.Equal(t, 3, len(a.nonterm2firsts))
	if !assert.Equal(t, 3, len(a.nonterm2firsts["A"])) {
		t.Error(a.nonterm2firsts["A"])
	}
	if !assert.Equal(t, 1, len(a.nonterm2firsts["B"])) {
		t.Error(a.nonterm2firsts["B"])
	}
	if !assert.Equal(t, 1, len(a.nonterm2firsts["C"])) {
		t.Error(a.nonterm2firsts["C"])
	}

	assert.True(t, a.HasFirst("A", "x"))
	assert.True(t, a.HasFirst("A", "y"))
	assert.True(t, a.HasFirst("A", "z"))
	assert.False(t, a.HasFirst("B", "x"))
	assert.True(t, a.HasFirst("B", "y"))
	assert.False(t, a.HasFirst("B", "z"))
	assert.False(t, a.HasFirst("C", "x"))
	assert.False(t, a.HasFirst("C", "y"))
	assert.True(t, a.HasFirst("C", "z"))

	assert.True(t, a.HasNilInFirst("A"))
	assert.True(t, a.HasNilInFirst("B"))
	assert.False(t, a.HasNilInFirst("C"))
}
