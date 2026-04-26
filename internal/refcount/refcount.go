package refcount

import (
	"os"
	"path/filepath"

	"github.com/jcombee/devenv/internal/global"
	"github.com/jcombee/devenv/internal/project"
)

// ActiveProjectsUsingService returns the names of all projects that list
// composeName as an active service in their state.yaml.
func ActiveProjectsUsingService(composeName string) ([]string, error) {
	entries, err := projectEntries()
	if err != nil {
		return nil, err
	}

	var names []string
	for _, name := range entries {
		state, err := project.LoadState(name)
		if err != nil {
			continue
		}
		if _, ok := state.Services[composeName]; ok && state.Running {
			names = append(names, name)
		}
	}
	return names, nil
}

// AllServiceUsage returns a map of composeName → []projectName for all running projects.
func AllServiceUsage() (map[string][]string, error) {
	entries, err := projectEntries()
	if err != nil {
		return nil, err
	}

	usage := map[string][]string{}
	for _, name := range entries {
		state, err := project.LoadState(name)
		if err != nil {
			continue
		}
		if !state.Running {
			continue
		}
		for composeName := range state.Services {
			usage[composeName] = append(usage[composeName], name)
		}
	}
	return usage, nil
}

func projectEntries() ([]string, error) {
	projectsDir := filepath.Join(global.Dir(), "projects")
	entries, err := os.ReadDir(projectsDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(project.MetaPath(e.Name())); err == nil {
			names = append(names, e.Name())
		}
	}
	return names, nil
}
