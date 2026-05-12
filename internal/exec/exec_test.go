package exec_test

import (
	"reflect"
	"testing"

	devexec "github.com/jcombee/devenv/internal/exec"
	"github.com/jcombee/devenv/internal/global"
	"github.com/jcombee/devenv/internal/project"
)

func TestBuildArgs_MySQL(t *testing.T) {
	ds := global.DockerSecrets{"mysql-8-0": {"root_password": "secret123"}}
	got := devexec.BuildArgs("mysql", "mysql-8-0", ds, project.Secrets{}, "myapp")
	want := []string{"mysql", "-u", "root", "-psecret123", "myapp"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgs_MariaDB(t *testing.T) {
	ds := global.DockerSecrets{"mariadb-latest": {"root_password": "pass456"}}
	got := devexec.BuildArgs("mariadb", "mariadb-latest", ds, project.Secrets{}, "shop")
	want := []string{"mysql", "-u", "root", "-ppass456", "shop"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgs_Percona(t *testing.T) {
	ds := global.DockerSecrets{"percona-8-0": {"root_password": "percpass"}}
	got := devexec.BuildArgs("percona", "percona-8-0", ds, project.Secrets{}, "blog")
	want := []string{"mysql", "-u", "root", "-ppercpass", "blog"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgs_Postgres(t *testing.T) {
	ds := global.DockerSecrets{"postgres-16": {"root_password": "pgpass"}}
	got := devexec.BuildArgs("postgres", "postgres-16", ds, project.Secrets{}, "api")
	want := []string{"env", "PGPASSWORD=pgpass", "psql", "-U", "postgres", "-d", "api"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgs_Mongo(t *testing.T) {
	ds := global.DockerSecrets{"mongo-latest": {"root_password": "mongopass"}}
	got := devexec.BuildArgs("mongo", "mongo-latest", ds, project.Secrets{}, "data")
	want := []string{"mongosh", "--username", "root", "--password", "mongopass"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgs_Redis(t *testing.T) {
	ps := project.Secrets{"redis": {"db_index": "3"}}
	got := devexec.BuildArgs("redis", "redis-latest", global.DockerSecrets{}, ps, "myapp")
	want := []string{"redis-cli", "-n", "3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgs_RedisNoIndex(t *testing.T) {
	got := devexec.BuildArgs("redis", "redis-latest", global.DockerSecrets{}, project.Secrets{}, "myapp")
	want := []string{"redis-cli"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgs_Unknown(t *testing.T) {
	got := devexec.BuildArgs("mailpit", "mailpit-latest", global.DockerSecrets{}, project.Secrets{}, "myapp")
	want := []string{"bash"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
