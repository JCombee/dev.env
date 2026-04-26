package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeGlobalServices(t *testing.T, h *Harness, content string) {
	t.Helper()
	dir := filepath.Join(h.Home, ".dev.env")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "services.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeProjectState(t *testing.T, h *Harness, projectName, content string) {
	t.Helper()
	dir := filepath.Join(h.Home, ".dev.env", "projects", projectName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// project.yaml required for refcount scanner
	if err := os.WriteFile(filepath.Join(dir, "project.yaml"), []byte("name: "+projectName+"\npath: /code/"+projectName+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "state.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestE2E_Status_NoServices(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")
	out, _ := h.MustRun("status")
	if !strings.Contains(out, "No services") {
		t.Errorf("expected 'No services' message, got: %s", out)
	}
}

func TestE2E_Status_ShowsRegisteredServices(t *testing.T) {
	h := NewHarness(t)
	writeGlobalServices(t, h, `services:
  mysql-8-0:
    image: mysql
    tag: "8.0"
  redis-latest:
    image: redis
    tag: latest
`)

	out, _ := h.MustRun("status")
	for _, want := range []string{"mysql:8.0", "redis:latest", "STATUS", "PROJECTS"} {
		if !strings.Contains(out, want) {
			t.Errorf("status output missing %q\n%s", want, out)
		}
	}
}

func TestE2E_Status_ShowsProjectUsage(t *testing.T) {
	h := NewHarness(t)
	writeGlobalServices(t, h, `services:
  mysql-8-0:
    image: mysql
    tag: "8.0"
`)
	writeProjectState(t, h, "api", `running: true
services:
  mysql-8-0: shared
acknowledged: []
`)

	out, _ := h.MustRun("status")
	if !strings.Contains(out, "api") {
		t.Errorf("expected project 'api' in status output\n%s", out)
	}
}

func TestE2E_Status_StoppedProjectNotShown(t *testing.T) {
	h := NewHarness(t)
	writeGlobalServices(t, h, `services:
  mysql-8-0:
    image: mysql
    tag: "8.0"
`)
	writeProjectState(t, h, "api", "running: true\nservices:\n  mysql-8-0: shared\nacknowledged: []\n")
	writeProjectState(t, h, "stopped-proj", "running: false\nservices:\n  mysql-8-0: shared\nacknowledged: []\n")

	out, _ := h.MustRun("status")
	if !strings.Contains(out, "api") {
		t.Errorf("expected running project 'api' in output\n%s", out)
	}
	if strings.Contains(out, "stopped-proj") {
		t.Errorf("stopped project should not appear in output\n%s", out)
	}
}

func TestE2E_Status_ShowsUnknownStatus(t *testing.T) {
	h := NewHarness(t)
	writeGlobalServices(t, h, `services:
  redis-latest:
    image: redis
    tag: latest
`)

	out, _ := h.MustRun("status")
	if !strings.Contains(out, "unknown") {
		t.Errorf("expected 'unknown' status before Phase 4 docker integration\n%s", out)
	}
}

func TestE2E_Status_MultipleProjects(t *testing.T) {
	h := NewHarness(t)
	writeGlobalServices(t, h, `services:
  mysql-8-0:
    image: mysql
    tag: "8.0"
`)
	writeProjectState(t, h, "api", "running: true\nservices:\n  mysql-8-0: shared\nacknowledged: []\n")
	writeProjectState(t, h, "shop", "running: true\nservices:\n  mysql-8-0: shared\nacknowledged: []\n")

	out, _ := h.MustRun("status")
	if !strings.Contains(out, "api") || !strings.Contains(out, "shop") {
		t.Errorf("expected both projects in output\n%s", out)
	}
}
