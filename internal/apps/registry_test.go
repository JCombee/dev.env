package apps_test

import (
	"strings"
	"testing"

	"github.com/jcombee/devenv/internal/apps"
	"github.com/jcombee/devenv/internal/global"
)

func TestFind_KnownApp(t *testing.T) {
	app, ok := apps.Find("litellm")
	if !ok {
		t.Fatal("expected litellm to be a known app")
	}
	if app.Name != "litellm" {
		t.Errorf("Name: got %q, want %q", app.Name, "litellm")
	}
	if app.DefaultPort != 4000 {
		t.Errorf("DefaultPort: got %d, want 4000", app.DefaultPort)
	}
	if app.Description == "" {
		t.Error("Description must not be empty — shown by `dev app list`")
	}
}

func TestFind_UnknownApp(t *testing.T) {
	if _, ok := apps.Find("nope"); ok {
		t.Error("expected unknown app to return false")
	}
}

func TestAll_EveryAppIsUsable(t *testing.T) {
	seen := map[string]bool{}
	for _, a := range apps.All {
		if seen[a.Name] {
			t.Errorf("duplicate app name %q", a.Name)
		}
		seen[a.Name] = true
		if a.Build == nil {
			t.Errorf("%s: Build must not be nil", a.Name)
		}
		if a.Info == nil {
			t.Errorf("%s: Info must not be nil", a.Name)
		}
		if a.DefaultPort == 0 {
			t.Errorf("%s: DefaultPort must be set", a.Name)
		}
	}
}

// buildLiteLLM returns the containers of the litellm app with secrets ensured.
func buildLiteLLM(t *testing.T, hostPort int) ([]apps.Container, global.DockerSecrets) {
	t.Helper()
	app, ok := apps.Find("litellm")
	if !ok {
		t.Fatal("litellm not registered")
	}
	ds := global.DockerSecrets{}
	apps.EnsureSecrets(app, ds)
	return app.Build(apps.BuildContext{HostPort: hostPort, Secrets: ds}), ds
}

func containerByName(containers []apps.Container, name string) (apps.Container, bool) {
	for _, c := range containers {
		if c.Name == name {
			return c, true
		}
	}
	return apps.Container{}, false
}

func TestLiteLLM_BuildsProxyAndDatabase(t *testing.T) {
	containers, _ := buildLiteLLM(t, 4000)
	if len(containers) != 2 {
		t.Fatalf("expected 2 containers, got %d", len(containers))
	}
	if _, ok := containerByName(containers, "litellm"); !ok {
		t.Error("missing litellm container")
	}
	if _, ok := containerByName(containers, "litellm-db"); !ok {
		t.Error("missing litellm-db container")
	}
}

func TestLiteLLM_ProxyPublishesHostPort(t *testing.T) {
	containers, _ := buildLiteLLM(t, 4123)
	proxy, _ := containerByName(containers, "litellm")
	if proxy.HostPort != 4123 {
		t.Errorf("HostPort: got %d, want 4123 (from BuildContext)", proxy.HostPort)
	}
	if proxy.Port != 4000 {
		t.Errorf("Port: got %d, want 4000 (container-internal)", proxy.Port)
	}
	if !strings.HasPrefix(proxy.Image, "ghcr.io/berriai/litellm") {
		t.Errorf("Image: got %q, want the litellm image", proxy.Image)
	}
}

func TestLiteLLM_DatabaseIsInternalOnly(t *testing.T) {
	containers, _ := buildLiteLLM(t, 4000)
	db, _ := containerByName(containers, "litellm-db")
	if db.HostPort != 0 {
		t.Errorf("HostPort: got %d, want 0 — the app's Postgres must not be published", db.HostPort)
	}
	if db.Image != "postgres" {
		t.Errorf("Image: got %q, want postgres", db.Image)
	}
	if db.Healthcheck == nil {
		t.Fatal("Healthcheck must be set so the proxy can wait for the database")
	}
	if len(db.Healthcheck.Test) == 0 {
		t.Error("Healthcheck.Test must not be empty")
	}
}

func TestLiteLLM_DatabaseProvisionsItself(t *testing.T) {
	containers, ds := buildLiteLLM(t, 4000)
	db, _ := containerByName(containers, "litellm-db")

	// The postgres image creates db + user from these on first boot; no dev-side provisioning.
	for k, want := range map[string]string{
		"POSTGRES_DB":       "litellm",
		"POSTGRES_USER":     "litellm",
		"POSTGRES_PASSWORD": ds["litellm-db"]["root_password"],
	} {
		if got := db.Env[k]; got != want {
			t.Errorf("Env[%s]: got %q, want %q", k, got, want)
		}
	}
	if len(db.NamedVolumes) == 0 {
		t.Error("database must use a named volume so data survives disable")
	}
	found := false
	for _, v := range db.Volumes {
		if strings.HasPrefix(v, db.NamedVolumes[0]+":") {
			found = true
		}
	}
	if !found {
		t.Errorf("named volume %q not mounted: %v", db.NamedVolumes[0], db.Volumes)
	}
}

func TestLiteLLM_ProxyEnvironment(t *testing.T) {
	containers, ds := buildLiteLLM(t, 4000)
	proxy, _ := containerByName(containers, "litellm")

	master := ds["litellm"]["master_key"]
	if !strings.HasPrefix(master, "sk-") {
		t.Errorf("master_key %q must be sk- prefixed", master)
	}
	if got := proxy.Env["LITELLM_MASTER_KEY"]; got != master {
		t.Errorf("LITELLM_MASTER_KEY: got %q, want %q", got, master)
	}
	if got := proxy.Env["LITELLM_SALT_KEY"]; got != ds["litellm"]["salt_key"] {
		t.Errorf("LITELLM_SALT_KEY not wired to the generated salt_key")
	}
	if got := proxy.Env["STORE_MODEL_IN_DB"]; got != "True" {
		t.Errorf("STORE_MODEL_IN_DB: got %q, want True — models are managed in the UI", got)
	}
	if got := proxy.Env["UI_USERNAME"]; got != "admin" {
		t.Errorf("UI_USERNAME: got %q, want admin", got)
	}
	if got := proxy.Env["UI_PASSWORD"]; got != ds["litellm"]["ui_password"] {
		t.Error("UI_PASSWORD not wired to the generated ui_password")
	}

	dbURL := proxy.Env["DATABASE_URL"]
	for _, want := range []string{"postgresql://litellm:", ds["litellm-db"]["root_password"], "@litellm-db:5432/litellm"} {
		if !strings.Contains(dbURL, want) {
			t.Errorf("DATABASE_URL %q missing %q", dbURL, want)
		}
	}
}

func TestLiteLLM_ProxyWaitsForDatabase(t *testing.T) {
	containers, _ := buildLiteLLM(t, 4000)
	proxy, _ := containerByName(containers, "litellm")
	if got := proxy.DependsOn["litellm-db"]; got != "service_healthy" {
		t.Errorf("DependsOn[litellm-db]: got %q, want service_healthy", got)
	}
}

func TestLiteLLM_NoConfigFileMounted(t *testing.T) {
	containers, _ := buildLiteLLM(t, 4000)
	proxy, _ := containerByName(containers, "litellm")
	for _, v := range proxy.Volumes {
		if strings.Contains(v, "config.yaml") {
			t.Errorf("litellm must run without a config file, got volume %q", v)
		}
	}
}

func TestEnsureSecrets_GeneratesOnceAndReuses(t *testing.T) {
	app, _ := apps.Find("litellm")
	ds := global.DockerSecrets{}

	apps.EnsureSecrets(app, ds)
	first := ds["litellm"]["master_key"]
	if first == "" {
		t.Fatal("master_key not generated")
	}
	if ds["litellm-db"]["root_password"] == "" {
		t.Fatal("litellm-db root_password not generated")
	}

	apps.EnsureSecrets(app, ds)
	if ds["litellm"]["master_key"] != first {
		t.Error("EnsureSecrets overwrote an existing secret — credentials must stay stable")
	}
}

func TestLiteLLM_Info(t *testing.T) {
	app, _ := apps.Find("litellm")
	ds := global.DockerSecrets{}
	apps.EnsureSecrets(app, ds)

	lines := app.Info(apps.BuildContext{HostPort: 4000, Secrets: ds})
	if len(lines) == 0 {
		t.Fatal("Info returned no lines")
	}
	var joined string
	for _, l := range lines {
		joined += l.Label + " " + l.Value + "\n"
	}
	for _, want := range []string{"http://127.0.0.1:4000", "/ui", ds["litellm"]["master_key"], "admin"} {
		if !strings.Contains(joined, want) {
			t.Errorf("Info output missing %q\n%s", want, joined)
		}
	}
}
