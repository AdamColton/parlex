package parlex

type strErr string

func (err strErr) Error() string { return string(err) }

const (
	ErrCouldNotLex    = strErr("Could Not Lex")
	ErrCouldNotParse  = strErr("Could Not Parse")
	ErrCouldNotReduce = strErr("Could Not Reduce")
	ErrBadGrammar     = strErr("Bad Grammar")
)
