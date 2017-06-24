package main

import (
	"bytes"
	"testing"
)

func TestJSON(t *testing.T) {
	s := `{
  "test": "test",
  "foo": "bar",
  "pi": [3,1,4,1,5],
  "null": null,
  "boolT": true,
  "boolF": false,
  "sub": {
    "one": 1,
    "two": 2
  },
  "nestedArr": [
    ["A","B","C"],
    [42,43,44],
    [
      {"name": "Adam"},
      {"name": "Maggie"},
      {"name": "Bea"}
    ]
  ]
}`
	tr := runner.Run(s)
	var buf bytes.Buffer
	prettyPrint(tr, &buf, "")
	if buf.String() != s {
		t.Error(buf.String())
	}
}

func TestCase(t *testing.T) {
	s := `{"test": true}`
	tr := runner.Run(s)
	var buf bytes.Buffer
	prettyPrint(tr, &buf, "")
	if s != buf.String() {
		t.Error(buf.String())
	}
}
