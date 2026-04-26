package cmd

import (
	"fmt"

	"github.com/jcombee/devenv/internal/global"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Initialize the ~/.dev.env/ directory structure",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSetup()
	},
}

func runSetup() error {
	if err := global.Setup(); err != nil {
		return fmt.Errorf("setup failed: %w", err)
	}
	fmt.Println("✓ ~/.dev.env/ ready")
	return nil
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
