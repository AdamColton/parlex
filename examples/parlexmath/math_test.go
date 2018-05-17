package parlexmath

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMath(t *testing.T) {
	tt := map[string]float64{
		"1+2+3":         6,
		"2*2+3":         7,
		"1+2*3":         7,
		"1+2-3":         0,
		"1*(2+3)":       5,
		"2*(2+3*2)-2*3": 10,
		"2*-3":          -6,
		"-1.5*4":        -6,
		"11--11":        22,
	}
	for str, expected := range tt {
		t.Run(str, func(t *testing.T) {
			actual, err := Eval(str)
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}
