package services_test

import (
	"testing"

	"github.com/jcombee/devenv/internal/services"
)

func TestFind_KnownService(t *testing.T) {
	svc, ok := services.Find("mysql")
	if !ok {
		t.Fatal("mysql not found")
	}
	if svc.Port != 3306 {
		t.Errorf("mysql port: got %d want 3306", svc.Port)
	}
	if !svc.CanBeShared {
		t.Error("mysql should be shareable")
	}
}

func TestFind_UnknownService(t *testing.T) {
	_, ok := services.Find("unknownimage")
	if ok {
		t.Error("should not find unknown service")
	}
}

func TestIsSupported(t *testing.T) {
	if !services.IsSupported("redis") {
		t.Error("redis should be supported")
	}
	if services.IsSupported("fakedb") {
		t.Error("fakedb should not be supported")
	}
}

func TestMemcached_DedicatedOnly(t *testing.T) {
	svc, ok := services.Find("memcached")
	if !ok {
		t.Fatal("memcached not found")
	}
	if svc.CanBeShared {
		t.Error("memcached should be dedicated-only")
	}
}

func TestAllServicesHavePort(t *testing.T) {
	for _, svc := range services.All {
		if svc.Port == 0 {
			t.Errorf("service %q has no port", svc.Image)
		}
	}
}
