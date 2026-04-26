package cmd

import (
	"fmt"
	"os"

	"github.com/jcombee/devenv/internal/global"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dev",
	Short: "DEV.ENV — shared Docker dev services for all your projects",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip auto-setup for the setup command itself
		if cmd.Name() == "setup" {
			return nil
		}
		if !global.IsSetUp() {
			if err := global.Setup(); err != nil {
				return fmt.Errorf("auto-setup failed: %w", err)
			}
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
