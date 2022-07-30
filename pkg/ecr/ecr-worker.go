package ecr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/digtux/laminar/pkg/cfg"
	"github.com/digtux/laminar/pkg/shared"
	"go.uber.org/zap"
	"sort"
	"strings"
)

type Worker struct {
	logger        *zap.SugaredLogger
	svcByRegistry map[string]EcrIFace
}

type EcrIFace interface {
	DescribeImagesPages(input *ecr.DescribeImagesInput, fn func(*ecr.DescribeImagesOutput, bool) bool) error
}

func New(logger *zap.SugaredLogger, regs []cfg.DockerRegistry) *Worker {
	return &Worker{
		logger:        logger,
		svcByRegistry: getSvcByRegistry(regs),
	}
}

func getSvcByRegistry(observedRegistries []cfg.DockerRegistry) map[string]EcrIFace {
	mySession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: "tooling",
	}))
	result := map[string]EcrIFace{}
	for _, reg := range observedRegistries {
		if result[reg.Reg] == nil {
			result[reg.Reg] = ecr.New(mySession, aws.NewConfig().WithRegion(reg.GetRegion()))
		}
	}
	return result
}

func (c *Worker) ScanAll(imageList []shared.DockerURI) (tagInfosByRepo map[string][]*shared.TagInfo, err error) {
	var uris []EcrURI
	uris, err = castAll(imageList)
	if err != nil {
		return nil, err
	}
	tagInfosByRepo = map[string][]*shared.TagInfo{}
	for _, image := range uris {
		// only scan if we're observing that registry
		if svc := c.svcByRegistry[image.Registry]; svc != nil {
			// only scan if we haven't scanned that repository yet
			if tagInfosByRepo[image.Repository] == nil {
				tagInfosByRepo[image.Repository], err = c.describeRepository(image, svc)
				c.sortByDateDesc(tagInfosByRepo[image.Repository])
			}
		} else {
			c.logger.Warnw("no service found for registry",
				"reg", image.Registry,
			)
		}
	}
	return
}

func castAll(imageList []shared.DockerURI) (result []EcrURI, err error) {
	result = make([]EcrURI, len(imageList))
	for i, image := range imageList {
		var ecrUri *EcrURI
		ecrUri, err = new(EcrURI).fromURI(image)
		if err != nil {
			return nil, err
		}
		result[i] = *ecrUri
	}
	return result, nil
}

func (c *Worker) sortByDateDesc(results []*shared.TagInfo) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].IsNewerThan(results[j])
	})
}

func (c *Worker) describeRepository(uri EcrURI, svc EcrIFace) (results []*shared.TagInfo, err error) {
	describeImageSettings := &ecr.DescribeImagesInput{
		RepositoryName: &uri.Repository,
	}
	err = svc.DescribeImagesPages(describeImageSettings, func(page *ecr.DescribeImagesOutput, lastPage bool) bool {
		// structure looks something like:
		// pages = [{tags: [tag1, tag2, tag3]}]
		// we merge it all into one big slice
		for _, details := range page.ImageDetails {
			results = append(results, func() []*shared.TagInfo {
				r1 := make([]*shared.TagInfo, len(details.ImageTags))
				for i, tag := range details.ImageTags {
					r1[i] = &shared.TagInfo{
						Image:   *details.RepositoryName,
						Hash:    strings.Split(*details.ImageDigest, ":")[1],
						Tag:     *tag,
						Created: *details.ImagePushedAt,
					}
				}
				return r1
			}()...)
		}
		return true
	})
	return
}
