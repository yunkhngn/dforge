package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yunkhngn/dforge/internal/docker"
)

func runShell(cl docker.Client, service string) error {
	return cl.Shell(service)
}

func init() {
	cmd := &cobra.Command{
		Use:   "shell <service>",
		Short: "Open a shell inside a container (bash, falls back to sh)",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runShell(docker.New(), args[0])
		},
	}
	RootCmd.AddCommand(cmd)
}
