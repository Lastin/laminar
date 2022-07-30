package fileOps

import (
	"fmt"
	"github.com/digtux/laminar/pkg/cfg"
	"github.com/digtux/laminar/pkg/shared"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/fs"
	"os"
	"path/filepath"
)

type Client struct {
	logger *zap.SugaredLogger
}

// GetFileList returns list of files matching specification, which might contain docker image URIs
func (c *Client) GetFileList(gitRepo cfg.GitRepo) (totalResult []string, err error) {
	c.logger.Info("updating Daemon's file list")
	relativeGitPath := gitRepo.GetRealPath()
	for _, path := range gitRepo.GetAllFilePaths() {
		var result []string
		result, err = c.findAllFiles(fmt.Sprintf("%s/%s", relativeGitPath, path))
		if err != nil {
			return
		}
		totalResult = append(totalResult, result...)
	}
	c.logger.Debugw("successfully found files in gitoperations",
		"count", len(totalResult),
	)
	return
}

// findAllFiles returns a slice containing paths to all files found in a directory
func (c *Client) findAllFiles(searchPath string) (result []string, err error) {
	// this function will handle each object inside the Walk()
	realPath := shared.GetFileAbsPath(searchPath, c.logger)
	err = filepath.Walk(realPath, func(path string, info fs.FileInfo, err error) error {
		if p := c.searchYamls(path, info); p != nil {
			result = append(result, *p)
		}
		return nil
	})
	if err != nil {
		err = errors.Wrap(err, "findAllFiles error: "+searchPath)
	}
	return
}

func (c *Client) searchYamls(path string, info os.FileInfo) *string {
	if info.Mode().IsRegular() {
		// TODO more expressive way to ignore certain files (in git) that users may want.. eg helm charts
		switch filepath.Ext(path) {
		case ".yml":
			return &path
		case ".yaml":
			return &path
		default:
			c.logger.Warnw("file not yaml, ignoring",
				"laminar.path", path)
		}
	}
	return nil
}
