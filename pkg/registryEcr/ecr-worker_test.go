package registryEcr

import (
	"github.com/digtux/laminar/pkg/cfg"
	"github.com/digtux/laminar/pkg/shared"
	"testing"
)

func TestClient_EcrWorker(t *testing.T) {
	logger := shared.GetLogger(true)
	w := New(logger, []cfg.DockerRegistry{
		{
			Reg:     "123456789012.dkr.ecr.eu-west-2.amazonaws.com/foo",
			Name:    "foo tooling",
			TimeOut: 0,
		},
	})
	w.ScanAll([]shared.DockerURI{
		{},
	})
}
