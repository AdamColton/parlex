package parser

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexer"
	"github.com/adamcolton/parlex/tree"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGpParsePR(t *testing.T) {
	lxr, err := lexer.New(`
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
	p := Packrat(grmr)
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

func TestParensPR(t *testing.T) {
	lxr, err := lexer.New(`
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
	p := Packrat(grmr)
	assert.NoError(t, err)
	pn := p.Parse(lxs)
	assert.NotNil(t, pn)
}

func TestLeftRecursion(t *testing.T) {
	lxr, err := lexer.New(`
    ( /\(/
    ) /\)/
    op /[+\-\*\/]/
    int /\d+/
    space /\s+/ -
  `)
	assert.NoError(t, err)
	grmr, err := grammar.New(`
    E -> E op E
      -> ( E )
      -> int
  `)
	assert.NoError(t, err)

	s := "5*(1+2)*3"
	lxs := lxr.Lex(s)
	p := Packrat(grmr)
	assert.NoError(t, err)
	pn := p.Parse(lxs)
	assert.NotNil(t, pn)

	if tpn, ok := pn.(*tree.PN); pn != nil && ok {
		expected, _ := tree.New(`
      E {
        E {
          int: '5'
        }
        op: '*'
        E {
          E {
            (: '('
            E {
              E {
                int: '1'
              }
              op: '+'
              E {
                int: '2'
              }
            }
            ): ')'
          }
          op: '*'
          E {
            int: '3'
          }
        }
      }
    `)
		if expected.String() != tpn.String() {
			t.Error(tpn.String())
		}
	} else {
		t.Error("Parse node should be of type *tree.PN")
	}
}

func TestCyclicRecursion(t *testing.T) {
	lxr, err := lexer.New(`
    + /\+/
    - /\-/
  `)
	assert.NoError(t, err)
	grmr, err := grammar.New(`
    A -> B +
      -> +
    B -> A A
      -> -
  `)
	assert.NoError(t, err)

	s := "+-++"
	lxs := lxr.Lex(s)
	p := Packrat(grmr)
	assert.NoError(t, err)
	var pn parlex.ParseNode
	ch := make(chan bool)
	go func() {
		pn = p.Parse(lxs)
		assert.NotNil(t, pn)
		ch <- true
	}()

	select {
	case <-ch:
	case <-time.After(time.Millisecond * 20):
		t.Error("Timeout")
	}

	expected, _ := tree.New(`
    A {
      B {
        A {
          +: '+'
        }
        A {
          B {
            -: '-'
          }
          +: '+'
        }
      }
      +: '+'
    }
  `)

	if expected.String() != pn.(*tree.PN).String() {
		t.Error(pn.(*tree.PN).String())
	}
}
