package compose

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/jcombee/devenv/internal/global"
)

// ProvisionDB creates the project-named database inside a running container.
// Retries for up to 60s to handle container startup lag.
// No-op for non-database images.
func ProvisionDB(composeFile, composeName, image, projectName string, ds global.DockerSecrets) error {
	switch image {
	case "mysql", "mariadb", "percona":
		return provisionMySQL(composeFile, composeName, projectName, ds)
	case "postgres":
		return provisionPostgres(composeFile, composeName, projectName, ds)
	}
	return nil
}

func provisionMySQL(composeFile, composeName, projectName string, ds global.DockerSecrets) error {
	pass := ds[composeName]["root_password"]
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", projectName)
	args := []string{
		"compose", "-f", composeFile, "exec", "-T", composeName,
		"mysql", "-u", "root", "-p" + pass, "-e", sql,
	}
	return retryExec(60*time.Second, args)
}

func provisionPostgres(composeFile, composeName, projectName string, ds global.DockerSecrets) error {
	pass := ds[composeName]["root_password"]
	sql := fmt.Sprintf(`CREATE DATABASE "%s"`, projectName)
	args := []string{
		"compose", "-f", composeFile, "exec", "-T", composeName,
		"env", "PGPASSWORD=" + pass, "psql", "-U", "postgres", "-c", sql,
	}
	return retryExec(60*time.Second, args)
}

func retryExec(timeout time.Duration, args []string) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		out, err := exec.Command("docker", args...).CombinedOutput()
		if err == nil || strings.Contains(string(out), "already exists") {
			return nil
		}
		lastErr = fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("provision timed out: %w", lastErr)
}
