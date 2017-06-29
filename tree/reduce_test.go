package tree

import (
	"github.com/adamcolton/parlex/symbol/stringsymbol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReduce(t *testing.T) {
	pn, _ := New(`
    E {
      T {
        P {
          (: '('
          E {
            T {
              int: '1'
            }
            op: '+'
            E {
              T {
                int: '2'
              }
            }
          }
          ): ')'
        }
      }
      op: '*'
      E {
        T {
          int: '3'
        }
      }
    }
  `)
	assert.NotNil(t, pn)

	reducer := Reducer{
		stringsymbol.Symbol("T"): PromoteSingleChild,
		stringsymbol.Symbol("E"): PromoteSingleChild,
		stringsymbol.Symbol("P"): func(node *PN) {
			node.PromoteChild(1)
		},
	}
	pn = reducer.Reduce(pn).(*PN)

	expected, _ := New(`
    E {
      E {
        (: '('
        int: '1'
        op: '+'
        int: '2'
        ): ')'
      }
      op: '*'
      int: '3'
    }
  `)
	assert.Equal(t, expected.String(), pn.String())
}
