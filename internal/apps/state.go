package apps

import (
	"os"

	"github.com/jcombee/devenv/internal/store"
)

// Entry is one app's record in apps.yaml.
type Entry struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port,omitempty"`
}

// State is the content of apps.yaml.
type State struct {
	Apps map[string]Entry `yaml:"apps"`
}

// Load reads apps.yaml. A missing or empty file yields an empty state.
func Load() (*State, error) {
	var st State
	err := store.Read(StatePath(), &st)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if st.Apps == nil {
		st.Apps = map[string]Entry{}
	}
	return &st, nil
}

// Save writes apps.yaml.
func (s *State) Save() error {
	return store.Write(StatePath(), s)
}

// IsEnabled reports whether the named app is enabled.
func (s *State) IsEnabled(name string) bool {
	return s.Apps[name].Enabled
}

// HostPort returns the app's published host port: the apps.yaml override when
// set, otherwise the app's default.
func (s *State) HostPort(a App) int {
	if p := s.Apps[a.Name].Port; p != 0 {
		return p
	}
	return a.DefaultPort
}

// Enabled returns the registered apps currently enabled, in registry order.
// Entries naming an app that is not registered are ignored.
func (s *State) Enabled() []App {
	var enabled []App
	for _, a := range All {
		if s.IsEnabled(a.Name) {
			enabled = append(enabled, a)
		}
	}
	return enabled
}

// Ports returns the host port of every enabled app, keyed by app name.
func (s *State) Ports() map[string]int {
	ports := map[string]int{}
	for _, a := range All {
		ports[a.Name] = s.HostPort(a)
	}
	return ports
}
