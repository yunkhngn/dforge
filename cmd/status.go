package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yunkhngn/dforge/internal/docker"
	"github.com/yunkhngn/dforge/internal/ui"
)

func formatStatus(statuses []docker.ServiceStatus) string {
	var b strings.Builder
	var ports []string
	for _, s := range statuses {
		state := s.State
		if s.Health != "" {
			state = s.Health
		}
		b.WriteString(fmt.Sprintf("%-14s %s\n", s.Name, state))
		ports = append(ports, s.Ports...)
	}
	if len(ports) > 0 {
		b.WriteString("\n" + ui.Title("Ports") + "\n\n")
		for _, p := range ports {
			b.WriteString(p + "\n")
		}
	}
	return b.String()
}

func init() {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display project status",
		RunE: func(c *cobra.Command, args []string) error {
			client := docker.New()
			statuses, err := client.Status()
			if err != nil {
				return err
			}
			fmt.Print(formatStatus(statuses))
			return nil
		},
	}
	RootCmd.AddCommand(cmd)
}
