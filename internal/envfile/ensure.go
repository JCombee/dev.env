package envfile

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jcombee/devenv/internal/global"
	"github.com/jcombee/devenv/internal/project"
	"github.com/jcombee/devenv/internal/store"
)

// EnsureProjectSecrets generates and stores per-project credentials for services
// that require them (redis DB index, typesense API key, soketi/reverb app credentials).
// Services that use global docker secrets (mysql, postgres, etc.) are no-ops here.
func EnsureProjectSecrets(image, composeName, projectName string, ps project.Secrets, ds global.DockerSecrets) error {
	switch image {
	case "redis":
		return ensureRedisIndex(projectName, ps)
	case "typesense":
		ensureField(ps, "typesense", "api_key", randHex(16))
	case "soketi":
		ensureField(ps, "soketi", "app_id", strconv.Itoa(rand.Intn(90000)+10000))
		ensureField(ps, "soketi", "app_key", randHex(16))
		ensureField(ps, "soketi", "app_secret", randHex(16))
	case "reverb":
		ensureField(ps, "reverb", "app_id", strconv.Itoa(rand.Intn(90000)+10000))
		ensureField(ps, "reverb", "app_key", randHex(16))
		ensureField(ps, "reverb", "app_secret", randHex(16))
	}
	return nil
}

// ensureRedisIndex assigns a unique Redis DB index (0–15) for this project
// by scanning all known project secrets files.
func ensureRedisIndex(projectName string, ps project.Secrets) error {
	if ps["redis"] != nil && ps["redis"]["db_index"] != "" {
		return nil
	}

	used, err := usedRedisIndices(projectName)
	if err != nil {
		return fmt.Errorf("scan redis indices: %w", err)
	}

	for i := 0; i <= 15; i++ {
		if !used[i] {
			if ps["redis"] == nil {
				ps["redis"] = map[string]string{}
			}
			ps["redis"]["db_index"] = strconv.Itoa(i)
			return nil
		}
	}
	// All 16 indices in use — fall back to index 0.
	if ps["redis"] == nil {
		ps["redis"] = map[string]string{}
	}
	ps["redis"]["db_index"] = "0"
	return nil
}

// usedRedisIndices returns the set of Redis DB indices already assigned to
// other projects by scanning ~/.dev.env/projects/*/secrets.yaml.
func usedRedisIndices(excludeProject string) (map[int]bool, error) {
	used := map[int]bool{}
	projectsDir := filepath.Join(global.Dir(), "projects")

	entries, err := os.ReadDir(projectsDir)
	if os.IsNotExist(err) {
		return used, nil
	}
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if !e.IsDir() || e.Name() == excludeProject {
			continue
		}
		var sec project.Secrets
		path := filepath.Join(projectsDir, e.Name(), "secrets.yaml")
		if err := store.Read(path, &sec); err != nil {
			continue
		}
		if sec["redis"] != nil {
			if idx, err := strconv.Atoi(sec["redis"]["db_index"]); err == nil {
				used[idx] = true
			}
		}
	}
	return used, nil
}

func ensureField(ps project.Secrets, image, field, generated string) {
	if ps[image] == nil {
		ps[image] = map[string]string{}
	}
	if ps[image][field] == "" {
		ps[image][field] = generated
	}
}

func randHex(n int) string {
	const chars = "0123456789abcdef"
	b := make([]byte, n*2)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
