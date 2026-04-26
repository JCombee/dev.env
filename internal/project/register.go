package project

import (
	"fmt"
	"os"

	"github.com/jcombee/devenv/internal/store"
)

// ErrNameCollision is returned when the project name is already used by a different path.
type ErrNameCollision struct {
	Name         string
	ExistingPath string
}

func (e *ErrNameCollision) Error() string {
	return fmt.Sprintf("project name %q already in use by: %s", e.Name, e.ExistingPath)
}

// Register ensures ~/.dev.env/projects/<name>/ exists and is associated with projectPath.
// Returns ErrNameCollision if the name is already registered to a different path.
func Register(name, projectPath string) error {
	dir := Dir(name)
	metaPath := MetaPath(name)

	if _, err := os.Stat(metaPath); err == nil {
		// Already registered — check for collision.
		var existing Meta
		if err := store.Read(metaPath, &existing); err != nil {
			return fmt.Errorf("read project.yaml: %w", err)
		}
		if existing.Path != projectPath {
			return &ErrNameCollision{Name: name, ExistingPath: existing.Path}
		}
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create project dir: %w", err)
	}

	meta := Meta{Name: name, Path: projectPath}
	if err := store.Write(metaPath, &meta); err != nil {
		return fmt.Errorf("write project.yaml: %w", err)
	}

	if err := store.Write(StatePath(name), &State{
		Services:     map[string]string{},
		Acknowledged: []string{},
	}); err != nil {
		return fmt.Errorf("write state.yaml: %w", err)
	}
	if err := store.Write(SecretsPath(name), &Secrets{}); err != nil {
		return fmt.Errorf("write secrets.yaml: %w", err)
	}

	return nil
}

// LoadState reads state.yaml for the named project.
func LoadState(name string) (*State, error) {
	var s State
	if err := store.Read(StatePath(name), &s); err != nil {
		return nil, err
	}
	if s.Services == nil {
		s.Services = map[string]string{}
	}
	return &s, nil
}

// SaveState writes state.yaml for the named project.
func SaveState(name string, s *State) error {
	return store.Write(StatePath(name), s)
}

// LoadSecrets reads secrets.yaml for the named project.
func LoadSecrets(name string) (Secrets, error) {
	var sec Secrets
	if err := store.Read(SecretsPath(name), &sec); err != nil {
		return nil, err
	}
	if sec == nil {
		sec = Secrets{}
	}
	return sec, nil
}

// SaveSecrets writes secrets.yaml for the named project.
func SaveSecrets(name string, sec Secrets) error {
	return store.Write(SecretsPath(name), sec)
}
