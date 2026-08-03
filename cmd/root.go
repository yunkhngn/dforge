package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "dforge",
	Short: "One command to bootstrap, inspect and manage Docker Compose projects",
}

func init() {
	RootCmd.PersistentFlags().Bool("yes", false, "assume yes for confirmation prompts")
	RootCmd.PersistentFlags().Bool("force", false, "bypass safety checks (overwrite/destructive)")
}

func Execute() error { return RootCmd.Execute() }
