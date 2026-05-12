package ports

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jcombee/devenv/internal/services"
)

// Result holds the resolved host port and whether the service must be dedicated.
type Result struct {
	Port      int
	Dedicated bool
}

// ParseVersionSuffix derives a 2-digit integer suffix from a semver-style tag.
// "8.0"→80, "5.7"→57, "7"→70, "16"→16, "8.11"→81, "v1.7"→17, "latest"/""→0.
func ParseVersionSuffix(tag string) int {
	tag = strings.TrimPrefix(tag, "v")
	parts := strings.SplitN(tag, ".", 2)
	major := parts[0]
	minor := "0"
	if len(parts) == 2 {
		minor = parts[1]
	}
	if _, err := strconv.Atoi(major); err != nil {
		return 0
	}
	combined := major + minor
	if len(combined) > 2 {
		combined = combined[:2]
	}
	n, err := strconv.Atoi(combined)
	if err != nil {
		return 0
	}
	return n
}

// Compute returns the deterministic host port for an image+tag pair.
// Formula: (basePort/100)*100 + versionSuffix.
func Compute(image, tag string) (int, error) {
	svc, ok := services.Find(image)
	if !ok {
		return 0, fmt.Errorf("unknown service: %s", image)
	}
	suffix := ParseVersionSuffix(tag)
	return (svc.Port/100)*100 + suffix, nil
}

// Resolve returns the host port and dedicated flag for a service, applying the
// priority chain: project override > global override > computed default.
// globalPorts is keyed by docker-compose service name (e.g. "mysql-8-0").
func Resolve(image, tag string, globalPorts map[string]int, projectPort *int) (Result, error) {
	if projectPort != nil {
		if *projectPort < 1 || *projectPort > 65535 {
			return Result{}, fmt.Errorf("service %q: port %d out of range (1–65535)", image, *projectPort)
		}
		return Result{Port: *projectPort, Dedicated: true}, nil
	}
	name := composeServiceName(image, tag)
	if p, ok := globalPorts[name]; ok {
		return Result{Port: p, Dedicated: false}, nil
	}
	p, err := Compute(image, tag)
	if err != nil {
		return Result{}, err
	}
	return Result{Port: p, Dedicated: false}, nil
}

// composeServiceName mirrors compose.ServiceName without importing that package.
func composeServiceName(image, tag string) string {
	if tag == "" {
		tag = "latest"
	}
	return image + "-" + strings.ReplaceAll(tag, ".", "-")
}
