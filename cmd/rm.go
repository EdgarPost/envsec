package cmd

import (
	"os"

	"github.com/EdgarPost/envsec/resolver"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm KEY",
	Short: "Remove an environment variable",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		result, err := resolver.Resolve(dir)
		if err != nil {
			return err
		}

		return st.Remove(result.ProjectKey, result.Subpath, args[0])
	},
}
