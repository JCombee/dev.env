package e2e_test

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestE2E_Exec_MySQL(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - mysql
`)

	h.MustRunFrom(proj, "start")

	// Reset call log so we only see exec calls.
	h.MustRunFrom(proj, "exec", "mysql")

	calls := h.DockerCalls()
	found := false
	for _, c := range calls {
		if len(c) >= 3 && c[0] == "compose" && contains(c, "exec") && contains(c, "mysql") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected docker compose exec call for mysql, got calls: %v", calls)
	}
}

func TestE2E_Exec_NotRunning(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - mysql
`)

	_, stderr, err := h.RunFrom(proj, "exec", "mysql")
	if err == nil {
		t.Error("expected non-zero exit when project not running")
	}
	if !strings.Contains(stderr, "not running") {
		t.Errorf("expected 'not running' in stderr, got: %s", stderr)
	}
}

func TestE2E_Exec_UnknownService(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - redis
`)

	h.MustRunFrom(proj, "start")

	_, stderr, err := h.RunFrom(proj, "exec", "mysql")
	if err == nil {
		t.Error("expected non-zero exit for unknown service")
	}
	if !strings.Contains(stderr, "mysql") {
		t.Errorf("expected 'mysql' in stderr, got: %s", stderr)
	}
}

func TestE2E_Exec_NoArgs(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - redis
  - mysql
`)

	h.MustRunFrom(proj, "start")

	stdout, _, err := h.RunFrom(proj, "exec")
	if err != nil {
		t.Fatalf("dev exec with no args should succeed, got: %v", err)
	}
	if !strings.Contains(stdout, "redis") || !strings.Contains(stdout, "mysql") {
		t.Errorf("expected service list in output, got: %s", stdout)
	}
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
