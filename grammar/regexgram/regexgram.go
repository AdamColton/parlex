package regexgram

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar"
	"github.com/adamcolton/parlex/lexeme"
	"github.com/adamcolton/parlex/lexer/simplelexer"
	"github.com/adamcolton/parlex/parser/packrat"
	"github.com/adamcolton/parlex/symbol/setsymbol"
	"github.com/adamcolton/parlex/tree"
	"strings"
)

const lexerProductions = `
  rarr     /->/
  symbol   /\w+/
  repeats  /\*/
  optional /\?/
  or       /\|/
  (        /\(/
  )        /\)/
  comment  /[\n\r]\s*\/\/[^\n\r]*/ -
  nl       /[\n\r]\s*/
  space    /[ \t]+/ -
`

const grammarProductions = `
  Grammar      -> NL Production Productions NL 
  Productions  -> Production Productions
               -> ContinueProd Productions
               -> 
  Production   -> nl symbol rarr Symbols
  ContinueProd -> nl rarr Symbols
  Symbols      -> Symbol Symbols
               ->
  Symbol       -> Group
               -> OrSymbol
               -> RepSymbol
               -> OptSymbol
               -> symbol
  Group        -> ( Symbols )
  OptSymbol    -> Group optional
               -> symbol optional
  RepSymbol    -> Group repeats
               -> symbol repeats
  OrSymbol     -> Group or MoreOr
               -> OptSymbol or MoreOr
               -> RepSymbol or MoreOr
               -> symbol or MoreOr
  MoreOr       -> Group or MoreOr
               -> OptSymbol or MoreOr
               -> RepSymbol or MoreOr
               -> symbol or MoreOr
               -> Group
               -> OptSymbol
               -> RepSymbol
               -> symbol
  NL           -> nl
               ->
`

var rdcr = tree.Reducer{
	"Grammar": tree.
		If(
			tree.ChildIs(0, "NL"), // Remove newlines at start of document
			tree.RemoveChild(0),
			nil,
		).
		If(
			tree.ChildIs(-1, "NL"), // Remove newlines at end of document
			tree.RemoveChild(-1),
			nil,
		).
		PromoteChildrenOf(1), // promote children of Productions
	"Productions": tree.
		PromoteChildrenOf(1), // promote children of Productions
	"Production": tree.
		RemoveChildren(0, 1). // remove new-line and rarr
		PromoteChildValue(0). // promote the non-terminal to be the production value
		PromoteChildrenOf(0), // replace Symbols with it's children
	"ContinueProd": tree.
		RemoveChildren(0, 0). // remove new-line and rarr
		PromoteChildrenOf(0), // replace Symbols with it's children
	"Symbols": tree.
		PromoteChildrenOf(-1), // Children of last will be either more symbols or nil
	"Symbol": tree.PromoteSingleChild,
	"OptSymbol": tree.
		RemoveChild(-1), // Remove ?
	"RepSymbol": tree.
		RemoveChild(-1), // Remove *
	"OrSymbol": tree.
		RemoveChild(1).       // Remove |
		PromoteChildrenOf(1), // promote the rest of the or condition
	"MoreOr": tree.
		RemoveChild(1).       // Remove |
		PromoteChildrenOf(1), // promote the rest of the or condition
	"Group": tree.
		RemoveChildren(0, -1). // Remove ( )
		PromoteChildrenOf(0),  // Promote the children to replace Group
}

var lxr = parlex.MustLexer(simplelexer.New(lexerProductions)).(*simplelexer.Lexer).
	InsertStart("nl", "\n")
var grmr = parlex.MustGrammar(grammar.New(grammarProductions))
var prsr = packrat.New(grmr)

var runner = parlex.New(lxr, prsr, rdcr)

type evalOp struct {
	grammar   *grammar.Grammar
	set       *setsymbol.Set
	rdcr      tree.Reducer
	stack     []*tree.PN
	nonterm   string
	bludgeons map[string][]string
	done      map[string]rules
}

func evalGrammar(node *tree.PN) (*grammar.Grammar, tree.Reducer) {
	op := &evalOp{
		grammar:   grammar.Empty(),
		set:       setsymbol.New(),
		rdcr:      tree.Reducer{},
		bludgeons: make(map[string][]string),
		done:      make(map[string]rules),
	}
	for _, c := range node.C {
		if c.Kind().String() == "Production" {
			op.nonterm = c.Value()
		} else {
			c.Lexeme.(*lexeme.Lexeme).V = op.nonterm
		}
		op.stack = append(op.stack, c)
	}

	for len(op.stack) > 0 {
		node := op.stack[0]
		op.stack = op.stack[1:]
		op.evalProd(node)
	}

	for nonterm, symbols := range op.bludgeons {
		op.rdcr[nonterm] = bludgeon(symbols)
	}

	return op.grammar, op.rdcr
}

func (op *evalOp) evalProd(node *tree.PN) {
	op.nonterm = node.Value()
	rs := rules{rule{}}
	for _, c := range node.C {
		rs = mergeRules(rs, op.evalSymbol(c))
	}

	for _, r := range rs {
		prod := op.ruleToProd(r)
		op.grammar.Add(op.set.Str(op.nonterm), prod)
	}
}

func (op *evalOp) evalSymbol(node *tree.PN) rules {
	switch node.Kind().String() {
	case "symbol":
		return rules{rule{node.Value()}}
	case "OptSymbol":
		return append(op.evalSymbol(node.C[0]), rule{})
	case "OrSymbol":
		var rs rules
		for _, c := range node.C {
			rs = append(rs, op.evalSymbol(c)...)
		}
		return rs
	case "Group":
		var rc ruleComb
		for _, c := range node.C {
			rc = append(rc, op.evalSymbol(c))
		}
		return rc.reduce()
	case "RepSymbol":
		return op.addRepeatAsProduction(node)
	}
	return nil
}

// addRepeatAsProduction creates two productions
// given:
// E*
// It adds
// E* -> E E*
//    ->
// And adds a rule to the reducer
func (op *evalOp) addRepeatAsProduction(node *tree.PN) rules {
	symName := op.getName(node)
	op.bludgeons[op.nonterm] = append(op.bludgeons[op.nonterm], symName)
	if rs, ok := op.done[symName]; ok {
		return rs
	}
	cp := tree.Clone(node)

	symNode := &tree.PN{
		Lexeme: &lexeme.Lexeme{
			K: op.set.Str("symbol"),
			V: symName,
		},
	}
	prod := &tree.PN{
		Lexeme: &lexeme.Lexeme{
			K: op.set.Str("Production"),
			V: symName,
		},
		C: append(cp.C, symNode),
	}
	for _, c := range prod.C {
		c.P = prod
	}
	op.stack = append(op.stack, prod)

	nilProd := &tree.PN{
		Lexeme: &lexeme.Lexeme{
			K: op.set.Str("Production"),
			V: symName,
		},
	}
	op.stack = append(op.stack, nilProd)

	rs := rules{rule{symName}}
	op.done[symName] = rs
	return rs
}

func bludgeon(symbols []string) func(*tree.PN) {
	return func(node *tree.PN) {
		for i := 0; i < len(node.C); i++ {
			if node.ChildAt(i, symbols...) {
				node.PromoteChildrenOf(i)
				i--
			}
		}
	}
}

func (op *evalOp) getName(node *tree.PN) string {
	switch node.Kind().String() {
	case "symbol":
		return node.Value()
	case "OptSymbol":
		return op.getName(node.C[0]) + "?"
	case "RepSymbol":
		return op.getName(node.C[0]) + "*"
	case "OrSymbol":
		var strs []string
		for _, c := range node.C {
			strs = append(strs, op.getName(c))
		}
		return strings.Join(strs, "|")
	case "Group":
		var strs []string
		for _, c := range node.C {
			strs = append(strs, op.getName(c))
		}
		return "(" + strings.Join(strs, "_") + ")"
	}

	return ""
}

func (op *evalOp) ruleToProd(r rule) *setsymbol.Production {
	prod := op.set.Production()
	for _, s := range r {
		prod.AddSymbols(op.set.Str(s))
	}
	return prod
}

type rule []string
type rules []rule
type ruleComb []rules

func mergeRules(a, b rules) rules {
	var out rules
	for _, ra := range a {
		for _, rb := range b {
			out = append(out, append(ra, rb...))
		}
	}
	return out
}

func (rc ruleComb) reduce() rules {
	if len(rc) == 0 {
		return nil
	}
	ra := rc[0]
	for _, rb := range rc[1:] {
		ra = mergeRules(ra, rb)
	}
	return ra
}

// New takes a grammar string and returns a grammar, reducer and error.
func New(grammarString string) (*grammar.Grammar, tree.Reducer, error) {
	parseTree, err := runner.Run(grammarString)
	if err != nil {
		return nil, nil, err
	}
	g, r := evalGrammar(parseTree.(*tree.PN))
	return g, r, nil
}

// Must returns a grammar and a reducer. If it fails to parse the grammar string
// it will panic.
func Must(grammarString string) (*grammar.Grammar, tree.Reducer) {
	g, r, err := New(grammarString)
	if err != nil {
		panic(err)
	}
	return g, r
}
