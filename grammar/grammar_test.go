package grammar

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/symbol/stringsymbol"
	"github.com/stretchr/testify/assert"
	"testing"
)

var A = stringsymbol.Symbol("A")
var x = stringsymbol.Symbol("x")

func TestGrammarString(t *testing.T) {
	g1, err := New(`
    A -> B C
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
	assert.True(t, g1.Productions(A) != nil)
	assert.True(t, g1.Productions(x) == nil)

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

	nilProd := g.Productions(stringsymbol.Symbol("NIL"))
	if assert.Equal(t, nilProd.Productions(), 1) {
		assert.Equal(t, nilProd.Production(0).Symbols(), 0)
	}
}
