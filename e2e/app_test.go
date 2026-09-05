package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// readHomeFile returns the contents of a file relative to the fake HOME.
func readHomeFile(t *testing.T, h *Harness, rel string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(h.Home, rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(data)
}

// dockerCalled reports whether any recorded docker call contains all the given args in order.
func dockerCalled(calls [][]string, want ...string) bool {
	for _, call := range calls {
		joined := " " + strings.Join(call, " ") + " "
		ok := true
		for _, w := range want {
			if !strings.Contains(joined, " "+w+" ") {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func TestE2E_App_ListShowsAvailableApps(t *testing.T) {
	h := NewHarness(t)
	out, _ := h.MustRun("app", "list")
	for _, want := range []string{"litellm", "APP", "STATUS", "PORT", "disabled"} {
		if !strings.Contains(out, want) {
			t.Errorf("app list output missing %q\n%s", want, out)
		}
	}
}

func TestE2E_App_EnableWritesComposeAndUps(t *testing.T) {
	h := NewHarness(t)
	out, _ := h.MustRun("app", "enable", "litellm")

	if !h.FileExists(".dev.env/apps/docker-compose.yml") {
		t.Fatal("apps compose file not written")
	}
	composeFile := readHomeFile(t, h, ".dev.env/apps/docker-compose.yml")
	for _, want := range []string{"litellm:", "litellm-db:", "4000:4000"} {
		if !strings.Contains(composeFile, want) {
			t.Errorf("apps compose missing %q\n%s", want, composeFile)
		}
	}

	calls := h.DockerCalls()
	if !dockerCalled(calls, "compose", "up", "litellm", "litellm-db") {
		t.Errorf("expected docker compose up for both containers, got %v", calls)
	}
	if !dockerCalled(calls, "-f", filepath.Join(h.Home, ".dev.env", "apps", "docker-compose.yml")) {
		t.Errorf("docker was not pointed at the apps compose file, got %v", calls)
	}

	state := readHomeFile(t, h, ".dev.env/apps.yaml")
	if !strings.Contains(state, "litellm") || !strings.Contains(state, "true") {
		t.Errorf("apps.yaml does not record litellm as enabled\n%s", state)
	}

	if !strings.Contains(out, "4000") {
		t.Errorf("enable output should show the proxy URL\n%s", out)
	}
}

func TestE2E_App_EnableGeneratesStableSecrets(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("app", "enable", "litellm")
	first := readHomeFile(t, h, ".dev.env/docker/secrets.yaml")

	for _, want := range []string{"litellm:", "master_key", "salt_key", "ui_password", "litellm-db:"} {
		if !strings.Contains(first, want) {
			t.Errorf("secrets.yaml missing %q\n%s", want, first)
		}
	}

	h.MustRun("app", "enable", "litellm")
	if second := readHomeFile(t, h, ".dev.env/docker/secrets.yaml"); second != first {
		t.Errorf("re-enabling rotated credentials\nbefore:\n%s\nafter:\n%s", first, second)
	}
}

func TestE2E_App_ListShowsRunningAfterEnable(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("app", "enable", "litellm")

	out, _ := h.MustRun("app", "list")
	if !strings.Contains(out, "running") {
		t.Errorf("expected litellm to show as running\n%s", out)
	}
}

func TestE2E_App_DisableStopsAndUnregisters(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("app", "enable", "litellm")
	h.MustRun("app", "disable", "litellm")

	if !dockerCalled(h.DockerCalls(), "compose", "stop", "litellm", "litellm-db") {
		t.Errorf("expected docker compose stop for the app containers, got %v", h.DockerCalls())
	}
	compose := readHomeFile(t, h, ".dev.env/apps/docker-compose.yml")
	if strings.Contains(compose, "litellm") {
		t.Errorf("disabled app must be removed from the compose file\n%s", compose)
	}
	if state := readHomeFile(t, h, ".dev.env/apps.yaml"); strings.Contains(state, "enabled: true") {
		t.Errorf("apps.yaml still marks an app enabled\n%s", state)
	}
	// Credentials survive so re-enabling restores the same app.
	if !h.FileExists(".dev.env/docker/secrets.yaml") {
		t.Error("secrets must be kept after disable")
	}
}

func TestE2E_App_DisableWhenNotEnabled(t *testing.T) {
	h := NewHarness(t)
	out, _, err := h.Run("app", "disable", "litellm")
	h.AssertExitCode(err, 0)
	if !strings.Contains(out, "not enabled") {
		t.Errorf("expected a 'not enabled' notice\n%s", out)
	}
}

func TestE2E_App_StartStop(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("app", "enable", "litellm")

	h.MustRun("app", "stop", "litellm")
	if !dockerCalled(h.DockerCalls(), "compose", "stop", "litellm") {
		t.Errorf("app stop did not stop the container, got %v", h.DockerCalls())
	}
	out, _ := h.MustRun("app", "list")
	if !strings.Contains(out, "stopped") {
		t.Errorf("expected status 'stopped' for an enabled but stopped app\n%s", out)
	}

	h.MustRun("app", "start")
	if !dockerCalled(h.DockerCalls(), "compose", "up", "litellm", "litellm-db") {
		t.Errorf("app start did not start enabled apps, got %v", h.DockerCalls())
	}

	// Still enabled after a stop/start cycle.
	if state := readHomeFile(t, h, ".dev.env/apps.yaml"); !strings.Contains(state, "enabled: true") {
		t.Errorf("stop/start must not change enabled state\n%s", state)
	}
}

func TestE2E_App_StartWithNoAppsEnabled(t *testing.T) {
	h := NewHarness(t)
	out, _, err := h.Run("app", "start")
	h.AssertExitCode(err, 0)
	if !strings.Contains(out, "No apps enabled") {
		t.Errorf("expected a 'No apps enabled' notice\n%s", out)
	}
}

func TestE2E_App_Info(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("app", "enable", "litellm")

	out, _ := h.MustRun("app", "info", "litellm")
	secrets := readHomeFile(t, h, ".dev.env/docker/secrets.yaml")
	for _, want := range []string{"http://127.0.0.1:4000", "/ui", "admin"} {
		if !strings.Contains(out, want) {
			t.Errorf("app info missing %q\n%s", want, out)
		}
	}
	if !strings.Contains(out, "sk-") {
		t.Errorf("app info should print the master key\n%s", out)
	}
	// The printed key must be the stored one.
	for _, line := range strings.Split(out, "\n") {
		if idx := strings.Index(line, "sk-"); idx >= 0 {
			key := strings.TrimSpace(line[idx:])
			if !strings.Contains(secrets, key) {
				t.Errorf("printed master key %q not found in secrets.yaml", key)
			}
		}
	}
}

func TestE2E_App_InfoWhenDisabled_ExitsOne(t *testing.T) {
	h := NewHarness(t)
	_, _, err := h.Run("app", "info", "litellm")
	h.AssertExitCode(err, 1)
}

func TestE2E_App_UnknownName_ExitsOne(t *testing.T) {
	h := NewHarness(t)
	for _, args := range [][]string{
		{"app", "enable", "nope"},
		{"app", "disable", "nope"},
		{"app", "start", "nope"},
		{"app", "stop", "nope"},
		{"app", "info", "nope"},
	} {
		_, stderr, err := h.Run(args...)
		h.AssertExitCode(err, 1)
		if !strings.Contains(stderr, "litellm") {
			t.Errorf("%v: error should list the known apps, got: %s", args, stderr)
		}
	}
}

func TestE2E_App_PortOverride(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("app", "enable", "litellm")

	appsYAML := filepath.Join(h.Home, ".dev.env", "apps.yaml")
	if err := os.WriteFile(appsYAML, []byte("apps:\n  litellm:\n    enabled: true\n    port: 4500\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	h.MustRun("app", "start", "litellm")
	compose := readHomeFile(t, h, ".dev.env/apps/docker-compose.yml")
	if !strings.Contains(compose, "4500:4000") {
		t.Errorf("port override from apps.yaml not applied\n%s", compose)
	}
	out, _ := h.MustRun("app", "info", "litellm")
	if !strings.Contains(out, "4500") {
		t.Errorf("app info should use the overridden port\n%s", out)
	}
}

func TestE2E_App_DoesNotAffectProjects(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("app", "enable", "litellm")

	// Apps are invisible to project-level status and never enter services.yaml.
	services := readHomeFile(t, h, ".dev.env/services.yaml")
	if strings.Contains(services, "litellm") {
		t.Errorf("apps must not be registered as project services\n%s", services)
	}
	out, _ := h.MustRun("status")
	if strings.Contains(out, "litellm") {
		t.Errorf("dev status must stay project-only\n%s", out)
	}
}
