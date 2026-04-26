package project_test

import (
	"errors"
	"os"
	"testing"

	"github.com/jcombee/devenv/internal/project"
	"github.com/jcombee/devenv/internal/store"
)

func setupHome(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)
	if err := os.MkdirAll(tmp+"/.dev.env/projects", 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestRegister_CreatesFiles(t *testing.T) {
	setupHome(t)
	if err := project.Register("myapp", "/code/myapp"); err != nil {
		t.Fatalf("Register: %v", err)
	}

	var meta project.Meta
	if err := store.Read(project.MetaPath("myapp"), &meta); err != nil {
		t.Fatalf("read project.yaml: %v", err)
	}
	if meta.Name != "myapp" || meta.Path != "/code/myapp" {
		t.Errorf("meta: %+v", meta)
	}

	state, err := project.LoadState("myapp")
	if err != nil {
		t.Fatalf("LoadState: %v", err)
	}
	if state.Services == nil {
		t.Error("state.Services should be initialised")
	}
}

func TestRegister_Idempotent(t *testing.T) {
	setupHome(t)
	for range 3 {
		if err := project.Register("myapp", "/code/myapp"); err != nil {
			t.Fatalf("Register: %v", err)
		}
	}
}

func TestRegister_NameCollision(t *testing.T) {
	setupHome(t)
	if err := project.Register("myapp", "/code/myapp"); err != nil {
		t.Fatal(err)
	}
	err := project.Register("myapp", "/code/other-app")
	if err == nil {
		t.Fatal("expected collision error")
	}
	var col *project.ErrNameCollision
	if !errors.As(err, &col) {
		t.Fatalf("expected ErrNameCollision, got %T: %v", err, err)
	}
	if col.ExistingPath != "/code/myapp" {
		t.Errorf("collision path: %q", col.ExistingPath)
	}
}
