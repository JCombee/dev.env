package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jcombee/devenv/internal/compose"
	"github.com/jcombee/devenv/internal/global"
	"github.com/jcombee/devenv/internal/refcount"
	"github.com/jcombee/devenv/internal/store"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show all containers and which projects use them",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStatus(compose.NewExecRunner())
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(runner compose.Runner) error {
	var svcFile global.Services
	svcPath := filepath.Join(global.Dir(), "services.yaml")
	if err := store.Read(svcPath, &svcFile); err != nil {
		return fmt.Errorf("read services.yaml: %w", err)
	}

	if len(svcFile.Services) == 0 {
		fmt.Println("No services registered. Run `dev start` in a project first.")
		return nil
	}

	// Query real docker status if compose file exists.
	statusMap := map[string]string{}
	sharedComposePath := compose.SharedPath()
	if _, err := os.Stat(sharedComposePath); err == nil {
		statuses, err := runner.PS(sharedComposePath)
		if err == nil {
			for _, s := range statuses {
				statusMap[s.Name] = s.State
			}
		}
	}

	usage, err := refcount.AllServiceUsage()
	if err != nil {
		return fmt.Errorf("read project states: %w", err)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"SERVICE", "STATUS", "PROJECTS"})

	for composeName, entry := range svcFile.Services {
		tag := entry.Tag
		if tag == "" {
			tag = "latest"
		}
		service := entry.Image + ":" + tag

		status := statusMap[composeName]
		if status == "" {
			status = "unknown"
		}

		projects := usage[composeName]
		t.AppendRow(table.Row{service, status, strings.Join(projects, ", ")})
	}

	t.Render()
	return nil
}
