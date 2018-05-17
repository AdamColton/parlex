package tree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConditional(t *testing.T) {
	pn, _ := New(`
    E {
      P {
        lp: "("
        num: "6"
        rp: ")"
      }
      foo:"bar"
      P {
        lp: "("
        num: "7"
        rp: ")"
      }
    }
  `)
  var hasThreeChildren Condition= func(node *PN) bool{
    return len(node.C) == 3
  }

  r := If(hasThreeChildren, PromoteChild(1), nil)
  r(pn)
  assert.NotNil(t,pn)

  assert.Equal(t, "bar", pn.Value())
}
