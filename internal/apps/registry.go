// Package apps holds the registry of global apps: containers that belong to
// the machine rather than to a project. They are enabled once with
// `dev app enable`, stay up independently of any project, and are never
// declared in a .dev.env.yaml.
//
// Adding an app means adding one App value to All — no changes to the
// commands or to compose generation are needed.
package apps

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/jcombee/devenv/internal/global"
)

// App describes a global app: its metadata, the credentials it needs, and how
// to turn those into containers.
type App struct {
	Name        string
	Description string
	// DefaultPort is the host port the app is published on unless apps.yaml
	// overrides it.
	DefaultPort int
	// Secrets are ensured in docker/secrets.yaml before Build is called.
	Secrets []SecretSpec
	// Build returns the containers making up the app.
	Build func(BuildContext) []Container
	// Info returns the URLs and credentials shown after `dev app enable` and
	// by `dev app info`.
	Info func(BuildContext) []InfoLine
}

// SecretSpec declares one credential of an app. Gen produces the value the
// first time the secret is needed; it is never regenerated afterwards.
type SecretSpec struct {
	ComposeName string
	Field       string
	Gen         func() string
}

// BuildContext carries everything Build and Info need to render an app.
type BuildContext struct {
	HostPort int
	Secrets  global.DockerSecrets
}

// Secret returns a stored credential, or "" when it has not been generated.
func (c BuildContext) Secret(composeName, field string) string {
	return c.Secrets[composeName][field]
}

// Container is one docker-compose service of an app.
type Container struct {
	Name  string
	Image string
	Tag   string
	// HostPort is the published port; 0 means the container is internal only.
	HostPort     int
	Port         int
	Env          map[string]string
	Volumes      []string
	NamedVolumes []string
	Command      []string
	DependsOn    map[string]string
	Healthcheck  *Healthcheck
}

// Healthcheck mirrors the docker-compose healthcheck block.
type Healthcheck struct {
	Test     []string
	Interval string
	Timeout  string
	Retries  int
}

// InfoLine is one label/value pair printed to the user.
type InfoLine struct {
	Label string
	Value string
}

// All registered global apps.
var All = []App{LiteLLM}

// Find returns the app with the given name, or false if it is not registered.
func Find(name string) (App, bool) {
	for _, a := range All {
		if a.Name == name {
			return a, true
		}
	}
	return App{}, false
}

// Names returns the names of all registered apps, for error messages.
func Names() []string {
	names := make([]string, 0, len(All))
	for _, a := range All {
		names = append(names, a.Name)
	}
	return names
}

// EnsureSecrets generates every credential the app declares that is not
// already stored. Existing values are never overwritten.
func EnsureSecrets(a App, ds global.DockerSecrets) {
	for _, s := range a.Secrets {
		global.EnsureDockerSecret(ds, s.ComposeName, s.Field, s.Gen())
	}
}

// randHex returns n random bytes hex-encoded.
func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("apps: generate random value: " + err.Error())
	}
	return hex.EncodeToString(b)
}
