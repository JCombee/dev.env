package secrets

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/jcombee/devenv/internal/project"
)

// GeneratePassword returns a random hex string of n bytes (2n chars).
func GeneratePassword(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("secrets: crypto/rand failed: " + err.Error())
	}
	return hex.EncodeToString(b)
}

// EnsureField sets secrets[image][field] to a generated password if not already set.
// Returns the (possibly existing) value. Callers must SaveSecrets after calling.
func EnsureField(sec project.Secrets, image, field string) string {
	if sec[image] == nil {
		sec[image] = map[string]string{}
	}
	if v, ok := sec[image][field]; ok && v != "" {
		return v
	}
	v := GeneratePassword(16)
	sec[image][field] = v
	return v
}

// EnsureStatic sets secrets[image][field] to value if not already set.
// Used for deterministic values like database names (project name) or index keys.
func EnsureStatic(sec project.Secrets, image, field, value string) string {
	if sec[image] == nil {
		sec[image] = map[string]string{}
	}
	if v, ok := sec[image][field]; ok && v != "" {
		return v
	}
	sec[image][field] = value
	return value
}
