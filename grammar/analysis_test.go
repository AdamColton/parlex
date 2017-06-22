package grammar

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnalyse(t *testing.T) {
	g, err := New(`
    AA -> B C
         x
    B -> y
    C -> z 
  `)
	assert.NoError(t, err)
	assert.NotNil(t, g)

	a := Analyse(g)
	assert.Equal(t, 3, len(a.nonterm2firsts))
	if !assert.Equal(t, 2, len(a.nonterm2firsts["AA"])) {
		t.Error(a.nonterm2firsts["AA"])
	}
	if !assert.Equal(t, 1, len(a.nonterm2firsts["B"])) {
		t.Error(a.nonterm2firsts["B"])
	}
	if !assert.Equal(t, 1, len(a.nonterm2firsts["C"])) {
		t.Error(a.nonterm2firsts["C"])
	}

	assert.True(t, a.HasFirst("AA", "x"))
	assert.True(t, a.HasFirst("AA", "y"))
	assert.False(t, a.HasFirst("AA", "z"))
	assert.False(t, a.HasFirst("B", "x"))
	assert.True(t, a.HasFirst("B", "y"))
	assert.False(t, a.HasFirst("B", "z"))
	assert.False(t, a.HasFirst("C", "x"))
	assert.False(t, a.HasFirst("C", "y"))
	assert.True(t, a.HasFirst("C", "z"))
}
