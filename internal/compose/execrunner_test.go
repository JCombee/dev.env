package compose_test

import (
	"testing"

	"github.com/jcombee/devenv/internal/compose"
)

func TestParsePS_JSONArray(t *testing.T) {
	data := []byte(`[{"Name":"mysql-8-0","State":"running","Status":"Up 2 hours"},{"Name":"redis-latest","State":"exited","Status":"Exited (0)"}]`)
	result, err := compose.ParsePS(data)
	if err != nil {
		t.Fatalf("ParsePS: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("want 2 entries, got %d", len(result))
	}
	if result[0].Name != "mysql-8-0" || result[0].State != "running" {
		t.Errorf("entry[0]: got %+v", result[0])
	}
	if result[1].Name != "redis-latest" || result[1].State != "exited" {
		t.Errorf("entry[1]: got %+v", result[1])
	}
}

func TestParsePS_JSONL(t *testing.T) {
	data := []byte(`{"Name":"mysql-8-0","State":"running","Status":"Up 2 hours"}
{"Name":"redis-latest","State":"exited","Status":"Exited (0)"}`)
	result, err := compose.ParsePS(data)
	if err != nil {
		t.Fatalf("ParsePS: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("want 2 entries, got %d", len(result))
	}
}

func TestParsePS_Empty(t *testing.T) {
	result, err := compose.ParsePS([]byte(""))
	if err != nil {
		t.Fatalf("ParsePS empty: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("want 0 entries, got %d", len(result))
	}
}
