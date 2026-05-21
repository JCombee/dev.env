package compose

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jcombee/devenv/internal/global"
)

// ImportArgs returns the container command args to pipe a SQL dump into the project database.
// Returns nil, false for images that don't support SQL import.
func ImportArgs(image, composeName, projectName string, ds global.DockerSecrets) ([]string, bool) {
	creds := ds[composeName]
	switch image {
	case "mysql", "mariadb", "percona":
		pass := creds["root_password"]
		return []string{"mysql", "-u", "root", "-p" + pass, projectName}, true
	case "postgres":
		pass := creds["root_password"]
		return []string{"env", "PGPASSWORD=" + pass, "psql", "-U", "postgres", "-d", projectName}, true
	}
	return nil, false
}

// ImportDB pipes filePath into the project database inside the running container.
func ImportDB(composeFile, composeName, image, projectName, filePath string, ds global.DockerSecrets) error {
	containerArgs, ok := ImportArgs(image, composeName, projectName, ds)
	if !ok {
		return fmt.Errorf("image %q does not support SQL import", image)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	args := append([]string{"compose", "-f", composeFile, "exec", "-T", composeName}, containerArgs...)
	cmd := exec.Command("docker", args...)
	cmd.Stdin = f
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
