package cfg

import (
	"github.com/digtux/laminar/pkg/shared"
	"github.com/gobwas/glob"
	"regexp"
	"strings"
)

// Global settings such as git commit user/email
type Global struct {
	GitUser     string      `yaml:"gitUser"`
	GitEmail    string      `yaml:"gitEmail"`
	GitMessage  interface{} `yaml:"gitMessage"`
	GitHubToken string      `yaml:"gitHubToken"`
}

// Config is the top level of config
type Config struct {
	DockerRegistries []DockerRegistry `yaml:"dockerRegistries"`
	GitRepos         []GitRepo        `yaml:"git"`
	Global           Global           `yaml:"global"`
}

// DockerRegistry contains info about the docker registries
type DockerRegistry struct {
	Reg     string `yaml:"reg"`
	Name    string `yaml:"name"`
	TimeOut int    `yaml:"timeOut,omitempty"`
	Profile string `yaml:"profile"`
}

func (d *DockerRegistry) GetRegistryProvider() shared.RegistryProvider {
	return shared.GetRegistryProvider(d.Reg)
}

func (d *DockerRegistry) GetRegion() string {
	return strings.Split(d.Reg, ".")[3]
}

type BlackList struct {
	Pattern string `yaml:"pattern"`
}

// Files to operate upon in a git repo
type Files struct {
	Path string `yaml:"path"`
}

// Update contains instructions about what to do with matching image
type Updates struct {
	PatternString PatternString `yaml:"pattern"`
	Files         []Files       `yaml:"files"`
	BlackList     []BlackList   `yaml:"blacklist"`
}

type PatternString struct {
	Match func(s string) bool
}

func (p *PatternString) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	s := new(string)
	err = unmarshal(s)
	if err == nil {
		parts := strings.Split(*s, ":")
		patternValue := parts[1]
		switch parts[0] {
		case "glob":
			p.Match = glob.MustCompile(patternValue).Match
		case "regex":
			p.Match = regexp.MustCompile(patternValue).MatchString
		}
	}
	return
}

type RemoteUpdates struct {
	Updates []Updates `yaml:"updates"`
}
