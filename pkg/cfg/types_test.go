package cfg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRepoPath(t *testing.T) {
	repo := &GitRepo{
		Branch: "foo:bar/baz",
		URL:    "https://qux.quux",
	}
	result := repo.GetRealPath()
	assert.Equal(t, "/tmp/laminar/https---qux.quux-foo-bar-baz", result)
}
