package apps

import (
	"path/filepath"

	"github.com/jcombee/devenv/internal/global"
)

// Dir is the directory holding the apps compose file.
func Dir() string {
	return filepath.Join(global.Dir(), "apps")
}

// StatePath is the path of apps.yaml.
func StatePath() string {
	return filepath.Join(global.Dir(), "apps.yaml")
}

// ComposePath is the path of the generated apps docker-compose.yml.
func ComposePath() string {
	return filepath.Join(Dir(), "docker-compose.yml")
}
