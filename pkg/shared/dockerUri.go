package shared

type DockerURI struct {
	string
	registryProvider *RegistryProvider
}

func (uri *DockerURI) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	s := new(string)
	err = unmarshal(s)
	if err == nil {
		uri.string = *s
	}
	return
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

func (uri *DockerURI) FromString(s string) {
	uri.string = s
}
