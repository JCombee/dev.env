package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jcombee/devenv/internal/store"
	"gopkg.in/yaml.v3"
)

const (
	File      = ".dev.env.yaml"
	LocalFile = ".dev.env.local.yaml"
)

type ServiceEntry struct {
	Image     string            `yaml:"image"`
	Tag       string            `yaml:"tag,omitempty"`
	Dedicated bool              `yaml:"dedicated,omitempty"`
	EnvMap    map[string]string `yaml:"env_map,omitempty"`
}

// UnmarshalYAML handles both the shorthand string form ("redis") and the full object form.
func (s *ServiceEntry) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		parts := strings.SplitN(value.Value, ":", 2)
		s.Image = parts[0]
		if len(parts) == 2 {
			s.Tag = parts[1]
		}
		return nil
	}
	type plain ServiceEntry
	return value.Decode((*plain)(s))
}

func (s ServiceEntry) ImageTag() string {
	tag := s.Tag
	if tag == "" {
		tag = "latest"
	}
	return s.Image + ":" + tag
}

type ProjectConfig struct {
	Project  string         `yaml:"project"`
	Type     string         `yaml:"type,omitempty"`
	Env      bool           `yaml:"env,omitempty"`
	Services []ServiceEntry `yaml:"services"`
}

// ServiceOverride holds per-service local overrides.
type ServiceOverride struct {
	Tag       string            `yaml:"tag,omitempty"`
	Dedicated *bool             `yaml:"dedicated,omitempty"`
	EnvMap    map[string]string `yaml:"env_map,omitempty"`
}

// LocalConfig holds optional overrides from .dev.env.local.yaml.
type LocalConfig struct {
	Project  string                     `yaml:"project,omitempty"`
	Type     string                     `yaml:"type,omitempty"`
	Env      *bool                      `yaml:"env,omitempty"`
	Services map[string]ServiceOverride `yaml:"services,omitempty"`
}

// Load reads .dev.env.yaml from dir and merges any .dev.env.local.yaml on top.
func Load(dir string) (*ProjectConfig, error) {
	var cfg ProjectConfig
	if err := store.Read(filepath.Join(dir, File), &cfg); err != nil {
		return nil, fmt.Errorf("read %s: %w", File, err)
	}

	localPath := filepath.Join(dir, LocalFile)
	if _, err := os.Stat(localPath); err == nil {
		var local LocalConfig
		if err := store.Read(localPath, &local); err != nil {
			return nil, fmt.Errorf("read %s: %w", LocalFile, err)
		}
		merge(&cfg, &local)
	}

	return &cfg, nil
}

// Save writes cfg to .dev.env.yaml in dir.
func Save(dir string, cfg *ProjectConfig) error {
	return store.Write(filepath.Join(dir, File), cfg)
}

// SaveLocal writes fields to .dev.env.local.yaml in dir (upsert).
func SaveLocal(dir string, local *LocalConfig) error {
	return store.Write(filepath.Join(dir, LocalFile), local)
}

func merge(cfg *ProjectConfig, local *LocalConfig) {
	if local.Project != "" {
		cfg.Project = local.Project
	}
	if local.Type != "" {
		cfg.Type = local.Type
	}
	if local.Env != nil {
		cfg.Env = *local.Env
	}
	for image, override := range local.Services {
		for i, svc := range cfg.Services {
			if svc.Image == image {
				if override.Tag != "" {
					cfg.Services[i].Tag = override.Tag
				}
				if override.Dedicated != nil {
					cfg.Services[i].Dedicated = *override.Dedicated
				}
				for k, v := range override.EnvMap {
					if cfg.Services[i].EnvMap == nil {
						cfg.Services[i].EnvMap = map[string]string{}
					}
					cfg.Services[i].EnvMap[k] = v
				}
			}
		}
	}
}
