package ecr

import (
	"fmt"
	"github.com/digtux/laminar/pkg/cache"
	"github.com/digtux/laminar/pkg/cfg"
	"github.com/digtux/laminar/pkg/shared"
	"regexp"
	"testing"
)

func TestClient_EcrWorker(t *testing.T) {
	logger := shared.GetLogger(true)
	w := New(logger, cache.Open(":memory", logger))
	w.ScanAll([]string{}, []cfg.DockerRegistry{
		{
			Reg:     "123456789012.dkr.ecr.eu-west-2.amazonaws.com/foo",
			Name:    "foo tooling",
			TimeOut: 0,
		},
	})
}

func TestCastEcr(t *testing.T) {
	//x, _ := new(EcrURI).fromURI("123456789012.dkr.ecr.eu-west-2.amazonaws.com/organisationX/repoX:feature-add-sftp-server-579b8a0-v.2.2.2")
	//assert.Equal(t, x, "")
	in := "123456789012.dkr.ecr.eu-west-2.amazonaws.com/xxxx"
	fmt.Println(regexp.MustCompile(`^[^.]+\.[^.]+\.ecr\.[^.]+\.amazonaws\.com/.+$`).MatchString(in))
}
