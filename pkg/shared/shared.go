package shared

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Misc helper functions

// GetFileAbsPath will expand on something such as ~/.ssh/my_id_rsa and return a string like /home/user/.ssh/my_id_rsa
func GetFileAbsPath(fileName string, log *zap.SugaredLogger) (result string) {

	if strings.HasPrefix(fileName, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		fileName = filepath.Join(dir, fileName[2:])
	}

	result, err := filepath.Abs(fileName)
	if err != nil {
		log.Fatalw("unable to determine path to a operations",
			"fileName", fileName,
			"error", err,
		)
	}

	return result
}

// IsDir will return true if the path is a directory
func IsDir(path string, log *zap.SugaredLogger) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Debugw("couldn't reading path",
			"path", path,
			"error", err,
		)
		return false
	}
	return fileInfo.IsDir()
}

// UniqueStrings takes an array of strings in, returns only the unique ones
func UniqueStrings(input []string) []string {
	sort.Strings(input)
	size := len(input)
	u := make([]string, 0, size)
	for i := 0; i < size; {
		current := input[i]
		u = append(u, current)
		for i < size && input[i] == current {
			i++
		}
	}
	return u
}

// GetLogger provides us with sugared logger
// switch between a vanilla Development or Production logging format (--debug)
// The only change from vanilla zap is the ProductionConfig outputs to stdout instead of stderr
func GetLogger(debug bool) (zapLog *zap.SugaredLogger) {
	// https://blog.sandipb.net/2018/05/02/using-zap-simple-use-cases/
	if debug {
		zapLogger, err := zap.NewDevelopment()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		sugar := zapLogger.Sugar()
		sugar.Debug("debug enabled")
		return sugar
	} else {
		// Override the default zap production Config a little
		// NewProductionConfig is json

		logConfig := zap.NewProductionConfig()
		// customise the "time" field to be ISO8601
		logConfig.EncoderConfig.TimeKey = "time"
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		// main message data into the key "msg"
		logConfig.EncoderConfig.MessageKey = "msg"

		// stdout+sterr into stdout
		logConfig.OutputPaths = []string{"stdout"}
		logConfig.ErrorOutputPaths = []string{"stdout"}
		zapLogger, err := logConfig.Build()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return zapLogger.Sugar()
	}
}

type RegistryProvider string

const (
	ECR     RegistryProvider = "ECR"       // Elastic Container Registry
	GAR                      = "GAR"       // Google Artifact Registry
	GCR                      = "GCR"       // Google Container Registry
	ACR                      = "ACR"       // Azure Container Registry
	DKH                      = "DockerHub" // Dockerhub Registry
	JFROG                    = "JFROG"     // JFrog Registry
	QUAY                     = "QUAYIO"    // Quay.io Container Registry
	UNKNOWN                  = "UNKNOWN"
)

var ecrRegex = regexp.MustCompile(`^[^.]+\.[^.]+\.ecr\.[^.]+\.amazonaws\.com/.+$`)

func GetRegistryProvider(s string) RegistryProvider {
	if ecrRegex.MatchString(s) {
		return ECR
	}
	return UNKNOWN
}

type DockerURI struct {
	string
	registryProvider *RegistryProvider
}

func (uri *DockerURI) GetRegistryProvider() RegistryProvider {
	if uri.registryProvider == nil {
		regProvider := GetRegistryProvider(uri.string)
		uri.registryProvider = &regProvider
	}
	return *uri.registryProvider
}

func (uri *DockerURI) String() string {
	return uri.string
}
