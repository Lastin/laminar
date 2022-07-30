package registry

import (
	"encoding/json"
	"github.com/digtux/laminar/pkg/cfg"
	"github.com/digtux/laminar/pkg/ecr"
	"github.com/digtux/laminar/pkg/shared"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"strings"
)

type Client struct {
	db     *buntdb.DB
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger, db *buntdb.DB) *Client {
	return &Client{
		logger: logger,
		db:     db,
	}
}

// Exec will check if we support that docker reg and then launch an appropriate worker
func (c *Client) Exec(registry cfg.DockerRegistry, imageList []string) {

	// grok will add some defaults lest the config doesn't include em
	registry = grokRegistrySettings(registry)
	c.logger.Debugw("DockerRegistry worker launching",
		"Registry", registry,
	)

	for _, image := range imageList {

		// Check if the image looks like an ECR image
		if strings.Contains(registry.Reg, "ecr") {
			ecr.EcrWorker(c.db, registry, imageList, c.logger)
			return
		}

		if strings.Contains(registry.Reg, "gcr") {
			GcrWorker(c.db, registry, imageList, c.logger)
			return
		}
	}

	c.logger.Fatal("unable to figure out which kind of registry you have")
}

// incase some fields are missing, lets set their defaults
func grokRegistrySettings(in cfg.DockerRegistry) cfg.DockerRegistry {

	// incase this is empty set it..
	if in.TimeOut == 0 {
		in.TimeOut = 30
	}
	return in
}

func (c *Client) CachedImagesToTagInfoListSpecificImage(
	imageString string,
	index string,
) (result []shared.TagInfo) {
	c.db.View(func(tx *buntdb.Tx) error {
		tx.Descend(index, func(key, val string) bool {
			// decode the data from the db
			x := JsonStringToTagInfo(val, c.logger)

			// if this image matches the imageString append it to the result
			if x.Image == imageString {
				result = append(result, x)
			}
			return true
		})
		return nil
	})
	//log.Debugw("searched DB for images",
	//	"image", imageString,
	//	"hits", len(result),
	//)
	return result
}

func JsonStringToTagInfo(s string, log *zap.SugaredLogger) shared.TagInfo {
	var data shared.TagInfo
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		log.Error("unmarshal error?")
		log.Fatal(err)
		return data
	}
	return data
}
