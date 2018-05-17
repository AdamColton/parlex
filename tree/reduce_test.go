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

func TestReplaceWithChild(t *testing.T) {
  pn, _ := New(`
    P {
      lp: "("
      num: "6"
      rp: ")"
    }
  `)
  pn.ReplaceWithChild(1)
  assert.Equal(t, "6", pn.Value())
  assert.Equal(t, "num", pn.Kind().String())
}

func TestChildIs(t *testing.T) {
  pn, _ := New(`
    P {
      lp: "("
      num: "6"
      rp: ")"
    }
  `)
  assert.True(t, pn.ChildIs(-1,"rp"))
  assert.False(t, pn.ChildIs(0,"rp"))
}

func TestPromoteChildrenOf(t *testing.T) {
  pn, _ := New(`
    E {
      P {
        lp: "("
        num: "6"
        rp: ")"
      }
      foo:"bar"
    }
  `)
  pn.PromoteChildrenOf(0)
  assert.Equal(t, "lp", pn.C[0].Kind().String())
  assert.Equal(t, "bar", pn.C[3].Value())
}

func TestPromoteGrandChildren(t *testing.T) {
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
  pn.PromoteGrandChildren()
  assert.Equal(t, "lp", pn.C[0].Kind().String())
  assert.Equal(t, "7", pn.C[4].Value())
}

func TestMerge(t *testing.T){
  pn, _ := New(`
    E {
      P {
        lp: "("
        num: "6"
        rp: ")"
      }
      foo {
        val: "bar"
      }
      P {
        lp: "("
        num: "7"
        rp: ")"
      }
      foo2: "foo2 val" {
        val: "bar2"
      }
    }
  `)

  r1 := Reducer{
    "P": RemoveChildren(0,1),
  }
  r1.Add("foo", PromoteChild(0))
  assert.True(t, r1.Can(pn.C[0]))
  assert.False(t, r1.Can(pn))

  r2 := Reducer{
    "P": ReplaceWithChild(0),
    "foo2": RemoveChild(0),
  }

  r := Merge(r1,r2)
  pn = r.Reduce(pn).(*PN)
  if assert.Len(t, pn.C, 4){
    assert.Equal(t, "6", pn.C[0].Value())
    assert.Equal(t, "bar", pn.C[1].Value())
    assert.Equal(t, "7", pn.C[2].Value())
    assert.Equal(t, "foo2 val", pn.C[3].Value())
  }
}

func TestRemoveAll(t *testing.T) {
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
  r := Reducer{
    "P": RemoveAll("lp","rp"),
  }
  pn = r.Reduce(pn).(*PN)
  assert.Len(t, pn.C[0].C, 1)
}