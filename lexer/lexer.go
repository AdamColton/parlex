package lexer

import (
	"errors"
	"fmt"
	"github.com/adamcolton/parlex"
	"regexp"
	"strings"
)

// Lexer implements parlex.Lexer. It can take a string and produce a slice of
// lexemes.
type Lexer struct {
	order   []parlex.Symbol
	rules   map[parlex.Symbol]*rule
	longest int
}

type rule struct {
	kind     parlex.Symbol
	re       *regexp.Regexp
	discard  bool
	priority int
}

func ruleFromLine(line string) (*rule, error) {
	m := lexStr.FindStringSubmatch(line)
	if len(m) == 4 {
		// if there is no regex, the word becomes the regex
		var re *regexp.Regexp
		var err error
		if m[2] == "" {
			re, err = regexp.Compile(m[1])
		} else {
			re, err = regexp.Compile(m[2])
		}
		if err != nil {
			return nil, err
		}
		return &rule{
			kind:    parlex.Symbol(m[1]),
			re:      re,
			discard: m[3] == "-",
		}, nil
	}
	return nil, nil
}

var lexStr = regexp.MustCompile(`([^\/\s]+)\s*(?:\/((?:[^\/\\]|(?:\\\/?))+)\/)?\s*(-?)`)

// ErrDuplicateKind will be thrown if a rule is duplicated. Instead, use regexp
// concatination.
var ErrDuplicateKind = errors.New("Duplicate Kind")

// New returns a new Lexer
func New(definitions ...string) (*Lexer, error) {
	rules := make(map[parlex.Symbol]*rule)
	var order []parlex.Symbol
	p := 0
	for _, definition := range definitions {
		for _, line := range strings.Split(definition, "\n") {
			rule, err := ruleFromLine(line)
			if err != nil {
				return nil, err
			}
			if rule != nil {
				if rules[rule.kind] != nil {
					return nil, ErrDuplicateKind
				}
				rule.priority = p
				p++
				rules[rule.kind] = rule
				order = append(order, rule.kind)
			}
		}
	}

	return &Lexer{
		rules:   rules,
		longest: -1,
		order:   order,
	}, nil
}

// Add a lexer rule
func (l *Lexer) Add(kind parlex.Symbol, re *regexp.Regexp, discard bool) error {
	if l.rules[kind] != nil {
		return ErrDuplicateKind
	}
	l.rules[kind] = &rule{
		kind:    kind,
		re:      re,
		discard: discard,
	}
	l.order = append(l.order, kind)
	return nil
}

// String exports the lexer as a string. The output of String can be used to
// make a copy of the lexer.
func (l *Lexer) String() string {
	if l.longest == -1 {
		for _, rule := range l.rules {
			if ln := rule.kind.Len(); ln > l.longest {
				l.longest = ln
			}
		}
	}

	format := fmt.Sprintf("%%-%ds %%s %%s", l.longest)
	segs := make([]string, len(l.order))
	for i, kind := range l.order {
		rule := l.rules[kind]
		d := ""
		if rule.discard {
			d = "-"
		}
		re := rule.re.String()
		if re == string(rule.kind) {
			re = ""
		} else {
			re = "/" + re + "/"
		}
		segs[i] = fmt.Sprintf(format, rule.kind, re, d)
	}
	return strings.Join(segs, "\n")
}

// Lex takes a string and produces a slice of lexemes that can be consumed by a
// parser.
func (l *Lexer) Lex(str string) []parlex.Lexeme {
	var lxs []parlex.Lexeme
	b := []byte(str)

	// find the first occurance of every lexer
	next := make(map[parlex.Symbol][]int)
	for kind, r := range l.rules {
		next[kind] = r.re.FindIndex(b)
	}

	errFlag := false
	errStart := 0

	for cur := 0; cur < len(b); {
		lx := &parlex.L{
			K: "",
			V: "",
		}
		lxEnd := cur
		lxP := -1
		discard := false

		// look in next for matches and take the longest one
		for kind, loc := range next {
			if loc != nil && loc[0] == cur {
				// longer always wins
				// if two are the same length, the lower priorty wins
				p := l.rules[kind].priority
				if loc[1] > lxEnd || (loc[1] == lxEnd && p < lxP) {
					lx.K = parlex.Symbol(kind)
					lx.V = string(b[loc[0]:loc[1]])
					lxEnd = loc[1]
					lxP = p
					discard = l.rules[kind].discard
				}
			}
		}

		if lxEnd == cur { // found no matches of any length
			if !errFlag { // if we're not already in an error state, start one
				errFlag = true
				errStart = cur
			}
			cur++
		} else {
			if errFlag { // if we were in an error state, resolve it by adding the error lexeme
				lxs = append(lxs, &parlex.L{
					K: "Error",
					V: string(b[errStart:cur]),
				})
				errFlag = false
			}
			if !discard {
				lxs = append(lxs, lx)
			}
			cur = lxEnd
		}

		for kind, loc := range next {
			if loc != nil && loc[0] <= cur {
				loc := l.rules[kind].re.FindIndex(b[cur:])
				if loc != nil {
					loc[0] += cur
					loc[1] += cur
				}
				next[kind] = loc
			}
		}
	}

	if errFlag { // if we were in an error state, resolve it by adding the error lexeme
		lxs = append(lxs, &parlex.L{
			K: "Error",
			V: string(b[errStart:]),
		})
	}

	return lxs
}
