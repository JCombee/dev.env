package global

import (
	"os"
	"path/filepath"

	"github.com/jcombee/devenv/internal/store"
)

// DockerSecrets stores admin/root credentials needed to start containers.
// Key: compose service name (e.g. "mysql-8-0"), value: credential fields.
type DockerSecrets map[string]map[string]string

func LoadDockerSecrets() (DockerSecrets, error) {
	path := filepath.Join(Dir(), "docker", "secrets.yaml")
	var ds DockerSecrets
	err := store.Read(path, &ds)
	if os.IsNotExist(err) {
		return DockerSecrets{}, nil
	}
	if err != nil {
		return nil, err
	}
	if ds == nil {
		ds = DockerSecrets{}
	}
	return ds, nil
}

func SaveDockerSecrets(ds DockerSecrets) error {
	path := filepath.Join(Dir(), "docker", "secrets.yaml")
	return store.Write(path, ds)
}

// EnsureDockerSecret ensures ds[composeName][field] is set, generating a
// random value if absent. Returns the (possibly existing) value.
func EnsureDockerSecret(ds DockerSecrets, composeName, field, generated string) string {
	if ds[composeName] == nil {
		ds[composeName] = map[string]string{}
	}
	if v, ok := ds[composeName][field]; ok && v != "" {
		return v
	}
	ds[composeName][field] = generated
	return generated
}
