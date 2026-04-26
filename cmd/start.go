package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/jcombee/devenv/internal/compose"
	"github.com/jcombee/devenv/internal/config"
	"github.com/jcombee/devenv/internal/global"
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
	if err := compose.WriteShared(svcFile.Services, dockerSecrets); err != nil {
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
