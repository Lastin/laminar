package fileOps

import (
	"go.uber.org/zap"
	"io/ioutil"
)

func ReadFile(filePath string, log *zap.SugaredLogger) ([]byte, string) {
	r, err := ioutil.ReadFile(filePath)
	stringContents := string(r)
	if err != nil {
		log.Fatal(err)
	}
	return r, stringContents
}
