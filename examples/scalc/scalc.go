package scalc

import (
	"fmt"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexer/simplelexer"
	"github.com/adamcolton/parlex/parser/packrat"
	"github.com/adamcolton/parlex/tree"
	"math"
	"strconv"
)

const lexerRules = `
  space /\s+/ -
  int   /(\+|-)?\d+/
  dec   /\.\d+/
  uop   /(--)/
  bop   /[\*\/+\-\^%]/
  sop   /(len)|(sum)|(avg)|(min)|(max)|(first)|(last)/
  smp   /(swap)|(drop)|(clear)/
  (     /\(/
  )     /\)/
`

const grammarRules = `
  Stack  -> Stack smp
         -> E Stack
         -> Stack P Stack
         ->
  E      -> Stack sop
         -> E uop
         -> E E bop
         -> Number
  Number -> int
         -> int dec
  P      -> ( Stack )
`

var rdcr = tree.Reducer{
	"Stack": stack,
	"E":     tree.PromoteChild(-1),
	"P":     tree.ReplaceWithChild(1),
}

func stack(node *tree.PN) {
	if node.ChildAt(-1, "Stack") {
		node.PromoteChildrenOf(-1)
	}
	if node.ChildAt(0, "Stack") {
		node.PromoteChildrenOf(0)
	}
	if node.ChildAt(-1, "smp") {
		node.PromoteChild(-1)
	} else {
		node.PromoteSingleChild()
	}
}

var lxr = parlex.MustLexer(simplelexer.New(lexerRules))
var grmr = parlex.MustGrammar(grammar.New(grammarRules))
var prsr = packrat.New(grmr)

type pfloat struct {
	v float64
	p int
}

func (p pfloat) String() string {
	f := fmt.Sprintf("%%.%df", p.p)
	return fmt.Sprintf(f, p.v)
}

func eval(node *tree.PN) []string {
	stack := evalStack(node)
	out := make([]string, len(stack))
	for i, s := range stack {
		out[i] = s.String()
	}
	return out
}

func evalStack(node *tree.PN) []pfloat {
	kind := node.Kind().String()

	if kind == "smp" {
		evalSmp(node)
		kind = "Stack"
	}

	switch kind {
	case "Stack":
		out := make([]pfloat, len(node.C))
		for i, ch := range node.C {
			out[i] = evalE(ch)
		}
		return out
	default:
		return []pfloat{evalE(node)}
	}
	return nil
}

func evalSmp(op *tree.PN) {
	switch op.Value() {
	case "swap":
		ln := len(op.C)
		if ln > 1 {
			op.C[ln-1], op.C[ln-2] = op.C[ln-2], op.C[ln-1]
		}
	case "drop":
		if len(op.C) > 0 {
			op.C = op.C[:len(op.C)-1]
		}
	case "clear":
		op.C = nil
	}
}

func evalE(node *tree.PN) pfloat {
	switch node.Kind().String() {
	case "Number":
		if c := node.Children(); c == 2 {
			c1 := node.C[1].Value()
			f, _ := strconv.ParseFloat(node.C[0].Value()+c1, 64)
			return pfloat{f, len(c1) - 1}
		} else if c == 1 {
			f, _ := strconv.ParseFloat(node.C[0].Value(), 64)
			return pfloat{f, 0}
		}
	case "uop":
		return evalUop(node.C[0], node)
	case "bop":
		return evalBop(node.C[0], node.C[1], node)
	case "sop":
		return evalSop(evalStack(node.C[0]), node)
	}
	return pfloat{}
}

func evalUop(a, op *tree.PN) pfloat {
	ae := evalE(a)
	switch op.Value() {
	case "--":
		ae.v = -ae.v
	}
	return ae
}

func evalBop(a, b, op *tree.PN) pfloat {
	ae := evalE(a)
	be := evalE(b)
	p := maxPrecision(ae, be)
	var v float64
	switch op.Value() {
	case "+":
		v = ae.v + be.v
	case "*":
		v = ae.v * be.v
	case "/":
		v = ae.v / be.v
	case "-":
		v = ae.v - be.v
	case "^":
		v = math.Pow(ae.v, be.v)
	case "%":
		v = math.Mod(ae.v, be.v)
	}

	return pfloat{v, p}
}

func evalSop(stack []pfloat, op *tree.PN) pfloat {
	var v pfloat
	switch op := op.Value(); op {
	case "sum", "avg":
		v.p = maxPrecision(stack...)
		for _, p := range stack {
			v.v += p.v
		}
		if op == "avg" && len(stack) > 0 {
			v.v /= float64(len(stack))
		}
	case "len":
		v.v = float64(len(stack))
	case "min":
		if len(stack) > 0 {
			v = stack[0]
			for _, p := range stack[1:] {
				if p.v < v.v {
					v = p
				}
			}
		}
	case "max":
		if len(stack) > 0 {
			v = stack[0]
			for _, p := range stack[1:] {
				if p.v > v.v {
					v = p
				}
			}
		}
	case "last":
		if len(stack) > 0 {
			v = stack[0]
		}
	case "first":
		if len(stack) > 0 {
			v = stack[len(stack)-1]
		}
	}

	return v
}

func maxPrecision(pfs ...pfloat) int {
	m := 0
	for _, p := range pfs {
		if p.p > m {
			m = p.p
		}
	}
	return m
}
