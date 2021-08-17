package pike

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadUint32(t *testing.T) {
	tt := map[string]struct {
		expected uint32
		slice    []byte
	}{
		"0": {
			expected: 0,
			slice:    []byte{0, 0, 0, 0},
		},
		"1": {
			expected: 1,
			slice:    []byte{1, 0, 0, 0},
		},
		"256": {
			expected: 256,
			slice:    []byte{0, 1, 0, 0},
		},
		"65536": {
			expected: 65536,
			slice:    []byte{0, 0, 1, 0},
		},
		"16777216": {
			expected: 16777216,
			slice:    []byte{0, 0, 0, 1},
		},
		"67305985": {
			expected: 1 + 2*256 + 3*65536 + 4*16777216,
			slice:    []byte{1, 2, 3, 4},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			w := wrapper{
				slice: tc.slice,
			}
			assert.Equal(t, tc.expected, w.idxUint32())
		})
	}
}
