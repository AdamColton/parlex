package parlexmath

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMath(t *testing.T) {
	expected := map[string]float64{
		"1+2+3":         6,
		"2*2+3":         7,
		"1+2*3":         7,
		"1+2-3":         0,
		"1*(2+3)":       5,
		"2*(2+3*2)-2*3": 10,
		"-2*-3":         6,
		"-1.5*-4":       6,
	}
	for str, i := range expected {
		ei, err := Eval(str)
		assert.NoError(t, err)
		if ei != i {
			t.Errorf("Got %d, expectd: %d\n", ei, i)
		}
	}
}
