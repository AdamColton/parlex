package parlex

// Run performs the lexing, parsing and reducing for an input
func Run(input string, lexer Lexer, parser Parser, reducer Reducer) ParseNode {
	lexemes := lexer.Lex(input)
	tree := parser.Parse(lexemes)
	if reducer != nil {
		tree = reducer.Reduce(tree)
	}
	return tree
}

// Runner holds a Lexer, Parser and Reducer and uses them to operate on an input
// string
type Runner struct {
	lexer   Lexer
	parser  Parser
	reducer Reducer
}

// New returns a new runner. The reducer can be nil.
func New(lexer Lexer, parser Parser, reducer Reducer) *Runner {
	return &Runner{
		lexer:   lexer,
		parser:  parser,
		reducer: reducer,
	}
}

// Run using the Parser, Lexer and Reducer in the Runner.
func (r *Runner) Run(input string) ParseNode {
	return Run(input, r.lexer, r.parser, r.reducer)
}
