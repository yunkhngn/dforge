package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yunkhngn/dforge/internal/env"
	"github.com/yunkhngn/dforge/internal/ui"
)

func formatEnvReport(r env.Report) string {
	var b strings.Builder
	if r.MissingFile {
		b.WriteString(ui.Warning(ui.WARN+" .env file not found") + "\n")
	}
	section := func(label string, names []string, warn bool) {
		if len(names) == 0 {
			return
		}
		b.WriteString(ui.Title(label) + "\n")
		for _, n := range names {
			line := "  " + n
			if warn {
				line = ui.Warning(line)
			}
			b.WriteString(line + "\n")
		}
	}
	section("Missing variables", r.Missing, true)
	section("Extra variables", r.Extra, true)
	section("Duplicate keys", r.Duplicates, true)
	section("Empty values", r.Empty, true)
	if len(r.Missing)+len(r.Extra)+len(r.Duplicates)+len(r.Empty) == 0 && !r.MissingFile {
		b.WriteString(ui.Success(ui.OK+" .env is valid") + "\n")
	}
	return b.String()
}

func init() {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Validate environment files (names only, never values)",
		RunE: func(c *cobra.Command, args []string) error {
			r, err := env.Compare(".env", ".env.example")
			if err != nil {
				return err
			}
			fmt.Print(formatEnvReport(r))
			return nil
		},
	}
	RootCmd.AddCommand(cmd)
}
