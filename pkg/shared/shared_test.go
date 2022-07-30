package shared

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUniqueStrings(t *testing.T) {
	result := UniqueStrings([]string{
		"a",
		"b",
		"c",
		"a",
		"c",
		"c",
	})
	assert.Equal(t, result, []string{"a", "b", "c"})
}
