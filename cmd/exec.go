package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jcombee/devenv/internal/compose"
	"github.com/jcombee/devenv/internal/config"
	devexec "github.com/jcombee/devenv/internal/exec"
	"github.com/jcombee/devenv/internal/global"
	"github.com/jcombee/devenv/internal/project"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:                "exec [service] [flags...]",
	Short:              "Open an interactive session inside a running container",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return runExecList()
		}
		return runExec(compose.NewExecRunner(), args[0], args[1:])
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}

func runExecList() error {
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
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("No services found. Run `dev start` first.")
			return nil
		}
		return fmt.Errorf("load project state: %w", err)
	}

	images := availableImages(cfg, state)
	if len(images) == 0 {
		fmt.Println("No services found. Run `dev start` first.")
		return nil
	}

	sort.Strings(images)
	fmt.Println("Available services:")
	for _, img := range images {
		fmt.Println(" ", img)
	}
	return nil
}

func runExec(runner compose.Runner, image string, extra []string) error {
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
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("project is not running, start it with `dev start`")
		}
		return fmt.Errorf("load project state: %w", err)
	}
	if !state.Running {
		return fmt.Errorf("project is not running, start it with `dev start`")
	}

	composeName, shareMode, found := findStateService(cfg, state, image)
	if !found {
		avail := availableImages(cfg, state)
		sort.Strings(avail)
		return fmt.Errorf("service %q not found in this project; available: %s", image, strings.Join(avail, ", "))
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
	ps, err := project.LoadSecrets(cfg.Project)
	if err != nil {
		return fmt.Errorf("load project secrets: %w", err)
	}

	args := devexec.BuildArgs(image, composeName, ds, ps, cfg.Project, extra)
	return runner.Exec(composeFile, composeName, args)
}

// findStateService searches state.Services for a compose name matching image.
// Returns composeName, shareMode ("shared"|"dedicated"), and whether found.
func findStateService(cfg *config.ProjectConfig, state *project.State, image string) (string, string, bool) {
	for composeName, mode := range state.Services {
		if mode == "shared" {
			if strings.HasPrefix(composeName, image+"-") {
				return composeName, mode, true
			}
		} else {
			prefix := cfg.Project + "-" + image + "-"
			if strings.HasPrefix(composeName, prefix) {
				return composeName, mode, true
			}
		}
	}
	return "", "", false
}

// availableImages returns the image names of services in this project's state.
func availableImages(cfg *config.ProjectConfig, state *project.State) []string {
	seen := map[string]bool{}
	for composeName, mode := range state.Services {
		var image string
		if mode == "shared" {
			// shared compose name = <image>-<tag>, find the matching config entry
			image = imageFromComposeName(cfg, composeName)
		} else {
			// dedicated compose name = <project>-<image>-<tag>
			prefix := cfg.Project + "-"
			if strings.HasPrefix(composeName, prefix) {
				image = imageFromComposeName(cfg, strings.TrimPrefix(composeName, prefix))
			}
		}
		if image != "" && !seen[image] {
			seen[image] = true
		}
	}
	var images []string
	for img := range seen {
		images = append(images, img)
	}
	return images
}

// imageFromComposeName finds the image name from .dev.env.yaml config that
// matches the given compose name (image-tag format).
func imageFromComposeName(cfg *config.ProjectConfig, composeName string) string {
	for _, svc := range cfg.Services {
		tag := svc.Tag
		if tag == "" {
			tag = "latest"
		}
		if composeName == compose.ServiceName(svc.Image, tag) {
			return svc.Image
		}
	}
	return ""
}
