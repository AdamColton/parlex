package simplelexer

import (
	"fmt"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/symbol/setsymbol"
	"regexp"
	"strings"
)

// DefaultErrorString is the value that will be assigned to any lex errors
var DefaultErrorString = "Error"

// New returns a new Lexer. It can be provided definitions, though most often it
// is only given a single definition. Each line will be one rule.
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

// InsertStart will insert a lexeme at the start of any results. This can be
// helpful to add a special lexeme to indicate the beginning or add something
// like a newline to make the format more consistent.
func (l *Lexer) InsertStart(kind, val string) *Lexer {
	l.insert.startKind, l.insert.startVal = kind, val
	return l
}

// InsertEnd will insert a lexeme at the end of any results. This can be helpful
// to add a special lexeme to indicate the end or add something like a newline
// to make the format more consistent.
func (l *Lexer) InsertEnd(kind, val string) *Lexer {
	l.insert.endKind, l.insert.endVal = kind, val
	return l
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
			return fmt.Errorf("Duplicate Kind: %s", l.set.ByIdx(r.kind).String())
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
