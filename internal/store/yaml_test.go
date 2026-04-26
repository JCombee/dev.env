package store_test

import (
	"path/filepath"
	"testing"

	"github.com/jcombee/devenv/internal/store"
)

type sample struct {
	Name  string `yaml:"name"`
	Value int    `yaml:"value"`
}

func TestReadWrite_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.yaml")
	in := sample{Name: "foo", Value: 42}

	if err := store.Write(path, &in); err != nil {
		t.Fatalf("Write: %v", err)
	}

	var out sample
	if err := store.Read(path, &out); err != nil {
		t.Fatalf("Read: %v", err)
	}

	if out.Name != in.Name || out.Value != in.Value {
		t.Errorf("round-trip mismatch: got %+v, want %+v", out, in)
	}
}

func TestRead_MissingFile(t *testing.T) {
	var s sample
	err := store.Read("/nonexistent/path.yaml", &s)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
