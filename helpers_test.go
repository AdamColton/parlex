package parlex

type symbol string

func (s symbol) String() string { return string(s) }

type production []symbol

func (p production) Symbols() int {
	return len(p)
}
func (p production) Symbol(i int) Symbol {
	if i < len(p) {
		return p[i]
	}
	return nil
}

func (p production) Iter() *ProductionIterator {
	return &ProductionIterator{
		Production: p,
	}
}

type productions []production

func (p productions) Productions() int {
	return len(p)
}

func (p productions) Production(i int) Production {
	if i < len(p) {
		return p[i]
	}
	return nil
}

func (p productions) Iter() *ProductionsIterator {
	return &ProductionsIterator{
		Productions: p,
	}
}

type lx struct {
	k    symbol
	v    string
	l, c int
}

func (l *lx) Kind() Symbol             { return l.k }
func (l *lx) Value() string            { return l.v }
func (l *lx) Pos() (line int, col int) { return l.l, l.c }

type testGrammar struct {
	order       []Symbol
	productions map[symbol]productions
	cur         symbol
}

func (tg *testGrammar) Productions(s Symbol) Productions {
	return tg.productions[s.(symbol)]
}

func (tg *testGrammar) NonTerminals() []Symbol {
	return tg.order
}

func (tg *testGrammar) new(s symbol) {
	tg.cur = s
	tg.order = append(tg.order, s)
}

func (tg *testGrammar) add(symbols ...symbol) {
	tg.productions[tg.cur] = append(tg.productions[tg.cur], production(symbols))
}

func (tg *testGrammar) reset() {
	tg.productions = make(map[symbol]productions)
	tg.order = nil
}

type testLexer struct{}

func (*testLexer) Lex(str string) []Lexeme { return nil }

var testErr = strErr("Test Error")

type testParser struct{}

func (*testParser) Parse([]Lexeme) ParseNode { return nil }
