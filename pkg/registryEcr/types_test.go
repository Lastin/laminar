package registryEcr

import (
	"github.com/digtux/laminar/pkg/shared"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEcrURI_fromURI(t *testing.T) {
	x := new(shared.DockerURI)
	x.FromString("123456789012.dkr.ecr.eu-west-2.amazonaws.com/organisationX/repoX:feature-add-sftp-server-579b8a0-v.2.2.2")
	result, _ := new(EcrURI).fromURI(x)
	assert.Equal(t, "", result)
}
