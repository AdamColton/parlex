package tree

import (
	"fmt"
	"github.com/adamcolton/parlex"
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
		parlex.Symbol("T"): PromoteSingleChild,
		parlex.Symbol("E"): PromoteSingleChild,
		parlex.Symbol("P"): func(node *PN) {
			node.PromoteChild(1)
		},
	}
	pn = reducer.Reduce(pn).(*PN)
	fmt.Println(pn)
}
