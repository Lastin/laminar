package operations

import (
	"bufio"
	"bytes"
	"github.com/labstack/gommon/log"
	"os"
	"path/filepath"
	"strings"

	"github.com/digtux/laminar/pkg/shared"
	"go.uber.org/zap"
)

type Client struct {
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger) *Client {
	return &Client{
		logger: logger,
	}
}

// FindFiles returns a slice containing paths to all files found in a directory
func (c *Client) FindFiles(searchPath string) []string {

	var result []string

	// this function will handle each object inside the Walk()
	var searchFunc = func(pathX string, infoX os.FileInfo, errX error) error {

		// check for errors
		if errX != nil {
			//log.Warnw("findFiles error",
			//	"path", pathX,
			//	"err", errX,
			//)
			return errX
		}

		if isFile, _ := isFile(pathX); isFile {
			c.logger.Debugw("findFiles found file",
				"fileName", infoX.Name(),
			)

			// TODO more expressive way to ignore certain files (in git) that users may want.. eg helm charts
			ext := filepath.Ext(pathX)
			switch ext {
			case ".yml":
				result = append(result, pathX)
			case ".yaml":
				result = append(result, pathX)
			default:
				c.logger.Warnw("file not yaml, ignoring",
					"laminar.path", pathX)
			}
		}

		return nil
	}

	realPath := shared.GetFileAbsPath(searchPath, c.logger)
	err := filepath.Walk(realPath, searchFunc)

	if err != nil {
		c.logger.Debugw("file error",
			"path", searchPath,
			"error", err,
		)
	}

	return result
}

// Search returns a slice of hits that match a string inside a operations
// The assumption is that this is only used against YAML files
// it should work on other types but YMMV
func (c *Client) Search(file string, searchString string) (matches []string) {
	pat := []byte(searchString)
	fp := shared.GetFileAbsPath(file, c.logger)
	f, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// start a scanner to search the operations
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if bytes.Contains(scanner.Bytes(), pat) {
			// if this matches we know the string is somewhere **within a line of text**
			// we should split that line of text (strings.Fields) and range over those to ensure that we
			// don't count the entire line as the actual hit
			// This should be enough for yaml (althoug I imagine it would also detect stuff in comments)
			// but it would be madness for a json operations for example..
			for _, field := range strings.Fields(scanner.Text()) {
				if bytes.Contains([]byte(field), pat) {
					matches = append(matches, field)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error(err)
	}
	return matches
}

// isFile will return true if the path is a normal file (not directory or link)
func isFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode().IsRegular(), nil
}
