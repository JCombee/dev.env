package compose

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/jcombee/devenv/internal/global"
)

// ContainerEnv returns the environment variables needed to start the container
// for image. Any required secrets are ensured in ds (caller must save ds after).
func ContainerEnv(image, composeName string, ds global.DockerSecrets) map[string]string {
	switch image {
	case "mysql", "mariadb", "percona":
		pass := global.EnsureDockerSecret(ds, composeName, "root_password", randHex(16))
		return map[string]string{
			"MYSQL_ROOT_PASSWORD": pass,
		}

	case "postgres":
		pass := global.EnsureDockerSecret(ds, composeName, "root_password", randHex(16))
		return map[string]string{
			"POSTGRES_PASSWORD": pass,
		}

	case "mongo":
		pass := global.EnsureDockerSecret(ds, composeName, "root_password", randHex(16))
		return map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": pass,
		}

	case "elasticsearch":
		return map[string]string{
			"discovery.type": "single-node",
			"ES_JAVA_OPTS":   "-Xms512m -Xmx512m",
			"xpack.security.enabled": "false",
		}

	case "opensearch":
		return map[string]string{
			"discovery.type":                      "single-node",
			"OPENSEARCH_JAVA_OPTS":                "-Xms512m -Xmx512m",
			"plugins.security.disabled":           "true",
			"DISABLE_INSTALL_DEMO_CONFIG":         "true",
		}

	case "meilisearch":
		key := global.EnsureDockerSecret(ds, composeName, "master_key", randHex(16))
		return map[string]string{
			"MEILI_MASTER_KEY": key,
		}

	case "minio":
		user := global.EnsureDockerSecret(ds, composeName, "root_user", randHex(8))
		pass := global.EnsureDockerSecret(ds, composeName, "root_password", randHex(16))
		return map[string]string{
			"MINIO_ROOT_USER":     user,
			"MINIO_ROOT_PASSWORD": pass,
		}

	case "rabbitmq":
		pass := global.EnsureDockerSecret(ds, composeName, "admin_password", randHex(16))
		return map[string]string{
			"RABBITMQ_DEFAULT_USER": "admin",
			"RABBITMQ_DEFAULT_PASS": pass,
		}
	}

	return nil
}

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("compose: crypto/rand failed: " + err.Error())
	}
	return hex.EncodeToString(b)
}
