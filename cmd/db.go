package cmd

import "github.com/spf13/cobra"

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations for the current project",
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
