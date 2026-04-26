package compose_test

import (
	"testing"

	"github.com/jcombee/devenv/internal/compose"
)

func TestServiceName(t *testing.T) {
	cases := []struct {
		image, tag, want string
	}{
		{"mysql", "8.0", "mysql-8-0"},
		{"redis", "latest", "redis-latest"},
		{"redis", "", "redis-latest"},
		{"elasticsearch", "8.11", "elasticsearch-8-11"},
		{"postgres", "15.2", "postgres-15-2"},
	}
	for _, c := range cases {
		got := compose.ServiceName(c.image, c.tag)
		if got != c.want {
			t.Errorf("ServiceName(%q,%q) = %q, want %q", c.image, c.tag, got, c.want)
		}
	}
}
