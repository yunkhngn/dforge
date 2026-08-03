package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yunkhngn/dforge/internal/docker"
)

func runLogs(cl docker.Client, service string, follow bool) error {
	return cl.Logs(service, follow)
}

func init() {
	cmd := &cobra.Command{
		Use:   "logs [service]",
		Short: "Show logs for a service",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			follow, _ := c.Flags().GetBool("follow")
			service := ""
			if len(args) == 1 {
				service = args[0]
			}
			return runLogs(docker.New(), service, follow)
		},
	}
	cmd.Flags().BoolP("follow", "f", false, "follow log output")
	RootCmd.AddCommand(cmd)
}
