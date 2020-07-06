package desultory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUniqueStrings(t *testing.T) {
	assert := assert.New(t)
	l := make(map[string]bool, 10)
	for i := 0; i < 10; i++ {
		s := GetUniqueString(1)
		assert.NotEmpty(s)
		assert.Len(s, 1)
		assert.NotContains(l, s)
		l[s] = true
	}
}
