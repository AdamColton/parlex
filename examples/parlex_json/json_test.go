package main

import (
	"bytes"
	"github.com/adamcolton/parlex/tree"
	"github.com/stretchr/testify/assert"
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

	tr1, err := runner.Run(s)
	assert.NoError(t, err)

	var buf bytes.Buffer
	prettyPrint(tr1, &buf, "")

	tr2, err := runner.Run(buf.String())
	assert.Equal(t, tr1.(*tree.PN).String(), tr2.(*tree.PN).String())
}

func TestCase(t *testing.T) {
	t.Skip()
	s := `{"test": true}`
	tr, err := runner.Run(s)
	assert.NotNil(t, err)
	var buf bytes.Buffer
	prettyPrint(tr, &buf, "")
	if s != buf.String() {
		t.Error(buf.String())
	}
}
