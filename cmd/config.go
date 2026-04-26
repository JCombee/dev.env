package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jcombee/devenv/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Read and write project configuration",
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all project configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConfigList()
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConfigGet(args[0])
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConfigSet(args[0], args[1])
	},
}

func init() {
	configCmd.AddCommand(configListCmd, configGetCmd, configSetCmd)
	rootCmd.AddCommand(configCmd)
}

func loadProjectConfig() (*config.ProjectConfig, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}
	cfg, err := config.Load(cwd)
	if err != nil {
		return nil, "", err
	}
	return cfg, cwd, nil
}

func runConfigList() error {
	cfg, _, err := loadProjectConfig()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendRow(table.Row{"project", cfg.Project})
	t.AppendRow(table.Row{"type", cfg.Type})
	t.AppendRow(table.Row{"env", fmt.Sprintf("%v", cfg.Env)})

	var svcParts []string
	for _, svc := range cfg.Services {
		tag := svc.Tag
		if tag == "" {
			tag = "latest"
		}
		part := svc.Image + ":" + tag
		if svc.Dedicated {
			part += " (dedicated)"
		}
		svcParts = append(svcParts, part)
	}
	t.AppendRow(table.Row{"services", strings.Join(svcParts, ", ")})
	t.Render()
	return nil
}

func runConfigGet(key string) error {
	cfg, _, err := loadProjectConfig()
	if err != nil {
		return err
	}

	val, err := getConfigValue(cfg, key)
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func runConfigSet(key, value string) error {
	cfg, cwd, err := loadProjectConfig()
	if err != nil {
		return err
	}

	if err := setConfigValue(cfg, key, value); err != nil {
		return err
	}

	if err := config.Save(cwd, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

func getConfigValue(cfg *config.ProjectConfig, key string) (string, error) {
	parts := strings.SplitN(key, ".", 3)
	switch parts[0] {
	case "project":
		return cfg.Project, nil
	case "type":
		return cfg.Type, nil
	case "env":
		return fmt.Sprintf("%v", cfg.Env), nil
	case "services":
		if len(parts) < 2 {
			return "", fmt.Errorf("usage: services.<image>[.<field>]")
		}
		svc := findService(cfg, parts[1])
		if svc == nil {
			return "", fmt.Errorf("service %q not found", parts[1])
		}
		if len(parts) == 2 {
			tag := svc.Tag
			if tag == "" {
				tag = "latest"
			}
			return svc.Image + ":" + tag, nil
		}
		switch parts[2] {
		case "tag":
			tag := svc.Tag
			if tag == "" {
				tag = "latest"
			}
			return tag, nil
		case "dedicated":
			return fmt.Sprintf("%v", svc.Dedicated), nil
		default:
			if strings.HasPrefix(parts[2], "env_map.") {
				mapKey := strings.TrimPrefix(parts[2], "env_map.")
				return svc.EnvMap[mapKey], nil
			}
			return "", fmt.Errorf("unknown service field %q", parts[2])
		}
	}
	return "", fmt.Errorf("unknown key %q", key)
}

func setConfigValue(cfg *config.ProjectConfig, key, value string) error {
	parts := strings.SplitN(key, ".", 3)
	switch parts[0] {
	case "project":
		cfg.Project = value
	case "type":
		cfg.Type = value
	case "env":
		switch strings.ToLower(value) {
		case "true", "1", "yes":
			cfg.Env = true
		case "false", "0", "no":
			cfg.Env = false
		default:
			return fmt.Errorf("env must be true or false")
		}
	case "services":
		if len(parts) < 3 {
			return fmt.Errorf("usage: services.<image>.<field>")
		}
		svc := findService(cfg, parts[1])
		if svc == nil {
			return fmt.Errorf("service %q not found", parts[1])
		}
		switch parts[2] {
		case "tag":
			svc.Tag = value
		case "dedicated":
			switch strings.ToLower(value) {
			case "true", "1", "yes":
				svc.Dedicated = true
			case "false", "0", "no":
				svc.Dedicated = false
			default:
				return fmt.Errorf("dedicated must be true or false")
			}
		default:
			if strings.HasPrefix(parts[2], "env_map.") {
				mapKey := strings.TrimPrefix(parts[2], "env_map.")
				if svc.EnvMap == nil {
					svc.EnvMap = map[string]string{}
				}
				svc.EnvMap[mapKey] = value
			} else {
				return fmt.Errorf("unknown service field %q", parts[2])
			}
		}
	default:
		return fmt.Errorf("unknown key %q", key)
	}
	return nil
}

// findService returns a pointer into cfg.Services for the given image name.
func findService(cfg *config.ProjectConfig, image string) *config.ServiceEntry {
	for i := range cfg.Services {
		if cfg.Services[i].Image == image {
			return &cfg.Services[i]
		}
	}
	return nil
}
