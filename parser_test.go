package parlex

import (
	"testing"

	"github.com/testify/assert"
)

func TestMustParser(t *testing.T) {
	p := &testParser{}
	assert.Equal(t, p, MustParser(p, nil))

	defer func() {
		assert.Equal(t, testErr, recover())
	}()
	MustParser(p, testErr)
}
