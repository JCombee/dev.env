package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeDevConfig writes a .dev.env.yaml in the given project dir.
func writeDevConfig(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".dev.env.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestE2E_Start_SharedService(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - mysql
`)

	out, _ := h.MustRunFrom(proj, "start")
	if !strings.Contains(out, "started") {
		t.Errorf("expected 'started' in output, got: %s", out)
	}

	calls := h.DockerCalls()
	foundUp := false
	for _, c := range calls {
		if len(c) >= 2 && c[0] == "compose" {
			for _, a := range c {
				if a == "up" {
					foundUp = true
				}
			}
		}
	}
	if !foundUp {
		t.Errorf("expected docker compose up call, got calls: %v", calls)
	}
}

func TestE2E_Start_UpdatesStateFile(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - redis
`)

	h.MustRunFrom(proj, "start")

	stateFile := filepath.Join(h.Home, ".dev.env", "projects", "myproject", "state.yaml")
	data, err := os.ReadFile(stateFile)
	if err != nil {
		t.Fatalf("state.yaml not found: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "running: true") {
		t.Errorf("state.yaml should have running:true\n%s", content)
	}
	if !strings.Contains(content, "redis") {
		t.Errorf("state.yaml should reference redis service\n%s", content)
	}
}

func TestE2E_Start_RegistersInServicesYAML(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - redis
`)

	h.MustRunFrom(proj, "start")

	svcFile := filepath.Join(h.Home, ".dev.env", "services.yaml")
	data, err := os.ReadFile(svcFile)
	if err != nil {
		t.Fatalf("services.yaml not found: %v", err)
	}
	if !strings.Contains(string(data), "redis") {
		t.Errorf("services.yaml should contain redis\n%s", string(data))
	}
}

func TestE2E_Stop_StopsSharedWhenNoOtherUsers(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - redis
`)

	h.MustRunFrom(proj, "start")
	out, _ := h.MustRunFrom(proj, "stop")

	if !strings.Contains(out, "stopped") {
		t.Errorf("expected 'stopped' in output, got: %s", out)
	}

	calls := h.DockerCalls()
	foundStop := false
	for _, c := range calls {
		for _, a := range c {
			if a == "stop" {
				foundStop = true
			}
		}
	}
	if !foundStop {
		t.Errorf("expected docker compose stop call, got: %v", calls)
	}
}

func TestE2E_Stop_KeepsSharedWhenOtherProjectRunning(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj1 := filepath.Join(h.Home, "project1")
	writeDevConfig(t, proj1, `project: project1
type: generic
services:
  - redis
`)
	proj2 := filepath.Join(h.Home, "project2")
	writeDevConfig(t, proj2, `project: project2
type: generic
services:
  - redis
`)

	h.MustRunFrom(proj1, "start")
	h.MustRunFrom(proj2, "start")

	// Stop project1 — redis should be kept (project2 still running).
	out, _ := h.MustRunFrom(proj1, "stop")
	if !strings.Contains(out, "kept running") {
		t.Errorf("expected 'kept running' message, got: %s", out)
	}
}

func TestE2E_Stop_UpdatesStateRunningFalse(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - redis
`)

	h.MustRunFrom(proj, "start")
	h.MustRunFrom(proj, "stop")

	stateFile := filepath.Join(h.Home, ".dev.env", "projects", "myproject", "state.yaml")
	data, err := os.ReadFile(stateFile)
	if err != nil {
		t.Fatalf("state.yaml not found: %v", err)
	}
	if !strings.Contains(string(data), "running: false") {
		t.Errorf("state.yaml should have running:false after stop\n%s", string(data))
	}
}

func TestE2E_Stop_AlreadyStopped(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - redis
`)

	h.MustRunFrom(proj, "start")
	h.MustRunFrom(proj, "stop")

	// Second stop should be a no-op, not an error.
	out, _, err := h.RunFrom(proj, "stop")
	h.AssertExitCode(err, 0)
	if !strings.Contains(out, "not running") {
		t.Errorf("expected 'not running' message, got: %s", out)
	}
}

func TestE2E_Start_DedicatedService(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - image: mysql
    tag: "8.0"
    dedicated: true
`)

	out, _ := h.MustRunFrom(proj, "start")
	if !strings.Contains(out, "dedicated") {
		t.Errorf("expected 'dedicated' in output, got: %s", out)
	}
}

func TestE2E_Start_NoConfig_ExitsOne(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "empty")
	if err := os.MkdirAll(proj, 0o755); err != nil {
		t.Fatal(err)
	}

	_, _, err := h.RunFrom(proj, "start")
	h.AssertExitCode(err, 1)
}

func TestE2E_Start_MySQL_ProvisionDB(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - mysql
`)

	out, _ := h.MustRunFrom(proj, "start")
	if !strings.Contains(out, `database "myproject" ready`) {
		t.Errorf("expected database ready message, got: %s", out)
	}

	calls := h.DockerCalls()
	foundExec := false
	for _, c := range calls {
		if len(c) >= 2 && c[0] == "compose" {
			for _, a := range c {
				if a == "exec" {
					foundExec = true
				}
			}
		}
	}
	if !foundExec {
		t.Errorf("expected docker compose exec call for DB provisioning, got calls: %v", calls)
	}
}

func TestE2E_Start_Postgres_ProvisionDB(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "pgproject")
	writeDevConfig(t, proj, `project: pgproject
type: generic
services:
  - postgres
`)

	out, _ := h.MustRunFrom(proj, "start")
	if !strings.Contains(out, `database "pgproject" ready`) {
		t.Errorf("expected database ready message, got: %s", out)
	}
}
