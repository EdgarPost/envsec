package cmd

import (
	"fmt"
	"os"

	"github.com/EdgarPost/envsec/resolver"
	"github.com/spf13/cobra"
)

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print resolved env file path(s)",
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

		for _, p := range paths {
			fmt.Fprintln(cmd.OutOrStdout(), p)
		}

		return nil
	},
}
