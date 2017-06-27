package topdown

import (
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexer/simplelexer"
	"github.com/adamcolton/parlex/tree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGpParse(t *testing.T) {
	lxr, err := simplelexer.New(`
    ( /\(/
    ) /\)/
    op /[+\-\*\/]/
    int /\d+/
    space /\s+/ -
  `)
	assert.NoError(t, err)
	grmr, err := grammar.New(`
    E -> T op E
      -> T
    T -> ( E )
      -> int
  `)
	assert.NoError(t, err)

	s := "1+2+3"
	lxs := lxr.Lex(s)
	p, err := New(grmr)
	assert.NoError(t, err)
	pn := p.Parse(lxs)
	if assert.NotNil(t, pn) {
		if tpn, ok := pn.(*tree.PN); ok {
			expected, _ := tree.New(`
        E {
          T {
            int: '1'
          }
          op: '+'
          E {
            T {
              int: '2'
            }
            op: '+'
            E {
              T {
                int: '3'
              }
            }
          }
        }
      `)
			assert.Equal(t, expected.String(), tpn.String())
		} else {
			t.Error("Parse node should be of type *tree.PN")
		}
	}
}

func TestParens(t *testing.T) {
	lxr, err := simplelexer.New(`
    ( /\(/
    ) /\)/
    op /[+\-\*\/]/
    int /\d+/
    space /\s+/ -
  `)
	assert.NoError(t, err)
	grmr, err := grammar.New(`
    E -> T op E
      -> T
    T -> P
      -> int
    P -> ( E )
  `)
	assert.NoError(t, err)

	s := "(1+2)*3"
	lxs := lxr.Lex(s)
	p, err := New(grmr)
	assert.NoError(t, err)
	pn := p.Parse(lxs)
	assert.NotNil(t, pn)
	//TODO: better assert
}

func TestNil(t *testing.T) {
	lxr, err := simplelexer.New(`
    ( /\(/
    ) /\)/
    op /[+\-\*\/]/
    int /\d+/
    space /\s+/
  `)
	assert.NoError(t, err)
	grmr, err := grammar.New(`
    E   -> T Gap op Gap E
        -> T
    T   -> P
        -> int
    P   -> ( Gap E Gap )
    Gap -> space Gap
        -> 
  `)
	assert.NoError(t, err)

	s := "( 1 + 2 )  *  3"
	lxs := lxr.Lex(s)
	p, err := New(grmr)
	assert.NoError(t, err)
	pn := p.Parse(lxs)
	assert.NotNil(t, pn)
}
