package apps

import "fmt"

// LiteLLM runs an OpenAI-compatible proxy in front of every LLM provider,
// backed by its own Postgres. It needs no config file: the container migrates
// its own schema on boot, and models, providers and API keys are managed in
// its admin UI and stored in that database.
var LiteLLM = App{
	Name:        "litellm",
	Description: "LLM proxy/gateway with admin UI",
	DefaultPort: 4000,
	Secrets: []SecretSpec{
		{ComposeName: "litellm", Field: "master_key", Gen: func() string { return "sk-" + randHex(16) }},
		{ComposeName: "litellm", Field: "salt_key", Gen: func() string { return randHex(16) }},
		{ComposeName: "litellm", Field: "ui_password", Gen: func() string { return randHex(8) }},
		{ComposeName: "litellm-db", Field: "root_password", Gen: func() string { return randHex(16) }},
	},
	Build: func(c BuildContext) []Container {
		dbPassword := c.Secret("litellm-db", "root_password")
		return []Container{
			{
				Name:     "litellm",
				Image:    "ghcr.io/berriai/litellm",
				Tag:      "main-stable",
				HostPort: c.HostPort,
				Port:     4000,
				Env: map[string]string{
					"LITELLM_MASTER_KEY": c.Secret("litellm", "master_key"),
					"LITELLM_SALT_KEY":   c.Secret("litellm", "salt_key"),
					"DATABASE_URL":       fmt.Sprintf("postgresql://litellm:%s@litellm-db:5432/litellm", dbPassword),
					"STORE_MODEL_IN_DB":  "True",
					"UI_USERNAME":        "admin",
					"UI_PASSWORD":        c.Secret("litellm", "ui_password"),
				},
				DependsOn: map[string]string{"litellm-db": "service_healthy"},
			},
			{
				Name:  "litellm-db",
				Image: "postgres",
				Tag:   "16",
				Port:  5432, // internal only — never published
				Env: map[string]string{
					"POSTGRES_DB":       "litellm",
					"POSTGRES_USER":     "litellm",
					"POSTGRES_PASSWORD": dbPassword,
				},
				Volumes:      []string{"litellm-db-data:/var/lib/postgresql/data"},
				NamedVolumes: []string{"litellm-db-data"},
				Healthcheck: &Healthcheck{
					Test:     []string{"CMD-SHELL", "pg_isready -U litellm"},
					Interval: "5s",
					Timeout:  "5s",
					Retries:  10,
				},
			},
		}
	},
	Info: func(c BuildContext) []InfoLine {
		base := fmt.Sprintf("http://127.0.0.1:%d", c.HostPort)
		return []InfoLine{
			{Label: "Proxy", Value: base},
			{Label: "Admin UI", Value: base + "/ui"},
			{Label: "Master key", Value: c.Secret("litellm", "master_key")},
			{Label: "UI login", Value: "admin / " + c.Secret("litellm", "ui_password")},
			{Label: "", Value: "Add your models and provider API keys in the admin UI."},
		}
	},
}
