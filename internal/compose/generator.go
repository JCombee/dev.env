package compose

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jcombee/devenv/internal/config"
	"github.com/jcombee/devenv/internal/global"
	"github.com/jcombee/devenv/internal/ports"
	"github.com/jcombee/devenv/internal/services"
	"gopkg.in/yaml.v3"
)

type composeFile struct {
	Services map[string]composeService `yaml:"services"`
	// Volumes declares named volumes; only global apps use them.
	Volumes map[string]struct{} `yaml:"volumes,omitempty"`
}

type composeService struct {
	Image       string                    `yaml:"image"`
	Ports       []string                  `yaml:"ports,omitempty"`
	Environment map[string]string         `yaml:"environment,omitempty"`
	Restart     string                    `yaml:"restart,omitempty"`
	Volumes     []string                  `yaml:"volumes,omitempty"`
	Command     []string                  `yaml:"command,omitempty"`
	DependsOn   map[string]dependsOnEntry `yaml:"depends_on,omitempty"`
	Healthcheck *composeHealthcheck       `yaml:"healthcheck,omitempty"`
}

type dependsOnEntry struct {
	Condition string `yaml:"condition"`
}

type composeHealthcheck struct {
	Test     []string `yaml:"test"`
	Interval string   `yaml:"interval,omitempty"`
	Timeout  string   `yaml:"timeout,omitempty"`
	Retries  int      `yaml:"retries,omitempty"`
}

// WriteShared regenerates ~/.dev.env/docker/docker-compose.yml from the
// full services.yaml map and current docker secrets.
// globalPorts overrides computed host ports by docker-compose service name; nil means use computed defaults.
func WriteShared(svcMap map[string]global.GlobalServiceEntry, ds global.DockerSecrets, globalPorts map[string]int) error {
	data, err := renderShared(svcMap, ds, globalPorts)
	if err != nil {
		return err
	}
	path := filepath.Join(global.Dir(), "docker", "docker-compose.yml")
	return os.WriteFile(path, data, 0o644)
}

// SharedPath returns the path to the shared docker-compose.yml.
func SharedPath() string {
	return filepath.Join(global.Dir(), "docker", "docker-compose.yml")
}

// DedicatedPath returns the path to a project's dedicated docker-compose.yml.
func DedicatedPath(projectName string) string {
	return filepath.Join(global.Dir(), "projects", projectName, "docker-compose.yml")
}

// WriteDedicated generates a docker-compose.yml for the dedicated services of
// one project and writes it to ~/.dev.env/projects/<name>/docker-compose.yml.
func WriteDedicated(projectName string, entries []config.ServiceEntry, ds global.DockerSecrets) error {
	data, err := renderDedicated(projectName, entries, ds)
	if err != nil {
		return err
	}
	path := DedicatedPath(projectName)
	return os.WriteFile(path, data, 0o644)
}

func renderShared(svcMap map[string]global.GlobalServiceEntry, ds global.DockerSecrets, globalPorts map[string]int) ([]byte, error) {
	cf := composeFile{Services: map[string]composeService{}}
	seen := map[int]string{}
	for composeName, entry := range svcMap {
		r, err := ports.Resolve(entry.Image, entry.Tag, globalPorts, nil)
		if err != nil {
			return nil, fmt.Errorf("resolve port for %s: %w", composeName, err)
		}
		if r.Port != 0 {
			if conflict, ok := seen[r.Port]; ok {
				return nil, fmt.Errorf("port conflict: %s and %s both resolve to host port %d", conflict, composeName, r.Port)
			}
			seen[r.Port] = composeName
		}
		cf.Services[composeName] = buildService(entry.Image, entry.Tag, composeName, r.Port, ds)
	}
	return yaml.Marshal(cf)
}

func renderDedicated(projectName string, entries []config.ServiceEntry, ds global.DockerSecrets) ([]byte, error) {
	cf := composeFile{Services: map[string]composeService{}}
	for _, entry := range entries {
		if !entry.Dedicated {
			continue
		}
		tag := entry.Tag
		if tag == "" {
			tag = "latest"
		}
		// Dedicated service name is project-prefixed to avoid cross-file conflicts.
		composeName := projectName + "-" + ServiceName(entry.Image, tag)
		r, err := ports.Resolve(entry.Image, tag, nil, entry.Port)
		if err != nil {
			return nil, fmt.Errorf("resolve port for %s: %w", composeName, err)
		}
		cf.Services[composeName] = buildService(entry.Image, tag, composeName, r.Port, ds)
	}
	return yaml.Marshal(cf)
}

func buildService(image, tag, composeName string, hostPort int, ds global.DockerSecrets) composeService {
	if tag == "" {
		tag = "latest"
	}
	svc, _ := services.Find(image)
	cs := composeService{
		Image:   image + ":" + tag,
		Restart: "unless-stopped",
	}
	if hostPort != 0 && svc.Port != 0 {
		cs.Ports = []string{fmt.Sprintf("%d:%d", hostPort, svc.Port)}
	}
	env := ContainerEnv(image, composeName, ds)
	if len(env) > 0 {
		cs.Environment = env
	}
	return cs
}
