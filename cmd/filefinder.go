package cmd

import (
	"fmt"
	"github.com/digtux/laminar/pkg/cfg"
	"github.com/digtux/laminar/pkg/common"
	"github.com/digtux/laminar/pkg/gitoperations"
	"github.com/digtux/laminar/pkg/operations"
	"go.uber.org/zap"
)

// FileFinder returns a list of files found in the gitoperations repo
// NOTE: updatePolicies with references to the specific gitRepo are required
// NOTE2: if it looks like a folder all files found underneath will be included
// TODO: how does this play with .gitoperations files? we probably want to filter+exclude them here
// TODO: might want to exclude tarballs and other kind of things also I guess..
func FileFinder(gitRepo cfg.GitRepo, log *zap.SugaredLogger) []string {

	// empty the list of files found
	fileList = nil

	var thisReposPaths []string

	// now loop though the UpdatePolicies and gather their files[].path values
	// only IF the gitRepo specifies the name of our gitRepo
	for _, update := range gitRepo.Updates {
		for _, p := range update.Files {
			log.Debugw("FileFinder searching for files",
				"path", p.Path,
				"gitRepo", gitRepo.URL,
				"branch", gitRepo.Branch,
			)

			thisReposPaths = append(thisReposPaths, p.Path)
		}
	}
	thisReposPaths = common.UniqueStrings(thisReposPaths)

	// get ready to add the discovered files to the slice
	for _, p := range thisReposPaths {
		// get the path of where the gitoperations repo is checked out
		relativeGitPath := gitoperations.GetRepoPath(gitRepo)
		// combine these
		realPath := fmt.Sprintf("%s/%s", relativeGitPath, p)

		// finally this will return all files found
		for _, x := range operations.FindFiles(realPath, log) {
			fileList = append(fileList, x)
		}

	}
	log.Debugw("successfully found files in gitoperations",
		"count", len(fileList),
	)
	return fileList
}
