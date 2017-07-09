package tree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReduce(t *testing.T) {
	pn, _ := New(`
    E {
      T {
        P {
          (: "("
          E {
            T {
              int: "1"
            }
            op: "+"
            E {
              T {
                int: "2"
              }
            }
          }
          ): ")"
        }
      }
      op: "*"
      E {
        T {
          int: "3"
        }
      }
    }
  `)
	assert.NotNil(t, pn)

	reducer := Reducer{
		"T": PromoteSingleChild,
		"E": PromoteSingleChild,
		"P": func(node *PN) {
			node.PromoteChild(1)
		},
	}
	pn = reducer.Reduce(pn).(*PN)

	expected, _ := New(`
    E {
      E {
        (: "("
        int: "1"
        op: "+"
        int: "2"
        ): ")"
      }
      op: "*"
      int: "3"
    }
  `)
	assert.Equal(t, expected.String(), pn.String())
}

func TestPromoteChildValue(t *testing.T) {
	pn, _ := New(`
    KeyVal {
      string: "name"
      colon: ":"
      string: "Adam"
    }
  `)
	pn.PromoteChildValue(0)
	assert.Equal(t, "name", pn.Value())
	assert.Len(t, pn.C, 2)
}

func TestRemoveChildValue(t *testing.T) {
	pn, _ := New(`
    KeyVal {
      string: "name"
      colon: ":"
      string: "Adam"
    }
  `)
	pn.RemoveChild(1)
	assert.Equal(t, "name", pn.C[0].Value())
	assert.Equal(t, "Adam", pn.C[1].Value())
	assert.Len(t, pn.C, 2)
}
