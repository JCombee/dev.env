package e2e_test

import (
	"testing"
)

func TestE2E_Setup_CreatesDirectoryTree(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	dirs := []string{
		".dev.env",
		".dev.env/docker",
		".dev.env/projects",
	}
	for _, d := range dirs {
		if !h.DirExists(d) {
			t.Errorf("expected dir %s to exist", d)
		}
	}

	files := []string{
		".dev.env/settings.yaml",
		".dev.env/services.yaml",
	}
	for _, f := range files {
		if !h.FileExists(f) {
			t.Errorf("expected file %s to exist", f)
		}
	}
}

func TestE2E_Setup_Idempotent(t *testing.T) {
	h := NewHarness(t)
	for range 3 {
		_, _, err := h.Run("setup")
		if err != nil {
			t.Fatalf("setup run failed: %v", err)
		}
	}
	if !h.DirExists(".dev.env") {
		t.Fatal("expected .dev.env to exist after repeated setup")
	}
}

func TestE2E_Setup_ExitZero(t *testing.T) {
	h := NewHarness(t)
	_, _, err := h.Run("setup")
	if err != nil {
		t.Fatalf("expected exit 0, got: %v", err)
	}
}

func TestE2E_AutoSetup_OnAnyCommand(t *testing.T) {
	h := NewHarness(t)
	// Run a command other than setup without prior setup — auto-setup should trigger.
	h.MustRun("version")
	if !h.DirExists(".dev.env") {
		t.Fatal("expected auto-setup to create .dev.env on first command")
	}
}
