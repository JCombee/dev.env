package project

// Meta is the content of project.yaml — disk location of the project.
type Meta struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// State is the content of state.yaml — runtime status of a project.
type State struct {
	Running      bool              `yaml:"running"`
	Services     map[string]string `yaml:"services"`    // compose-name → "shared"|"dedicated"
	Acknowledged []string          `yaml:"acknowledged"` // unsupported images the user has accepted
}

// Secrets is the content of secrets.yaml — generated credentials, never overwritten.
// Outer key: image name (e.g. "mysql"). Inner map: credential fields.
type Secrets map[string]map[string]string
