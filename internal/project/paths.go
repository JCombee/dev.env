package project

import (
	"path/filepath"

	"github.com/jcombee/devenv/internal/global"
)

// Dir returns ~/.dev.env/projects/<name>.
func Dir(name string) string {
	return filepath.Join(global.Dir(), "projects", name)
}

// MetaPath returns the path to project.yaml for the named project.
func MetaPath(name string) string {
	return filepath.Join(Dir(name), "project.yaml")
}

// StatePath returns the path to state.yaml for the named project.
func StatePath(name string) string {
	return filepath.Join(Dir(name), "state.yaml")
}

// SecretsPath returns the path to secrets.yaml for the named project.
func SecretsPath(name string) string {
	return filepath.Join(Dir(name), "secrets.yaml")
}
