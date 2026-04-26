package global

import "gopkg.in/yaml.v3"

type Settings struct {
	DefaultType string `yaml:"default_type"`
}

// GlobalServiceEntry is one entry in services.yaml.
// The map key is the docker-compose service name (e.g. "mysql-8-0").
type GlobalServiceEntry struct {
	Image string `yaml:"image"`
	Tag   string `yaml:"tag"`
}

type Services struct {
	Services map[string]GlobalServiceEntry `yaml:"services"`
}

// UnmarshalYAML handles both the old list format (services: []) written by
// early builds and the current map format (services: {}).
func (s *Services) UnmarshalYAML(value *yaml.Node) error {
	// Walk the mapping node manually to find the "services" key.
	for i := 0; i+1 < len(value.Content); i += 2 {
		if value.Content[i].Value != "services" {
			continue
		}
		child := value.Content[i+1]
		if child.Kind == yaml.SequenceNode {
			// Old format: empty list — initialise as empty map and return.
			s.Services = map[string]GlobalServiceEntry{}
			return nil
		}
	}
	type plain Services
	return value.Decode((*plain)(s))
}
