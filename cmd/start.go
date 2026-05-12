package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/jcombee/devenv/internal/compose"
	"github.com/jcombee/devenv/internal/config"
	"github.com/jcombee/devenv/internal/envfile"
	"github.com/jcombee/devenv/internal/global"
	"github.com/jcombee/devenv/internal/ports"
	"github.com/jcombee/devenv/internal/project"
	"github.com/jcombee/devenv/internal/store"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start services for this project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStart(compose.NewExecRunner())
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart(runner compose.Runner) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Load(cwd)
	if err != nil {
		return err
	}

	// Register project — handle name collision interactively.
	if err := registerWithCollisionPrompt(cfg, cwd); err != nil {
		return err
	}

	// Load global state.
	svcFile, err := loadServicesFile()
	if err != nil {
		return err
	}
	dockerSecrets, err := global.LoadDockerSecrets()
	if err != nil {
		return err
	}
	settings, err := global.ReadSettings()
	if err != nil {
		return err
	}

	sharedComposePath := compose.SharedPath()
	state := &project.State{
		Running:      true,
		Services:     map[string]string{},
		Acknowledged: []string{},
	}

	var sharedToStart []string
	var dedicatedEntries []config.ServiceEntry
	servicesChanged := false

	for _, svcEntry := range cfg.Services {
		tag := svcEntry.Tag
		if tag == "" {
			tag = "latest"
		}
		composeName := compose.ServiceName(svcEntry.Image, tag)

		if svcEntry.Dedicated {
			dedicatedEntries = append(dedicatedEntries, svcEntry)
			dedicatedComposeName := cfg.Project + "-" + composeName
			state.Services[dedicatedComposeName] = "dedicated"
		} else {
			if _, exists := svcFile.Services[composeName]; !exists {
				svcFile.Services[composeName] = global.GlobalServiceEntry{
					Image: svcEntry.Image,
					Tag:   tag,
				}
				servicesChanged = true
			}
			compose.ContainerEnv(svcEntry.Image, composeName, dockerSecrets)
			sharedToStart = append(sharedToStart, composeName)
			state.Services[composeName] = "shared"
		}
	}

	if servicesChanged {
		if err := saveServicesFile(svcFile); err != nil {
			return fmt.Errorf("save services.yaml: %w", err)
		}
		if err := global.SaveDockerSecrets(dockerSecrets); err != nil {
			return fmt.Errorf("save docker secrets: %w", err)
		}
	}

	// Always regenerate shared compose — ensures new secrets are reflected.
	if err := compose.WriteShared(svcFile.Services, dockerSecrets, settings.Ports); err != nil {
		return fmt.Errorf("write docker-compose.yml: %w", err)
	}

	// Start shared services.
	if len(sharedToStart) > 0 {
		if err := runner.Up(sharedComposePath, sharedToStart...); err != nil {
			return fmt.Errorf("docker compose up: %w", err)
		}
		for _, name := range sharedToStart {
			entry := svcFile.Services[name]
			tag := entry.Tag
			if tag == "" {
				tag = "latest"
			}
			fmt.Printf("✓ %s:%s  started (shared)\n", entry.Image, tag)
		}
	}

	// Start dedicated services.
	if len(dedicatedEntries) > 0 {
		if err := compose.WriteDedicated(cfg.Project, dedicatedEntries, dockerSecrets); err != nil {
			return fmt.Errorf("write dedicated compose: %w", err)
		}
		dedicatedComposePath := compose.DedicatedPath(cfg.Project)
		var dedicatedNames []string
		for _, e := range dedicatedEntries {
			tag := e.Tag
			if tag == "" {
				tag = "latest"
			}
			dedicatedNames = append(dedicatedNames, cfg.Project+"-"+compose.ServiceName(e.Image, tag))
		}
		if err := runner.Up(dedicatedComposePath, dedicatedNames...); err != nil {
			return fmt.Errorf("docker compose up (dedicated): %w", err)
		}
		for _, e := range dedicatedEntries {
			tag := e.Tag
			if tag == "" {
				tag = "latest"
			}
			fmt.Printf("✓ %s:%s  started (dedicated)\n", e.Image, tag)
		}
	}

	if err := global.SaveDockerSecrets(dockerSecrets); err != nil {
		return fmt.Errorf("save docker secrets: %w", err)
	}

	// Ensure per-project secrets (redis DB index, soketi/reverb app creds, etc.).
	projSecrets, err := project.LoadSecrets(cfg.Project)
	if err != nil {
		return fmt.Errorf("load project secrets: %w", err)
	}
	for _, svcEntry := range cfg.Services {
		tag := svcEntry.Tag
		if tag == "" {
			tag = "latest"
		}
		composeName := compose.ServiceName(svcEntry.Image, tag)
		if err := envfile.EnsureProjectSecrets(svcEntry.Image, composeName, cfg.Project, projSecrets, dockerSecrets); err != nil {
			return fmt.Errorf("ensure project secrets for %s: %w", svcEntry.Image, err)
		}
	}
	if err := project.SaveSecrets(cfg.Project, projSecrets); err != nil {
		return fmt.Errorf("save project secrets: %w", err)
	}

	// Write .env when enabled.
	if cfg.Env {
		allVars := map[string]string{}
		for _, svcEntry := range cfg.Services {
			tag := svcEntry.Tag
			if tag == "" {
				tag = "latest"
			}
			composeName := compose.ServiceName(svcEntry.Image, tag)
			r, _ := ports.Resolve(svcEntry.Image, tag, settings.Ports, svcEntry.Port)
			svcVars := envfile.ServiceVars(svcEntry.Image, composeName, r.Port, cfg.Project, dockerSecrets, projSecrets)
			svcVars = envfile.ApplyEnvMap(svcVars, svcEntry.EnvMap)
			for k, v := range svcVars {
				allVars[k] = v
			}
		}

		envPath := filepath.Join(cwd, ".env")
		existing, parseErr := envfile.ParseEnv(envPath)

		if os.IsNotExist(parseErr) {
			if err := envfile.Write(envPath, allVars); err != nil {
				return fmt.Errorf("write .env: %w", err)
			}
			fmt.Println("✓ .env written")
		} else if parseErr != nil {
			return fmt.Errorf("read .env: %w", parseErr)
		} else {
			newKeys, changedKeys := envfile.Diff(allVars, existing)
			userVars := envfile.UserVars(allVars, existing)

			if len(newKeys) == 0 && len(changedKeys) == 0 {
				fmt.Println("✓ .env up to date")
			} else {
				envfile.PrintDiff(newKeys, changedKeys, allVars, existing)
				update, err := confirmPrompt("Update .env file?")
				if err != nil {
					return err
				}
				if update {
					if err := envfile.WriteMerged(envPath, allVars, userVars); err != nil {
						return fmt.Errorf("write .env: %w", err)
					}
					fmt.Println("✓ .env updated")
				} else {
					fmt.Println("~ .env skipped")
				}
			}
		}
	}

	if err := project.SaveState(cfg.Project, state); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	return nil
}

func registerWithCollisionPrompt(cfg *config.ProjectConfig, cwd string) error {
	err := project.Register(cfg.Project, cwd)
	if err == nil {
		return nil
	}

	var col *project.ErrNameCollision
	if !errors.As(err, &col) {
		return err
	}

	fmt.Printf("! Project name %q is already in use by: %s\n\n", cfg.Project, col.ExistingPath)
	fmt.Println("  Using the same name means sharing the same database, credentials, and state.")
	fmt.Println("  If intentional (e.g. monorepo), continue. Otherwise choose a new name.")
	fmt.Println()

	var newName string
	if err := huh.NewInput().
		Title("New project name").
		Placeholder(cfg.Project).
		Value(&newName).
		Run(); err != nil {
		return err
	}
	if newName == "" {
		newName = cfg.Project
	}

	cfg.Project = newName
	local := &config.LocalConfig{Project: newName}
	if err := config.SaveLocal(cwd, local); err != nil {
		return fmt.Errorf("save .dev.env.local.yaml: %w", err)
	}

	return project.Register(cfg.Project, cwd)
}

func loadServicesFile() (*global.Services, error) {
	var svcFile global.Services
	svcPath := filepath.Join(global.Dir(), "services.yaml")
	if err := store.Read(svcPath, &svcFile); err != nil {
		return nil, fmt.Errorf("read services.yaml: %w", err)
	}
	if svcFile.Services == nil {
		svcFile.Services = map[string]global.GlobalServiceEntry{}
	}
	return &svcFile, nil
}

func saveServicesFile(svcFile *global.Services) error {
	svcPath := filepath.Join(global.Dir(), "services.yaml")
	return store.Write(svcPath, svcFile)
}

// confirmPrompt asks a yes/no question. Uses huh in a real terminal;
// falls back to reading a line from stdin when not connected to a TTY.
func confirmPrompt(title string) (bool, error) {
	fi, err := os.Stdin.Stat()
	isTTY := err == nil && (fi.Mode()&os.ModeCharDevice) != 0
	if !isTTY {
		var line string
		fmt.Scanln(&line)
		line = strings.ToLower(strings.TrimSpace(line))
		return line == "y" || line == "yes", nil
	}
	var v bool
	if err := huh.NewConfirm().Title(title).Value(&v).Run(); err != nil {
		return false, err
	}
	return v, nil
}
