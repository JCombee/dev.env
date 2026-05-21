package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestE2E_DBImport_MySQL_GoldenPath(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - mysql
`)
	h.MustRunFrom(proj, "start")

	writeSQLDump(t, proj, "dump.sql", "INSERT INTO test VALUES (1);\n")

	out, _ := h.MustRunFrom(proj, "db", "import", "dump.sql")
	if !strings.Contains(out, "Import complete") {
		t.Errorf("expected 'Import complete' in output, got: %s", out)
	}

	calls := h.DockerCalls()
	if !hasExecCall(calls) {
		t.Errorf("expected docker compose exec call for import, got calls: %v", calls)
	}
}

func TestE2E_DBImport_Postgres_GoldenPath(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "pgproject")
	writeDevConfig(t, proj, `project: pgproject
type: generic
services:
  - postgres
`)
	h.MustRunFrom(proj, "start")

	writeSQLDump(t, proj, "dump.sql", "CREATE TABLE test (id INT);\n")

	out, _ := h.MustRunFrom(proj, "db", "import", "dump.sql")
	if !strings.Contains(out, "Import complete") {
		t.Errorf("expected 'Import complete' in output, got: %s", out)
	}
}

func TestE2E_DBImport_FileNotFound(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - mysql
`)
	h.MustRunFrom(proj, "start")

	_, _, err := h.RunFrom(proj, "db", "import", "missing.sql")
	h.AssertExitCode(err, 1)
}

func TestE2E_DBImport_NoDBService(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "cacheonlyproject")
	writeDevConfig(t, proj, `project: cacheonlyproject
type: generic
services:
  - redis
`)
	h.MustRunFrom(proj, "start")

	writeSQLDump(t, proj, "dump.sql", "SELECT 1;\n")

	_, _, err := h.RunFrom(proj, "db", "import", "dump.sql")
	h.AssertExitCode(err, 1)
}

func TestE2E_DBImport_MultipleDBServices(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "multidb")
	writeDevConfig(t, proj, `project: multidb
type: generic
services:
  - mysql
  - postgres
`)
	h.MustRunFrom(proj, "start")

	writeSQLDump(t, proj, "dump.sql", "SELECT 1;\n")

	_, _, err := h.RunFrom(proj, "db", "import", "dump.sql")
	h.AssertExitCode(err, 1)
}

func TestE2E_DBImport_ProjectNotRunning(t *testing.T) {
	h := NewHarness(t)
	h.MustRun("setup")

	proj := filepath.Join(h.Home, "myproject")
	writeDevConfig(t, proj, `project: myproject
type: generic
services:
  - mysql
`)

	writeSQLDump(t, proj, "dump.sql", "SELECT 1;\n")

	// No dev start — no state file.
	_, _, err := h.RunFrom(proj, "db", "import", "dump.sql")
	h.AssertExitCode(err, 1)
}

// writeSQLDump writes a SQL file in dir with the given content.
func writeSQLDump(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// hasExecCall reports whether any docker compose exec call appears in calls.
func hasExecCall(calls [][]string) bool {
	for _, c := range calls {
		if len(c) < 2 || c[0] != "compose" {
			continue
		}
		for _, a := range c {
			if a == "exec" {
				return true
			}
		}
	}
	return false
}
