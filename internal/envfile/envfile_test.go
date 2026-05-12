package envfile_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jcombee/devenv/internal/envfile"
	"github.com/jcombee/devenv/internal/global"
	"github.com/jcombee/devenv/internal/project"
)

// --- ServiceVars tests ---

func TestServiceVars_MySQL(t *testing.T) {
	ds := global.DockerSecrets{"mysql-8-0": {"root_password": "secret123"}}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("mysql", "mysql-8-0", 3380, "myproject", ds, ps)

	expect := map[string]string{
		"DB_HOST":     "127.0.0.1",
		"DB_PORT":     "3380",
		"DB_DATABASE": "myproject",
		"DB_USERNAME": "root",
		"DB_PASSWORD": "secret123",
	}
	assertVars(t, vars, expect)
}

func TestServiceVars_MariaDB(t *testing.T) {
	ds := global.DockerSecrets{"mariadb-latest": {"root_password": "pw"}}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("mariadb", "mariadb-latest", 3300, "proj", ds, ps)
	if vars["DB_USERNAME"] != "root" {
		t.Errorf("DB_USERNAME: got %q want root", vars["DB_USERNAME"])
	}
	if vars["DB_PASSWORD"] != "pw" {
		t.Errorf("DB_PASSWORD: got %q want pw", vars["DB_PASSWORD"])
	}
}

func TestServiceVars_Percona(t *testing.T) {
	ds := global.DockerSecrets{"percona-latest": {"root_password": "pw"}}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("percona", "percona-latest", 3300, "proj", ds, ps)
	if vars["DB_PASSWORD"] != "pw" {
		t.Errorf("DB_PASSWORD: got %q want pw", vars["DB_PASSWORD"])
	}
}

func TestServiceVars_Postgres(t *testing.T) {
	ds := global.DockerSecrets{"postgres-16": {"root_password": "pgpass"}}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("postgres", "postgres-16", 5416, "myproject", ds, ps)

	expect := map[string]string{
		"DB_HOST":     "127.0.0.1",
		"DB_PORT":     "5416",
		"DB_DATABASE": "myproject",
		"DB_USERNAME": "postgres",
		"DB_PASSWORD": "pgpass",
	}
	assertVars(t, vars, expect)
}

func TestServiceVars_Mongo(t *testing.T) {
	ds := global.DockerSecrets{"mongo-latest": {"root_password": "mongopw"}}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("mongo", "mongo-latest", 27000, "myproject", ds, ps)

	if vars["MONGODB_HOST"] != "127.0.0.1" {
		t.Errorf("MONGODB_HOST: got %q", vars["MONGODB_HOST"])
	}
	if vars["MONGODB_PORT"] != "27000" {
		t.Errorf("MONGODB_PORT: got %q", vars["MONGODB_PORT"])
	}
	if vars["MONGODB_DATABASE"] != "myproject" {
		t.Errorf("MONGODB_DATABASE: got %q", vars["MONGODB_DATABASE"])
	}
	if vars["MONGODB_USERNAME"] != "root" {
		t.Errorf("MONGODB_USERNAME: got %q", vars["MONGODB_USERNAME"])
	}
	if vars["MONGODB_PASSWORD"] != "mongopw" {
		t.Errorf("MONGODB_PASSWORD: got %q", vars["MONGODB_PASSWORD"])
	}
	uri := vars["MONGODB_URI"]
	if !strings.Contains(uri, "mongopw") || !strings.Contains(uri, "27000") {
		t.Errorf("MONGODB_URI malformed: %q", uri)
	}
}

func TestServiceVars_Redis(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{"redis": {"db_index": "3"}}
	vars := envfile.ServiceVars("redis", "redis-latest", 6300, "myproject", ds, ps)

	expect := map[string]string{
		"REDIS_HOST": "127.0.0.1",
		"REDIS_PORT": "6300",
		"REDIS_DB":   "3",
	}
	assertVars(t, vars, expect)
}

func TestServiceVars_Memcached(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("memcached", "memcached-latest", 11200, "proj", ds, ps)
	assertVars(t, vars, map[string]string{
		"MEMCACHED_HOST": "127.0.0.1",
		"MEMCACHED_PORT": "11200",
	})
}

func TestServiceVars_Elasticsearch(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("elasticsearch", "elasticsearch-8-11", 9281, "myproject", ds, ps)

	if vars["ELASTICSEARCH_HOST"] != "http://127.0.0.1:9281" {
		t.Errorf("ELASTICSEARCH_HOST: got %q", vars["ELASTICSEARCH_HOST"])
	}
	if vars["ELASTICSEARCH_INDEX"] != "myproject" {
		t.Errorf("ELASTICSEARCH_INDEX: got %q", vars["ELASTICSEARCH_INDEX"])
	}
}

func TestServiceVars_OpenSearch(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("opensearch", "opensearch-latest", 9200, "proj", ds, ps)
	if _, ok := vars["ELASTICSEARCH_HOST"]; !ok {
		t.Error("opensearch should emit ELASTICSEARCH_HOST")
	}
}

func TestServiceVars_Meilisearch(t *testing.T) {
	ds := global.DockerSecrets{"meilisearch-latest": {"master_key": "masterkey123"}}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("meilisearch", "meilisearch-latest", 7700, "myproject", ds, ps)
	assertVars(t, vars, map[string]string{
		"MEILISEARCH_HOST":  "http://127.0.0.1:7700",
		"MEILISEARCH_KEY":   "masterkey123",
		"MEILISEARCH_INDEX": "myproject",
	})
}

func TestServiceVars_Typesense(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{"typesense": {"api_key": "tskey"}}
	vars := envfile.ServiceVars("typesense", "typesense-latest", 8108, "myproject", ds, ps)
	assertVars(t, vars, map[string]string{
		"TYPESENSE_HOST":       "127.0.0.1",
		"TYPESENSE_PORT":       "8108",
		"TYPESENSE_PROTOCOL":   "http",
		"TYPESENSE_API_KEY":    "tskey",
		"TYPESENSE_COLLECTION": "myproject",
	})
}

func TestServiceVars_Solr(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("solr", "solr-latest", 8983, "myproject", ds, ps)
	assertVars(t, vars, map[string]string{
		"SOLR_HOST": "127.0.0.1",
		"SOLR_PORT": "8983",
		"SOLR_CORE": "myproject",
	})
}

func TestServiceVars_Cassandra(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("cassandra", "cassandra-latest", 9042, "my-project", ds, ps)
	assertVars(t, vars, map[string]string{
		"CASSANDRA_HOST":     "127.0.0.1",
		"CASSANDRA_PORT":     "9042",
		"CASSANDRA_KEYSPACE": "my_project",
	})
}

func TestServiceVars_RabbitMQ(t *testing.T) {
	ds := global.DockerSecrets{"rabbitmq-latest": {"admin_password": "rabbitpw"}}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("rabbitmq", "rabbitmq-latest", 5600, "myproject", ds, ps)
	assertVars(t, vars, map[string]string{
		"RABBITMQ_HOST":     "127.0.0.1",
		"RABBITMQ_PORT":     "5600",
		"RABBITMQ_VHOST":    "myproject",
		"RABBITMQ_USERNAME": "admin",
		"RABBITMQ_PASSWORD": "rabbitpw",
	})
}

func TestServiceVars_Kafka(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("kafka", "kafka-latest", 9092, "myproject", ds, ps)
	assertVars(t, vars, map[string]string{
		"KAFKA_BROKERS":      "127.0.0.1:9092",
		"KAFKA_TOPIC_PREFIX": "myproject",
	})
}

func TestServiceVars_MinIO(t *testing.T) {
	ds := global.DockerSecrets{"minio-latest": {"root_user": "miniouser", "root_password": "miniopw"}}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("minio", "minio-latest", 9000, "myproject", ds, ps)
	assertVars(t, vars, map[string]string{
		"MINIO_ENDPOINT":   "http://127.0.0.1:9000",
		"MINIO_ACCESS_KEY": "miniouser",
		"MINIO_SECRET_KEY": "miniopw",
		"MINIO_BUCKET":     "myproject",
		"MINIO_USE_SSL":    "false",
	})
}

func TestServiceVars_Soketi(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{"soketi": {"app_id": "123", "app_key": "key", "app_secret": "secret"}}
	vars := envfile.ServiceVars("soketi", "soketi-latest", 6001, "myproject", ds, ps)
	assertVars(t, vars, map[string]string{
		"PUSHER_HOST":       "127.0.0.1",
		"PUSHER_PORT":       "6001",
		"PUSHER_SCHEME":     "http",
		"PUSHER_APP_ID":     "123",
		"PUSHER_APP_KEY":    "key",
		"PUSHER_APP_SECRET": "secret",
	})
}

func TestServiceVars_Reverb(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{"reverb": {"app_id": "456", "app_key": "rkey", "app_secret": "rsecret"}}
	vars := envfile.ServiceVars("reverb", "reverb-latest", 8080, "myproject", ds, ps)
	assertVars(t, vars, map[string]string{
		"REVERB_HOST":       "127.0.0.1",
		"REVERB_PORT":       "8080",
		"REVERB_SCHEME":     "http",
		"REVERB_APP_ID":     "456",
		"REVERB_APP_KEY":    "rkey",
		"REVERB_APP_SECRET": "rsecret",
	})
}

func TestServiceVars_Mailpit(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("mailpit", "mailpit-latest", 1025, "proj", ds, ps)
	assertVars(t, vars, map[string]string{
		"MAIL_HOST":   "127.0.0.1",
		"MAIL_PORT":   "1025",
		"MAIL_MAILER": "smtp",
	})
}

func TestServiceVars_Mailhog(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("mailhog", "mailhog-latest", 1025, "proj", ds, ps)
	if vars["MAIL_MAILER"] != "smtp" {
		t.Errorf("MAIL_MAILER: got %q want smtp", vars["MAIL_MAILER"])
	}
}

func TestServiceVars_UnknownImage_ReturnsEmpty(t *testing.T) {
	ds := global.DockerSecrets{}
	ps := project.Secrets{}
	vars := envfile.ServiceVars("unknownimage", "unknownimage-latest", 9999, "proj", ds, ps)
	if len(vars) != 0 {
		t.Errorf("expected empty vars for unknown image, got %v", vars)
	}
}

// --- ApplyEnvMap tests ---

func TestApplyEnvMap_RenamesCanonicalKeys(t *testing.T) {
	vars := map[string]string{
		"DB_HOST":     "127.0.0.1",
		"DB_DATABASE": "myproject",
		"DB_PASSWORD": "secret",
	}
	envMap := map[string]string{
		"DB_HOST":     "DATABASE_HOST",
		"DB_DATABASE": "DATABASE_NAME",
	}
	result := envfile.ApplyEnvMap(vars, envMap)

	if _, ok := result["DB_HOST"]; ok {
		t.Error("DB_HOST should have been renamed and removed")
	}
	if result["DATABASE_HOST"] != "127.0.0.1" {
		t.Errorf("DATABASE_HOST: got %q want 127.0.0.1", result["DATABASE_HOST"])
	}
	if result["DATABASE_NAME"] != "myproject" {
		t.Errorf("DATABASE_NAME: got %q want myproject", result["DATABASE_NAME"])
	}
	// Unmapped key passes through.
	if result["DB_PASSWORD"] != "secret" {
		t.Errorf("DB_PASSWORD should pass through unchanged, got %q", result["DB_PASSWORD"])
	}
}

func TestApplyEnvMap_NilMap_NoChange(t *testing.T) {
	vars := map[string]string{"DB_HOST": "127.0.0.1"}
	result := envfile.ApplyEnvMap(vars, nil)
	if result["DB_HOST"] != "127.0.0.1" {
		t.Error("nil env_map should not change vars")
	}
}

func TestApplyEnvMap_EmptyMap_NoChange(t *testing.T) {
	vars := map[string]string{"REDIS_HOST": "127.0.0.1"}
	result := envfile.ApplyEnvMap(vars, map[string]string{})
	if result["REDIS_HOST"] != "127.0.0.1" {
		t.Error("empty env_map should not change vars")
	}
}

func TestApplyEnvMap_UnknownCanonical_Ignored(t *testing.T) {
	vars := map[string]string{"DB_HOST": "127.0.0.1"}
	envMap := map[string]string{"NONEXISTENT": "SOMETHING"}
	result := envfile.ApplyEnvMap(vars, envMap)
	if result["DB_HOST"] != "127.0.0.1" {
		t.Error("unknown canonical key in env_map should not affect existing vars")
	}
}

// --- Write tests ---

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	vars := map[string]string{"DB_HOST": "127.0.0.1", "DB_PORT": "3380"}

	if err := envfile.Write(path, vars); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf(".env not created: %v", err)
	}
}

func TestWrite_ContainsAllVars(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	vars := map[string]string{
		"DB_HOST":     "127.0.0.1",
		"DB_PORT":     "3380",
		"DB_PASSWORD": "secret",
	}

	if err := envfile.Write(path, vars); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	content := string(data)

	for k, v := range vars {
		expected := k + "=" + v
		if !strings.Contains(content, expected) {
			t.Errorf(".env missing %q\n%s", expected, content)
		}
	}
}

func TestWrite_AlphabeticalOrder(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	vars := map[string]string{"Z_VAR": "z", "A_VAR": "a", "M_VAR": "m"}

	if err := envfile.Write(path, vars); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	content := string(data)

	aIdx := strings.Index(content, "A_VAR")
	mIdx := strings.Index(content, "M_VAR")
	zIdx := strings.Index(content, "Z_VAR")
	if !(aIdx < mIdx && mIdx < zIdx) {
		t.Errorf("vars not in alphabetical order:\n%s", content)
	}
}

func TestWrite_Overwrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := envfile.Write(path, map[string]string{"OLD": "value"}); err != nil {
		t.Fatal(err)
	}
	if err := envfile.Write(path, map[string]string{"NEW": "value"}); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	content := string(data)
	if strings.Contains(content, "OLD=") {
		t.Error(".env should be overwritten, old vars should not remain")
	}
	if !strings.Contains(content, "NEW=value") {
		t.Error(".env should contain new vars after overwrite")
	}
}

func TestWrite_HasHeaderComment(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := envfile.Write(path, map[string]string{"X": "y"}); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if !strings.HasPrefix(string(data), "#") {
		t.Error(".env should start with a header comment")
	}
}

// --- EnsureProjectSecrets tests ---

func TestEnsureProjectSecrets_Redis_AssignsIndex(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	ps := project.Secrets{}
	ds := global.DockerSecrets{}

	if err := envfile.EnsureProjectSecrets("redis", "redis-latest", "myproject", ps, ds); err != nil {
		t.Fatal(err)
	}
	idx, ok := ps["redis"]["db_index"]
	if !ok || idx == "" {
		t.Error("redis db_index not assigned")
	}
}

func TestEnsureProjectSecrets_Redis_IdempotentIndex(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	ps := project.Secrets{"redis": {"db_index": "5"}}
	ds := global.DockerSecrets{}

	if err := envfile.EnsureProjectSecrets("redis", "redis-latest", "myproject", ps, ds); err != nil {
		t.Fatal(err)
	}
	if ps["redis"]["db_index"] != "5" {
		t.Errorf("existing db_index should not be overwritten, got %q", ps["redis"]["db_index"])
	}
}

func TestEnsureProjectSecrets_Typesense_GeneratesKey(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	ps := project.Secrets{}
	ds := global.DockerSecrets{}

	if err := envfile.EnsureProjectSecrets("typesense", "typesense-latest", "proj", ps, ds); err != nil {
		t.Fatal(err)
	}
	key, ok := ps["typesense"]["api_key"]
	if !ok || key == "" {
		t.Error("typesense api_key not generated")
	}
}

func TestEnsureProjectSecrets_Soketi_GeneratesCredentials(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	ps := project.Secrets{}
	ds := global.DockerSecrets{}

	if err := envfile.EnsureProjectSecrets("soketi", "soketi-latest", "proj", ps, ds); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"app_id", "app_key", "app_secret"} {
		if ps["soketi"][field] == "" {
			t.Errorf("soketi %s not generated", field)
		}
	}
}

func TestEnsureProjectSecrets_Reverb_GeneratesCredentials(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	ps := project.Secrets{}
	ds := global.DockerSecrets{}

	if err := envfile.EnsureProjectSecrets("reverb", "reverb-latest", "proj", ps, ds); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"app_id", "app_key", "app_secret"} {
		if ps["reverb"][field] == "" {
			t.Errorf("reverb %s not generated", field)
		}
	}
}

func TestEnsureProjectSecrets_MySQL_NoOp(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	ps := project.Secrets{}
	ds := global.DockerSecrets{}

	// MySQL uses global docker secrets — EnsureProjectSecrets should not add anything to ps.
	if err := envfile.EnsureProjectSecrets("mysql", "mysql-8-0", "proj", ps, ds); err != nil {
		t.Fatal(err)
	}
	if len(ps) != 0 {
		t.Errorf("mysql should not modify project secrets, got %v", ps)
	}
}

// --- helpers ---

func assertVars(t *testing.T, got, want map[string]string) {
	t.Helper()
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s: got %q want %q", k, got[k], v)
		}
	}
}
