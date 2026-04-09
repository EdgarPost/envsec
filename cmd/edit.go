package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/EdgarPost/envsec/resolver"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open env file in $EDITOR",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		result, err := resolver.Resolve(dir)
		if err != nil {
			return err
		}

		paths, err := st.Path(result.ProjectKey, result.Subpath)
		if err != nil {
			return err
		}
		if len(paths) == 0 {
			return fmt.Errorf("no env file found — run 'envsec init' first")
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		// Edit the most specific file (last in the list)
		target := paths[len(paths)-1]
		c := exec.Command(editor, target)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}
