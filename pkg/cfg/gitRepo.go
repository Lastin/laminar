package cfg

import (
	"fmt"
	"github.com/digtux/laminar/pkg/shared"
	"github.com/digtux/laminar/pkg/yamlOps"
	"regexp"
)

// GitRepo which laminar operates on
type GitRepo struct {
	URL               string    `yaml:"url"`
	Branch            string    `yaml:"branch"`
	Username          string    `yaml:"username"`
	Token             string    `yaml:"token"`
	Key               string    `yaml:"key"`
	PollFreq          int       `yaml:"pollFreq"`
	Name              string    `yaml:"name"`
	RemoteConfig      bool      `yaml:"remoteConfig"` // propogate []Updates from remote git ".laminar.yaml" ?
	Updates           []Updates `yaml:"updates,omitempty"`
	PreCommitCommands []string  `yaml:"preCommitCommands,omitempty"`
	//PostChange   []PostChanges `yaml:"postChange"`
}

func (repo *GitRepo) GetTotalPathsSize() (total int) {
	for _, update := range repo.Updates {
		total += len(update.Files)
	}
	return
}

// GetAllFilePaths Returns all file paths combined, and unique
func (repo *GitRepo) GetAllFilePaths() []string {
	repoPaths := make([]string, repo.GetTotalPathsSize())
	i := 0
	// now loop though the UpdatePolicies and gather their files[].path values
	for _, update := range repo.Updates {
		for _, p := range update.Files {
			repoPaths[i] = p.Path
			i++
		}
	}
	return shared.UniqueStrings(repoPaths)
}

func (repo *GitRepo) GetRealPath() string {
	r := regexp.MustCompile("[/:]")
	return fmt.Sprintf(
		"/tmp/laminar/%s-%s",
		r.ReplaceAllString(repo.URL, "-"),
		r.ReplaceAllString(repo.Branch, "-"),
	)
}

// GetFileList returns list of files matching specification, which might contain docker image URIs
func (repo *GitRepo) GetFileList() (totalResult []string, err error) {
	for _, path := range repo.GetAllFilePaths() {
		var result []string
		result, err = yamlOps.FindYamlsInPath(fmt.Sprintf("%s/%s", repo.GetRealPath(), path))
		if err != nil {
			return
		}
		totalResult = append(totalResult, result...)
	}
	return
}

func (repo *GitRepo) SetUpdates(remoteUpdates *RemoteUpdates) {
	repo.Updates = make([]Updates, len(remoteUpdates.Updates))
	// now assemble that list for this run
	for _, update := range remoteUpdates.Updates {
		repo.Updates = append(repo.Updates, update)
	}
}
