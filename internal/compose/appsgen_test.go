package compose_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jcombee/devenv/internal/apps"
	"github.com/jcombee/devenv/internal/compose"
	"github.com/jcombee/devenv/internal/global"
	"gopkg.in/yaml.v3"
)

func setupAppsHome(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	if err := global.Setup(); err != nil {
		t.Fatal(err)
	}
}

// writeLiteLLM renders the litellm app and returns the generated compose file.
func writeLiteLLM(t *testing.T, ports map[string]int) (string, global.DockerSecrets) {
	t.Helper()
	app, ok := apps.Find("litellm")
	if !ok {
		t.Fatal("litellm not registered")
	}
	ds := global.DockerSecrets{}
	apps.EnsureSecrets(app, ds)

	if err := compose.WriteApps([]apps.App{app}, ports, ds); err != nil {
		t.Fatalf("WriteApps: %v", err)
	}
	data, err := os.ReadFile(compose.AppsPath())
	if err != nil {
		t.Fatalf("read apps compose: %v", err)
	}
	return string(data), ds
}

func TestAppsPath_SeparateFromShared(t *testing.T) {
	setupAppsHome(t)
	want := filepath.Join(global.Dir(), "apps", "docker-compose.yml")
	if compose.AppsPath() != want {
		t.Errorf("AppsPath: got %q, want %q", compose.AppsPath(), want)
	}
	if compose.AppsPath() == compose.SharedPath() {
		t.Error("apps must not share the project compose file")
	}
}

func TestWriteApps_RendersLiteLLMStack(t *testing.T) {
	setupAppsHome(t)
	content, ds := writeLiteLLM(t, nil)

	for _, want := range []string{
		"litellm:",
		"litellm-db:",
		"postgres:16",
		"4000:4000",
		"restart: unless-stopped",
		"STORE_MODEL_IN_DB",
		ds["litellm"]["master_key"],
	} {
		if !strings.Contains(content, want) {
			t.Errorf("apps compose missing %q\n%s", want, content)
		}
	}
}

func TestWriteApps_DatabaseNotPublished(t *testing.T) {
	setupAppsHome(t)
	content, _ := writeLiteLLM(t, nil)

	var cf struct {
		Services map[string]struct {
			Ports []string `yaml:"ports"`
		} `yaml:"services"`
	}
	if err := yaml.Unmarshal([]byte(content), &cf); err != nil {
		t.Fatalf("generated compose is not valid YAML: %v\n%s", err, content)
	}
	if len(cf.Services["litellm-db"].Ports) != 0 {
		t.Errorf("litellm-db must publish no host port, got %v", cf.Services["litellm-db"].Ports)
	}
	if len(cf.Services["litellm"].Ports) != 1 {
		t.Errorf("litellm must publish exactly one port, got %v", cf.Services["litellm"].Ports)
	}
}

func TestWriteApps_EmitsVolumesDependsOnAndHealthcheck(t *testing.T) {
	setupAppsHome(t)
	content, _ := writeLiteLLM(t, nil)

	var cf struct {
		Services map[string]struct {
			Volumes     []string          `yaml:"volumes"`
			DependsOn   map[string]any    `yaml:"depends_on"`
			Healthcheck map[string]any    `yaml:"healthcheck"`
			Environment map[string]string `yaml:"environment"`
		} `yaml:"services"`
		Volumes map[string]any `yaml:"volumes"`
	}
	if err := yaml.Unmarshal([]byte(content), &cf); err != nil {
		t.Fatalf("generated compose is not valid YAML: %v\n%s", err, content)
	}

	db := cf.Services["litellm-db"]
	if len(db.Volumes) == 0 {
		t.Error("litellm-db has no volumes")
	}
	if len(db.Healthcheck) == 0 {
		t.Error("litellm-db has no healthcheck")
	}
	if _, ok := cf.Volumes["litellm-db-data"]; !ok {
		t.Errorf("top-level volumes must declare litellm-db-data, got %v", cf.Volumes)
	}
	if _, ok := cf.Services["litellm"].DependsOn["litellm-db"]; !ok {
		t.Error("litellm must depend on litellm-db")
	}
}

func TestWriteApps_PortOverride(t *testing.T) {
	setupAppsHome(t)
	content, _ := writeLiteLLM(t, map[string]int{"litellm": 4500})

	if !strings.Contains(content, "4500:4000") {
		t.Errorf("port override not applied\n%s", content)
	}
	if strings.Contains(content, "4000:4000") {
		t.Errorf("default port still used despite override\n%s", content)
	}
}

func TestWriteApps_NoAppsWritesEmptyFile(t *testing.T) {
	setupAppsHome(t)
	if err := compose.WriteApps(nil, nil, global.DockerSecrets{}); err != nil {
		t.Fatalf("WriteApps with no apps: %v", err)
	}
	data, err := os.ReadFile(compose.AppsPath())
	if err != nil {
		t.Fatalf("read apps compose: %v", err)
	}
	if strings.Contains(string(data), "litellm") {
		t.Errorf("disabled app still present in compose file\n%s", data)
	}
}

func TestWriteApps_PortConflict(t *testing.T) {
	setupAppsHome(t)
	app, _ := apps.Find("litellm")
	twin := app
	twin.Name = "litellm-twin"

	err := compose.WriteApps([]apps.App{app, twin}, nil, global.DockerSecrets{})
	if err == nil {
		t.Fatal("expected an error when two apps resolve to the same host port")
	}
	if !strings.Contains(err.Error(), "port conflict") {
		t.Errorf("error should mention a port conflict, got: %v", err)
	}
}

func TestWriteShared_UnaffectedByAppFields(t *testing.T) {
	setupAppsHome(t)
	svcMap := map[string]global.GlobalServiceEntry{
		"redis-latest": {Image: "redis", Tag: "latest"},
	}
	if err := compose.WriteShared(svcMap, global.DockerSecrets{}, nil); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(compose.SharedPath())
	if err != nil {
		t.Fatal(err)
	}
	for _, unwanted := range []string{"volumes:", "depends_on:", "healthcheck:", "command:"} {
		if strings.Contains(string(data), unwanted) {
			t.Errorf("shared compose gained %q — app-only keys must stay omitempty\n%s", unwanted, data)
		}
	}
}
