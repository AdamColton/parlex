package grammar

import (
	"github.com/adamcolton/parlex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGrammarString(t *testing.T) {
	g1, err := New(`
    AA -> B C
         x
    B -> y
    C -> z 
  `)
	assert.NoError(t, err)
	assert.NotNil(t, g1)

	g2, err := New(g1.String())
	assert.NoError(t, err)
	assert.Equal(t, g1.String(), g2.String())

	assert.Equal(t, g1.String(), parlex.GrammarString(g1))
}

func TestNil(t *testing.T) {
	g, err := New(`
    A   -> B C
           x
    B   -> Y
    C   -> z
        -> NIL
    Y   -> A
    NIL ->
  `)
	assert.NoError(t, err)
	assert.NotNil(t, g)

	nilProd := g.Productions("NIL")
	if assert.Len(t, nilProd, 1) {
		assert.Len(t, nilProd[0], 0)
	}
}
