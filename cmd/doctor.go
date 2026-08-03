package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yunkhngn/dforge/internal/doctor"
	"github.com/yunkhngn/dforge/internal/ui"
)

func formatDoctorReport(r doctor.Report) string {
	var b strings.Builder
	for _, res := range r.Results {
		switch {
		case res.Passed:
			b.WriteString(ui.Success(ui.OK+" "+res.Title) + "\n")
		case res.Severity == doctor.Error:
			b.WriteString(ui.Error(ui.ERR+" "+res.Title) + "\n")
		default:
			b.WriteString(ui.Warning(ui.WARN+" "+res.Title) + "\n")
		}
	}
	b.WriteString("\n" + ui.Title(fmt.Sprintf("Score: %d/100", r.Score)) + "\n")
	return b.String()
}

func init() {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Analyze Docker configuration and score it",
		RunE: func(c *cobra.Command, args []string) error {
			state := doctor.CollectState(".")
			rep := doctor.Run(state, doctor.DefaultChecks())
			fmt.Print(formatDoctorReport(rep))
			return nil
		},
	}
	RootCmd.AddCommand(cmd)
}
