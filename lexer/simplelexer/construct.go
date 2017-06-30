package simplelexer

import (
	"errors"
	"fmt"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/symbol/setsymbol"
	"regexp"
	"strings"
)

// ErrDuplicateKind will be thrown if a rule is duplicated. Instead, use regexp
// concatination.
var ErrDuplicateKind = errors.New("Duplicate Kind")

var DefaultErrorString = "Error"

// New returns a new Lexer
func New(definitions ...string) (*Lexer, error) {
	l := &Lexer{
		compare: lengthThenPriority,
		Error:   DefaultErrorString,
		set:     setsymbol.New(),
	}
	for _, definition := range definitions {
		for _, line := range strings.Split(definition, "\n") {
			r, err := l.ruleFromLine(line)
			if err != nil {
				return nil, err
			}
			if r == nil {
				continue
			}
			err = l.addRule(r)
			if err != nil {
				return nil, err
			}
		}
	}

	return l, nil
}

var lexStr = regexp.MustCompile(`([^\/\s]+)\s*(?:\/((?:[^\/\\]|(?:\\\/?))+)\/)?\s*(-?)`)

func (l *Lexer) ruleFromLine(line string) (*rule, error) {
	m := lexStr.FindStringSubmatch(line)
	if len(m) != 4 {
		return nil, nil
	}
	i := 2
	if m[i] == "" {
		// if there is no regex, the word becomes the regex
		i = 1
	}
	re, err := regexp.Compile(m[i])
	if err != nil {
		return nil, err
	}
	return &rule{
		kind:    l.set.Str(m[1]).Idx(),
		re:      re,
		discard: m[3] == "-",
	}, nil
}

// Add a lexer rule
func (l *Lexer) Add(kind parlex.Symbol, re *regexp.Regexp, discard bool) error {
	return l.addRule(&rule{
		kind:    l.set.Symbol(kind).Idx(),
		re:      re,
		discard: discard,
	})
}

func (l *Lexer) addRule(r *rule) error {
	if r.kind < len(l.rules) {
		if l.rules[r.kind] != nil {
			return ErrDuplicateKind
		}
	} else {
		l.rules = append(l.rules, make([]*rule, 1+r.kind-len(l.rules))...)
	}

	r.priority = l.priorityCounter
	l.priorityCounter++
	l.rules[r.kind] = r
	l.order = append(l.order, r.kind)
	return nil
}

// String exports the lexer as a string. The output of String can be used to
// make a copy of the lexer.
func (l *Lexer) String() string {
	var longest int
	for _, rule := range l.rules {
		if ln := parlex.SymLen(l.set.ByIdx(rule.kind)); ln > longest {
			longest = ln
		}
	}

	format := fmt.Sprintf("%%-%ds %%s %%s", longest)
	lines := make([]string, len(l.order))
	for i, kind := range l.order {
		rule := l.rules[kind]
		d := ""
		if rule.discard {
			d = "-"
		}
		str := l.set.ByIdx(kind).String()
		re := rule.re.String()
		if re == str {
			re = ""
		} else {
			re = "/" + re + "/"
		}
		lines[i] = fmt.Sprintf(format, str, re, d)
	}
	return strings.Join(lines, "\n")
}
