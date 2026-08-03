package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yunkhngn/dforge/internal/docker"
	"github.com/yunkhngn/dforge/internal/ui"
)

func runClean(cl docker.Client, force, assumeYes bool) error {
	if !force && !assumeYes {
		ok, err := ui.Confirm("Remove dangling images, unused volumes, networks, and build cache?", false)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("clean aborted")
		}
	}
	if err := cl.Clean(force); err != nil {
		return err
	}
	fmt.Println(ui.Success(ui.OK + " cleaned Docker resources"))
	return nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean unused Docker resources",
		RunE: func(c *cobra.Command, args []string) error {
			force, _ := c.Flags().GetBool("force")
			yes, _ := c.Flags().GetBool("yes")
			return runClean(docker.New(), force, yes)
		},
	}
	cmd.Flags().BoolP("force", "f", false, "force clean without confirmation")
	cmd.Flags().BoolP("yes", "y", false, "assume yes for confirmation prompt")
	RootCmd.AddCommand(cmd)
}
