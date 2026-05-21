package compose_test

import (
	"reflect"
	"testing"

	"github.com/jcombee/devenv/internal/compose"
	"github.com/jcombee/devenv/internal/global"
)

func TestImportArgs_MySQL(t *testing.T) {
	ds := global.DockerSecrets{"mysql-8-0": {"root_password": "secret"}}
	got, ok := compose.ImportArgs("mysql", "mysql-8-0", "myapp", ds)
	if !ok {
		t.Fatal("expected ok=true for mysql")
	}
	want := []string{"mysql", "-u", "root", "-psecret", "myapp"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestImportArgs_MariaDB(t *testing.T) {
	ds := global.DockerSecrets{"mariadb-latest": {"root_password": "pass"}}
	got, ok := compose.ImportArgs("mariadb", "mariadb-latest", "shop", ds)
	if !ok {
		t.Fatal("expected ok=true for mariadb")
	}
	want := []string{"mysql", "-u", "root", "-ppass", "shop"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestImportArgs_Percona(t *testing.T) {
	ds := global.DockerSecrets{"percona-8-0": {"root_password": "percpass"}}
	got, ok := compose.ImportArgs("percona", "percona-8-0", "blog", ds)
	if !ok {
		t.Fatal("expected ok=true for percona")
	}
	want := []string{"mysql", "-u", "root", "-ppercpass", "blog"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestImportArgs_Postgres(t *testing.T) {
	ds := global.DockerSecrets{"postgres-16": {"root_password": "pgpass"}}
	got, ok := compose.ImportArgs("postgres", "postgres-16", "api", ds)
	if !ok {
		t.Fatal("expected ok=true for postgres")
	}
	want := []string{"env", "PGPASSWORD=pgpass", "psql", "-U", "postgres", "-d", "api"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestImportArgs_Unsupported(t *testing.T) {
	ds := global.DockerSecrets{}
	_, ok := compose.ImportArgs("redis", "redis-latest", "myapp", ds)
	if ok {
		t.Error("expected ok=false for unsupported image")
	}
}

func TestImportArgs_EmptyPassword(t *testing.T) {
	ds := global.DockerSecrets{"mysql-latest": {}}
	got, ok := compose.ImportArgs("mysql", "mysql-latest", "proj", ds)
	if !ok {
		t.Fatal("expected ok=true even with empty password")
	}
	want := []string{"mysql", "-u", "root", "-p", "proj"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
