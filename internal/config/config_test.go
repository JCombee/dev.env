package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jcombee/devenv/internal/config"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoad_Basic(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, config.File, `
project: myapp
type: laravel
env: true
services:
  - image: mysql
    tag: "8.0"
  - redis
`)
	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Project != "myapp" {
		t.Errorf("project: got %q want %q", cfg.Project, "myapp")
	}
	if cfg.Type != "laravel" {
		t.Errorf("type: got %q", cfg.Type)
	}
	if !cfg.Env {
		t.Error("env should be true")
	}
	if len(cfg.Services) != 2 {
		t.Fatalf("services: got %d want 2", len(cfg.Services))
	}
	if cfg.Services[0].Image != "mysql" || cfg.Services[0].Tag != "8.0" {
		t.Errorf("service[0]: %+v", cfg.Services[0])
	}
	// shorthand: redis with no tag
	if cfg.Services[1].Image != "redis" {
		t.Errorf("service[1] image: got %q want redis", cfg.Services[1].Image)
	}
}

func TestLoad_ShorthandService(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, config.File, "project: x\nservices:\n  - redis\n")
	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Services[0].Image != "redis" {
		t.Errorf("got image %q", cfg.Services[0].Image)
	}
}

func TestLoad_LocalOverride(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, config.File, "project: myapp\ntype: laravel\nservices:\n  - image: mysql\n    tag: \"8.0\"\n")
	writeFile(t, dir, config.LocalFile, "project: myapp-local\nservices:\n  mysql:\n    tag: \"8.4\"\n")

	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Project != "myapp-local" {
		t.Errorf("project override: got %q", cfg.Project)
	}
	if cfg.Services[0].Tag != "8.4" {
		t.Errorf("service tag override: got %q", cfg.Services[0].Tag)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(t.TempDir())
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.ProjectConfig{
		Project: "test",
		Type:    "node",
		Env:     true,
		Services: []config.ServiceEntry{
			{Image: "postgres", Tag: "15"},
		},
	}
	if err := config.Save(dir, cfg); err != nil {
		t.Fatal(err)
	}
	loaded, err := config.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Project != cfg.Project {
		t.Errorf("project: got %q want %q", loaded.Project, cfg.Project)
	}
	if loaded.Services[0].Image != "postgres" {
		t.Errorf("service image: got %q", loaded.Services[0].Image)
	}
}
