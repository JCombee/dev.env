package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jcombee/devenv/internal/apps"
	"github.com/jcombee/devenv/internal/compose"
	"github.com/jcombee/devenv/internal/global"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage global apps that run independently of any project",
}

var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available global apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAppList(compose.NewExecRunner())
	},
}

var appEnableCmd = &cobra.Command{
	Use:   "enable <app>",
	Short: "Enable a global app and start it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAppEnable(compose.NewExecRunner(), args[0])
	},
}

var appDisableCmd = &cobra.Command{
	Use:   "disable <app>",
	Short: "Stop a global app and disable it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAppDisable(compose.NewExecRunner(), args[0])
	},
}

var appStartCmd = &cobra.Command{
	Use:   "start [app]",
	Short: "Start enabled global apps",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAppStart(compose.NewExecRunner(), firstArg(args))
	},
}

var appStopCmd = &cobra.Command{
	Use:   "stop [app]",
	Short: "Stop enabled global apps without disabling them",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAppStop(compose.NewExecRunner(), firstArg(args))
	},
}

var appInfoCmd = &cobra.Command{
	Use:   "info <app>",
	Short: "Show an app's URLs and credentials",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAppInfo(args[0])
	},
}

func init() {
	appCmd.AddCommand(appListCmd, appEnableCmd, appDisableCmd, appStartCmd, appStopCmd, appInfoCmd)
	rootCmd.AddCommand(appCmd)
}

func firstArg(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

// findApp resolves an app name, returning an error listing the known apps.
func findApp(name string) (apps.App, error) {
	app, ok := apps.Find(name)
	if !ok {
		return apps.App{}, fmt.Errorf("unknown app %q (available: %s)", name, strings.Join(apps.Names(), ", "))
	}
	return app, nil
}

// appContainers returns the compose service names of one app.
func appContainers(app apps.App, st *apps.State, ds global.DockerSecrets) []string {
	return compose.AppContainerNames(app, st.HostPort(app), ds)
}

// writeAppsCompose regenerates the apps compose file from the current state.
func writeAppsCompose(st *apps.State, ds global.DockerSecrets) error {
	if err := compose.WriteApps(st.Enabled(), st.Ports(), ds); err != nil {
		return fmt.Errorf("write apps compose file: %w", err)
	}
	return nil
}

func runAppList(runner compose.Runner) error {
	st, err := apps.Load()
	if err != nil {
		return fmt.Errorf("read apps.yaml: %w", err)
	}

	ds, err := global.LoadDockerSecrets()
	if err != nil {
		return fmt.Errorf("read docker secrets: %w", err)
	}

	// Query real docker status if the apps compose file exists.
	statusMap := map[string]string{}
	if _, err := os.Stat(compose.AppsPath()); err == nil {
		if statuses, err := runner.PS(compose.AppsPath()); err == nil {
			for _, s := range statuses {
				statusMap[s.Name] = s.State
			}
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"APP", "STATUS", "PORT", "DESCRIPTION"})

	for _, app := range apps.All {
		status := "disabled"
		if st.IsEnabled(app.Name) {
			status = "stopped"
			for _, name := range appContainers(app, st, ds) {
				if statusMap[name] == "running" {
					status = "running"
					break
				}
			}
		}
		t.AppendRow(table.Row{app.Name, status, st.HostPort(app), app.Description})
	}

	t.Render()
	return nil
}

func runAppEnable(runner compose.Runner, name string) error {
	app, err := findApp(name)
	if err != nil {
		return err
	}

	ds, err := global.LoadDockerSecrets()
	if err != nil {
		return fmt.Errorf("read docker secrets: %w", err)
	}
	apps.EnsureSecrets(app, ds)
	if err := global.SaveDockerSecrets(ds); err != nil {
		return fmt.Errorf("save docker secrets: %w", err)
	}

	st, err := apps.Load()
	if err != nil {
		return fmt.Errorf("read apps.yaml: %w", err)
	}
	entry := st.Apps[app.Name]
	entry.Enabled = true
	st.Apps[app.Name] = entry
	if err := st.Save(); err != nil {
		return fmt.Errorf("save apps.yaml: %w", err)
	}

	if err := writeAppsCompose(st, ds); err != nil {
		return err
	}
	if err := runner.Up(compose.AppsPath(), appContainers(app, st, ds)...); err != nil {
		return fmt.Errorf("start %s: %w", app.Name, err)
	}

	fmt.Printf("✓ %s  started\n\n", app.Name)
	printAppInfo(app, st.HostPort(app), ds)
	return nil
}

func runAppDisable(runner compose.Runner, name string) error {
	app, err := findApp(name)
	if err != nil {
		return err
	}

	st, err := apps.Load()
	if err != nil {
		return fmt.Errorf("read apps.yaml: %w", err)
	}
	if !st.IsEnabled(app.Name) {
		fmt.Printf("~ app %q is not enabled\n", app.Name)
		return nil
	}

	ds, err := global.LoadDockerSecrets()
	if err != nil {
		return fmt.Errorf("read docker secrets: %w", err)
	}

	if err := runner.Stop(compose.AppsPath(), appContainers(app, st, ds)...); err != nil {
		return fmt.Errorf("stop %s: %w", app.Name, err)
	}

	entry := st.Apps[app.Name]
	entry.Enabled = false
	st.Apps[app.Name] = entry
	if err := st.Save(); err != nil {
		return fmt.Errorf("save apps.yaml: %w", err)
	}
	if err := writeAppsCompose(st, ds); err != nil {
		return err
	}

	fmt.Printf("✓ %s  stopped and disabled\n", app.Name)
	return nil
}

func runAppStart(runner compose.Runner, name string) error {
	st, ds, targets, err := appTargets(name)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		fmt.Println("No apps enabled. Run `dev app enable <app>` first.")
		return nil
	}

	// Regenerate first so edits to apps.yaml (e.g. a port override) take effect.
	if err := writeAppsCompose(st, ds); err != nil {
		return err
	}

	for _, app := range targets {
		if err := runner.Up(compose.AppsPath(), appContainers(app, st, ds)...); err != nil {
			return fmt.Errorf("start %s: %w", app.Name, err)
		}
		fmt.Printf("✓ %s  started\n", app.Name)
	}
	return nil
}

func runAppStop(runner compose.Runner, name string) error {
	st, ds, targets, err := appTargets(name)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		fmt.Println("No apps enabled.")
		return nil
	}

	for _, app := range targets {
		if err := runner.Stop(compose.AppsPath(), appContainers(app, st, ds)...); err != nil {
			return fmt.Errorf("stop %s: %w", app.Name, err)
		}
		fmt.Printf("✓ %s  stopped\n", app.Name)
	}
	return nil
}

// appTargets resolves the apps a start/stop call applies to: the named app when
// given (and enabled), otherwise every enabled app.
func appTargets(name string) (*apps.State, global.DockerSecrets, []apps.App, error) {
	st, err := apps.Load()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read apps.yaml: %w", err)
	}
	ds, err := global.LoadDockerSecrets()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read docker secrets: %w", err)
	}

	if name == "" {
		return st, ds, st.Enabled(), nil
	}

	app, err := findApp(name)
	if err != nil {
		return nil, nil, nil, err
	}
	if !st.IsEnabled(app.Name) {
		return nil, nil, nil, fmt.Errorf("app %q is not enabled — run `dev app enable %s` first", app.Name, app.Name)
	}
	return st, ds, []apps.App{app}, nil
}

func runAppInfo(name string) error {
	app, err := findApp(name)
	if err != nil {
		return err
	}
	st, err := apps.Load()
	if err != nil {
		return fmt.Errorf("read apps.yaml: %w", err)
	}
	if !st.IsEnabled(app.Name) {
		return fmt.Errorf("app %q is not enabled — run `dev app enable %s` first", app.Name, app.Name)
	}
	ds, err := global.LoadDockerSecrets()
	if err != nil {
		return fmt.Errorf("read docker secrets: %w", err)
	}

	printAppInfo(app, st.HostPort(app), ds)
	return nil
}

func printAppInfo(app apps.App, hostPort int, ds global.DockerSecrets) {
	for _, line := range app.Info(apps.BuildContext{HostPort: hostPort, Secrets: ds}) {
		if line.Label == "" {
			fmt.Printf("  %s\n", line.Value)
			continue
		}
		fmt.Printf("  %-11s %s\n", line.Label, line.Value)
	}
}
