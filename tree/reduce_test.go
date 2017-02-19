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
		parlex.Symbol("T"): ReduceSingleChild,
		parlex.Symbol("E"): ReduceSingleChild,
		parlex.Symbol("P"): func(node *PN) {
			PromoteChild(node, 1)
		},
	}
	reducer.Reduce(pn)
	fmt.Println(pn)
}
