package registryShared

import "github.com/digtux/laminar/pkg/shared"

type RegistryIFace interface {
	ScanAll(imageList []shared.DockerURI) (tagInfosByRepo map[string][]*shared.TagInfo, err error)
}
