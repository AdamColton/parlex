package main

import (
	"fmt"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexer"
	"github.com/adamcolton/parlex/parser"
	"github.com/adamcolton/parlex/tree"
)

const lexerRules = `
  space  /\s+/ -
  number /\d*\.?\d+/
  string /(\"([^\"\\]|(\\.))*\")/
  ,
  {
  }
  [      /\[/
  ]
  :
`
const grammarRules = `
  Value         -> string
                -> number
                -> Array
                -> Object
  Object        -> { ObjectDef }
  ObjectDef     -> KeyVal , ObjectDef
                -> KeyVal
  KeyVal        -> Key : Value
  Key           -> string
  Array         -> [ ArrayContents ]
  ArrayContents -> Value , ArrayContents
                -> Value
`

func reduceList(node *tree.PN) {
	node.RemoveChild(1)       // remove comma
	node.PromoteChildrenOf(1) // promote recursive portion of list
}

func reduceListWrapper(node *tree.PN) {
	node.RemoveChildren(0, -1)  // remove wrapper
	node.PromoteGrandChildren() // promote rows in ObjectDef
}

var reducer = tree.Reducer{
	"Value":         tree.PromoteSingleChild,
	"ObjectDef":     reduceList,
	"Object":        reduceListWrapper,
	"KeyVal":        tree.RemoveChild(1), // remove :
	"Key":           tree.PromoteChildValue(0),
	"ArrayContents": reduceList,
	"Array":         reduceListWrapper,
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	lxr, err := lexer.New(lexerRules)
	check(err)

	grmr, err := grammar.New(grammarRules)
	check(err)

	p, err := parser.TopDown(grmr)
	check(err)

	s := `
    {
      "test":"test",
      "foo": "bar",
      "pi":[3,1,4,1,5],
      "sub":{
        "one":1,
        "two":2
      },
      "nestedArr":[
        ["A","B","C"],
        [42,43,44],
        [{"name":"Adam"}, {"name":"Maggie"}, {"name":"Bea"}]
      ]
    }
  `

	tr := parlex.Run(s, lxr, p, reducer)
	fmt.Println(tr)
}
