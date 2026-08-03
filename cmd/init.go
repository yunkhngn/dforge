package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yunkhngn/dforge/internal/detect"
	"github.com/yunkhngn/dforge/internal/fsutil"
	"github.com/yunkhngn/dforge/internal/generator"
	"github.com/yunkhngn/dforge/internal/ui"
)

func runInit(root string, chosen detect.Framework, assumeYes bool) error {
	files, err := generator.Render(chosen)
	if err != nil {
		return err
	}

	// Pre-check every target so we never leave a half-generated project:
	// either all files are written or none are.
	if !assumeYes {
		var existing []string
		for name := range files {
			if fsutil.Exists(filepath.Join(root, name)) {
				existing = append(existing, name)
			}
		}
		if len(existing) > 0 {
			sort.Strings(existing)
			return fmt.Errorf("refusing to overwrite existing files: %s (use --force to overwrite)",
				strings.Join(existing, ", "))
		}
	}

	for name, content := range files {
		path := filepath.Join(root, name)
		if err := fsutil.WriteWithBackup(path, []byte(content)); err != nil {
			return err
		}
		fmt.Println(ui.Success(ui.OK + " wrote " + name))
	}
	return nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Detect the project and generate Docker files",
		RunE: func(c *cobra.Command, args []string) error {
			force, _ := c.Flags().GetBool("force")
			yes, _ := c.Flags().GetBool("yes")
			flagFw, _ := c.Flags().GetString("framework")

			cands, err := detect.Detect(".")
			if err != nil {
				return err
			}
			chosen, err := pickFramework(cands, detect.Framework(flagFw))
			if err != nil {
				return err
			}
			return runInit(".", chosen, force || yes)
		},
	}
	cmd.Flags().String("framework", "", "override detected framework")
	RootCmd.AddCommand(cmd)
}

func pickFramework(cands []detect.Candidate, override detect.Framework) (detect.Framework, error) {
	if override != "" {
		return override, nil
	}
	switch len(cands) {
	case 0:
		return "", fmt.Errorf("no framework detected; pass --framework")
	case 1:
		return cands[0].Framework, nil
	default:
		labels := make([]string, len(cands))
		for i, c := range cands {
			labels[i] = string(c.Framework)
		}
		idx, err := ui.Select("Select framework:", labels)
		if err != nil {
			return "", err
		}
		return cands[idx].Framework, nil
	}
}
