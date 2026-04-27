package ports_test

import (
	"testing"

	"github.com/jcombee/devenv/internal/ports"
)

func TestParseVersionSuffix(t *testing.T) {
	cases := []struct {
		tag  string
		want int
	}{
		{"8.0", 80},
		{"5.7", 57},
		{"7", 70},
		{"16", 16},
		{"8.11", 81},
		{"8.4", 84},
		{"latest", 0},
		{"", 0},
		{"abc", 0},
		{"v1.7", 17},
	}
	for _, tc := range cases {
		got := ports.ParseVersionSuffix(tc.tag)
		if got != tc.want {
			t.Errorf("ParseVersionSuffix(%q) = %d, want %d", tc.tag, got, tc.want)
		}
	}
}

func TestCompute(t *testing.T) {
	cases := []struct {
		image string
		tag   string
		want  int
	}{
		{"mysql", "8.0", 3380},
		{"mysql", "5.7", 3357},
		{"mysql", "latest", 3300},
		{"postgres", "16", 5416},
		{"postgres", "latest", 5400},
		{"redis", "7", 6370},
		{"redis", "latest", 6300},
		{"elasticsearch", "8.11", 9281},
		{"mongo", "latest", 27000},
		{"mongo", "7", 27070},
		{"rabbitmq", "3.12", 5631},
		{"memcached", "latest", 11200},
	}
	for _, tc := range cases {
		got, err := ports.Compute(tc.image, tc.tag)
		if err != nil {
			t.Errorf("Compute(%q, %q) unexpected error: %v", tc.image, tc.tag, err)
			continue
		}
		if got != tc.want {
			t.Errorf("Compute(%q, %q) = %d, want %d", tc.image, tc.tag, got, tc.want)
		}
	}
}

func TestCompute_UnknownImage(t *testing.T) {
	_, err := ports.Compute("unknownimage", "latest")
	if err == nil {
		t.Error("expected error for unknown image, got nil")
	}
}

func TestResolve_ComputedDefault(t *testing.T) {
	result, err := ports.Resolve("mysql", "8.0", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Port != 3380 {
		t.Errorf("port: got %d want 3380", result.Port)
	}
	if result.Dedicated {
		t.Error("dedicated should be false for computed default")
	}
}

func TestResolve_GlobalOverride(t *testing.T) {
	globalPorts := map[string]int{"mysql-8-0": 13380}
	result, err := ports.Resolve("mysql", "8.0", globalPorts, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Port != 13380 {
		t.Errorf("port: got %d want 13380", result.Port)
	}
	if result.Dedicated {
		t.Error("dedicated should be false for global override")
	}
}

func TestResolve_ProjectOverride(t *testing.T) {
	p := 13306
	result, err := ports.Resolve("mysql", "8.0", nil, &p)
	if err != nil {
		t.Fatal(err)
	}
	if result.Port != 13306 {
		t.Errorf("port: got %d want 13306", result.Port)
	}
	if !result.Dedicated {
		t.Error("dedicated should be true when project port override is set")
	}
}

func TestResolve_ProjectOverrideWinsOverGlobal(t *testing.T) {
	globalPorts := map[string]int{"mysql-8-0": 13380}
	p := 23306
	result, err := ports.Resolve("mysql", "8.0", globalPorts, &p)
	if err != nil {
		t.Fatal(err)
	}
	if result.Port != 23306 {
		t.Errorf("port: got %d want 23306", result.Port)
	}
	if !result.Dedicated {
		t.Error("dedicated should be true when project port override is set")
	}
}

func TestResolve_InvalidPort(t *testing.T) {
	cases := []int{0, -1, 65536, 100000}
	for _, p := range cases {
		pp := p
		_, err := ports.Resolve("mysql", "8.0", nil, &pp)
		if err == nil {
			t.Errorf("expected error for invalid port %d", p)
		}
	}
}

func TestResolve_UnknownImage(t *testing.T) {
	_, err := ports.Resolve("unknownimage", "latest", nil, nil)
	if err == nil {
		t.Error("expected error for unknown image")
	}
}
