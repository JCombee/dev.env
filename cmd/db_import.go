package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jcombee/devenv/internal/compose"
	"github.com/jcombee/devenv/internal/config"
	"github.com/jcombee/devenv/internal/global"
	"github.com/jcombee/devenv/internal/project"
	"github.com/spf13/cobra"
)

var dbImportCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import a SQL dump into the project database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDBImport(args[0])
	},
}

func init() {
	dbCmd.AddCommand(dbImportCmd)
}

func runDBImport(filePath string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(cwd, filePath)
	}

	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("file %q not found", filepath.Base(filePath))
	}

	cfg, err := config.Load(cwd)
	if err != nil {
		return err
	}

	var dbServices []config.ServiceEntry
	for _, svc := range cfg.Services {
		if isDBImage(svc.Image) {
			dbServices = append(dbServices, svc)
		}
	}
	if len(dbServices) == 0 {
		return fmt.Errorf("no supported database service configured")
	}
	if len(dbServices) > 1 {
		return fmt.Errorf("multiple database services configured — only one supported for import")
	}
	svcEntry := dbServices[0]

	state, err := project.LoadState(cfg.Project)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("project is not running — run `dev start` first")
		}
		return fmt.Errorf("load project state: %w", err)
	}
	if !state.Running {
		return fmt.Errorf("project is not running — run `dev start` first")
	}

	composeName, shareMode, found := findStateService(cfg, state, svcEntry.Image)
	if !found {
		return fmt.Errorf("container for %q not found in project state — run `dev start` first", svcEntry.Image)
	}

	var composeFile string
	if shareMode == "shared" {
		composeFile = compose.SharedPath()
	} else {
		composeFile = compose.DedicatedPath(cfg.Project)
	}

	ds, err := global.LoadDockerSecrets()
	if err != nil {
		return fmt.Errorf("load docker secrets: %w", err)
	}

	tag := svcEntry.Tag
	if tag == "" {
		tag = "latest"
	}
	fmt.Printf("Importing %s into %q (%s:%s)...\n", filepath.Base(filePath), cfg.Project, svcEntry.Image, tag)

	if err := compose.ImportDB(composeFile, composeName, svcEntry.Image, cfg.Project, filePath, ds); err != nil {
		return fmt.Errorf("import: %w", err)
	}

	fmt.Println("✓ Import complete")
	return nil
}
