package registryCore

import (
	"github.com/digtux/laminar/pkg/cfg"
	"github.com/digtux/laminar/pkg/registryEcr"
	"github.com/digtux/laminar/pkg/registryShared"
	"github.com/digtux/laminar/pkg/shared"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

type Client struct {
	db              *buntdb.DB
	logger          *zap.SugaredLogger
	registries      []cfg.DockerRegistry
	providerClients map[shared.RegistryProvider]registryShared.RegistryIFace
}

func New(logger *zap.SugaredLogger, regs []cfg.DockerRegistry, db *buntdb.DB) *Client {
	c := &Client{
		logger:     logger,
		db:         db,
		registries: regs,
	}
	c.providerClients = c.initialiseClients(regs)
	return c
}

func (c *Client) initialiseClients(regs []cfg.DockerRegistry) map[shared.RegistryProvider]registryShared.RegistryIFace {
	result := map[shared.RegistryProvider]registryShared.RegistryIFace{}
	for _, reg := range regs {
		if reg.GetRegistryProvider() == shared.ECR {
			// initialise once
			if result[shared.ECR] == nil {
				result[shared.ECR] = registryEcr.New(c.logger, c.registries)
			}
		}
	}
	return result
}

func (c *Client) ScanAll(images []shared.DockerURI) (err error) {
	grouped := groupByKind(images)
	for regProvider, group := range grouped {
		if client := c.providerClients[regProvider]; client != nil {
			_, err = client.ScanAll(group)
			if err != nil {
				return
			}
		} else {
			c.logger.Warnw("not client for registry provider",
				"registry_provider", regProvider,
			)
		}
	}
	return
}

func (c *Client) writeAllToCache(tagInfosByRepo map[string][]*shared.TagInfo) (err error) {
	for _, tagInfoList := range tagInfosByRepo {
		for _, tagInfo := range tagInfoList {
			err = tagInfo.WriteToCache(c.db)
			if err != nil {
				return
			}
		}
	}
	return
}

func groupByKind(images []shared.DockerURI) map[shared.RegistryProvider][]shared.DockerURI {
	result := map[shared.RegistryProvider][]shared.DockerURI{}
	for _, image := range images {
		result[image.GetRegistryProvider()] = append(result[image.GetRegistryProvider()], image)
	}
	return result
}
