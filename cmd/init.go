package cmd

import (
	"fmt"
	"os"

	"github.com/EdgarPost/envsec/resolver"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Register current directory as a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		result, err := resolver.Resolve(dir)
		if err != nil {
			return err
		}

		if err := st.Init(result.ProjectKey, result.Subpath); err != nil {
			return err
		}

		paths, _ := st.Path(result.ProjectKey, result.Subpath)
		if len(paths) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "Initialized %s\n", paths[len(paths)-1])
		}

		return nil
	},
}
