package cmd

import (
	"fmt"
	"strings"

	"github.com/jcombee/devenv/internal/compose"
	"github.com/jcombee/devenv/internal/config"
	"github.com/jcombee/devenv/internal/project"
	"github.com/jcombee/devenv/internal/refcount"
	"github.com/spf13/cobra"
	"os"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop services for this project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStop(compose.NewExecRunner())
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func runStop(runner compose.Runner) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Load(cwd)
	if err != nil {
		return err
	}

	state, err := project.LoadState(cfg.Project)
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	if !state.Running {
		fmt.Println("~ project is not running")
		return nil
	}

	// Mark project as stopped before scanning refcounts so counts are accurate.
	state.Running = false
	if err := project.SaveState(cfg.Project, state); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	sharedComposePath := compose.SharedPath()
	dedicatedComposePath := compose.DedicatedPath(cfg.Project)

	for composeName, mode := range state.Services {
		switch mode {
		case "dedicated":
			if err := runner.Stop(dedicatedComposePath, composeName); err != nil {
				return fmt.Errorf("stop %s: %w", composeName, err)
			}
			fmt.Printf("✓ %s  stopped (dedicated)\n", composeName)

		case "shared":
			users, err := refcount.ActiveProjectsUsingService(composeName)
			if err != nil {
				return fmt.Errorf("refcount %s: %w", composeName, err)
			}
			if len(users) == 0 {
				if err := runner.Stop(sharedComposePath, composeName); err != nil {
					return fmt.Errorf("stop %s: %w", composeName, err)
				}
				fmt.Printf("✓ %s  stopped\n", composeName)
			} else {
				fmt.Printf("~ %s  kept running (in use by: %s)\n", composeName, strings.Join(users, ", "))
			}
		}
	}

	return nil
}
