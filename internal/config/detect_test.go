package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jcombee/devenv/internal/config"
)

func TestDetectType_Laravel(t *testing.T) {
	dir := t.TempDir()
	content := `{"require": {"laravel/framework": "^10.0"}}`
	os.WriteFile(filepath.Join(dir, "composer.json"), []byte(content), 0o644)
	if got := config.DetectType(dir); got != "laravel" {
		t.Errorf("got %q want laravel", got)
	}
}

func TestDetectType_Node(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{}`), 0o644)
	if got := config.DetectType(dir); got != "node" {
		t.Errorf("got %q want node", got)
	}
}

func TestDetectType_Unknown(t *testing.T) {
	if got := config.DetectType(t.TempDir()); got != "" {
		t.Errorf("got %q want empty", got)
	}
}

func TestDetectType_LaravelTakesPrecedence(t *testing.T) {
	dir := t.TempDir()
	// both files present — laravel wins
	os.WriteFile(filepath.Join(dir, "composer.json"), []byte(`{"require":{"laravel/framework":"^10"}}`), 0o644)
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{}`), 0o644)
	if got := config.DetectType(dir); got != "laravel" {
		t.Errorf("got %q want laravel", got)
	}
}
