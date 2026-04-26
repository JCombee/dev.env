package secrets_test

import (
	"testing"

	"github.com/jcombee/devenv/internal/project"
	"github.com/jcombee/devenv/internal/secrets"
)

func TestGeneratePassword_Length(t *testing.T) {
	p := secrets.GeneratePassword(16)
	if len(p) != 32 { // 16 bytes = 32 hex chars
		t.Errorf("length: got %d want 32", len(p))
	}
}

func TestGeneratePassword_Unique(t *testing.T) {
	a := secrets.GeneratePassword(16)
	b := secrets.GeneratePassword(16)
	if a == b {
		t.Error("two generated passwords should not be equal")
	}
}

func TestEnsureField_GeneratesOnce(t *testing.T) {
	sec := project.Secrets{}
	v1 := secrets.EnsureField(sec, "mysql", "password")
	v2 := secrets.EnsureField(sec, "mysql", "password")
	if v1 != v2 {
		t.Error("EnsureField should return same value on second call")
	}
	if v1 == "" {
		t.Error("generated value should not be empty")
	}
}

func TestEnsureStatic(t *testing.T) {
	sec := project.Secrets{}
	v1 := secrets.EnsureStatic(sec, "mysql", "database", "myapp")
	v2 := secrets.EnsureStatic(sec, "mysql", "database", "other") // should not overwrite
	if v1 != "myapp" || v2 != "myapp" {
		t.Errorf("EnsureStatic: got %q / %q, want myapp/myapp", v1, v2)
	}
}
