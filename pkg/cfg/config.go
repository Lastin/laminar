package cfg

import (
	"github.com/pkg/errors"
	"reflect"

	"gopkg.in/yaml.v1"
)

// ParseConfig will read a config and infer some defaults if they're omitted (one day)
func ParseConfig(data []byte) (Config, error) {

	var yamlConfig Config
	var empty Config

	_ = yaml.Unmarshal(data, &yamlConfig)

	// lets return an error if an yaml.Unmarshal returned no new data
	if reflect.DeepEqual(yamlConfig, empty) {
		err := errors.New("no data was loaded")
		return yamlConfig, err
	}
	return yamlConfig, nil
}

// ParseUpdates will read the .laminar.yaml from a repo and return its RemoteUpdates
func ParseUpdates(data []byte) (*RemoteUpdates, error) {
	yamlUpdates := new(RemoteUpdates)
	err := yaml.Unmarshal(data, &yamlUpdates)
	if err != nil {
		return nil, errors.Wrap(err, "ParseUpdates failed")
	}
	return yamlUpdates, err
}
