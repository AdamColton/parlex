// Package stacklexer provides a more advanced lexer. See readme for full
// details.
package stacklexer

import (
	"errors"
	"fmt"
	"github.com/adamcolton/parlex/lexeme"
	"github.com/adamcolton/parlex/symbol/setsymbol"
	"regexp"
	"strconv"
	"strings"
)

// StackLexer is defined as a set of sublexer and will lex a string using a
// stack of lexers to provide more power than a simple lexer.
type StackLexer struct {
	lexers  map[string]*subLexer
	start   *subLexer
	set     *setsymbol.Set
	Error   string
	compare func(e1, p1, e2, p2 int) bool
	insert  struct {
		startKind string
		startVal  string
		endKind   string
		endVal    string
	}
}

type subLexer struct {
	*StackLexer
	order           []int
	rules           []*rule
	priorityCounter int
	name            string
	inheirit        struct {
		from []*subLexer
		by   []*subLexer
	}
}

type rule struct {
	kind       int
	re         *regexp.Regexp
	discard    bool
	priority   int
	push       string
	pop        int
	submatches []submatch
}

type submatch struct {
	section int
	str     string
}

var intRe = regexp.MustCompile(`\d+`)

func parseSubmatches(str string) []submatch {
	if str == "" {
		return nil
	}
	segs := strings.Split(str, "|")
	ms := make([]submatch, len(segs))
	for i, s := range segs {
		if intRe.MatchString(s) {
			sec, _ := strconv.Atoi(s)
			ms[i] = submatch{section: sec}
		} else {
			ms[i] = submatch{section: -1, str: s}
		}
	}
	return ms
}

// ErrCyclic will be returned if a stack lexers form a cyclic inheritance
var ErrCyclic = errors.New("Cyclic Inheritance")

var subParserDef = regexp.MustCompile(`==\s*([a-zA-Z_][a-zA-Z_0-9]*)\s*(?:==)?\s*\n`)
var subParserLine = regexp.MustCompile(`([^\/\s]+)\s*(?:\/((?:[^\/\\]|(?:\\\/?))+)\/\s*(?:\(((?:[^\\\)]|(?:\\[^\n]))*)\))?)?\s*((?:\^+)|(?:[a-zA-Z_][a-zA-Z_0-9]*))?\s*(-?)`)

// New will try to parse any definitions it is given. If parsing fails,
// *Stacklexer will be nil and error returned. If the definition parses
// successfully, a *StackLexer is returned and error is nil.
func New(definitions ...string) (*StackLexer, error) {
	l := &StackLexer{
		lexers:  make(map[string]*subLexer),
		set:     setsymbol.New(),
		Error:   "Error",
		compare: lengthThenPriority,
	}
	subLexerDefs := make(map[string]string)
	for _, definition := range definitions {
		ms := subParserDef.FindAllStringSubmatchIndex(definition, -1)
		cur := struct {
			name string
			idx  int
		}{
			idx: -1,
		}
		for _, m := range ms {
			if cur.idx != -1 {
				subLexerDefs[cur.name] = definition[cur.idx:m[0]]
			}
			cur.idx = m[1]
			cur.name = definition[m[2]:m[3]]
			sl := &subLexer{
				StackLexer: l,
				name:       cur.name,
			}
			l.lexers[cur.name] = sl
			if l.start == nil {
				l.start = sl
			}
		}
		if cur.idx != -1 {
			subLexerDefs[cur.name] = definition[cur.idx:]
		}
	}

	done := make(map[string]bool)
	stack := make(map[string]bool)
	for _, sl := range l.lexers {
		sl.parse(subLexerDefs, done, stack)
	}

	return l, nil
}

// Must will parse the definitions given and panic if there is an error.
func Must(definitions ...string) *StackLexer {
	sl, err := New(definitions...)
	if err != nil {
		panic(err)
	}
	return sl
}

// InsertStart will insert a lexeme at the start of any results. This can be
// helpful to add a special lexeme to indicate the beginning or add something
// like a newline to make the format more consistent.
func (l *StackLexer) InsertStart(kind, val string) *StackLexer {
	l.insert.startKind, l.insert.startVal = kind, val
	return l
}

// InsertEnd will insert a lexeme at the end of any results. This can be helpful
// to add a special lexeme to indicate the end or add something like a newline
// to make the format more consistent.
func (l *StackLexer) InsertEnd(kind, val string) *StackLexer {
	l.insert.endKind, l.insert.endVal = kind, val
	return l
}

func (sl *subLexer) parse(defs map[string]string, done, stack map[string]bool) error {
	if stack[sl.name] {
		return ErrCyclic
	}
	if done[sl.name] {
		return nil
	}
	done[sl.name] = true
	stack[sl.name] = true
	var err error
	for _, line := range strings.Split(defs[sl.name], "\n") {
		if err != nil {
			break
		}
		m := subParserLine.FindStringSubmatch(line)
		if len(m) != 6 {
			//TODO: if line is not empty, give warning or error
			continue
		}

		if from, found := sl.lexers[m[1]]; found {
			err = from.addHeir(sl, defs, done, stack)
			continue
		}

		err = sl.addMatch(m)
	}
	stack[sl.name] = false
	return err
}

func (sl *subLexer) addMatch(m []string) error {
	i := 2
	if m[i] == "" {
		// if there is no regex, the word becomes the regex
		i = 1
	}
	re, err := regexp.Compile(m[i])
	if err != nil {
		return err
	}

	push := m[4]
	pop := 0
	for _, r := range push {
		if r != '^' {
			pop = 0
			break
		}
		pop++
	}
	if pop > 0 {
		push = ""
	}

	r := rule{
		kind:       sl.set.Str(m[1]).Idx(),
		re:         re,
		push:       push,
		pop:        pop,
		discard:    m[5] == "-",
		submatches: parseSubmatches(m[3]),
	}
	sl.addRule(r)
	return nil
}

func (sl *subLexer) addRule(r rule) {
	ln := len(sl.rules)
	if r.kind >= ln {
		sl.rules = append(sl.rules, make([]*rule, 1+r.kind-ln)...)
	}
	if sl.rules[r.kind] != nil {
		return
	}
	r.priority = sl.priorityCounter
	sl.priorityCounter++
	sl.rules[r.kind] = &r
	sl.order = append(sl.order, r.kind)
	for _, heir := range sl.inheirit.by {
		heir.addRule(r)
	}
}

func (sl *subLexer) addHeir(heir *subLexer, defs map[string]string, done, stack map[string]bool) error {
	err := sl.parse(defs, done, stack)
	if err != nil {
		return err
	}

	sl.inheirit.by = append(sl.inheirit.by, heir)
	heir.inheirit.from = append(heir.inheirit.from, sl)

	for _, r := range sl.rules {
		if r != nil {
			heir.addRule(*r)
		}
	}
	return nil
}

type errLexeme struct {
	*lexeme.Lexeme
}

func (e *errLexeme) Error() string {
	return fmt.Sprintf("Lex Error %d:%d) %s", e.L, e.C, e.Value())
}

// ByLength sets the lexer to choose the longest match and use priority to
// decide a tie. This is the default.
func (l *StackLexer) ByLength() *StackLexer {
	l.compare = lengthThenPriority
	return l
}

// ByPriority sets the lexer to choose the highest priority match and use the
// length to decide a tie.
func (l *StackLexer) ByPriority() *StackLexer {
	l.compare = priorityThenLength
	return l
}

func priorityThenLength(e1, p1, e2, p2 int) bool {
	return p2 == -1 || p1 < p2 || (p1 == p2 && e1 > e2)
}
func lengthThenPriority(e1, p1, e2, p2 int) bool {
	return e1 > e2 || (e1 == e2 && (p2 == -1 || p1 < p2))
}
