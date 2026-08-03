package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yunkhngn/dforge/internal/compose"
	"github.com/yunkhngn/dforge/internal/ui"
)

func runRemove(path, name string) error {
	f, err := compose.Load(path)
	if err != nil {
		return err
	}
	if err := f.RemoveService(name); err != nil {
		return err
	}
	if err := f.Save(); err != nil {
		return err
	}
	fmt.Println(ui.Success(ui.OK + " removed " + name))
	return nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "remove <service>",
		Short: "Remove a service from compose.yaml",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runRemove("compose.yaml", args[0])
		},
	}
	RootCmd.AddCommand(cmd)
}
