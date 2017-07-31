package grammar

import (
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

	assert.True(t, g1.Productions(A) != nil)
	assert.True(t, g1.Productions(x) == nil)

	assert.Equal(t, "x", g1.Productions(A).Production(1).Symbol(0).String())

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
	assert.NotEqual(t, -1, g.set.Idx(stringsymbol.Symbol("NIL")))
	nilProd := g.Productions(stringsymbol.Symbol("NIL"))
	if assert.Equal(t, nilProd.Productions(), 1) {
		assert.Equal(t, nilProd.Production(0).Symbols(), 0)
	}
}

func TestBasic(t *testing.T) {
	grmr, err := New(`
    E -> E op E
      -> ( E )
      -> int
  `)
	assert.Len(t, grmr.NonTerminals(), 1)
	assert.NoError(t, err)
}
