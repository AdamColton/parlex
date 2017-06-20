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
  op1 /[\*\/]/
  op2 /[\+\-]/
  ( /\(/
  ) /\)/
`
const grammarRules = `
  E  -> E op1 E
  	 -> E2
  E2 -> E op2 E
     -> ( E )
     -> number
`

var reducer = tree.Reducer{
	"E": tree.PromoteSingleChild,
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

	p, err := parser.Packrat(grmr)
	check(err)

	tr := parlex.Run("1+2+3", lxr, p, reducer)
	fmt.Println(tr)

	tr = parlex.Run("1*2+3", lxr, p, reducer)
	fmt.Println(tr)
}
