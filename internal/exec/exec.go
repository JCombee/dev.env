package exec

import (
	"github.com/jcombee/devenv/internal/global"
	"github.com/jcombee/devenv/internal/project"
)

// BuildArgs returns the command arguments to run inside the container for the
// given image, with credentials and project context pre-filled.
func BuildArgs(image, composeName string, ds global.DockerSecrets, ps project.Secrets, projectName string) []string {
	creds := ds[composeName]

	switch image {
	case "mysql", "mariadb", "percona":
		pass := creds["root_password"]
		return []string{"mysql", "-u", "root", "-p" + pass, projectName}

	case "postgres":
		pass := creds["root_password"]
		return []string{"env", "PGPASSWORD=" + pass, "psql", "-U", "postgres", "-d", projectName}

	case "mongo":
		pass := creds["root_password"]
		return []string{"mongosh", "--username", "root", "--password", pass}

	case "redis":
		if ps["redis"] != nil {
			if idx := ps["redis"]["db_index"]; idx != "" {
				return []string{"redis-cli", "-n", idx}
			}
		}
		return []string{"redis-cli"}
	}

	return []string{"bash"}
}
