package simplelexer

import (
	"errors"
	"fmt"
	"github.com/adamcolton/parlex"
	"regexp"
	"strings"
)

// ErrDuplicateKind will be thrown if a rule is duplicated. Instead, use regexp
// concatination.
var ErrDuplicateKind = errors.New("Duplicate Kind")

// New returns a new Lexer
func New(definitions ...string) (*Lexer, error) {
	l := &Lexer{
		rules:   make(map[parlex.Symbol]*rule),
		compare: lengthThenPriority,
		Error:   "Error",
	}
	for _, definition := range definitions {
		for _, line := range strings.Split(definition, "\n") {
			r, err := ruleFromLine(line)
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

func ruleFromLine(line string) (*rule, error) {
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
		kind:    parlex.Symbol(m[1]),
		re:      re,
		discard: m[3] == "-",
	}, nil
}

// Add a lexer rule
func (l *Lexer) Add(kind parlex.Symbol, re *regexp.Regexp, discard bool) error {
	return l.addRule(&rule{
		kind:    kind,
		re:      re,
		discard: discard,
	})
}

func (l *Lexer) addRule(r *rule) error {
	if l.rules[r.kind] != nil {
		return ErrDuplicateKind
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
		if ln := rule.kind.Len(); ln > longest {
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
		re := rule.re.String()
		if re == string(rule.kind) {
			re = ""
		} else {
			re = "/" + re + "/"
		}
		lines[i] = fmt.Sprintf(format, rule.kind, re, d)
	}
	return strings.Join(lines, "\n")
}
