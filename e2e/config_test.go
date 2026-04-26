package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleConfig = `project: myapp
type: laravel
env: true
services:
  - image: mysql
    tag: "8.0"
  - image: redis
`

func writeConfig(t *testing.T, h *Harness, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(h.Home, ".dev.env.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestE2E_ConfigList(t *testing.T) {
	h := NewHarness(t)
	writeConfig(t, h, sampleConfig)

	out, _ := h.MustRunFrom(h.Home, "config", "list")
	for _, want := range []string{"myapp", "laravel", "mysql", "redis"} {
		if !strings.Contains(out, want) {
			t.Errorf("config list output missing %q\noutput: %s", want, out)
		}
	}
}

func TestE2E_ConfigGet(t *testing.T) {
	h := NewHarness(t)
	writeConfig(t, h, sampleConfig)

	out, _ := h.MustRunFrom(h.Home, "config", "get", "type")
	if strings.TrimSpace(out) != "laravel" {
		t.Errorf("config get type: got %q want laravel", strings.TrimSpace(out))
	}
}

func TestE2E_ConfigGet_ServiceTag(t *testing.T) {
	h := NewHarness(t)
	writeConfig(t, h, sampleConfig)

	out, _ := h.MustRunFrom(h.Home, "config", "get", "services.mysql.tag")
	if strings.TrimSpace(out) != "8.0" {
		t.Errorf("config get services.mysql.tag: got %q want 8.0", strings.TrimSpace(out))
	}
}

func TestE2E_ConfigSet(t *testing.T) {
	h := NewHarness(t)
	writeConfig(t, h, sampleConfig)

	h.MustRunFrom(h.Home, "config", "set", "type", "node")
	out, _ := h.MustRunFrom(h.Home, "config", "get", "type")
	if strings.TrimSpace(out) != "node" {
		t.Errorf("after set type=node, got %q", strings.TrimSpace(out))
	}
}

func TestE2E_ConfigSet_ServiceTag(t *testing.T) {
	h := NewHarness(t)
	writeConfig(t, h, sampleConfig)

	h.MustRunFrom(h.Home, "config", "set", "services.mysql.tag", "8.4")
	out, _ := h.MustRunFrom(h.Home, "config", "get", "services.mysql.tag")
	if strings.TrimSpace(out) != "8.4" {
		t.Errorf("after set mysql tag=8.4, got %q", strings.TrimSpace(out))
	}
}

// --- error propagation tests: all must exit 1, not 0 ---

func TestE2E_ConfigGet_UnknownKey_ExitsOne(t *testing.T) {
	h := NewHarness(t)
	writeConfig(t, h, sampleConfig)

	_, _, err := h.RunFrom(h.Home, "config", "get", "doesnotexist")
	h.AssertExitCode(err, 1)
}

func TestE2E_ConfigGet_MissingConfigFile_ExitsOne(t *testing.T) {
	h := NewHarness(t)
	// No .dev.env.yaml written — must propagate read error as exit 1.
	_, _, err := h.RunFrom(h.Home, "config", "get", "type")
	h.AssertExitCode(err, 1)
}

func TestE2E_ConfigSet_InvalidEnvValue_ExitsOne(t *testing.T) {
	h := NewHarness(t)
	writeConfig(t, h, sampleConfig)

	_, _, err := h.RunFrom(h.Home, "config", "set", "env", "notabool")
	h.AssertExitCode(err, 1)
}

func TestE2E_ConfigSet_UnknownService_ExitsOne(t *testing.T) {
	h := NewHarness(t)
	writeConfig(t, h, sampleConfig)

	_, _, err := h.RunFrom(h.Home, "config", "set", "services.postgres.tag", "15")
	h.AssertExitCode(err, 1)
}
