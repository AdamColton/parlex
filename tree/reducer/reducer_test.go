package reducer

import (
	//"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/tree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStuff(t *testing.T) {
	str := `
Value       PromoteSingleChild
Object      RemoveChildren(0, -1)               // remove { }
Array       RemoveChildren(0, -1)               // remove [ ]
KeyVal      PromoteChildValue(0).RemoveChild(0) // Promote key, remove :
// another comment 1
// another comment 2
MoreVals    ReplaceWithChild(1)
MoreKeyVals ReplaceWithChild(1)
Fooo        If(
				ChildIs(0, "string"), // Conditional
				PromoteSingleChild, // then
				RemoveChildren(0) // else
            )
`

	rdcr, err := Parse(str)
	assert.NoError(t, err)

	assert.NotNil(t, rdcr["Array"])

	pn1, err := tree.New(`
		KeyVal {
			key: "Foo"
			colon: ":"
			val: "Bar"
	  }
  `)
	assert.NoError(t, err)
	pn1 = rdcr.Reduce(pn1).(*tree.PN)
	assert.Equal(t, "KeyVal", pn1.Kind().String())
	assert.Equal(t, "Foo", pn1.Value())
	if assert.Len(t, pn1.C, 1) {
		assert.Equal(t, "val", pn1.C[0].Kind().String())
		assert.Equal(t, "Bar", pn1.C[0].Value())
	}
}
