package global

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
