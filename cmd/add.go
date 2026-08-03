package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yunkhngn/dforge/internal/compose"
	"github.com/yunkhngn/dforge/internal/services"
	"github.com/yunkhngn/dforge/internal/ui"
)

func runAdd(path, name string) error {
	svc, ok := services.Get(name)
	if !ok {
		return fmt.Errorf("unknown service %q (available: %v)", name, services.Names())
	}
	f, err := compose.Load(path)
	if err != nil {
		return err
	}
	if err := f.AddService(svc); err != nil {
		return err
	}
	if err := f.Save(); err != nil {
		return err
	}
	fmt.Println(ui.Success(ui.OK + " added " + name))
	return nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "add <service>",
		Short: "Add a predefined service to compose.yaml",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return runAdd("compose.yaml", args[0])
		},
	}
	RootCmd.AddCommand(cmd)
}
