package grammar

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveLeftRecursion(t *testing.T) {
	grmr, err := New(`
    E -> E op E
      -> ( E )
      -> int
  `)
	assert.NoError(t, err)
	noRecur := RemoveLeftRecursion(grmr)
	expected, err := New(`
    E  -> ( E ) E'
       -> int E'
    E' -> op E E'
       -> 
  `)
	assert.NoError(t, err)
	assert.Equal(t, expected.String(), noRecur.(*Grammar).String())

	grmr, err = New(`
    A -> B C
    B -> x
      ->
    C -> A
      -> y
  `)
	assert.NoError(t, err)
	noRecur = RemoveLeftRecursion(grmr)
	expected, err = New(`
    A  -> B C
    B  -> x
       -> 
    C  -> x C C'
       -> y C'
    C' -> 
    `)
	assert.NoError(t, err)
	assert.Equal(t, expected.String(), noRecur.(*Grammar).String())
}
