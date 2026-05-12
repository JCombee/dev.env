package compose_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jcombee/devenv/internal/compose"
	"github.com/jcombee/devenv/internal/config"
	"github.com/jcombee/devenv/internal/global"
)

func TestWriteShared_CreatesFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	if err := global.Setup(); err != nil {
		t.Fatal(err)
	}

	svcMap := map[string]global.GlobalServiceEntry{
		"mysql-8-0":    {Image: "mysql", Tag: "8.0"},
		"redis-latest": {Image: "redis", Tag: "latest"},
	}
	ds := global.DockerSecrets{}

	if err := compose.WriteShared(svcMap, ds, nil); err != nil {
		t.Fatalf("WriteShared: %v", err)
	}

	data, err := os.ReadFile(compose.SharedPath())
	if err != nil {
		t.Fatalf("read compose file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "mysql:8.0") {
		t.Errorf("compose file missing mysql:8.0\n%s", content)
	}
	if !strings.Contains(content, "redis:latest") {
		t.Errorf("compose file missing redis:latest\n%s", content)
	}
}

func TestWriteDedicated_ProjectPrefixed(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	if err := global.Setup(); err != nil {
		t.Fatal(err)
	}
	projDir := filepath.Join(global.Dir(), "projects", "myproject")
	if err := os.MkdirAll(projDir, 0o755); err != nil {
		t.Fatal(err)
	}

	entries := []config.ServiceEntry{
		{Image: "mysql", Tag: "8.0", Dedicated: true},
	}
	ds := global.DockerSecrets{}

	if err := compose.WriteDedicated("myproject", entries, ds); err != nil {
		t.Fatalf("WriteDedicated: %v", err)
	}

	data, err := os.ReadFile(compose.DedicatedPath("myproject"))
	if err != nil {
		t.Fatalf("read dedicated compose: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "myproject-mysql-8-0") {
		t.Errorf("dedicated compose missing project-prefixed name\n%s", content)
	}
	if !strings.Contains(content, "mysql:8.0") {
		t.Errorf("dedicated compose missing image tag\n%s", content)
	}
}

func TestWriteShared_Idempotent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	if err := global.Setup(); err != nil {
		t.Fatal(err)
	}

	svcMap := map[string]global.GlobalServiceEntry{
		"redis-latest": {Image: "redis", Tag: "latest"},
	}
	ds := global.DockerSecrets{}

	if err := compose.WriteShared(svcMap, ds, nil); err != nil {
		t.Fatal(err)
	}
	if err := compose.WriteShared(svcMap, ds, nil); err != nil {
		t.Fatalf("second WriteShared: %v", err)
	}
}
