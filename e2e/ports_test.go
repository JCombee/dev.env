package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestE2E_Ports_ComputedPortInComposeFile(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - image: mysql
    tag: "8.0"
`)

	h.MustRunFrom(proj, "start")

	composeFile := filepath.Join(h.Home, ".dev.env", "docker", "docker-compose.yml")
	data, err := os.ReadFile(composeFile)
	if err != nil {
		t.Fatalf("docker-compose.yml not found: %v", err)
	}
	if !strings.Contains(string(data), "3380:3306") {
		t.Errorf("expected computed port mapping 3380:3306 in compose file, got:\n%s", string(data))
	}
}

func TestE2E_Ports_GlobalOverrideInComposeFile(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	settingsPath := filepath.Join(h.Home, ".dev.env", "settings.yaml")
	if err := os.WriteFile(settingsPath, []byte("default_type: generic\nports:\n  mysql-8-0: 13380\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - image: mysql
    tag: "8.0"
`)

	h.MustRunFrom(proj, "start")

	composeFile := filepath.Join(h.Home, ".dev.env", "docker", "docker-compose.yml")
	data, err := os.ReadFile(composeFile)
	if err != nil {
		t.Fatalf("docker-compose.yml not found: %v", err)
	}
	if !strings.Contains(string(data), "13380:3306") {
		t.Errorf("expected global override port mapping 13380:3306 in compose file, got:\n%s", string(data))
	}
}

func TestE2E_Ports_ProjectOverride_ForcesDedicated(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - image: mysql
    tag: "8.0"
`)
	if err := os.WriteFile(filepath.Join(proj, ".dev.env.local.yaml"), []byte("services:\n  mysql:\n    port: 23306\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, _ := h.MustRunFrom(proj, "start")
	if !strings.Contains(out, "dedicated") {
		t.Errorf("expected 'dedicated' in output when project port override is set, got: %s", out)
	}

	// dedicated container compose file lives under projects/<name>/
	composeFile := filepath.Join(h.Home, ".dev.env", "projects", "myproject", "docker-compose.yml")
	data, err := os.ReadFile(composeFile)
	if err != nil {
		t.Fatalf("project docker-compose.yml not found: %v", err)
	}
	if !strings.Contains(string(data), "23306:3306") {
		t.Errorf("expected project port mapping 23306:3306 in dedicated compose file, got:\n%s", string(data))
	}
}

func TestE2E_Ports_Collision_ExitsNonZero(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	// Force a collision via global overrides — two different services, same host port.
	settingsPath := filepath.Join(h.Home, ".dev.env", "settings.yaml")
	if err := os.WriteFile(settingsPath, []byte("ports:\n  mysql-8-0: 9999\n  redis-latest: 9999\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - image: mysql
    tag: "8.0"
  - redis
`)

	_, stderr, err := h.RunFrom(proj, "start")
	h.AssertExitCode(err, 1)
	if !strings.Contains(stderr, "port conflict") {
		t.Errorf("expected 'port conflict' in stderr, got: %s", stderr)
	}
}
