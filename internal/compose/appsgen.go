package compose

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jcombee/devenv/internal/apps"
	"github.com/jcombee/devenv/internal/global"
	"gopkg.in/yaml.v3"
)

// AppsPath returns the path to the global apps docker-compose.yml. It is
// deliberately separate from SharedPath so project commands can never start or
// stop an app.
func AppsPath() string {
	return filepath.Join(global.Dir(), "apps", "docker-compose.yml")
}

// WriteApps regenerates the apps compose file from the enabled apps.
// hostPorts overrides an app's default host port, keyed by app name.
func WriteApps(enabled []apps.App, hostPorts map[string]int, ds global.DockerSecrets) error {
	data, err := renderApps(enabled, hostPorts, ds)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(AppsPath()), 0o755); err != nil {
		return err
	}
	return os.WriteFile(AppsPath(), data, 0o644)
}

// AppContainerNames returns the compose service names of one app.
func AppContainerNames(a apps.App, hostPort int, ds global.DockerSecrets) []string {
	containers := a.Build(apps.BuildContext{HostPort: hostPort, Secrets: ds})
	names := make([]string, 0, len(containers))
	for _, c := range containers {
		names = append(names, c.Name)
	}
	return names
}

func renderApps(enabled []apps.App, hostPorts map[string]int, ds global.DockerSecrets) ([]byte, error) {
	cf := composeFile{Services: map[string]composeService{}}
	seen := map[int]string{}

	for _, a := range enabled {
		hostPort := a.DefaultPort
		if p, ok := hostPorts[a.Name]; ok && p != 0 {
			hostPort = p
		}
		for _, c := range a.Build(apps.BuildContext{HostPort: hostPort, Secrets: ds}) {
			if c.HostPort != 0 {
				if conflict, ok := seen[c.HostPort]; ok {
					return nil, fmt.Errorf("port conflict: %s and %s both resolve to host port %d", conflict, c.Name, c.HostPort)
				}
				seen[c.HostPort] = c.Name
			}
			cf.Services[c.Name] = buildAppService(c)
			for _, v := range c.NamedVolumes {
				if cf.Volumes == nil {
					cf.Volumes = map[string]struct{}{}
				}
				cf.Volumes[v] = struct{}{}
			}
		}
	}

	return yaml.Marshal(cf)
}

func buildAppService(c apps.Container) composeService {
	tag := c.Tag
	if tag == "" {
		tag = "latest"
	}
	cs := composeService{
		Image:   c.Image + ":" + tag,
		Restart: "unless-stopped",
		Volumes: c.Volumes,
		Command: c.Command,
	}
	if c.HostPort != 0 && c.Port != 0 {
		cs.Ports = []string{fmt.Sprintf("%d:%d", c.HostPort, c.Port)}
	}
	if len(c.Env) > 0 {
		cs.Environment = c.Env
	}
	if len(c.DependsOn) > 0 {
		cs.DependsOn = map[string]dependsOnEntry{}
		for name, condition := range c.DependsOn {
			cs.DependsOn[name] = dependsOnEntry{Condition: condition}
		}
	}
	if c.Healthcheck != nil {
		cs.Healthcheck = &composeHealthcheck{
			Test:     c.Healthcheck.Test,
			Interval: c.Healthcheck.Interval,
			Timeout:  c.Healthcheck.Timeout,
			Retries:  c.Healthcheck.Retries,
		}
	}
	return cs
}
