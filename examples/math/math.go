package main

import (
	"fmt"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexer/simplelexer"
	"github.com/adamcolton/parlex/parser/packrat"
	"github.com/adamcolton/parlex/symbol/stringsymbol"
	"github.com/adamcolton/parlex/tree"
	"os"
	"strconv"
	"strings"
)

const lexerRules = `
  space  /\s+/ -
  number /\d*\.?\d+/
  op1 /[\*\/]/
  op2 /[\+\-]/
  ( /\(/
  ) /\)/
`

var lxr = parlex.MustLexer(simplelexer.New(lexerRules))

var (
	nt_E     = stringsymbol.Symbol("E")
	nt_P     = stringsymbol.Symbol("P")
	t_space  = stringsymbol.Symbol("space")
	t_number = stringsymbol.Symbol("number")
	t_op1    = stringsymbol.Symbol("op1")
	t_op2    = stringsymbol.Symbol("op2")
)

const grammarRules = `
  E -> E op2 E
    -> E op1 E
    -> number
    -> op2 number
    -> P
  P -> ( E )
`

var grmr = parlex.MustGrammar(grammar.New(grammarRules))
var prsr = packrat.New(grmr)

var reducer = tree.Reducer{
	nt_E: func(node *tree.PN) {
		if !node.PromoteSingleChild() {
			node.PromoteChild(1)
		}
	},
	nt_P: tree.ReplaceWithChild(1),
}

var runner = parlex.New(lxr, prsr, reducer)

func main() {
	fmt.Println(eval(runner.Run(strings.Join(os.Args[1:], " "))))
}

func eval(node parlex.ParseNode) float64 {
	switch node.Kind() {
	case t_number:
		i, _ := strconv.ParseFloat(node.Value(), 64)
		return i
	case t_op1, t_op2:
		a := eval(node.Child(0))
		b := eval(node.Child(1))
		switch node.Value() {
		case "+":
			return a + b
		case "*":
			return a * b
		case "/":
			return a / b
		case "-":
			return a - b
		}
	}
	return 0
}
