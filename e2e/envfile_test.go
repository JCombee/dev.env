package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestE2E_Env_WrittenWhenEnabled(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
env: true
services:
  - image: mysql
    tag: "8.0"
`)

	h.MustRunFrom(proj, "start")

	envPath := filepath.Join(proj, ".env")
	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf(".env not written: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "DB_HOST=127.0.0.1") {
		t.Errorf(".env missing DB_HOST\n%s", content)
	}
	if !strings.Contains(content, "DB_PORT=3380") {
		t.Errorf(".env missing DB_PORT=3380\n%s", content)
	}
	if !strings.Contains(content, "DB_DATABASE=myproject") {
		t.Errorf(".env missing DB_DATABASE\n%s", content)
	}
	if !strings.Contains(content, "DB_PASSWORD=") {
		t.Errorf(".env missing DB_PASSWORD\n%s", content)
	}
}

func TestE2E_Env_NotWrittenWhenDisabled(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
env: false
services:
  - mysql
`)

	h.MustRunFrom(proj, "start")

	envPath := filepath.Join(proj, ".env")
	if _, err := os.Stat(envPath); err == nil {
		t.Error(".env should not be written when env: false")
	}
}

func TestE2E_Env_DefaultNotWritten(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	// No env: field — defaults to false.
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - redis
`)

	h.MustRunFrom(proj, "start")

	envPath := filepath.Join(proj, ".env")
	if _, err := os.Stat(envPath); err == nil {
		t.Error(".env should not be written when env is not set")
	}
}

func TestE2E_Env_Redis(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
env: true
services:
  - redis
`)

	h.MustRunFrom(proj, "start")

	data, err := os.ReadFile(filepath.Join(proj, ".env"))
	if err != nil {
		t.Fatalf(".env not written: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "REDIS_HOST=127.0.0.1") {
		t.Errorf(".env missing REDIS_HOST\n%s", content)
	}
	if !strings.Contains(content, "REDIS_PORT=") {
		t.Errorf(".env missing REDIS_PORT\n%s", content)
	}
	if !strings.Contains(content, "REDIS_DB=") {
		t.Errorf(".env missing REDIS_DB\n%s", content)
	}
}

func TestE2E_Env_EnvMapRenamesVars(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
env: true
services:
  - image: mysql
    tag: "8.0"
    env_map:
      DB_HOST: DATABASE_HOST
      DB_DATABASE: DATABASE_NAME
`)

	h.MustRunFrom(proj, "start")

	data, err := os.ReadFile(filepath.Join(proj, ".env"))
	if err != nil {
		t.Fatalf(".env not written: %v", err)
	}
	content := string(data)

	if strings.Contains(content, "DB_HOST=") {
		t.Error("DB_HOST should have been renamed, not present in .env")
	}
	if !strings.Contains(content, "DATABASE_HOST=127.0.0.1") {
		t.Errorf("DATABASE_HOST not found in .env\n%s", content)
	}
	if strings.Contains(content, "DB_DATABASE=") {
		t.Error("DB_DATABASE should have been renamed, not present in .env")
	}
	if !strings.Contains(content, "DATABASE_NAME=myproject") {
		t.Errorf("DATABASE_NAME not found in .env\n%s", content)
	}
	// Unmapped var passes through.
	if !strings.Contains(content, "DB_PASSWORD=") {
		t.Errorf("DB_PASSWORD should pass through unchanged\n%s", content)
	}
}

func TestE2E_Env_MultipleServices(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
env: true
services:
  - image: mysql
    tag: "8.0"
  - redis
`)

	h.MustRunFrom(proj, "start")

	data, err := os.ReadFile(filepath.Join(proj, ".env"))
	if err != nil {
		t.Fatalf(".env not written: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "DB_HOST=") {
		t.Errorf(".env missing mysql vars\n%s", content)
	}
	if !strings.Contains(content, "REDIS_HOST=") {
		t.Errorf(".env missing redis vars\n%s", content)
	}
}

func TestE2E_Env_OutputIndicatesWritten(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
env: true
services:
  - mysql
`)

	out, _ := h.MustRunFrom(proj, "start")
	if !strings.Contains(out, ".env") {
		t.Errorf("expected .env mention in output, got: %s", out)
	}
}

func TestE2E_Env_NoDiffSkipsWizard(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
env: true
services:
  - mysql
`)

	// First run: creates .env.
	h.MustRunFrom(proj, "start")

	// Second run: same config, nothing changed — should print "up to date".
	out, _ := h.MustRunFrom(proj, "start")
	if !strings.Contains(out, "up to date") {
		t.Errorf("expected 'up to date' on second run with no changes, got: %s", out)
	}
}

func TestE2E_Env_UserVarsPreservedOnUpdate(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
env: true
services:
  - mysql
`)

	// Write a .env with a user var before dev start.
	envPath := filepath.Join(proj, ".env")
	if err := os.WriteFile(envPath, []byte("APP_KEY=base64:abc123\nDB_PORT=9999\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Confirm update (DB_PORT will change from 9999 → computed).
	out, _, _ := h.RunFromWithInput(proj, "y\n", "start")
	_ = out

	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf(".env not found after update: %v", err)
	}
	content := string(data)

	// User var must survive.
	if !strings.Contains(content, "APP_KEY=base64:abc123") {
		t.Errorf("user var APP_KEY not preserved\n%s", content)
	}
	// Dev var must be updated.
	if !strings.Contains(content, "DB_HOST=127.0.0.1") {
		t.Errorf("dev var DB_HOST missing after update\n%s", content)
	}
}

func TestE2E_Env_DeclineSkipsWrite(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
env: true
services:
  - mysql
`)

	// Existing .env with a different port — will trigger diff.
	envPath := filepath.Join(proj, ".env")
	original := "# custom\nDB_PORT=9999\n"
	if err := os.WriteFile(envPath, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	// Decline the update.
	h.RunFromWithInput(proj, "n\n", "start")

	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf(".env not found: %v", err)
	}
	// File should be unchanged.
	if string(data) != original {
		t.Errorf(".env was modified despite declining\ngot:\n%s\nwant:\n%s", string(data), original)
	}
}
