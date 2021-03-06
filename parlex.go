package parlex

// Run performs the lexing, parsing and reducing for an input
func Run(input string, lexer Lexer, parser Parser, reducers ...Reducer) (ParseNode, error) {
	lexemes := lexer.Lex(input)
	if lexemes == nil {
		return nil, ErrCouldNotLex
	}
	errs := LexErrors(lexemes)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	parseTree := parser.Parse(lexemes)
	if parseTree == nil {
		return nil, ErrCouldNotParse
	}

	for _, reducer := range reducers {
		if reducer != nil {
			parseTree = reducer.Reduce(parseTree)
			if parseTree == nil {
				return nil, ErrCouldNotReduce
			}
		}
	}

	return parseTree, nil
}

// Runner holds a Lexer, Parser and Reducer and uses them to operate on an input
// string
type Runner struct {
	lexer    Lexer
	parser   Parser
	reducers []Reducer
}

// New returns a new runner. The reducer can be nil.
func New(lexer Lexer, parser Parser, reducers ...Reducer) *Runner {
	return &Runner{
		lexer:    lexer,
		parser:   parser,
		reducers: reducers,
	}
}

// Run using the Parser, Lexer and Reducer in the Runner.
func (r *Runner) Run(input string) (ParseNode, error) {
	return Run(input, r.lexer, r.parser, r.reducers...)
}
