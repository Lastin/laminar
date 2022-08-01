package yamlOps

import (
	"github.com/digtux/laminar/pkg/shared"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path/filepath"
)

func IsYaml(path string, info os.FileInfo) bool {
	if info.Mode().IsRegular() {
		switch filepath.Ext(path) {
		case ".yml":
			fallthrough
		case ".yaml":
			return true
		}
	}
	return false
}

// FindYamlsInPath returns a slice containing paths to all files found in a directory
func FindYamlsInPath(searchPath string) (result []string, err error) {
	// this function will handle each object inside the Walk()
	realPath, err := shared.GetFileAbsPath(searchPath)
	if err == nil {
		err = filepath.Walk(realPath, func(path string, info fs.FileInfo, err error) error {
			if IsYaml(path, info) {
				result = append(result, path)
			}
			return err
		})
	}
	if err != nil {
		err = errors.Wrap(err, "findAllFiles error: "+searchPath)
	}
	return
}
