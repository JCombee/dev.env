package global

import (
	"os"
	"path/filepath"

	"github.com/jcombee/devenv/internal/store"
)

func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".dev.env")
}

func Setup() error {
	base := Dir()
	dirs := []string{
		base,
		filepath.Join(base, "docker"),
		filepath.Join(base, "projects"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}

	settingsPath := filepath.Join(base, "settings.yaml")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		if err := store.Write(settingsPath, &Settings{}); err != nil {
			return err
		}
	}

	servicesPath := filepath.Join(base, "services.yaml")
	if _, err := os.Stat(servicesPath); os.IsNotExist(err) {
		if err := store.Write(servicesPath, &Services{Services: map[string]GlobalServiceEntry{}}); err != nil {
			return err
		}
	}

	return nil
}

func IsSetUp() bool {
	_, err := os.Stat(Dir())
	return err == nil
}
