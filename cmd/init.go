package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/jcombee/devenv/internal/config"
	"github.com/jcombee/devenv/internal/services"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize DEV.ENV for the current project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg := &config.ProjectConfig{}

	// Step 1: project name
	if err := huh.NewInput().
		Title("Project name").
		Value(&cfg.Project).
		Run(); err != nil {
		return err
	}

	// Step 2: project type
	detected := config.DetectType(cwd)
	typeOptions := []huh.Option[string]{
		huh.NewOption("laravel", "laravel"),
		huh.NewOption("node", "node"),
		huh.NewOption("(none)", ""),
	}
	if detected != "" {
		fmt.Printf("Detected project type: %s\n", detected)
		cfg.Type = detected
	}
	if err := huh.NewSelect[string]().
		Title("Project type").
		Options(typeOptions...).
		Value(&cfg.Type).
		Run(); err != nil {
		return err
	}

	// Step 3: select services
	var selectedImages []string
	serviceOptions := make([]huh.Option[string], len(services.All))
	for i, svc := range services.All {
		serviceOptions[i] = huh.NewOption(svc.Image, svc.Image)
	}
	if err := huh.NewMultiSelect[string]().
		Title("Select services").
		Options(serviceOptions...).
		Value(&selectedImages).
		Run(); err != nil {
		return err
	}

	// Step 4: per-service configuration
	for _, image := range selectedImages {
		svc, _ := services.Find(image)

		entry := config.ServiceEntry{
			Image: image,
			Tag:   svc.DefaultTag,
		}

		if err := huh.NewInput().
			Title(fmt.Sprintf("%s tag", image)).
			Placeholder(svc.DefaultTag).
			Value(&entry.Tag).
			Run(); err != nil {
			return err
		}
		if entry.Tag == "" {
			entry.Tag = svc.DefaultTag
		}

		if svc.CanBeShared {
			modeOptions := []huh.Option[bool]{
				huh.NewOption("shared", false),
				huh.NewOption("dedicated", true),
			}
			if err := huh.NewSelect[bool]().
				Title(fmt.Sprintf("%s mode", image)).
				Options(modeOptions...).
				Value(&entry.Dedicated).
				Run(); err != nil {
				return err
			}
		} else {
			entry.Dedicated = true
		}

		cfg.Services = append(cfg.Services, entry)
	}

	// Step 5: .env generation
	if err := huh.NewConfirm().
		Title("Enable .env generation?").
		Value(&cfg.Env).
		Run(); err != nil {
		return err
	}

	// Write config
	if err := config.Save(cwd, cfg); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Printf("✓ Wrote %s\n", strings.Join([]string{cwd, config.File}, "/"))
	return nil
}
