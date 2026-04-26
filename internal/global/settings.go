package global

type Settings struct {
	DefaultType string `yaml:"default_type,omitempty"`
}

type Services struct {
	Services []ServiceEntry `yaml:"services,omitempty"`
}

type ServiceEntry struct {
	Image   string `yaml:"image"`
	Tag     string `yaml:"tag"`
	Name    string `yaml:"name"`
	Shared  bool   `yaml:"shared"`
}
