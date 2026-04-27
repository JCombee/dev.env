package global_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jcombee/devenv/internal/global"
)

// runSetup calls global.Setup with a temp dir as the home directory.
func runSetup(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp) // Windows fallback
	if err := global.Setup(); err != nil {
		t.Fatalf("Setup() error: %v", err)
	}
	return tmp
}

func TestSetup_CreatesDirectories(t *testing.T) {
	tmp := runSetup(t)
	base := filepath.Join(tmp, ".dev.env")

	dirs := []string{base, filepath.Join(base, "docker"), filepath.Join(base, "projects")}
	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			t.Errorf("expected dir %s to exist", d)
		}
	}
}

func TestSetup_CreatesSettingsAndServices(t *testing.T) {
	tmp := runSetup(t)
	base := filepath.Join(tmp, ".dev.env")

	files := []string{
		filepath.Join(base, "settings.yaml"),
		filepath.Join(base, "services.yaml"),
	}
	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", f)
		}
	}
}

func TestSetup_Idempotent(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)

	for i := range 3 {
		if err := global.Setup(); err != nil {
			t.Fatalf("Setup() call %d error: %v", i+1, err)
		}
	}
}

func TestSetup_StubFilesAreValidYAML(t *testing.T) {
	tmp := runSetup(t)
	base := filepath.Join(tmp, ".dev.env")

	cases := map[string]string{
		"settings.yaml": "default_type:",
		"services.yaml": "services:",
	}
	for name, want := range cases {
		data, err := os.ReadFile(filepath.Join(base, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		content := string(data)
		if !containsString(content, want) {
			t.Errorf("%s: want YAML key %q, got:\n%s", name, want, content)
		}
	}
}

func containsString(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := range len(s) - len(sub) + 1 {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestSettings_PortsRoundTrip(t *testing.T) {
	tmp := runSetup(t)
	base := filepath.Join(tmp, ".dev.env")
	settingsPath := filepath.Join(base, "settings.yaml")

	content := "default_type: laravel\nports:\n  mysql-8-0: 13380\n  redis-7-0: 16370\n"
	if err := os.WriteFile(settingsPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	settings, err := global.ReadSettings()
	if err != nil {
		t.Fatal(err)
	}
	if settings.Ports["mysql-8-0"] != 13380 {
		t.Errorf("mysql-8-0 port: got %d want 13380", settings.Ports["mysql-8-0"])
	}
	if settings.Ports["redis-7-0"] != 16370 {
		t.Errorf("redis-7-0 port: got %d want 16370", settings.Ports["redis-7-0"])
	}
}

func TestIsSetUp(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)

	if global.IsSetUp() {
		t.Fatal("IsSetUp() should be false before Setup()")
	}
	if err := global.Setup(); err != nil {
		t.Fatal(err)
	}
	if !global.IsSetUp() {
		t.Fatal("IsSetUp() should be true after Setup()")
	}
}
