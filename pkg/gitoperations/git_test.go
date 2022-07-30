package gitoperations

import (
	"github.com/digtux/laminar/pkg/cfg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRepoPath(t *testing.T) {
	result := GetRepoPath(cfg.GitRepo{
		Branch: "foo:bar/baz",
		URL:    "https://qux.quux",
	})
	assert.Equal(t, "/tmp/https---qux.quux-foo-bar-baz", result)
}
