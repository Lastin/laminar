package ecr

import (
	"errors"
	"fmt"
	"github.com/digtux/laminar/pkg/shared"
	"regexp"
)

type EcrURI struct {
	Registry   string
	Region     string
	Repository string
	Tag        string
}

func (eap *EcrURI) fromURI(uri shared.DockerURI) (*EcrURI, error) {
	r := regexp.MustCompile(
		`(?P<registry>` + //registry start
			`(?P<account>[^.]+)\.` +
			`(?P<sub2>[^.]+)\.` +
			`(?P<sub1>[^.]+)\.` +
			`(?P<region>[^.]+)\.` +
			`(amazonaws\.com)/` +
			`(?P<org>[^/]+)` +
			`)/` + //registry end
			`(?P<service>[^:]+):` +
			`(?P<tag>.+)`,
	)
	groups := r.FindStringSubmatch(uri.String())
	if groups == nil {
		return nil, errors.New("could not parse uri" + uri.String())
	}
	//// EG "112233445566.dkr.ecr.eu-west-2.amazonaws.com/acmecorp/app-name"
	eap.Registry = groups[r.SubexpIndex("registry")]
	eap.Region = groups[r.SubexpIndex("region")]
	eap.Repository = fmt.Sprintf("%s/%s",
		groups[r.SubexpIndex("org")],
		groups[r.SubexpIndex("service")],
	)
	eap.Tag = groups[r.SubexpIndex("tag")]
	return eap, nil
}
