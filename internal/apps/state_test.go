package apps_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jcombee/devenv/internal/apps"
	"github.com/jcombee/devenv/internal/global"
)

func setupHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	if err := global.Setup(); err != nil {
		t.Fatal(err)
	}
	return home
}

func TestPaths_UnderGlobalDir(t *testing.T) {
	setupHome(t)
	if want := filepath.Join(global.Dir(), "apps"); apps.Dir() != want {
		t.Errorf("Dir: got %q, want %q", apps.Dir(), want)
	}
	if want := filepath.Join(global.Dir(), "apps.yaml"); apps.StatePath() != want {
		t.Errorf("StatePath: got %q, want %q", apps.StatePath(), want)
	}
}

func TestSetup_CreatesAppsDirAndState(t *testing.T) {
	setupHome(t)
	info, err := os.Stat(apps.Dir())
	if err != nil || !info.IsDir() {
		t.Fatalf("expected %s to exist as a directory: %v", apps.Dir(), err)
	}
	if _, err := os.Stat(apps.StatePath()); err != nil {
		t.Fatalf("expected apps.yaml stub: %v", err)
	}
}

func TestLoad_MissingFileIsEmpty(t *testing.T) {
	setupHome(t)
	if err := os.Remove(apps.StatePath()); err != nil {
		t.Fatal(err)
	}
	st, err := apps.Load()
	if err != nil {
		t.Fatalf("Load with no file: %v", err)
	}
	if len(st.Apps) != 0 {
		t.Errorf("expected empty state, got %v", st.Apps)
	}
	if st.IsEnabled("litellm") {
		t.Error("no app should be enabled in an empty state")
	}
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	setupHome(t)
	st, err := apps.Load()
	if err != nil {
		t.Fatal(err)
	}
	st.Apps["litellm"] = apps.Entry{Enabled: true, Port: 4123}
	if err := st.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	data, err := os.ReadFile(apps.StatePath())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "litellm") {
		t.Errorf("apps.yaml missing litellm entry:\n%s", data)
	}

	reloaded, err := apps.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reloaded.IsEnabled("litellm") {
		t.Error("litellm should be enabled after reload")
	}
	if got := reloaded.Apps["litellm"].Port; got != 4123 {
		t.Errorf("Port: got %d, want 4123", got)
	}
}

func TestHostPort_OverrideBeatsDefault(t *testing.T) {
	setupHome(t)
	app, _ := apps.Find("litellm")

	st, err := apps.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got := st.HostPort(app); got != app.DefaultPort {
		t.Errorf("HostPort without override: got %d, want %d", got, app.DefaultPort)
	}

	st.Apps["litellm"] = apps.Entry{Enabled: true, Port: 4500}
	if got := st.HostPort(app); got != 4500 {
		t.Errorf("HostPort with override: got %d, want 4500", got)
	}
}

func TestEnabled_ReturnsRegisteredAppsOnly(t *testing.T) {
	setupHome(t)
	st, err := apps.Load()
	if err != nil {
		t.Fatal(err)
	}
	st.Apps["litellm"] = apps.Entry{Enabled: true}
	st.Apps["ghost-app"] = apps.Entry{Enabled: true} // not in the registry
	st.Apps["disabled-one"] = apps.Entry{Enabled: false}

	enabled := st.Enabled()
	if len(enabled) != 1 || enabled[0].Name != "litellm" {
		t.Errorf("Enabled: got %v, want [litellm]", enabled)
	}
}
